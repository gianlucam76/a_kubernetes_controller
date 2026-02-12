package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *InventoryReconciler) requeueInventoryForPods(
	ctx context.Context, o client.Object,
) []reconcile.Request {
	// Get all Inventory instances in the object namespace

	items, ok := r.inventories[o.GetNamespace()]
	if !ok {
		return nil
	}

	result := make([]reconcile.Request, len(items))
	for i := range items {
		result[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: o.GetNamespace(),
				Name:      items[i],
			},
		}
	}

	return result
}
