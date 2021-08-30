// Copyright 2021 Monoskope Authors
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
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func localContextConfig() (*rest.Config, error) {
	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	return config, err
}

// NewClient creates a new kubernetes.Clientset using the
// in-cluster config with a fallback to the local kubeconfig.
func NewClient() (client.Client, error) {
	// try using in cluster config
	config, err := rest.InClusterConfig()

	if err == rest.ErrNotInCluster {
		// if not in cluster try local context
		config, err = localContextConfig()
	}

	// if any of the above failed, return error
	if err != nil {
		return nil, err
	}

	cl, err := client.New(config, client.Options{})
	if err != nil {
		return nil, err
	}

	return cl, nil
}
