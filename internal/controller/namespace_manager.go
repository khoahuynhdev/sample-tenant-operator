package controller

import (
	"context"

	multitenancyv1 "codereliant.io/tenant/api/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	tenantOperatorAnnotation = "tenant-operator"
)

func (r *TenantReconciler) EnsureNamespace(ctx context.Context, tenant *multitenancyv1.Tenant, namespaceName string) error {
	log := log.FromContext(ctx)
	// define a namespace object
	namespace := &corev1.Namespace{}

	// attempt to get the namespace with provided name
	err := r.Get(ctx, client.ObjectKey{Namespace: namespaceName}, namespace)
	if err != nil {
		// if the namespace doesn't exist, create it
		if apierrors.IsNotFound(err) {
			log.Info("Creating Namespace", "namespace", namespaceName)
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespaceName,
					Annotations: map[string]string{
						"adminEmail": tenant.Spec.AdminEmail,
						"managed-by": tenantOperatorAnnotation,
					},
				},
			}
			// attempt to create the namespace
			if err = r.Create(ctx, namespace); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		log.Info("Namespace already existed", "namespace", namespace)
	}
	return nil
}
