/*
Copyright 2023.

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

package controller

import (
	"context"
	"encoding/json"
	"os"

	"github.com/blacktailed/test-kubebuilder.git/pkg/common"
	"github.com/blacktailed/test-kubebuilder.git/pkg/slack"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "k8s.io/api/core/v1"
)

// EventReconciler reconciles a Event object
type EventReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=example.test.kubebuiler,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.test.kubebuiler,resources=events/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.test.kubebuiler,resources=events/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Event object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *EventReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// TODO(user): your logic here
	event := &v1.Event{}
	pod := &v1.Pod{}
	// event := &metav1.Event{}

	if err := r.Get(ctx, req.NamespacedName, event); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// if event.Type == "Warning" || event.Reason == "BackOff" {
	if event.Namespace == "test" {
		podName := event.InvolvedObject.Name

		// Pod 이름으로 Pod 리소스 조회
		podKey := types.NamespacedName{Name: podName, Namespace: event.Namespace}
		if err := r.Get(ctx, podKey, pod); err != nil {
			if client.IgnoreNotFound(err) != nil {
				return ctrl.Result{}, err
			}
			// Pod를 찾지 못한 경우
			l.Info("Pod not found", "Pod Name: ", podName)
		} else {
			// Pod의 Yaml 정보를 Json을 변환
			podYaml, err := json.Marshal(pod)
			if err != nil {
				return ctrl.Result{}, err
			}
			l.Info("Event", "Name", event.Name, "Namespace", event.Namespace, "Reason", event.Reason, "Pod Name", podName, "Pod YAML", string(podYaml))

			//slack message test
			content := &common.SlackMsg{
				WebhookURL: os.Getenv("SLACK_WEBHOOK_URL"),
				Text:       podYaml,
			}
			slack.TestMess(content)
		}
		// podSpec := pod.Items
		// l.Info("Event", "Name", event.Name, "Namespace", event.Namespace, "Reason", event.Reason, "Object Meta", event.InvolvedObject, "Pod Name", podName, "Pod List", pod.ObjectMeta)
	}

	// if err := r.Update(ctx, pod); err != nil {
	// }
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		For(&v1.Event{}).
		// Watches(&v1.Pod{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
