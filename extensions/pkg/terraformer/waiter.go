// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package terraformer

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/gardener/pkg/utils/retry"
)

// WaitForCleanEnvironment waits until no Terraform Pod(s) exist for the current instance of the Terraformer.
func (t *terraformer) WaitForCleanEnvironment(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, t.deadlineCleaning)
	defer cancel()

	var (
		err     error
		podList = &corev1.PodList{}
	)

	t.logger.Info("Waiting for clean environment")
	err = retry.UntilTimeout(ctx, 5*time.Second, 2*time.Minute, func(ctx context.Context) (done bool, err error) {
		podList, err = t.listPods(ctx)
		if err != nil {
			return retry.SevereError(err)
		}

		if len(podList.Items) > 0 {
			t.logger.Info("Waiting until all Terraformer pods have been cleaned up")
			return retry.MinorError(fmt.Errorf("at least one terraformer pod still exists: %s", podList.Items[0].Name))
		}

		return retry.Ok()
	})

	if err == context.DeadlineExceeded && len(podList.Items) > 0 {
		t.logger.Info("Fetching logs of Terraformer pods as waiting for clean environment timed out")
		for _, pod := range podList.Items {
			podLogger := t.logger.WithValues("pod", client.ObjectKey{Namespace: t.namespace, Name: pod.Name})
			podLogs, err := t.retrievePodLogs(ctx, podLogger, &pod)

			if err != nil {
				podLogger.Error(err, "Could not retrieve logs of Terraformer pod, "+pod.Name)
				continue
			}
			podLogger.Info("Logs of Terraformer pod, " + pod.Name + ": " + podLogs)
		}
	}

	return err
}

type podStatus byte

const (
	podStatusSucceeded podStatus = iota
	podStatusFailure
	podStatusCreationTimeout
)

// waitForPod waits for the Terraform Pod to be completed (either successful or failed).
// It checks the Pod status field to identify the state.
func (t *terraformer) waitForPod(ctx context.Context, logger logr.Logger, pod *corev1.Pod) (podStatus, string) {
	var (
		status             = podStatusFailure
		terminationMessage = ""
		log                = logger.WithValues("pod", client.ObjectKeyFromObject(pod))
	)

	timeoutCtx, cancel := context.WithTimeout(ctx, t.deadlinePod)
	defer cancel()

	log.Info("Waiting for Terraformer pod to be completed...")
	_ = retry.Until(timeoutCtx, 5*time.Second, func(ctx context.Context) (bool, error) {
		if err := t.client.Get(ctx, client.ObjectKeyFromObject(pod), pod); err != nil {
			if apierrors.IsNotFound(err) {
				log.Info("Terraformer pod disappeared unexpectedly, somebody must have manually deleted it")
				return retry.Ok()
			}

			log.Error(err, "Error retrieving pod")
			return retry.SevereError(err)
		}

		var (
			phase             = pod.Status.Phase
			containerStatuses = pod.Status.ContainerStatuses
		)

		if len(containerStatuses) > 0 {
			switch phase {
			case corev1.PodPending:
				// Check whether the Pod has been created successfully
				if containerStateWaiting := containerStatuses[0].State.Waiting; containerStateWaiting != nil && containerStateWaiting.Reason == "ContainerCreating" {
					if podAge := time.Now().UTC().Sub(pod.CreationTimestamp.Time.UTC()); podAge > t.deadlinePodCreation {
						status = podStatusCreationTimeout
						log.Info("Timeout creating pod")
						return retry.Ok()
					}
				}

			case corev1.PodSucceeded, corev1.PodFailed:
				// Check whether the Pod has been executed successfully
				if containerStateTerminated := containerStatuses[0].State.Terminated; containerStateTerminated != nil {
					if containerStateTerminated.ExitCode == 0 {
						status = podStatusSucceeded
					}
					terminationMessage = containerStateTerminated.Message
				}
				return retry.Ok()
			}
		}

		log.Info("Waiting for Terraformer pod to be completed, pod hasn't finished yet", "phase", phase, "len-of-containerstatuses", len(containerStatuses))
		return retry.MinorError(fmt.Errorf("pod was not successful: phase=%s, len-of-containerstatuses=%d", phase, len(containerStatuses)))
	})

	return status, terminationMessage
}
