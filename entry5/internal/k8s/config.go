package k8s

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	TownLabel = "town"
)

// GetClientSet returns a kubernetes clientset from any found kubeconfig
func GetClientSet() *kubernetes.Clientset {
	// load kubeconfig
	cfg, err := rest.InClusterConfig()
	if err != nil {
		// fall back to local kubeconfig
		homeDir, _ := os.UserHomeDir()
		kubeconfig := filepath.Join(homeDir, ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			// panic as we can't continue without a clientset and we should be able to use in-cluster/default config
			panic(err)
		}
	}

	// create clientset
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		// panic as we can't continue without a clientset
		panic(err)
	}
	return clientset
}
