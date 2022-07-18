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
	"path"
	"strings"
	"sync"

	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/go-git/go-git/v5"
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

	r.log.V(logger.DebugLevel).Info("Reconciling users...")
	if err := r.reconcileUsers(ctx); err != nil {
		return fmt.Errorf("error reconciling users: %w", err)
	}

	return nil
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

	// get all clusters accessible for the user
	clusters, err := r.clusterAccesses.GetClustersAccessibleByUserId(ctx, user.ID())
	if err != nil {
		return err
	}

	username := strings.Split(user.Email, "@")[0]
	for _, cluster := range clusters {
		path := path.Join(r.config.LocalDirectory, cluster.Cluster.Name, fmt.Sprintf("%s.yaml", username))

		// Remove bindings for deleted users
		if user.IsDeleted() {
			r.log.V(logger.DebugLevel).Info("User is deleted. Removing bindings", "user", user.Email)
			if _, err := os.Stat(path); err == nil {
				if err := os.Remove(path); err != nil {
					r.log.Error(err, "failed to remove file", "file", path, "user", user.Email, "cluster", cluster.Cluster.Name)
				}
			}
		}
	}
	return nil
}
