/*
Copyright 2026.

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
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configv1alpha1 "my.domain/inventory/api/v1alpha1"
)

// InventoryReconciler reconciles a Inventory object
type InventoryReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Logger logr.Logger

	// key: namespace, value: list names
	inventories map[string][]string
}

// +kubebuilder:rbac:groups=config.my.domain,resources=inventories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=config.my.domain,resources=inventories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=config.my.domain,resources=inventories/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile
func (r *InventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	logger.V(1).Info("Reconciling")
	// Fecth the clusterSummary instance
	inventory := &configv1alpha1.Inventory{}
	if err := r.Get(ctx, req.NamespacedName, inventory); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		logger.Error(err, "Failed to fetch inventory")
		return reconcile.Result{}, fmt.Errorf(
			"failed to fetch inventory %s: %w",
			req.NamespacedName, err,
		)
	}

	if !inventory.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, inventory, logger)
	}

	return r.reconcileNormal(ctx, inventory, logger)
}

func (r *InventoryReconciler) reconcileDelete(ctx context.Context,
	inventory *configv1alpha1.Inventory, logger logr.Logger) (ctrl.Result, error) {
	// cleanup

	// remove instance from inventories

	// remove finalizer

	return ctrl.Result{}, nil
}

func (r *InventoryReconciler) reconcileNormal(ctx context.Context,
	inventory *configv1alpha1.Inventory, logger logr.Logger) (ctrl.Result, error) {
	// TODO: reports errors

	// add finalizier
	if !controllerutil.ContainsFinalizer(inventory, configv1alpha1.InventoryFinalizer) {
		if err := r.addFinalizer(ctx, inventory); err != nil {
			logger.V(1).Error(err, "failed to add finalizer")
			return reconcile.Result{}, err
		}
	}

	_, ok := r.inventories[inventory.Namespace]
	if !ok {
		r.inventories[inventory.Namespace] = make([]string, 0)
	}
	r.inventories[inventory.Namespace] = append(r.inventories[inventory.Namespace], inventory.Name)

	listOptions := []client.ListOption{
		client.InNamespace(inventory.Namespace),
	}
	pods := corev1.PodList{}
	err := r.List(ctx, &pods, listOptions...)
	if err != nil {
		logger.V(1).Error(err, "failed to list pods")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	key := configv1alpha1.DefaultAnnotation
	if inventory.Spec.AnnotationKey != "" {
		key = inventory.Spec.AnnotationKey
	}

	podsOwner := make(map[string][]string, 0)
	// walk pods
	for i := range pods.Items {
		if pods.Items[i].Annotations != nil {
			owner, ok := pods.Items[i].Annotations[key]
			if ok {
				if podsOwner[owner] == nil {
					podsOwner[owner] = make([]string, 0)
				}
				podsOwner[owner] = append(podsOwner[owner], pods.Items[i].Name)
			}
		}
	}

	// TODO: move this logic to fuinction. add ownerReference

	// create/update configMap
	configMapName := inventory.Name
	configMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Namespace: inventory.Namespace, Name: configMapName}, configMap)
	if err != nil {
		if apierrors.IsNotFound(err) {
			configMap.Namespace = inventory.Namespace
			configMap.Name = configMapName
			configMap.Data = map[string]string{}
			for owner := range podsOwner {
				configMap.Data[owner] = strings.Join(podsOwner[owner], ",")
			}
			return ctrl.Result{}, r.Create(ctx, configMap)
		}
		return ctrl.Result{}, err
	}

	configMap.Data = map[string]string{}
	for owner := range podsOwner {
		configMap.Data[owner] = strings.Join(podsOwner[owner], ",")
	}
	return ctrl.Result{}, r.Update(ctx, configMap)
}

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO: react to pod and configMap changes with predicates
	_, err := ctrl.NewControllerManagedBy(mgr).
		For(&configv1alpha1.Inventory{},
			builder.WithPredicates(
				InventoryPredicate{Logger: r.Logger.WithName("inventoryPredicate")}),
		).
		Watches(&corev1.Pod{},
			handler.EnqueueRequestsFromMapFunc(r.requeueInventoryForPods),
			builder.WithPredicates(
				PodPredicates(mgr.GetLogger().WithValues("predicate", "podpredicates")),
			),
		).
		Named("inventory").
		Build(r)
	if err != nil {
		return err
	}

	r.inventories = make(map[string][]string)
	return nil
}

func (r *InventoryReconciler) addFinalizer(ctx context.Context, inventory *configv1alpha1.Inventory) error {
	controllerutil.AddFinalizer(inventory, configv1alpha1.InventoryFinalizer)
	return r.Update(ctx, inventory)
}
