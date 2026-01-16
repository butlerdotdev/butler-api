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

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

// ControlPlaneExposureMode defines how the control plane is exposed to clients.
// This determines the network path for kubectl and other API clients to reach
// the tenant cluster's Kubernetes API server.
// +kubebuilder:validation:Enum=Gateway;LoadBalancer;NodePort
type ControlPlaneExposureMode string

const (
	// ControlPlaneExposureModeGateway exposes the control plane via Gateway API
	// with SNI-based TLS routing. Multiple tenant control planes share a single
	// Gateway IP, differentiated by hostname (e.g., tenant-a.k8s.example.com).
	//
	// This is the most IP-efficient option, recommended for environments with
	// limited IP addresses or many tenant clusters.
	//
	// Requirements:
	// - Gateway API CRDs installed
	// - Gateway controller (e.g., Cilium Gateway API)
	// - DNS wildcard record pointing to Gateway IP
	// - ButlerConfig.spec.controlPlane.gateway configured
	ControlPlaneExposureModeGateway ControlPlaneExposureMode = "Gateway"

	// ControlPlaneExposureModeLoadBalancer exposes the control plane via a
	// dedicated LoadBalancer Service. Each tenant cluster gets its own external
	// IP address assigned by MetalLB or a cloud provider.
	//
	// This is the simplest option but consumes one IP per tenant cluster.
	// Good for environments with ample IP addresses or few clusters.
	ControlPlaneExposureModeLoadBalancer ControlPlaneExposureMode = "LoadBalancer"

	// ControlPlaneExposureModeNodePort exposes the control plane via a NodePort
	// Service. Clients connect to any management cluster node IP on the
	// allocated NodePort.
	//
	// Useful for edge deployments, air-gapped environments, or when LoadBalancer
	// and Gateway options are unavailable.
	ControlPlaneExposureModeNodePort ControlPlaneExposureMode = "NodePort"
)

// ProviderReference references a ProviderConfig resource.
// Used when the provider may be in a different namespace.
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
// This is the most common reference type within Butler.
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

// TeamResourceLimits defines resource quotas and restrictions for a Team.
// This allows platform administrators to control resource consumption per team.
type TeamResourceLimits struct {
	// MaxClusters is the maximum number of TenantClusters this team can create.
	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxClusters *int32 `json:"maxClusters,omitempty"`

	// MaxNodesPerCluster is the maximum worker nodes per cluster.
	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxNodesPerCluster *int32 `json:"maxNodesPerCluster,omitempty"`

	// MaxTotalNodes is the maximum total worker nodes across all clusters.
	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxTotalNodes *int32 `json:"maxTotalNodes,omitempty"`

	// MaxCPUCores is the maximum total CPU cores across all clusters.
	// +optional
	MaxCPUCores *resource.Quantity `json:"maxCPUCores,omitempty"`

	// MaxMemory is the maximum total memory across all clusters.
	// +optional
	MaxMemory *resource.Quantity `json:"maxMemory,omitempty"`

	// MaxStorage is the maximum total storage across all clusters.
	// +optional
	MaxStorage *resource.Quantity `json:"maxStorage,omitempty"`

	// DefaultNodeCount is the default worker count for new clusters.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=3
	DefaultNodeCount *int32 `json:"defaultNodeCount,omitempty"`

	// DefaultCPUPerNode is the default CPU cores per worker node.
	// +optional
	DefaultCPUPerNode *resource.Quantity `json:"defaultCPUPerNode,omitempty"`

	// DefaultMemoryPerNode is the default memory per worker node.
	// +optional
	DefaultMemoryPerNode *resource.Quantity `json:"defaultMemoryPerNode,omitempty"`

	// AllowedKubernetesVersions restricts which K8s versions can be used.
	// If empty, all supported versions are allowed.
	// +optional
	AllowedKubernetesVersions []string `json:"allowedKubernetesVersions,omitempty"`

	// AllowedProviders restricts which ProviderConfigs can be used.
	// If empty, all providers the team has access to are allowed.
	// +optional
	AllowedProviders []string `json:"allowedProviders,omitempty"`

	// AllowedAddons restricts which addons can be installed.
	// If empty, all addons are allowed.
	// +optional
	AllowedAddons []string `json:"allowedAddons,omitempty"`

	// DeniedAddons explicitly denies certain addons.
	// Takes precedence over AllowedAddons.
	// +optional
	DeniedAddons []string `json:"deniedAddons,omitempty"`
}

