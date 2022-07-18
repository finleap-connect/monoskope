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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"gopkg.in/yaml.v2"
)

var (
	DefaultTimeout        = 60 * time.Second
	DefaultUsernamePrefix = "oidc:"
)

// Config is the configuration for the GitRepoReconciler.
type Config struct {
	Repositories   []*GitRepository      `yaml:"repositories"`
	Mappings       []*ClusterRoleMapping `yaml:"mappings"`
	UsernamePrefix string                `yaml:"usernamePrefix"` // UsernamePrefix is prepended to usernames to prevent clashes with existing names (such as system: users). For example, the value oidc: will create usernames like oidc:jane.doe. Defaults to oidc:.
}

type ReconcilerConfig struct {
	LocalDirectory string
	UsernamePrefix string
	Mappings       []*ClusterRoleMapping
}

func NewReconcilerConfig(localDirectory, usernamePrefix string, mappings []*ClusterRoleMapping) *ReconcilerConfig {
	return &ReconcilerConfig{localDirectory, usernamePrefix, mappings}
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
	// BasicAuthPath is an optional field to specify the file containing credentials to authenticate towards a Git repository over HTTPS using basic access authentication.
	BasicAuthPath string `yaml:"basicAuthPath"`
	// SSHAuthPath is an optional field to specify the file containing credentials to authenticate towards a Git repository over SSH. With the respective private key of the SSH key pair, and the host keys of the Git repository.
	SSHAuthPath string `yaml:"sshAuthPath"`
	// cloneOptions are the parsed settings
	cloneOptions *git.CloneOptions
}

// ClusterRoleMapping is a mapping from m8 roles to ClusterRole's in a K8s cluster
type ClusterRoleMapping struct {
	Scope           es.Scope `yaml:"scope"`
	Role            es.Role  `yaml:"role"`
	ClusterRoleName string   `yaml:"clusterRoleName"`
}

// NewConfigFromFile creates a new GitRepoReconcilerConfig from a given yaml file
func NewConfigFromFile(data []byte) (*Config, error) {
	// Unmarshal
	conf := &Config{}
	err := yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	// Set default values
	if len(conf.UsernamePrefix) == 0 {
		conf.UsernamePrefix = DefaultUsernamePrefix
	}

	for _, repo := range conf.Repositories {
		// Check required fields are set
		if err := conf.parseCloneOptions(repo); err != nil {
			return conf, err
		}

		if repo.Interval == nil {
			return nil, ErrIntervalIsRequired
		}

		if len(repo.Branch) == 0 {
			return nil, ErrBranchIsRequired
		}

		// Set default values
		if repo.Timeout == nil {
			repo.Timeout = &DefaultTimeout
		}
	}

	return conf, nil
}

// configureBasicAuth reads the file containing the basic auth information and unmarshal's it's content into the clone options given.
func configureBasicAuth(repo *GitRepository, cloneOptions *git.CloneOptions) error {
	// read file
	data, err := ioutil.ReadFile(repo.BasicAuthPath)
	if err != nil {
		return err
	}

	// unmarshal
	basicAuth := &GitBasicAuth{}
	err = yaml.Unmarshal(data, basicAuth)
	if err != nil {
		return err
	}

	// set clone options auth
	cloneOptions.Auth = &http.BasicAuth{
		Username: basicAuth.Username,
		Password: basicAuth.Password,
	}

	return nil
}

// configureSSHAuth reads the file containing the ssh auth information and unmarshal's it's content into the clone options given.
func configureSSHAuth(repo *GitRepository, cloneOptions *git.CloneOptions) error {
	// read file
	data, err := ioutil.ReadFile(repo.SSHAuthPath)
	if err != nil {
		return err
	}

	// unmarshal
	sshAuth := &GitSSHAuth{}
	err = yaml.Unmarshal(data, sshAuth)
	if err != nil {
		return err
	}

	// set clone options auth
	if _, err := os.Stat(sshAuth.PrivateKeyPath); err != nil {
		return fmt.Errorf("read file %s failed: %w", sshAuth.PrivateKeyPath, err)
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", sshAuth.PrivateKeyPath, sshAuth.Password)
	if err != nil {
		return err
	}
	cloneOptions.Auth = publicKeys

	return nil
}

// parseCloneOptions parses the configuration using the git library to validate.
func (c *Config) parseCloneOptions(repo *GitRepository) error {
	cloneOptions := &git.CloneOptions{
		URL:           repo.URL,
		ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
		SingleBranch:  true,
		NoCheckout:    false,
		Depth:         1,
	}

	// Set CA
	if len(repo.CA) != 0 {
		if data, err := base64.StdEncoding.DecodeString(repo.CA); err != nil {
			return err
		} else {
			cloneOptions.CABundle = data
		}
	}

	// Configure basic auth optionally
	if len(repo.BasicAuthPath) != 0 {
		if err := configureBasicAuth(repo, cloneOptions); err != nil {
			return err
		}
	}

	// Configure ssh auth
	if len(repo.SSHAuthPath) != 0 {
		if err := configureSSHAuth(repo, cloneOptions); err != nil {
			return err
		}
	}

	if err := cloneOptions.Validate(); err != nil {
		return err
	}
	repo.cloneOptions = cloneOptions

	return nil
}
