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

// TeamResourceLimits defines resource quotas and restrictions for a Team.
// This is separate from ResourceLimits in butlerconfig_types.go which defines
// platform-wide defaults. TeamResourceLimits includes additional fields for
// feature restrictions that are team-specific.
type TeamResourceLimits struct {
	// ====== Cluster Limits ======

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

	// ====== Compute Limits ======

	// MaxCPUCores is the maximum total CPU cores across all clusters.
	// +optional
	MaxCPUCores *resource.Quantity `json:"maxCPUCores,omitempty"`

	// MaxMemory is the maximum total memory across all clusters.
	// +optional
	MaxMemory *resource.Quantity `json:"maxMemory,omitempty"`

	// MaxStorage is the maximum total storage across all clusters.
	// +optional
	MaxStorage *resource.Quantity `json:"maxStorage,omitempty"`

	// ====== Per-Cluster Defaults ======

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

	// ====== Feature Restrictions ======

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

	// ====== Utilization Percentages ======

	// ClusterUtilization is percentage of MaxClusters used.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	ClusterUtilization *int32 `json:"clusterUtilization,omitempty"`

	// NodeUtilization is percentage of MaxTotalNodes used.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	NodeUtilization *int32 `json:"nodeUtilization,omitempty"`

	// CPUUtilization is percentage of MaxCPUCores used.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	CPUUtilization *int32 `json:"cpuUtilization,omitempty"`

	// MemoryUtilization is percentage of MaxMemory used.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	MemoryUtilization *int32 `json:"memoryUtilization,omitempty"`
}

// Kubernetes recommended labels.
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

	// LabelTenant identifies the tenant cluster.
	LabelTenant = "butler.butlerlabs.dev/tenant"

	// LabelSourceNamespace indicates the source namespace for generated resources.
	LabelSourceNamespace = "butler.butlerlabs.dev/source-namespace"

	// LabelSourceName indicates the source name for generated resources.
	LabelSourceName = "butler.butlerlabs.dev/source-name"

	// LabelNetworkPool identifies the NetworkPool associated with a resource.
	LabelNetworkPool = "butler.butlerlabs.dev/network-pool"

	// LabelProviderConfig identifies the ProviderConfig associated with a resource.
	LabelProviderConfig = "butler.butlerlabs.dev/provider-config"

	// LabelWorkspaceOwner identifies the owner of a workspace (hashed email).
	LabelWorkspaceOwner = "butler.butlerlabs.dev/workspace-owner"

	// LabelAllocationType identifies the IP allocation type (loadbalancer, nodes).
	LabelAllocationType = "butler.butlerlabs.dev/allocation-type"
)

// Butler-specific annotations.
const (
	// AnnotationDescription provides a human-readable description.
	AnnotationDescription = "butler.butlerlabs.dev/description"

	// AnnotationCreatedBy indicates who created the resource.
	AnnotationCreatedBy = "butler.butlerlabs.dev/created-by"

	// AnnotationConnect signals the controller to create/tear down the SSH service.
	AnnotationConnect = "butler.butlerlabs.dev/connect"

	// AnnotationConnectTime records when the SSH service was created.
	AnnotationConnectTime = "butler.butlerlabs.dev/connect-time"
)

// Finalizers.
const (
	// FinalizerTeam is the finalizer for Team resources.
	FinalizerTeam = "butler.butlerlabs.dev/team"

	// FinalizerTenantCluster is the finalizer for TenantCluster resources.
	FinalizerTenantCluster = "butler.butlerlabs.dev/tenantcluster"

	// FinalizerTenantAddon is the finalizer for TenantAddon resources.
	FinalizerTenantAddon = "butler.butlerlabs.dev/tenantaddon"

	// FinalizerUser is the finalizer for User resources.
	FinalizerUser = "butler.butlerlabs.dev/user"

	// FinalizerNetworkPool is the finalizer for NetworkPool resources.
	FinalizerNetworkPool = "butler.butlerlabs.dev/networkpool"

	// FinalizerIPAllocation is the finalizer for IPAllocation resources.
	FinalizerIPAllocation = "butler.butlerlabs.dev/ipallocation"

	// FinalizerProviderConfig is the finalizer for ProviderConfig resources.
	FinalizerProviderConfig = "butler.butlerlabs.dev/providerconfig"

	// FinalizerWorkspace is the finalizer for Workspace resources.
	FinalizerWorkspace = "butler.butlerlabs.dev/workspace"
)

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

	// ReasonQuotaExceeded indicates a resource quota was exceeded.
	ReasonQuotaExceeded = "QuotaExceeded"

	// ReasonPoolExhausted indicates a NetworkPool has no available IPs.
	ReasonPoolExhausted = "PoolExhausted"

	// ReasonAllocationFailed indicates an IP allocation failed.
	ReasonAllocationFailed = "AllocationFailed"

	// ReasonProviderAccessDenied indicates the team does not have access to the provider.
	ReasonProviderAccessDenied = "ProviderAccessDenied"

	// ReasonNetworkNotReady indicates the network is not ready.
	ReasonNetworkNotReady = "NetworkNotReady"

	// ReasonCredentialsInvalid indicates provider credentials are invalid.
	ReasonCredentialsInvalid = "CredentialsInvalid"

	// ReasonNetworkReachable indicates the network is reachable.
	ReasonNetworkReachable = "NetworkReachable"

	// ReasonPoolAvailable indicates the network pool has capacity.
	ReasonPoolAvailable = "PoolAvailable"
)
