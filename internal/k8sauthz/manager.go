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
	"time"

	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/git"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type Manager struct {
	log                     logger.Logger
	gitClient               *git.GitClient
	reconciler              *GitRepoReconciler
	userRepository          repositories.UserRepository
	clusterAccessRepository repositories.ClusterAccessRepository
	eg                      errgroup.Group
	quitChannels            []chan struct{}
}

func NewManager(userRepository repositories.UserRepository, clusterAccessRepository repositories.ClusterAccessRepository) *Manager {
	m := &Manager{
		log:                     logger.WithName("GitRepoManager"),
		userRepository:          userRepository,
		clusterAccessRepository: clusterAccessRepository,
	}
	m.userRepository.RegisterObserver(m)
	return m
}

func (m *Manager) Run(ctx context.Context, conf *Config) error {
	m.log.Info("Starting reconciliation loops...")
	gitClient, err := git.NewGitClient(conf.Repository)
	if err != nil {
		return err
	}
	m.gitClient = gitClient

	// Clone repo
	m.log.Info("Cloning repo...", "url", conf.Repository.URL)
	if err := gitClient.Clone(ctx); err != nil {
		return err
	}

	m.log.Info("Configuring reconciler...", "url", conf.Repository.URL)
	m.reconciler = NewGitRepoReconciler(conf, m.userRepository, m.clusterAccessRepository, gitClient)

	// initial reconcile
	if err := m.reconciler.Reconcile(ctx); err != nil {
		m.log.Error(err, "Failed running reconciliation loop.")
	}

	// schedule reconcile loop
	ticker := time.NewTicker(*conf.Interval)
	quit := make(chan struct{})
	m.quitChannels = append(m.quitChannels, quit)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.eg.Go(func() error {
					err := m.reconciler.Reconcile(ctx)
					if err != nil {
						m.log.Error(err, "Failed running reconciliation loop.")
					}
					return err
				})
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

func (m *Manager) Notify(ctx context.Context, u *projections.User) {
	m.log.V(logger.DebugLevel).Info("Received notification from repo for user.", "user", u.Email)
	if err := m.reconciler.ReconcileUser(ctx, u); err != nil {
		m.log.Error(err, "Failed to reconcile user.")
	}
}

func (m *Manager) Close() error {
	m.log.Info("Stopping reconciliation loops...")
	for _, c := range m.quitChannels {
		close(c)
	}

	m.log.Info("Waiting reconciliations to finish...")
	if err := m.eg.Wait(); err != nil {
		m.log.Error(err, "Encountered errors while reconciling.")
	}

	m.log.Info("Cleaning up...")
	if err := m.gitClient.Close(); err != nil {
		m.log.Error(err, "Encountered errors cleaning up.")
	}

	return m.eg.Wait()
}
