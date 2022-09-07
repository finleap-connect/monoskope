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
	"strings"
	"sync"

	api_projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/git"
	mk8s "github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/finleap-connect/monoskope/pkg/logger"
	gogit "github.com/go-git/go-git/v5"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	defaultDirectoryMode = 0755
)

// GitRepoReconciler reconciles the resources within the target repo to match the expected state.
type GitRepoReconciler struct {
	log             logger.Logger
	config          *Config
	users           repositories.UserRepository
	clusterAccesses repositories.ClusterAccessRepository
	gitClient       *git.GitClient
	dir             string
	mutex           sync.Mutex
}

// NewGitRepoReconciler creates a new GitRepoReconciler configured via the given config.
func NewGitRepoReconciler(
	config *Config,
	userRepo repositories.UserRepository,
	clusterAccessRepo repositories.ClusterAccessRepository,
	gitClient *git.GitClient,
) *GitRepoReconciler {
	return &GitRepoReconciler{logger.WithName("GitRepoReconciler"), config, userRepo, clusterAccessRepo, gitClient, filepath.Join(gitClient.GetLocalDirectory(), config.SubDir), sync.Mutex{}}
}

func (r *GitRepoReconciler) Reconcile(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.log.Info("Started reconciling...")

	r.log.Info("Pulling latest changes..")
	if err := r.gitClient.Pull(ctx); err != nil {
		return err
	}

	r.log.Info("Reconciling...")
	if err := r.reconcileUsers(ctx); err != nil {
		return fmt.Errorf("error reconciling users: %w", err)
	}

	r.log.Info("Pushing changes to git repo...")
	if err := r.gitClient.Push(ctx); err != nil && err != gogit.NoErrAlreadyUpToDate {
		r.log.Error(err, "Reconciling finished with errors.")
		return err
	}
	return nil
}

func (r *GitRepoReconciler) ReconcileUser(ctx context.Context, user *projections.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.log.Info("Start reconciling user...", "user", user.Email)

	r.log.Info("Pulling latest changes..")
	if err := r.gitClient.Pull(ctx); err != nil {
		return err
	}

	r.log.Info("Reconciling user...")
	if err := r.reconcileUser(ctx, user); err != nil {
		return fmt.Errorf("error reconciling user: %w", err)
	}

	r.log.Info("Pushing changes to git repo...")
	if err := r.gitClient.Push(ctx); err != nil {
		r.log.Error(err, "Reconciling finished with errors.")
		return err
	}
	return nil
}

// removeAll cleans up a directory keeping all hidden files and directories
func removeAll(path string) error {
	matches, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return fmt.Errorf("failed to clean local dir: %w", err)
	}
	for _, match := range matches {
		if strings.HasPrefix(filepath.Base(match), ".") {
			continue
		}
		if err := os.RemoveAll(match); err != nil && err != os.ErrNotExist {
			return fmt.Errorf("failed to clean local dir: %w", err)
		}
	}
	return nil
}

func (r *GitRepoReconciler) reconcileUsers(ctx context.Context) error {
	// Clean
	if err := removeAll(r.dir); err != nil {
		return err
	}
	if err := r.gitClient.AddAll(ctx); err != nil {
		return err
	}

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
		return fmt.Errorf("failed to get cluster access: %w", err)
	}

	// sanitize user name
	sanitizedName, err := mk8s.GetK8sName(user.Name)
	if err != nil {
		return err
	}

	// Do not handle deleted users
	if user.IsDeleted() {
		return nil
	}

	// Create/reconcile bindings for existing users
	return r.createClusterRolesForUser(ctx, user, sanitizedName, clusterAccesses)
}

func (r *GitRepoReconciler) createClusterRolesForUser(ctx context.Context, user *projections.User, sanitizedName string, clusterAccesses []*api_projections.ClusterAccessV2) error {
	r.log.V(logger.DebugLevel).Info("Reconciling bindings...", "user", user.Email, "clusterAccesses", len(clusterAccesses))

	for _, clusterAccess := range clusterAccesses {
		r.log.V(logger.DebugLevel).Info("Reconciling binding...", "user", user.Email, "cluster", clusterAccess.Cluster.Name)
		path := filepath.Join(r.dir, clusterAccess.Cluster.Name, sanitizedName)

		// Create user sub dir
		if err := os.MkdirAll(path, defaultDirectoryMode); err != nil {
			r.log.Error(err, "Failed to create path", "path", path)
			return fmt.Errorf("failed to mkdir: %w", err)
		}

		// Reconcile bindings for existing users
		for _, clusterAccessRole := range clusterAccess.ClusterRoles {
			if clusterRole := r.config.getClusterRoleMapping(clusterAccessRole.Scope.String(), clusterAccessRole.Role); clusterRole != "" {
				if err := r.createClusterRoleBinding(ctx, path, clusterRole, user, clusterAccess.Cluster, sanitizedName); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *GitRepoReconciler) createClusterRoleBinding(ctx context.Context, dir, clusterRoleName string, user *projections.User, cluster *api_projections.Cluster, sanitizedName string) error {
	filePath := filepath.Join(dir, fmt.Sprintf("%s.yaml", clusterRoleName))
	relFilePath, err := filepath.Rel(r.gitClient.GetLocalDirectory(), filePath)
	if err != nil {
		return fmt.Errorf("failed to generate relative file path: %w", err)
	}

	r.log.V(logger.DebugLevel).Info("Creating cluster role binding...", "path", filePath)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file `%s`: %w", filePath, err)
	}

	// create new cluster role binding for user and cluster role
	crb := mk8s.NewClusterRoleBinding(clusterRoleName, sanitizedName, r.config.UsernamePrefix, map[string]string{
		"user":    user.Email,
		"cluster": cluster.Name,
	})

	err = new(printers.YAMLPrinter).PrintObj(crb, file)
	if err != nil {
		return fmt.Errorf("failed to print cluster role binding as yaml: %w", err)
	}

	r.log.V(logger.DebugLevel).Info("Committing cluster role binding...", "path", filePath)
	err = r.gitClient.AddAndCommit(ctx, relFilePath, fmt.Sprintf("Add ClusterRoleBinding for user %s", user.Email))
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}