// TeamResourceUsage shows current resource consumption for a Team.
// Used for monitoring and quota enforcement.
type TeamResourceUsage struct {
	// Clusters is the number of TenantClusters.
	// +optional
	Clusters int32 `json:"clusters,omitempty"`

	// TotalNodes is the total number of worker nodes.
	// +optional
	TotalNodes int32 `json:"totalNodes,omitempty"`

	// TotalCPU is the total CPU cores allocated.
	// +optional
	TotalCPU *resource.Quantity `json:"totalCPU,omitempty"`

	// TotalMemory is the total memory allocated.
	// +optional
	TotalMemory *resource.Quantity `json:"totalMemory,omitempty"`

	// TotalStorage is the total storage allocated.
	// +optional
	TotalStorage *resource.Quantity `json:"totalStorage,omitempty"`

	// ClusterUtilization is percentage of MaxClusters used (0-100).
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	ClusterUtilization *int32 `json:"clusterUtilization,omitempty"`

	// NodeUtilization is percentage of MaxTotalNodes used (0-100).
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	NodeUtilization *int32 `json:"nodeUtilization,omitempty"`

	// CPUUtilization is percentage of MaxCPUCores used (0-100).
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	CPUUtilization *int32 `json:"cpuUtilization,omitempty"`

	// MemoryUtilization is percentage of MaxMemory used (0-100).
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	MemoryUtilization *int32 `json:"memoryUtilization,omitempty"`
}

// Kubernetes recommended labels for resource organization.
// See: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
const (
	// LabelManagedBy indicates the tool managing the resource.
	// Uses the Kubernetes standard label for interoperability with
	// kubectl, Helm, ArgoCD, Prometheus, and other CNCF tools.
	LabelManagedBy = "app.kubernetes.io/managed-by"
)

// Butler-specific labels for resource tracking and multi-tenancy.
const (
	// LabelTeam identifies the team that owns a resource.
	LabelTeam = "butler.butlerlabs.dev/team"

	// LabelTenant identifies the tenant cluster a resource belongs to.
	LabelTenant = "butler.butlerlabs.dev/tenant"

	// LabelSourceNamespace indicates the source namespace for generated resources.
	LabelSourceNamespace = "butler.butlerlabs.dev/source-namespace"

	// LabelSourceName indicates the source name for generated resources.
	LabelSourceName = "butler.butlerlabs.dev/source-name"
)

// Butler-specific annotations.
const (
	// AnnotationDescription provides a human-readable description.
	AnnotationDescription = "butler.butlerlabs.dev/description"

	// AnnotationCreatedBy indicates who created the resource (user, system, etc.).
	AnnotationCreatedBy = "butler.butlerlabs.dev/created-by"
)

// Finalizers used by Butler controllers for cleanup.
const (
	// FinalizerTeam ensures Team cleanup before deletion.
	FinalizerTeam = "butler.butlerlabs.dev/team"

	// FinalizerTenantCluster ensures TenantCluster cleanup before deletion.
	FinalizerTenantCluster = "butler.butlerlabs.dev/tenantcluster"

	// FinalizerTenantAddon ensures TenantAddon cleanup before deletion.
	FinalizerTenantAddon = "butler.butlerlabs.dev/tenantaddon"

	// FinalizerUser ensures User cleanup before deletion.
	FinalizerUser = "butler.butlerlabs.dev/user"
)

// Condition types following Kubernetes API conventions.
// See: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md
const (
	// ConditionTypeReady indicates the resource is ready for use.
	ConditionTypeReady = "Ready"

	// ConditionTypeProgressing indicates the resource is making progress toward Ready.
	ConditionTypeProgressing = "Progressing"

	// ConditionTypeDegraded indicates the resource is in a degraded state.
	ConditionTypeDegraded = "Degraded"
)

// Condition reasons used across Butler resources.
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

	// ReasonReconciling indicates active reconciliation in progress.
	ReasonReconciling = "Reconciling"

	// ReasonValidationFailed indicates validation failed.
	ReasonValidationFailed = "ValidationFailed"

	// ReasonQuotaExceeded indicates a resource quota was exceeded.
	ReasonQuotaExceeded = "QuotaExceeded"

	// ReasonGatewayNotConfigured indicates Gateway mode was requested but Gateway is not configured.
	ReasonGatewayNotConfigured = "GatewayNotConfigured"
)
