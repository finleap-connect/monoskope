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

package k8sauthzreactor

import (
	"fmt"
	"os"
	"time"

	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"gopkg.in/yaml.v2"
)

var (
	DefaultTimeout = 60 * time.Second
)

// GitRepoReconcilerConfig is the configuration for the GitRepoReconciler.
type GitRepoReconcilerConfig struct {
	Repositories []*GitRepository      `yaml:"repositories"`
	Mappings     []*ClusterRoleMapping `yaml:"mappings"`
	cloneOptions *git.CloneOptions
}

// GitBasicAuth is used to authenticate towards a Git repository over HTTPS using basic access authentication.
type GitBasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// GitSSHAuth is used to authenticate towards a Git repository over SSH. With the respective private key of the SSH key pair, and the host keys of the Git repository.
type GitSSHAuth struct {
	PrivateKeyPath string `yaml:"privateKeyPath"`
	Password       string `yaml:"password"`
}

// GitRepository is configuration to connect to a git repository.
type GitRepository struct {
	// URL is a required field that specifies the HTTP/S or SSH address of the Git repository.
	URL string `yaml:"url"`
	// CA is an optional field to specify the Certificate Authority to trust while connecting with a git repository over HTTPS. If not specified OS CA's are used.
	CA string `yaml:"caCert"`
	// Branch is a required field that specifies the branch of the repository to use.
	Branch string `yaml:"branch"`
	// Internal is a required field that specifies the interval at which the Git repository must be fetched.
	Interval *time.Duration `yaml:"interval"`
	// Timeout is an optional field to specify a timeout for Git operations like cloning. Defaults to 60s.
	Timeout *time.Duration `yaml:"timeout"`
	// AllClusters is an optional field to specify if the RBAC for all clusters should be managed. Defaults to false.
	AllClusters bool `yaml:"allClusters"`
	// Clusters is an optional field to specify a list of clusters for which the RBAC should be managed.
	Clusters []string `yaml:"clusters"`
	// BasicAuth is an optional field to specify credentials to authenticate towards a Git repository over HTTPS using basic access authentication.
	BasicAuth *GitBasicAuth `yaml:"basicAuth"`
	// SSHAuth is an optional field to specify credentials to authenticate towards a Git repository over SSH. With the respective private key of the SSH key pair, and the host keys of the Git repository.
	SSHAuth *GitSSHAuth `yaml:"sshAuth"`
}

// ClusterRoleMapping is a mapping from m8 roles to ClusterRole's in a K8s cluster
type ClusterRoleMapping struct {
	Scope           es.Scope `yaml:"scope"`
	Role            es.Role  `yaml:"role"`
	ClusterRoleName string   `yaml:"clusterRoleName"`
}

// NewConfigFromFile creates a new GitRepoReconcilerConfig from a given yaml file
func NewConfigFromFile(data []byte) (*GitRepoReconcilerConfig, error) {
	// Unmarshal
	conf := &GitRepoReconcilerConfig{}
	err := yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	for _, repo := range conf.Repositories {
		// Check required fields are set
		if err := conf.parseCloneOptions(repo); err != nil {
			return conf, err
		}

		if repo.Interval == nil {
			return nil, ErrIntervalIsRequired
		}

		// Set default values
		if repo.Timeout == nil {
			repo.Timeout = &DefaultTimeout
		}
	}

	return conf, nil
}

// parseCloneOptions parses the configuration using the git library to validate.
func (c *GitRepoReconcilerConfig) parseCloneOptions(repo *GitRepository) error {
	cloneOptions := &git.CloneOptions{
		URL: repo.URL,
	}

	// Configure basic auth optionally
	if repo.BasicAuth != nil {
		cloneOptions.Auth = &http.BasicAuth{
			Username: repo.BasicAuth.Username,
			Password: repo.BasicAuth.Password,
		}
	}

	// Configure ssh auth
	if repo.SSHAuth != nil {
		_, err := os.Stat(repo.SSHAuth.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("read file %s failed: %w", repo.SSHAuth.PrivateKeyPath, err)
		}

		publicKeys, err := ssh.NewPublicKeysFromFile("git", repo.SSHAuth.PrivateKeyPath, repo.SSHAuth.Password)
		if err != nil {
			return err
		}
		cloneOptions.Auth = publicKeys
	}

	if err := cloneOptions.Validate(); err != nil {
		return err
	}
	c.cloneOptions = cloneOptions
	return nil
}
