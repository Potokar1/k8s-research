package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListPods returns the names of pods in a namespace with the given label
func ListPods(ctx context.Context, namespace string, labelValue string) ([]string, error) {
	clientset := GetClientSet()

	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: TownLabel + "=" + labelValue,
	})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}
	return podNames, nil
}
