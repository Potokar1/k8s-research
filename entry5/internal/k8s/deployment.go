package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListDeployments returns the the unique set of deployment labels in a namespace
func ListDeployments(ctx context.Context, namespace string) ([]string, error) {
	clientset := GetClientSet()

	deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	deploymentTownLabels := make(map[string]struct{})
	for _, deployment := range deployments.Items {
		// only add unique labels
		if _, ok := deploymentTownLabels[deployment.Labels[TownLabel]]; !ok {
			deploymentTownLabels[deployment.Labels[TownLabel]] = struct{}{}
		}
	}

	var deploymentNames []string
	for deployment := range deploymentTownLabels {
		deploymentNames = append(deploymentNames, deployment)
	}
	return deploymentNames, nil
}
