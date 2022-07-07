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
	"time"

	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Repositories []GitRepository      `yaml:"repositories"`
	Mappings     []ClusterRoleMapping `yaml:"mappings"`
}

// GitCredentialsBasicAccessAuth is used to authenticate towards a Git repository over HTTPS using basic access authentication.
type GitCredentialsBasicAccessAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// GitSSHAuth is used to authenticate towards a Git repository over SSH. With the respective private key of the SSH key pair, and the host keys of the Git repository.
type GitSSHAuth struct {
	Identity   string `yaml:"identity"`
	KnownHosts string `yaml:"knownHosts"`
}

// GitRepository is configuration to connect to a git repository.
type GitRepository struct {
	// URL is a required field that specifies the HTTP/S or SSH address of the Git repository.
	URL string `yaml:"url"`
	// CA is the Certificate Authority to trust while connecting with a Git repository over HTTPS.
	CA string `yaml:"caCert"`
	// Branch is the branch of the repository to use.
	Branch string `yaml:"branch"`
	// Internal is a required field that specifies the interval at which the Git repository must be fetched.
	Interval time.Duration `yaml:"interval"`
	// Timeout is an optional field to specify a timeout for Git operations like cloning.
	Timeout time.Duration `yaml:"timeout"`
	// allClusters specifies if the RBAC for all clusters should be managed.
	AllClusters bool `yaml:"allClusters"`
	// Clusters specifies a list of clusters for which the RBAC should be managed.
	Clusters []string `yaml:"clusters"`
	// CredentialsBasicAccessAuth is used to authenticate towards a Git repository over HTTPS using basic access authentication.
	CredentialsBasicAccessAuth *GitCredentialsBasicAccessAuth `yaml:"credentialsBasicAccessAuth"`
	// SSHAuth is used to authenticate towards a Git repository over SSH. With the respective private key of the SSH key pair, and the host keys of the Git repository.
	SSHAuth *GitSSHAuth `yaml:"sshAuth"`
}

// ClusterRoleMapping is a mapping from m8 roles to ClusterRole's in a K8s cluster
type ClusterRoleMapping struct {
	Scope           es.Scope `yaml:"scope"`
	Role            es.Role  `yaml:"role"`
	ClusterRoleName string   `yaml:"clusterRoleName"`
}

// NewConfigFromFile creates a new Config from a given yaml file
func NewConfigFromFile(data []byte) (*Config, error) {
	conf := &Config{}
	err := yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
