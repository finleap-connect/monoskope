// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8sauthz

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	mk8s "github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"k8s.io/cli-runtime/pkg/printers"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	rbac "k8s.io/api/rbac/v1"
)

const (
	defaultDirectoryMode = 0755
)

// GitRepoReconciler reconciles the resources within the target repo to match the expected state.
type GitRepoReconciler struct {
	log             logger.Logger
	config          *ReconcilerConfig
	users           repositories.UserRepository
	clusterAccesses repositories.ClusterAccessRepository
	gitRepo         *git.Repository
	mutex           sync.Mutex
}

// NewGitRepoReconciler creates a new GitRepoReconciler configured via the given config.
func NewGitRepoReconciler(
	config *ReconcilerConfig,
	userRepo repositories.UserRepository,
	clusterAccessRepo repositories.ClusterAccessRepository,
	gitRepo *git.Repository,
) *GitRepoReconciler {
	remote, _ := gitRepo.Remote("origin")
	return &GitRepoReconciler{logger.WithName("GitRepoReconciler").WithValues("remote", remote.Config().URLs), config, userRepo, clusterAccessRepo, gitRepo, sync.Mutex{}}
}

func (r *GitRepoReconciler) Reconcile(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.log.V(logger.DebugLevel).Info("Started reconciling..")

	r.log.V(logger.DebugLevel).Info("Fetching git repo..")
	if err := r.gitRepo.Fetch(&git.FetchOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error fetching latest changes from repository: %w", err)
	}

	w, err := r.gitRepo.Worktree()
	if err != nil {
		return err
	}

	r.log.V(logger.DebugLevel).Info("Reconciling users...")
	if err := r.reconcileUsers(ctx, w); err != nil {
		return fmt.Errorf("error reconciling users: %w", err)
	}

	return nil
}

func (r *GitRepoReconciler) reconcileUsers(ctx context.Context, w *git.Worktree) error {
	// Get all users including deleted ones
	users, err := r.users.AllWith(ctx, true)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// reconcile each user
	for _, user := range users {
		if err := r.reconcileUser(ctx, user, w); err != nil {
			return err
		}
	}

	return nil
}

func (r *GitRepoReconciler) reconcileUser(ctx context.Context, user *projections.User, w *git.Worktree) error {
	r.log.V(logger.DebugLevel).Info("Reconciling user...", "user", user.Email)

	// get all clusterAccesses accessible for the user
	clusterAccesses, err := r.clusterAccesses.GetClustersAccessibleByUserId(ctx, user.ID())
	if err != nil {
		return err
	}

	sanitizedName, err := mk8s.GetK8sName(user.Name)
	if err != nil {
		return err
	}
	for _, clusterAccess := range clusterAccesses {
		path := filepath.Join(r.config.LocalDirectory, clusterAccess.Cluster.Name, sanitizedName)

		if user.IsDeleted() {
			// Remove bindings for deleted users
			r.log.V(logger.DebugLevel).Info("User is deleted. Removing bindings...", "user", user.Email)
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				if err := os.Remove(path); err != nil {
					r.log.Error(err, "Failed to remove path", "path", path, "user", user.Email, "cluster", clusterAccess.Cluster.Name)
				}
			}
		} else {
			// Create user sub dir
			if err := os.MkdirAll(path, defaultDirectoryMode); err != nil {
				r.log.Error(err, "Failed to create path", "path", path)
				return err
			}

			// Reconcile bindings for existing users
			r.log.V(logger.DebugLevel).Info("User exists. Reconciling bindings...", "user", user.Email)
			for _, role := range clusterAccess.Roles {
				if clusterRole, ok := r.config.Mappings[role]; ok {
					if err := r.createClusterRoleBinding(ctx, path, clusterRole, sanitizedName, w); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (r *GitRepoReconciler) createClusterRoleBinding(ctx context.Context, dir, clusterRoleName, sanitizedUsername string, w *git.Worktree) error {
	filePath := filepath.Join(dir, fmt.Sprintf("%s.yaml", clusterRoleName))
	r.log.V(logger.DebugLevel).Info("Creating cluster role binding...", "path", filePath)

	crb := new(rbac.ClusterRoleBinding)
	crb.Subjects = append(crb.Subjects, rbac.Subject{
		Name: fmt.Sprintf("%s%s", r.config.UsernamePrefix, sanitizedUsername),
		Kind: "User",
	})
	crb.ObjectMeta.Name = fmt.Sprintf("%s-%s", sanitizedUsername, clusterRoleName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	err = new(printers.YAMLPrinter).PrintObj(crb, file)
	if err != nil {
		return err
	}
	r.log.V(logger.DebugLevel).Info("Committing cluster role binding...", "clusterRole", clusterRoleName, "path", filePath)

	_, err = w.Add(filePath)
	if err != nil {
		return err
	}

	commit, err := w.Commit(fmt.Sprintf("Add ClusterRoleBinding for user %s", sanitizedUsername), &git.CommitOptions{
		Author: &object.Signature{
			Name:  users.GitRepoReconcilerUser.User.Name,
			Email: users.GitRepoReconcilerUser.User.Email,
			When:  time.Now().UTC(),
		},
	})
	if err != nil {
		return err
	}

	_, err = r.gitRepo.CommitObject(commit)
	if err != nil {
		return err
	}

	return nil
}
