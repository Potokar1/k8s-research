package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
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

// ListAllPodsInNamespace returns a list of all pods in a given namespace
func ListAllPodsInNamespace(ctx context.Context, namespace string) ([]string, error) {
	pods := GetAllPodsInNamespace(ctx, namespace) // this is a helper function to get all pods in a namespace

	// return pod names
	var allPods []string
	for _, pod := range pods {
		// append the pod name to the list
		allPods = append(allPods, pod.Name)
	}

	return allPods, nil
}

// GetAllPodsInNamespace returns a slice of all pods in a given namespace
func GetAllPodsInNamespace(ctx context.Context, namespace string) []v1.Pod {
	// get clientset
	clientset := GetClientSet()

	// list pods
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []v1.Pod{}
	}

	// return pod items
	return pods.Items
}
