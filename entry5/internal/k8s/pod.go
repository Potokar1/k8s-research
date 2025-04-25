package k8s

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

type podModifiedEvent struct {
	PodName   string            `json:"pod_name"`
	Inventory map[string]string `json:"inventory"`
}

// WatchPods watches for changes to pod(s) given filters such as namespace and label
func WatchPods(ctx context.Context, namespace string, labelValue string) (chan podModifiedEvent, error) {
	clientset := GetClientSet()

	listOpts := metav1.ListOptions{}
	if labelValue != "" {
		listOpts.LabelSelector = TownLabel + "=" + labelValue
	}

	watch, err := clientset.CoreV1().Pods(namespace).Watch(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	// Create a channel to send pod modified events
	podModifiedChan := make(chan podModifiedEvent)
	go func() {
		defer close(podModifiedChan)
		for event := range watch.ResultChan() {
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue // skip if the object is not a Pod
			}
			switch event.Type {
			case "MODIFIED":
				podModifiedChan <- podModifiedEvent{
					PodName:   pod.Name,
					Inventory: pod.Annotations,
				}
			}
		}
	}()
	return podModifiedChan, nil
}

func PatchPod(ctx context.Context, namespace string, podName string, inventory map[string]string) error {
	clientset := GetClientSet()

	data := map[string]any{
		"metadata": map[string]any{
			"annotations": inventory,
		},
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	strictFieldValidation := "Strict"
	// Patch the pod with the provided data
	patchOptions := metav1.PatchOptions{
		FieldValidation: strictFieldValidation,
	}
	_, err = clientset.CoreV1().Pods(namespace).Patch(ctx, podName, types.MergePatchType, dataBytes, patchOptions)
	return err
}
