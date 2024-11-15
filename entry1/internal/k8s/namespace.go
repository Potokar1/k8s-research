package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListNamespaces returns a list of namespaces
func ListNamespaces(ctx context.Context) ([]string, error) {
	// get clientset
	clientset := GetClientSet()

	// list namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// return namespaces
	var ns []string
	for _, namespace := range namespaces.Items {
		ns = append(ns, namespace.Name)
	}
	return ns, nil
}
