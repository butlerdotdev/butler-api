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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MultiTenancyMode defines how multi-tenancy is enforced across the platform.
// +kubebuilder:validation:Enum=Enforced;Optional;Disabled
type MultiTenancyMode string

const (
	// MultiTenancyModeEnforced requires all TenantClusters to belong to a Team.
	// Recommended for enterprise deployments with multiple teams.
	MultiTenancyModeEnforced MultiTenancyMode = "Enforced"

	// MultiTenancyModeOptional allows Teams but doesn't require them.
	// Useful during migration or for mixed deployment models.
	MultiTenancyModeOptional MultiTenancyMode = "Optional"

	// MultiTenancyModeDisabled disables Team functionality entirely.
	// Simplest mode for demos and single-user deployments.
	MultiTenancyModeDisabled MultiTenancyMode = "Disabled"
)

// ButlerConfigSpec defines the desired state of ButlerConfig.
// ButlerConfig is a singleton that configures platform-wide settings.
type ButlerConfigSpec struct {
	// MultiTenancy configures how multi-tenancy is handled.
	// +optional
	MultiTenancy MultiTenancyConfig `json:"multiTenancy,omitempty"`

	// DefaultNamespace is the namespace for TenantClusters when not using Teams.
	// +kubebuilder:default="butler-tenants"
	// +optional
	DefaultNamespace string `json:"defaultNamespace,omitempty"`

	// DefaultProviderConfigRef references the default ProviderConfig.
	// +optional
	DefaultProviderConfigRef *LocalObjectReference `json:"defaultProviderConfigRef,omitempty"`

	// DefaultTeamLimits are the default resource limits for new Teams.
	// +optional
	DefaultTeamLimits *ResourceLimits `json:"defaultTeamLimits,omitempty"`

	// DefaultAddonVersions specifies default versions for platform addons.
	// +optional
	DefaultAddonVersions *AddonVersions `json:"defaultAddonVersions,omitempty"`

	// ControlPlane configures platform-wide control plane exposure settings.
	// +optional
	ControlPlane *PlatformControlPlaneConfig `json:"controlPlane,omitempty"`
}

// MultiTenancyConfig configures multi-tenancy behavior.
type MultiTenancyConfig struct {
	// Mode determines how multi-tenancy is enforced.
	// +kubebuilder:default="Disabled"
	// +optional
	Mode MultiTenancyMode `json:"mode,omitempty"`
}

// ResourceLimits defines resource limits for Teams.
type ResourceLimits struct {
	// MaxClusters is the maximum number of TenantClusters a Team can create.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MaxClusters *int32 `json:"maxClusters,omitempty"`

	// MaxWorkersPerCluster is the maximum workers per TenantCluster.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MaxWorkersPerCluster *int32 `json:"maxWorkersPerCluster,omitempty"`

	// MaxTotalCPU is the maximum total CPU cores across all clusters.
	// +optional
	MaxTotalCPU *resource.Quantity `json:"maxTotalCPU,omitempty"`

	// MaxTotalMemory is the maximum total memory across all clusters.
	// +optional
	MaxTotalMemory *resource.Quantity `json:"maxTotalMemory,omitempty"`

	// MaxTotalStorage is the maximum total storage across all clusters.
	// +optional
	MaxTotalStorage *resource.Quantity `json:"maxTotalStorage,omitempty"`
}

// AddonVersions specifies default versions for Butler-managed addons.
type AddonVersions struct {
	// Cilium is the default Cilium CNI version.
	// +optional
	Cilium string `json:"cilium,omitempty"`

	// MetalLB is the default MetalLB load balancer version.
	// +optional
	MetalLB string `json:"metallb,omitempty"`

	// CertManager is the default cert-manager version.
	// +optional
	CertManager string `json:"certManager,omitempty"`

	// Longhorn is the default Longhorn storage version.
	// +optional
	Longhorn string `json:"longhorn,omitempty"`

	// Traefik is the default Traefik ingress controller version.
	// +optional
	Traefik string `json:"traefik,omitempty"`

	// FluxCD is the default Flux CD version for GitOps.
	// +optional
	FluxCD string `json:"fluxcd,omitempty"`
}

