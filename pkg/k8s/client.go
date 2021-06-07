package k8s

import (
	"context"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// K8sClient is a Kubernetes client
type K8sClient interface {
	Get(ctx context.Context, k8sObject runtime.Object) error
	Create(ctx context.Context, k8sObject runtime.Object) error
}

type k8sClient struct {
	client client.Client
}

func localContextConfig() (*rest.Config, error) {
	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	return config, err
}

// NewClient creates a new kubernetes.Clientset using the
// in-cluster config with a fallback to the local kubeconfig.
func NewClient() (K8sClient, error) {
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

	return &k8sClient{
		client: cl,
	}, nil
}

func (c *k8sClient) Get(ctx context.Context, k8sObject runtime.Object) error {
	// Get namespaced name from object
	objectKey, err := client.ObjectKeyFromObject(k8sObject)
	if err != nil {
		return err
	}

	// Get object
	err = c.client.Get(ctx, objectKey, k8sObject)
	if err != nil {
		return err
	}

	return nil
}

func (c *k8sClient) Create(ctx context.Context, k8sObject runtime.Object) error {
	return c.client.Create(ctx, k8sObject)
}
