package git

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitClient struct {
	config         *GitConfig
	log            logger.Logger
	localDirectory string
	repo           *git.Repository
}

func NewGitClient(config *GitConfig) (*GitClient, error) {
	dir, err := os.MkdirTemp("", "m8-k8sauthz")
	if err != nil {
		return nil, err
	}
	return &GitClient{config: config, localDirectory: dir, log: logger.WithName("git-client")}, nil
}

func (c *GitClient) GetLocalDirectory() string {
	return c.localDirectory
}

func (c *GitClient) Clone(ctx context.Context) error {
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

func (c *GitClient) Pull(ctx context.Context) error {
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

func (c *GitClient) AddAll(ctx context.Context) error {
	w, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	return nil
}

func (c *GitClient) Add(ctx context.Context, filePath string) error {
	w, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	if err := w.AddWithOptions(&git.AddOptions{Path: filePath}); err != nil {
		return fmt.Errorf("failed to add file `%s`: %w", filePath, err)
	}

	return nil
}

func (c *GitClient) Commit(ctx context.Context, msg string) error {
	w, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	commit, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  users.GitRepoReconcilerUser.User.Name,
			Email: users.GitRepoReconcilerUser.User.Email,
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

func (c *GitClient) AddAndCommit(ctx context.Context, filePath, msg string) error {
	if err := c.Add(ctx, filePath); err != nil {
		return err
	}

	if err := c.Commit(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (c *GitClient) AddAllAndCommit(ctx context.Context, msg string) error {
	if err := c.AddAll(ctx); err != nil {
		return err
	}

	if err := c.Commit(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (c *GitClient) Push(ctx context.Context) error {
	po, err := c.config.getPushOptions()
	if err != nil {
		return err
	}
	if err := c.repo.PushContext(ctx, po); err != nil {
		return err
	}
	return nil
}
