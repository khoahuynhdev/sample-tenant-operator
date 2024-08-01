/*
Copyright 2024.

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

package v1

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:object:generate=false

type TenantValidator struct {
	client.Client
}

// ValidateDelete implements admission.CustomValidator.
func (v *TenantValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	panic("unimplemented")
}

// ValidateUpdate implements admission.CustomValidator.
func (v *TenantValidator) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (warnings admission.Warnings, err error) {
	panic("unimplemented")
}

func (v *TenantValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	// read the object, make sure it's Tenant type
	tenant, ok := obj.(*Tenant)
	if !ok {
		return nil, fmt.Errorf("unexpected object type, expected Tenant")
	}
	// get the list of existing namespaces
	var namespaces corev1.NamespaceList
	if err := v.List(ctx, &namespaces); err != nil {
		return nil, fmt.Errorf("failed to list namespaces %v", err)
	}

	for _, ns := range tenant.Spec.Namespaces {
		// check if namespace already exists
		if namespaceExists(namespaces, ns) {
			return nil, fmt.Errorf("namespace %s already exists", ns)
		}
	}
	return nil, nil
}

func namespaceExists(namespace corev1.NamespaceList, ns string) bool {
	for _, n := range namespace.Items {
		if n.ObjectMeta.Name == ns {
			return true
		}
	}
	return false
}

// log is for logging in this package.
var tenantlog = logf.Log.WithName("tenant-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Tenant) SetupWebhookWithManager(mgr ctrl.Manager) error {
	validator := &TenantValidator{mgr.GetClient()}
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(validator).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-multitenancy-codereliant-io-v1-tenant,mutating=false,failurePolicy=fail,sideEffects=None,groups=multitenancy.codereliant.io,resources=tenants,verbs=create;update,versions=v1,name=vtenant.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Tenant{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Tenant) ValidateCreate() (admission.Warnings, error) {
	tenantlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Tenant) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	tenantlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Tenant) ValidateDelete() (admission.Warnings, error) {
	tenantlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
