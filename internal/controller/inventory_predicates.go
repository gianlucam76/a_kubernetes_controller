package controller

import (
	"reflect"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	configv1alpha1 "my.domain/inventory/api/v1alpha1"
)

type InventoryPredicate struct {
	Logger logr.Logger
}

func (p InventoryPredicate) Create(e event.CreateEvent) bool {
	// Always reconcile on creation
	return true
}

func (p InventoryPredicate) Update(e event.UpdateEvent) bool {
	newInventory := e.ObjectNew.(*configv1alpha1.Inventory)
	oldInventory := e.ObjectOld.(*configv1alpha1.Inventory)
	log := p.Logger.WithValues("predicate", "updateInventoryPredicate",
		"Inventory", newInventory.Name,
	)

	if oldInventory == nil {
		log.V(1).Info("old is nil")
		return true
	}

	if !reflect.DeepEqual(oldInventory.Spec, newInventory.Spec) {
		log.V(1).Info("Spec changed")
		return true
	}

	return false
}

func (p InventoryPredicate) Delete(e event.DeleteEvent) bool {
	// Always reconcile on deletion
	return true
}

func (p InventoryPredicate) Generic(e event.GenericEvent) bool {
	// Ignore generic
	return false
}

func PodPredicates(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldPod := e.ObjectOld.(*corev1.Pod)
			newPod := e.ObjectNew.(*corev1.Pod)

			if oldPod == nil {
				return true
			}

			if !reflect.DeepEqual(oldPod.Annotations, newPod.Annotations) {
				logger.V(1).Info("pod annotation changed. reconcile")
				return true
			}

			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			pod := e.Object.(*corev1.Pod)
			return pod.Annotations != nil
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}
