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

	api_projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	mk8s "github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"k8s.io/cli-runtime/pkg/printers"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	return &GitRepoReconciler{logger.WithName("GitRepoReconciler"), config, userRepo, clusterAccessRepo, gitRepo, sync.Mutex{}}
}

func (r *GitRepoReconciler) Reconcile(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.log.Info("Started reconciling...")

	r.log.Info("Fetching git repo..")
	if err := r.gitRepo.Fetch(&git.FetchOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error fetching latest changes from repository: %w", err)
	}

	r.log.Info("Reconciling users...")
	if err := r.reconcileUsers(ctx); err != nil {
		return fmt.Errorf("error reconciling users: %w", err)
	}

	r.log.Info("Pushing changes to git repo...")
	err := r.gitRepo.PushContext(ctx, &git.PushOptions{})
	if err == nil {
		r.log.Info("Reconciling finished.")

	} else {
		r.log.Error(err, "Reconciling finished with errors.")
	}
	return err
}

func (r *GitRepoReconciler) reconcileUsers(ctx context.Context) error {
	// Get all users including deleted ones
	users, err := r.users.AllWith(ctx, true)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// reconcile each user
	for _, user := range users {
		if err := r.reconcileUser(ctx, user); err != nil {
			return err
		}
	}

	return nil
}

func (r *GitRepoReconciler) reconcileUser(ctx context.Context, user *projections.User) error {
	r.log.V(logger.DebugLevel).Info("Reconciling user...", "user", user.Email)

	// get all clusterAccesses accessible for the user
	clusterAccesses, err := r.clusterAccesses.GetClustersAccessibleByUserIdV2(ctx, user.ID())
	if err != nil {
		return err
	}

	// sanitize user name
	sanitizedName, err := mk8s.GetK8sName(user.Name)
	if err != nil {
		return err
	}

	// Remove bindings for deleted users
	if user.IsDeleted() {
		return r.removeClusterRolesForUser(user, sanitizedName, clusterAccesses)
	}

	// Create/reconcile bindings for existing users
	return r.createClusterRolesForUser(ctx, user, sanitizedName, clusterAccesses)
}

func (r *GitRepoReconciler) createClusterRolesForUser(ctx context.Context, user *projections.User, sanitizedName string, clusterAccesses []*api_projections.ClusterAccessV2) error {
	r.log.V(logger.DebugLevel).Info("Reconciling bindings...", "user", user.Email)

	for _, clusterAccess := range clusterAccesses {
		path := filepath.Join(r.config.LocalDirectory, clusterAccess.Cluster.Name, sanitizedName)

		// Create user sub dir
		if err := os.MkdirAll(path, defaultDirectoryMode); err != nil {
			r.log.Error(err, "Failed to create path", "path", path)
			return err
		}

		// Reconcile bindings for existing users
		for _, clusterAccessRole := range clusterAccess.ClusterRoles {
			if clusterRole := getClusterRoleMapping(r.config.Mappings, clusterAccessRole.Scope.String(), clusterAccessRole.Role); clusterRole != "" {
				if err := r.createClusterRoleBinding(ctx, path, clusterRole, user, sanitizedName); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *GitRepoReconciler) removeClusterRolesForUser(user *projections.User, sanitizedName string, clusterAccesses []*api_projections.ClusterAccessV2) error {
	r.log.V(logger.DebugLevel).Info("User is deleted. Removing bindings...", "user", user.Email)

	w, err := r.gitRepo.Worktree()
	if err != nil {
		return err
	}

	for _, clusterAccess := range clusterAccesses {
		path := filepath.Join(r.config.LocalDirectory, clusterAccess.Cluster.Name, sanitizedName)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			// remove folder with bindings if present
			if err := os.Remove(path); err != nil {
				r.log.Error(err, "Failed to remove path", "path", path)
				return err
			}

			r.log.V(logger.DebugLevel).Info("Committing removal of cluster role bindings...", "path", path)
			if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
				return err
			}

			// commit
			commit, err := w.Commit(fmt.Sprintf("Removing ClusterRoleBindings for user %s", user.Email), &git.CommitOptions{
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
		}
	}
	return nil
}

func (r *GitRepoReconciler) createClusterRoleBinding(ctx context.Context, dir, clusterRoleName string, user *projections.User, sanitizedName string) error {
	filePath := filepath.Join(dir, fmt.Sprintf("%s.yaml", clusterRoleName))
	relFilePath, err := filepath.Rel(r.config.LocalDirectory, filePath)
	if err != nil {
		return err
	}

	r.log.V(logger.DebugLevel).Info("Creating cluster role binding...", "path", filePath)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// create new cluster role binding for user and clusterrole
	crb := mk8s.NewClusterRoleBinding(clusterRoleName, sanitizedName, r.config.UsernamePrefix, map[string]string{
		"user": user.Email,
	})

	err = new(printers.YAMLPrinter).PrintObj(crb, file)
	if err != nil {
		return err
	}

	r.log.V(logger.DebugLevel).Info("Committing cluster role binding...", "path", filePath)
	w, err := r.gitRepo.Worktree()
	if err != nil {
		return err
	}

	if err := w.AddWithOptions(&git.AddOptions{Path: relFilePath}); err != nil {
		return err
	}

	// commit
	commit, err := w.Commit(fmt.Sprintf("Add ClusterRoleBinding for user %s", user.Email), &git.CommitOptions{
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
