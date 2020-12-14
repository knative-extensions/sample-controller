/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package simpledeployment

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"

	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
	"knative.dev/sample-controller/pkg/apis/samples"
	samplesv1alpha1 "knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
	simpledeploymentreconciler "knative.dev/sample-controller/pkg/client/injection/reconciler/samples/v1alpha1/simpledeployment"
)

// podOwnerLabelKey is the key to a label that points to the owner (creator) of the
// pod, allowing us to easily list all pods a single SimpleDeployment created.
const podOwnerLabelKey = samples.GroupName + "/podOwner"

// Reconciler implements simpledeploymentreconciler.Interface for
// SimpleDeployment resources.
type Reconciler struct {
	kubeclient kubernetes.Interface
	podLister  corev1listers.PodLister
}

// Check that our Reconciler implements Interface
var _ simpledeploymentreconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, d *samplesv1alpha1.SimpleDeployment) reconciler.Event {
	// This logger has all the context necessary to identify which resource is being reconciled.
	logger := logging.FromContext(ctx)

	// Get all the pods created by the current SimpleDeployment. The result is read from
	// cache (via the lister).
	selector := labels.SelectorFromSet(labels.Set{
		podOwnerLabelKey: d.Name,
	})
	existingPods, err := r.podLister.Pods(d.Namespace).List(selector)
	if err != nil {
		return fmt.Errorf("failed to list existing pods: %w", err)
	}
	logger.Infof("Found %d pods in total", len(existingPods))

	// Find out which pods have the current image and which ones are outdated.
	currentPods, outdatedPods := partitionPods(existingPods, d.Spec.Image)

	// Remove all outdated pods. They don't represent a state we want to be in.
	logger.Infof("Deleting %d outdated pods", len(outdatedPods))
	for _, pod := range outdatedPods {
		if err := r.kubeclient.CoreV1().Pods(d.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("failed to delete pod: %w", err)
		}
	}

	currentCount := int32(len(currentPods))
	if currentCount < d.Spec.Replicas {
		// We don't have as many replicas as we should, so create the remaining ones.
		toCreate := d.Spec.Replicas - currentCount
		logger.Infof("Got %d existing pods, want %d -> creating %d", currentCount, d.Spec.Replicas, toCreate)
		pod := makePod(d)
		for i := int32(0); i < toCreate; i++ {
			if _, err := r.kubeclient.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
				return fmt.Errorf("failed to create pod: %w", err)
			}
		}
	} else if currentCount > d.Spec.Replicas {
		// We have too many replicas, so remove the ones that are too much.
		toDelete := currentCount - d.Spec.Replicas
		logger.Infof("Got %d existing pods, want %d -> removing %d", currentCount, d.Spec.Replicas, toDelete)
		for i := int32(0); i < toDelete; i++ {
			if err := r.kubeclient.CoreV1().Pods(d.Namespace).Delete(ctx, currentPods[i].Name, metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("failed to delete pod: %w", err)
			}
		}
	}

	// Surface the readiness of the pods we've launched.
	var readyPods int32
	for _, p := range currentPods {
		if isPodReady(p) {
			readyPods++
		}
	}

	d.Status.ReadyReplicas = readyPods
	if readyPods >= d.Spec.Replicas {
		d.Status.MarkPodsReady()
	} else {
		d.Status.MarkPodsNotReady(d.Spec.Replicas - readyPods)
	}

	return nil
}

// makePod generates a simple pod to be created in the given namespace with the given
// image.
func makePod(d *samplesv1alpha1.SimpleDeployment) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    d.Namespace,
			GenerateName: d.Name + "-",
			Labels: map[string]string{
				// The label allows for easy querying of all the pods created.
				podOwnerLabelKey: d.Name,
			},
			// The OwnerReference makes sure the pods get removed automatically once the
			// SimpleDeployment is removed.
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(d)},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "user-container",
				Image: d.Spec.Image,
			}},
		},
	}
}

// partitionPods returns a list of pods that have the correct image set and a list of
// pods with a wrong (potentially old) image.
func partitionPods(pods []*corev1.Pod, wantImage string) ([]*corev1.Pod, []*corev1.Pod) {
	var current []*corev1.Pod
	var outdated []*corev1.Pod

	for _, pod := range pods {
		if pod.Spec.Containers[0].Image == wantImage {
			current = append(current, pod)
		} else {
			outdated = append(outdated, pod)
		}
	}

	return current, outdated
}

// isPodReady returns whether or not the given pod is ready.
func isPodReady(p *corev1.Pod) bool {
	if p.Status.Phase == corev1.PodRunning && p.DeletionTimestamp == nil {
		for _, cond := range p.Status.Conditions {
			if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
				return true
			}
		}
	}
	return false
}
