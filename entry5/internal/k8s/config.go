package k8s

import (
	"log/slog"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	TownLabel = "town"
)

// GetClientSet returns a kubernetes clientset from any found kubeconfig
func GetClientSet() *kubernetes.Clientset {
	// load kubeconfig

	kubeconfig := ""
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// slog error
		slog.Error("no home directory found, using in-cluster config or default config if no in-cluster config found")
	} else {
		kubeconfig = filepath.Join(homeDir, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// panic as we can't continue without a clientset and we should be able to use in-cluster/default config
		panic(err)
	}

	// create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}
