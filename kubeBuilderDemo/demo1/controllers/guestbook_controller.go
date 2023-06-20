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

package controllers

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "GoDevEach/api/v1"
)

// GuestbookReconciler reconciles a Guestbook object
type GuestbookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.demo1.com,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.demo1.com,resources=guestbooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.demo1.com,resources=guestbooks/finalizers,verbs=update

//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Guestbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *GuestbookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.WithName("GuestBook")
	//logger := r.Log.WithValues("GuestBook", req.NamespacedName)

	instance := &webappv1.Guestbook{}
	if err := r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	logger.Info(fmt.Sprintf("0. %+v", instance))

	// TODO(user): your logic here
	//logger.Info(fmt.Sprintf("1. %v", req))
	//logger.Info(fmt.Sprintf("2. %s", debug.Stack()))

	// 从crd的资源(cr)中获取ConfigMap的数据
	configMapData := instance.Spec.ConfigMap1
	logger.Info(fmt.Sprintf("1. %+v", configMapData))
	cname := "my-configmap"
	cm := newCMForCR(instance, cname)

	// 尝试获取已存在的 ConfigMap 对象
	existingConfigMap := &corev1.ConfigMap{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cname,
		Namespace: req.Namespace,
	}, existingConfigMap)

	if err != nil {
		if errors.IsNotFound(err) {
			// 如果 ConfigMap 不存在，则创建它
			err = r.Client.Create(context.TODO(), cm)
			if err != nil {
				// 处理错误
				logger.Error(err, fmt.Sprintf("!!!!!!!!!"))
			}
		} else {
			// 处理其他错误
			logger.Error(err, fmt.Sprintf("!!!!!!!!!"))
		}
	} else {
		// 如果 ConfigMap 已存在，则更新它的标签和数据
		existingConfigMap.Labels = cm.Labels
		existingConfigMap.Data = cm.Data

		err = r.Client.Update(context.TODO(), existingConfigMap)
		if err != nil {
			// 处理错误
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GuestbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Guestbook{}).
		Complete(r)
}

func newCMForCR(cr *webappv1.Guestbook, cname string) *corev1.ConfigMap {

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cname,
			Namespace: cr.Namespace,
			Labels:    cr.GetLabels(),
		},
		Data: map[string]string{
			"key1": cr.Spec.ConfigMap1.Key1.String(),
			"key2": cr.Spec.ConfigMap1.Key2.String(),
		},
	}
	return configMap // 创建一个新的 ConfigMap 对象
}
