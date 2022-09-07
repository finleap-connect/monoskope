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
	"os"
	"time"

	"github.com/finleap-connect/monoskope/pkg/git"
	"github.com/finleap-connect/monoskope/pkg/logger"

	"gopkg.in/yaml.v2"
)

const (
	DefaultTimeout        = 60 * time.Second
	DefaultInterval       = 10 * time.Minute
	DefaultUsernamePrefix = "oidc:"
)

type ClusterRoleMapping struct {
	Scope       string `yaml:"scope"`
	Role        string `yaml:"role"`
	ClusterRole string `yaml:"clusterRole"`
}

// Config is the configuration for the GitRepoReconciler.
type Config struct {
	log logger.Logger
	// Internal is a required field that specifies the interval at which the Git repository must be fetched.
	Interval *time.Duration `yaml:"interval"`
	// Repository is the git config to use
	Repository *git.GitConfig `yaml:"repository"`
	// Mappings define which k8s role in m8 leads to which cluster role within clusters
	Mappings []*ClusterRoleMapping `yaml:"mappings"`
	// UsernamePrefix is prepended to usernames to prevent clashes with existing names (such as system: users). For example, the value oidc: will create usernames like oidc:jane.doe. Defaults to oidc:.
	UsernamePrefix string `yaml:"usernamePrefix"`
	// AllClusters is an optional field to specify if the RBAC for all clusters should be managed. Defaults to false.
	AllClusters bool `yaml:"allClusters"`
	// Clusters is an optional field to specify a list of clusters for which the RBAC should be managed.
	Clusters []string `yaml:"clusters"`
	// SubDir is the relative path within the repo where to reconcile yamls
	SubDir string `yaml:"subdir"`
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
	conf.log = logger.WithName("config")

	if err := conf.setDefaults(); err != nil {
		return nil, err
	}

	return conf, nil
}

// setDefaults sets the default values on the configuration
func (conf *Config) setDefaults() error {
	if len(conf.UsernamePrefix) == 0 {
		conf.UsernamePrefix = DefaultUsernamePrefix
	}
	if conf.Interval == nil {
		interval := DefaultInterval
		conf.Interval = &interval
	}
	if conf.Repository.Timeout == nil {
		timeout := DefaultTimeout
		conf.Repository.Timeout = &timeout
	}
	return nil
}

func (conf *Config) getClusterRoleMapping(scope, role string) string {
	for _, m := range conf.Mappings {
		if m.Scope == scope && m.Role == role {
			return m.ClusterRole
		}
	}
	return ""
}
