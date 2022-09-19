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
	MonoskopeDomain        = "monoskope.io"
	ClusterRoleBindingKind = "ClusterRoleBinding"
	ClusterRoleKind        = "ClusterRole"
)

// NewClusterRoleBinding creates a new K8s API ClusterRoleBinding resource
func NewClusterRoleBinding(clusterRoleName, userName, userNamePrefix string, labels map[string]string) *rbac.ClusterRoleBinding {
	// prefix labels with MonoskopeDomain
	labelsWithPrefix := make(map[string]string)
	for k, v := range labels {
		labelsWithPrefix[fmt.Sprintf("%s/%s", MonoskopeDomain, k)] = v
	}
	return &rbac.ClusterRoleBinding{
		TypeMeta: api.TypeMeta{
			Kind:       ClusterRoleBindingKind,
			APIVersion: rbac.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: api.ObjectMeta{
			Name:   fmt.Sprintf("%s:%s", clusterRoleName, userName),
			Labels: labelsWithPrefix,
		},
		RoleRef: rbac.RoleRef{
			Kind:     ClusterRoleKind,
			Name:     clusterRoleName,
			APIGroup: rbac.SchemeGroupVersion.Group,
		},
		Subjects: []rbac.Subject{
			{
				Name: fmt.Sprintf("%s%s", userNamePrefix, userName),
				Kind: rbac.UserKind,
			},
		},
	}
}