// PlatformControlPlaneConfig defines platform-level control plane settings.
type PlatformControlPlaneConfig struct {
	// DefaultExposureMode is the default exposure mode for new TenantClusters.
	// TenantClusters can override this in their spec.
	// If not specified, defaults to LoadBalancer for backward compatibility.
	//
	// Options:
	// - Gateway: SNI-based routing via Gateway API (most IP-efficient)
	// - LoadBalancer: Dedicated LoadBalancer IP per cluster
	// - NodePort: NodePort service per cluster
	// +kubebuilder:default="LoadBalancer"
	// +optional
	DefaultExposureMode ControlPlaneExposureMode `json:"defaultExposureMode,omitempty"`

	// Gateway configures Gateway API exposure settings for control planes.
	// Required when DefaultExposureMode is Gateway or any TenantCluster uses Gateway mode.
	// Butler creates and manages the Gateway resource.
	// +optional
	Gateway *GatewayConfig `json:"gateway,omitempty"`
}

// GatewayConfig defines Gateway API configuration for control plane exposure.
// When configured, Butler manages a Gateway resource that routes traffic to
// tenant control planes based on SNI (Server Name Indication) hostname.
// This allows multiple tenant API servers to share a single external IP.
type GatewayConfig struct {
	// Domain is the base domain for control plane hostnames.
	// TenantClusters are exposed as {cluster-name}.{domain}.
	// Example: "k8s.example.com" results in hostnames like "tenant-a.k8s.example.com"
	//
	// DNS must be configured with a wildcard record pointing to the Gateway address:
	//   *.k8s.example.com -> <Gateway IP>
	//
	// Required when using Gateway exposure mode.
	// +kubebuilder:validation:MinLength=1
	Domain string `json:"domain"`

	// GatewayName is the name of the Gateway resource Butler manages.
	// Butler creates this Gateway and maintains ownership of it.
	// +kubebuilder:default="butler-control-plane"
	// +optional
	GatewayName string `json:"gatewayName,omitempty"`

	// GatewayNamespace is the namespace where the Gateway resource is created.
	// Should be a namespace where the Gateway controller has permissions.
	// +kubebuilder:default="butler-system"
	// +optional
	GatewayNamespace string `json:"gatewayNamespace,omitempty"`

	// GatewayClassName is the GatewayClass to use for the managed Gateway.
	// Must reference an existing GatewayClass in the cluster.
	// Common values: "cilium" (recommended), "istio", "envoy-gateway"
	// +kubebuilder:default="cilium"
	// +optional
	GatewayClassName string `json:"gatewayClassName,omitempty"`

	// Annotations are additional annotations to apply to the Gateway resource.
	// Use for Gateway controller-specific configuration (e.g., load balancer settings).
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ButlerConfigStatus defines the observed state of ButlerConfig.
type ButlerConfigStatus struct {
	// Conditions represent the latest available observations of the config's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ObservedGeneration is the spec generation that was last reconciled.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// TeamCount is the current number of Team resources.
	// +optional
	TeamCount int32 `json:"teamCount,omitempty"`

	// ClusterCount is the current number of TenantCluster resources.
	// +optional
	ClusterCount int32 `json:"clusterCount,omitempty"`

	// Gateway contains the status of the managed Gateway resource.
	// Only populated when Gateway exposure mode is configured.
	// +optional
	Gateway *GatewayStatus `json:"gateway,omitempty"`
}

// GatewayStatus contains the status of the Butler-managed Gateway resource.
type GatewayStatus struct {
	// Ready indicates the Gateway is programmed and ready to accept traffic.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Address is the Gateway's external address (IP or hostname).
	// This is the address that DNS wildcard records should point to.
	// +optional
	Address string `json:"address,omitempty"`

	// ListenerCount is the number of active listeners on the Gateway.
	// Should be 2 when healthy (kube-apiserver on 6443, konnectivity on 8132).
	// +optional
	ListenerCount int32 `json:"listenerCount,omitempty"`

	// TenantCount is the number of TenantClusters using this Gateway.
	// +optional
	TenantCount int32 `json:"tenantCount,omitempty"`

	// Message provides additional information about the Gateway status.
	// +optional
	Message string `json:"message,omitempty"`
}

// ButlerConfig condition types.
const (
	// ButlerConfigConditionGatewayReady indicates the managed Gateway is ready.
	ButlerConfigConditionGatewayReady = "GatewayReady"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=bc
// +kubebuilder:printcolumn:name="Mode",type="string",JSONPath=".spec.multiTenancy.mode",description="Multi-tenancy mode"
// +kubebuilder:printcolumn:name="Exposure",type="string",JSONPath=".spec.controlPlane.defaultExposureMode",description="Default CP exposure"
// +kubebuilder:printcolumn:name="Teams",type="integer",JSONPath=".status.teamCount",description="Number of teams"
// +kubebuilder:printcolumn:name="Clusters",type="integer",JSONPath=".status.clusterCount",description="Number of clusters"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ButlerConfig is the Schema for the butlerconfigs API.
// It is a singleton resource that configures platform-wide Butler settings.
// Only one ButlerConfig named "butler" should exist in the cluster.
type ButlerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ButlerConfigSpec   `json:"spec,omitempty"`
	Status ButlerConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ButlerConfigList contains a list of ButlerConfig.
type ButlerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ButlerConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ButlerConfig{}, &ButlerConfigList{})
}

