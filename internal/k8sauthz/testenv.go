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

	"github.com/finleap-connect/monoskope/pkg/git"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type TestEnv struct {
	tempDir       string
	repoOriginDir string
	gitClient     *git.GitClient
}

func NewTestEnv() (*TestEnv, error) {
	env := &TestEnv{}

	// Temp dir to clone the repository
	dir, err := os.MkdirTemp("", "m8-git-repo-reconciler")
	if err != nil {
		return nil, err
	}
	env.tempDir = dir
	env.repoOriginDir = filepath.Join(dir, "origin")

	r, err := gogit.PlainInit(env.repoOriginDir, false)
	if env.err(err) != nil {
		return nil, err
	}

	f, err := os.Create(filepath.Join(env.repoOriginDir, ".gitignore"))
	if env.err(err) != nil {
		return nil, err
	}
	f.Close()

	fRelName, err := filepath.Rel(env.repoOriginDir, f.Name())
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
	_, err = wt.Commit("init", &gogit.CommitOptions{
		Author: &object.Signature{
			Name: "testenv",
		},
	})
	if env.err(err) != nil {
		return nil, err
	}

	gitConfig, err := git.NewGitConfig(env.repoOriginDir, &git.GitAuthor{Name: "test", Email: "test@monoskope.io"})
	if env.err(err) != nil {
		return nil, err
	}

	gitClient, err := git.NewGitClient(gitConfig)
	if env.err(err) != nil {
		return nil, err
	}

	err = gitClient.Clone(context.Background())
	if env.err(err) != nil {
		return nil, err
	}

	env.gitClient = gitClient
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
	if err := filepath.Walk(env.gitClient.GetLocalDirectory(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path)
			return nil
		}); err != nil {
		return err
	}
	if err := os.RemoveAll(env.tempDir); err != nil {
		return err
	}
	if err := env.gitClient.Close(); err != nil {
		return err
	}
	return nil
}
