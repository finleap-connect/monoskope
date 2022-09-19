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

package git

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitClient is a small wrapper around go-git library methods.
// It provides clone/pull/push options as needed and has a simpler interface.
type GitClient struct {
	config         *GitConfig
	log            logger.Logger
	localDirectory string
	repo           *git.Repository
}

// NewGitClient creates a new go-git library wrapper
func NewGitClient(config *GitConfig) (*GitClient, error) {
	dir, err := os.MkdirTemp("", "m8-k8sauthz")
	if err != nil {
		return nil, err
	}
	return &GitClient{config: config, localDirectory: dir, log: logger.WithName("git-client")}, nil
}

// GetLocalDirectory returns the directory where the repo is cloned to
func (c *GitClient) GetLocalDirectory() string {
	return c.localDirectory
}

// Clone clones the configured repo
func (c *GitClient) Clone(ctx context.Context) error {
	var cancel context.CancelFunc
	if c.config.Timeout != nil {
		ctx, cancel = context.WithTimeout(ctx, *c.config.Timeout)
		defer cancel()
	}

	cloneOptions, err := c.config.getCloneOptions()
	if err != nil {
		return err
	}

	repo, err := git.PlainCloneContext(ctx, c.localDirectory, false, cloneOptions)
	if err != nil {
		return err
	}
	c.repo = repo
	return nil
}

// Pull pulls the latest changes of the configured repo
func (c *GitClient) Pull(ctx context.Context) error {
	var cancel context.CancelFunc
	if c.config.Timeout != nil {
		ctx, cancel = context.WithTimeout(ctx, *c.config.Timeout)
		defer cancel()
	}

	w, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	po, err := c.config.getPullOptions()
	if err != nil {
		return err
	}

	if err := w.PullContext(ctx, po); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error pulling latest changes from repository: %w", err)
	}
	return nil
}

// AddAll stages all changes in the working directory
func (c *GitClient) AddAll(_ context.Context) error {
	w, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	for file, s := range status {
		if s.Worktree == git.Deleted {
			_, err = w.Add(file)
			if err != nil {
				return err
			}
		}
	}

	c.log.V(logger.DebugLevel).Info("Status after add all.", "status", status.String())

	return nil
}

// Add stages the given file
func (c *GitClient) Add(_ context.Context, filePath string) error {
	w, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	if err := w.AddWithOptions(&git.AddOptions{Path: filePath}); err != nil {
		return fmt.Errorf("failed to add file `%s`: %w", filePath, err)
	}

	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	c.log.V(logger.DebugLevel).Info("Status after add.", "filePath", filePath, "status", status.String())

	return nil
}

// Commit commits all changes with the given commit message
func (c *GitClient) Commit(_ context.Context, msg string) error {
	w, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	// Don't create empty commits
	if status.IsClean() {
		return nil
	}

	c.log.V(logger.DebugLevel).Info("Creating commit.", "message", msg, "status", status.String())
	commit, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  c.config.Author.Name,
			Email: c.config.Author.Email,
			When:  time.Now().UTC(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	_, err = c.repo.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("failed to commit object: %w", err)
	}

	return nil
}

// AddAndCommit stages the given file an creates a commit
func (c *GitClient) AddAndCommit(ctx context.Context, filePath, msg string) error {
	if err := c.Add(ctx, filePath); err != nil {
		return err
	}

	if err := c.Commit(ctx, msg); err != nil {
		return err
	}

	return nil
}

// AddAndCommit stages all changes an creates a commit
func (c *GitClient) AddAllAndCommit(ctx context.Context, msg string) error {
	if err := c.AddAll(ctx); err != nil {
		return err
	}

	if err := c.Commit(ctx, msg); err != nil {
		return err
	}

	return nil
}

// Push pushes all outstanding changes
func (c *GitClient) Push(ctx context.Context) error {
	var cancel context.CancelFunc
	if c.config.Timeout != nil {
		ctx, cancel = context.WithTimeout(ctx, *c.config.Timeout)
		defer cancel()
	}

	po, err := c.config.getPushOptions()
	if err != nil {
		return err
	}
	if err := c.repo.PushContext(ctx, po); err != nil {
		return err
	}
	return nil
}

// Close cleans up the clone directory
func (c *GitClient) Close() error {
	return os.RemoveAll(c.localDirectory)
}