// GetGatewayDomain returns the configured gateway domain.
// Returns empty string if gateway is not configured.
func (bc *ButlerConfig) GetGatewayDomain() string {
	if bc.Spec.ControlPlane != nil && bc.Spec.ControlPlane.Gateway != nil {
		return bc.Spec.ControlPlane.Gateway.Domain
	}
	return ""
}

// GetGatewayName returns the gateway name with default fallback.
func (bc *ButlerConfig) GetGatewayName() string {
	if bc.Spec.ControlPlane != nil && bc.Spec.ControlPlane.Gateway != nil && bc.Spec.ControlPlane.Gateway.GatewayName != "" {
		return bc.Spec.ControlPlane.Gateway.GatewayName
	}
	return "butler-control-plane"
}

// GetGatewayNamespace returns the gateway namespace with default fallback.
func (bc *ButlerConfig) GetGatewayNamespace() string {
	if bc.Spec.ControlPlane != nil && bc.Spec.ControlPlane.Gateway != nil && bc.Spec.ControlPlane.Gateway.GatewayNamespace != "" {
		return bc.Spec.ControlPlane.Gateway.GatewayNamespace
	}
	return "butler-system"
}

// GetGatewayClassName returns the gateway class name with default fallback.
func (bc *ButlerConfig) GetGatewayClassName() string {
	if bc.Spec.ControlPlane != nil && bc.Spec.ControlPlane.Gateway != nil && bc.Spec.ControlPlane.Gateway.GatewayClassName != "" {
		return bc.Spec.ControlPlane.Gateway.GatewayClassName
	}
	return "cilium"
}

// GetDefaultExposureMode returns the default exposure mode with fallback.
func (bc *ButlerConfig) GetDefaultExposureMode() ControlPlaneExposureMode {
	if bc.Spec.ControlPlane != nil && bc.Spec.ControlPlane.DefaultExposureMode != "" {
		return bc.Spec.ControlPlane.DefaultExposureMode
	}
	return ControlPlaneExposureModeLoadBalancer
}

// IsGatewayConfigured returns true if gateway configuration is present and valid.
func (bc *ButlerConfig) IsGatewayConfigured() bool {
	return bc.Spec.ControlPlane != nil &&
		bc.Spec.ControlPlane.Gateway != nil &&
		bc.Spec.ControlPlane.Gateway.Domain != ""
}
