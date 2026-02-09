package controller

import (
	"reflect"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/event"

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
