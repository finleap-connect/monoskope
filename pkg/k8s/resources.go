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

package k8s

import (
	"fmt"

	rbac "k8s.io/api/rbac/v1"
	api "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	BASE_DOMAIN = "monoskope.io"
)

// NewClusterRoleBinding creates a new K8s API ClusterRoleBinding resource
func NewClusterRoleBinding(clusterRoleName, userName, userNamePrefix string, labels map[string]string) *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		TypeMeta: api.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name:   fmt.Sprintf("%s:%s", clusterRoleName, userName),
			Labels: labels,
		},
		RoleRef: rbac.RoleRef{
			Kind:     "ClusterRole", // ClusterRoleBindings can only refer to ClusterRole
			Name:     clusterRoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: []rbac.Subject{
			{
				Name: fmt.Sprintf("%s%s", userNamePrefix, userName),
				Kind: "User",
			},
		},
	}
}
