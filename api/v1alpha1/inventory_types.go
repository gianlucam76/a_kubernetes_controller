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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	InventoryFinalizer = "my.domain/inventory"

	InventoryKind = "Inventory"
)

const (
	// This is the annotation used by default by the Inventory controller
	// to decide who owns a pod
	DefaultAnnotation = "owned-by"
)

// InventorySpec defines the desired state of Inventory
type InventorySpec struct {
	// By default inventory controller looks for annotation "owned-by" to
	// decide who owns it. DefaultAnnotation if defined, overrides that.
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:validation:MaxLength=50
	// +optional
	AnnotationKey string `json:"annotationKey,omitempty"`
}

type UserData struct {
	Username string `json:"username,omitempty"`

	Pods []string `json:"pods,omitempty"`
}

// InventoryStatus defines the observed state of Inventory.
type InventoryStatus struct {
	// FailureMessage reports any error the controller hit
	// *optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	Data []UserData `json:"data,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=inv
// +kubebuilder:subresource:status

// Inventory is the Schema for the inventories API
type Inventory struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Inventory
	// +required
	Spec InventorySpec `json:"spec"`

	// status defines the observed state of Inventory
	// +optional
	Status InventoryStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// InventoryList contains a list of Inventory
type InventoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Inventory `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Inventory{}, &InventoryList{})
}
