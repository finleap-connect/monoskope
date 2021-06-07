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
