package k8sauthz

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

type TestEnv struct {
	tempDir string
	repoDir string
	gitRepo *git.Repository
}

func NewTestEnv() (*TestEnv, error) {
	env := &TestEnv{}

	// Temp dir to clone the repository
	dir, err := os.MkdirTemp("", "m8-git-repo-reconciler")
	if err != nil {
		return nil, err
	}
	env.tempDir = dir
	env.repoDir = filepath.Join(dir, "repo")
	repoOriginDir := filepath.Join(dir, "origin")

	r, err := git.PlainInit(repoOriginDir, false)
	if env.err(err) != nil {
		return nil, err
	}

	f, err := os.Create(filepath.Join(repoOriginDir, ".gitignore"))
	if env.err(err) != nil {
		return nil, err
	}
	f.Close()

	fRelName, err := filepath.Rel(repoOriginDir, f.Name())
	if env.err(err) != nil {
		return nil, err
	}

	wt, err := r.Worktree()
	if env.err(err) != nil {
		return nil, err
	}
	_, err = wt.Add(fRelName)
	if env.err(err) != nil {
		return nil, err
	}
	_, err = wt.Commit("init", &git.CommitOptions{})
	if env.err(err) != nil {
		return nil, err
	}

	r, err = git.PlainClone(env.repoDir, false, &git.CloneOptions{
		URL: repoOriginDir,
	})
	if env.err(err) != nil {
		return nil, err
	}
	env.gitRepo = r
	return env, nil
}

func (env *TestEnv) err(err error) error {
	if err != nil {
		_ = env.Shutdown()
		return err
	}
	return nil
}

func (env *TestEnv) Shutdown() error {
	os.RemoveAll(env.tempDir) // clean up

	return nil
}
