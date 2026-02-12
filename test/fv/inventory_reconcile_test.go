package fv_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	configv1alpha1 "my.domain/inventory/api/v1alpha1"
)

var _ = Describe("Inventory", func() {
	It("Reconciliation updates ConfigMap", Label("FV"), func() {
		inventory := &configv1alpha1.Inventory{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "demo",
				Name:      "my-inventory",
			},
		}

		Expect(k8sClient.Create(context.TODO(), inventory)).To(Succeed())

		configMap := &corev1.ConfigMap{}
		Eventually(func() bool {
			err := k8sClient.Get(context.TODO(),
				types.NamespacedName{Name: inventory.Name, Namespace: inventory.Namespace},
				configMap)
			if err != nil {
				return false
			}
			return len(configMap.Data) == 0
		}, timeout, pollingInterval).Should(BeTrue())

		Consistently(func() bool {
			err := k8sClient.Get(context.TODO(),
				types.NamespacedName{Name: inventory.Name, Namespace: inventory.Namespace},
				configMap)
			if err != nil {
				return false
			}
			return len(configMap.Data) == 0
		}, timeout, pollingInterval).Should(BeTrue())

		pod := &corev1.Pod{} // TODO: add all fields including annotation
		Expect(k8sClient.Create(context.TODO(), pod)).To(Succeed())

		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: inventory.Name, Namespace: inventory.Namespace},
			configMap)).To(Succeed())
		Expect(len(configMap.Data)).To(Equal(1))

		currentInventory := &configv1alpha1.Inventory{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Namespace: inventory.Namespace, Name: inventory.Name},
			currentInventory)).To(Succeed())

		Expect(k8sClient.Delete(context.TODO(), currentInventory)).To(Succeed())

		err := k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: inventory.Name, Namespace: inventory.Namespace},
			configMap)
		Expect(err).ToNot(BeNil())
		Expect(apierrors.IsNotFound(err)).To(BeTrue())
	})
})
