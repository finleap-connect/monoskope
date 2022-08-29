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
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"gopkg.in/yaml.v2"
)

const (
	DefaultTimeout        = 60 * time.Second
	DefaultUsernamePrefix = "oidc:"

	AuthTypeBasic               = "basic"
	AuthTypeBasicSuffixUsername = ".basic.username"
	AuthTypeBasicSuffixPassword = ".basic.password"

	AuthTypeSSH                 = "ssh"
	AuthTypeSSHSuffixPrivateKey = ".ssh.privateKey"
	AuthTypeSSHSuffixPassword   = ".ssh.password"
)

type ClusterRoleMapping struct {
	Scope       string
	Role        string
	ClusterRole string
}

// Config is the configuration for the GitRepoReconciler.
type Config struct {
	Repositories   []*GitRepository      `yaml:"repositories"`
	Mappings       []*ClusterRoleMapping `yaml:"mappings"`
	UsernamePrefix string                `yaml:"usernamePrefix"` // UsernamePrefix is prepended to usernames to prevent clashes with existing names (such as system: users). For example, the value oidc: will create usernames like oidc:jane.doe. Defaults to oidc:.
}

type ReconcilerConfig struct {
	RootDirectory  string
	SubPath        string
	UsernamePrefix string
	Mappings       []*ClusterRoleMapping `yaml:"mappings"`
}

func NewReconcilerConfig(rootDir, subPath, usernamePrefix string, mappings []*ClusterRoleMapping) *ReconcilerConfig {
	return &ReconcilerConfig{rootDir, subPath, usernamePrefix, mappings}
}

type GitAuth struct {
	Type      string `yaml:"type"`
	EnvPrefix string `yaml:"envPrefix"`
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
	// SubDir is the relative path within the repo where to reconcile yamls
	SubDir string  `yaml:"subdir"`
	Auth   GitAuth `yaml:"auth"`
	// cloneOptions are the parsed settings
	cloneOptions *git.CloneOptions
}

// NewConfigFromFile creates a new GitRepoReconcilerConfig from a given yaml file path
func NewConfigFromFilePath(name string) (*Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewConfigFromFile(data)
}

// NewConfigFromFile creates a new GitRepoReconcilerConfig from a given yaml file
func NewConfigFromFile(data []byte) (*Config, error) {
	// Unmarshal
	conf := &Config{}

	if err := yaml.Unmarshal(data, conf); err != nil {
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
			return conf, ErrIntervalIsRequired
		}

		// Set default values
		if repo.Timeout == nil {
			timeout := DefaultTimeout
			repo.Timeout = &timeout
		}
	}

	return conf, nil
}

func getClusterRoleMapping(mappings []*ClusterRoleMapping, scope, role string) string {
	for _, m := range mappings {
		if m.Scope == scope && m.Role == role {
			return m.ClusterRole
		}
	}
	return ""
}

// configureBasicAuth reads the file containing the basic auth information and unmarshal's it's content into the clone options given.
func configureBasicAuth(repo *GitRepository, cloneOptions *git.CloneOptions) error {
	// get env
	username := os.Getenv(fmt.Sprintf("%s%s", repo.Auth.EnvPrefix, AuthTypeBasicSuffixUsername))
	password := os.Getenv(fmt.Sprintf("%s%s", repo.Auth.EnvPrefix, AuthTypeBasicSuffixPassword))

	// set clone options auth
	cloneOptions.Auth = &http.BasicAuth{
		Username: username,
		Password: password,
	}

	return nil
}

// configureSSHAuth reads the file containing the ssh auth information and unmarshal's it's content into the clone options given.
func configureSSHAuth(repo *GitRepository, cloneOptions *git.CloneOptions) error {
	// get env
	privateKeyBase64 := os.Getenv(fmt.Sprintf("%s%s", repo.Auth.EnvPrefix, AuthTypeSSHSuffixPrivateKey))
	password := os.Getenv(fmt.Sprintf("%s%s", repo.Auth.EnvPrefix, AuthTypeSSHSuffixPassword))

	// set clone options auth
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("failed to create temp file to write private key to: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	if err := os.WriteFile(tmpFile.Name(), []byte(privateKey), 0600); err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", tmpFile.Name(), password)
	if err != nil {
		return err
	}
	cloneOptions.Auth = publicKeys

	return nil
}

// parseCloneOptions parses the configuration using the git library to validate.
func (c *Config) parseCloneOptions(repo *GitRepository) error {
	cloneOptions := &git.CloneOptions{
		URL:          repo.URL,
		SingleBranch: true,
		Depth:        1,
	}
	if repo.Branch != "" {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(repo.Branch)
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
	if repo.Auth.Type == AuthTypeBasic {
		if err := configureBasicAuth(repo, cloneOptions); err != nil {
			return err
		}
	}

	// Configure ssh auth
	if repo.Auth.Type == AuthTypeSSH {
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
