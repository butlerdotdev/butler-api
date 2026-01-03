/*
Copyright 2026 The Butler Authors.

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

// ProviderReference references a ProviderConfig resource.
type ProviderReference struct {
	// Name is the name of the ProviderConfig resource.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Namespace is the namespace of the ProviderConfig resource.
	// If not specified, the namespace of the referencing resource is used.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// SecretReference references a Secret resource.
type SecretReference struct {
	// Name is the name of the Secret.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Namespace is the namespace of the Secret.
	// If not specified, the namespace of the referencing resource is used.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Key is the key within the Secret to reference.
	// If not specified, the entire Secret data is used.
	// +optional
	Key string `json:"key,omitempty"`
}

// LocalObjectReference references a resource in the same namespace.
type LocalObjectReference struct {
	// Name is the name of the resource.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// NamespacedObjectReference references a resource in any namespace.
type NamespacedObjectReference struct {
	// Name is the name of the resource.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Namespace is the namespace of the resource.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Namespace string `json:"namespace"`
}

// Condition types following Kubernetes API conventions.
// See: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties
const (
	// ConditionTypeReady indicates the resource is ready for use.
	ConditionTypeReady = "Ready"

	// ConditionTypeProgressing indicates the resource is making progress toward Ready.
	ConditionTypeProgressing = "Progressing"

	// ConditionTypeDegraded indicates the resource is in a degraded state.
	ConditionTypeDegraded = "Degraded"
)

// Condition reasons for MachineRequest.
const (
	// ReasonPending indicates the request is waiting to be processed.
	ReasonPending = "Pending"

	// ReasonCreating indicates the resource is being created.
	ReasonCreating = "Creating"

	// ReasonCreated indicates the resource was successfully created.
	ReasonCreated = "Created"

	// ReasonRunning indicates the resource is running.
	ReasonRunning = "Running"

	// ReasonWaitingForIP indicates waiting for IP address assignment.
	ReasonWaitingForIP = "WaitingForIP"

	// ReasonFailed indicates the operation failed.
	ReasonFailed = "Failed"

	// ReasonDeleting indicates the resource is being deleted.
	ReasonDeleting = "Deleting"

	// ReasonDeleted indicates the resource was deleted.
	ReasonDeleted = "Deleted"

	// ReasonProviderError indicates an error from the infrastructure provider.
	ReasonProviderError = "ProviderError"

	// ReasonInvalidConfiguration indicates invalid configuration.
	ReasonInvalidConfiguration = "InvalidConfiguration"

	// ReasonReady indicates the resource is ready.
	ReasonReady = "Ready"

	// ReasonWaitingForDependencies indicates waiting for dependencies.
	ReasonWaitingForDependencies = "WaitingForDependencies"

	// ReasonReconciling indicates active reconciliation.
	ReasonReconciling = "Reconciling"

	// ReasonValidationFailed indicates validation failed.
	ReasonValidationFailed = "ValidationFailed"
)

// Butler-specific labels.
const (
	// LabelTeam identifies the team that owns a resource.
	LabelTeam = "butler.butlerlabs.dev/team"

	// LabelTenant identifies the tenant cluster.
	LabelTenant = "butler.butlerlabs.dev/tenant"

	// LabelManagedBy indicates the resource is managed by Butler.
	LabelManagedBy = "butler.butlerlabs.dev/managed-by"

	// LabelSourceNamespace indicates the source namespace for generated resources.
	LabelSourceNamespace = "butler.butlerlabs.dev/source-namespace"

	// LabelSourceName indicates the source name for generated resources.
	LabelSourceName = "butler.butlerlabs.dev/source-name"
)

// Butler-specific annotations.
const (
	// AnnotationDescription provides a human-readable description.
	AnnotationDescription = "butler.butlerlabs.dev/description"

	// AnnotationCreatedBy indicates who created the resource.
	AnnotationCreatedBy = "butler.butlerlabs.dev/created-by"
)

// Finalizers.
const (
	// FinalizerTeam is the finalizer for Team resources.
	FinalizerTeam = "butler.butlerlabs.dev/team"

	// FinalizerTenantCluster is the finalizer for TenantCluster resources.
	FinalizerTenantCluster = "butler.butlerlabs.dev/tenantcluster"

	// FinalizerTenantAddon is the finalizer for TenantAddon resources.
	FinalizerTenantAddon = "butler.butlerlabs.dev/tenantaddon"
)
