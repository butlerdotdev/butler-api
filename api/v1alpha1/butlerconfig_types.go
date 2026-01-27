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

// MultiTenancyMode defines how multi-tenancy is enforced.
// +kubebuilder:validation:Enum=Enforced;Optional;Disabled
type MultiTenancyMode string

const (
	// MultiTenancyModeEnforced requires all TenantClusters to belong to a Team.
	// Teams must exist before TenantClusters can be created.
	// This is the recommended mode for enterprise deployments.
	MultiTenancyModeEnforced MultiTenancyMode = "Enforced"

	// MultiTenancyModeOptional allows Teams but doesn't require them.
	// TenantClusters can exist in the default namespace without a Team.
	MultiTenancyModeOptional MultiTenancyMode = "Optional"

	// MultiTenancyModeDisabled disables Team functionality.
	// All TenantClusters exist in the default namespace.
	// This is the simplest mode for demos and single-user deployments.
	MultiTenancyModeDisabled MultiTenancyMode = "Disabled"
)

// ButlerConfigSpec defines the desired state of ButlerConfig.
type ButlerConfigSpec struct {
	// MultiTenancy configures how multi-tenancy is handled.
	// +optional
	MultiTenancy MultiTenancyConfig `json:"multiTenancy,omitempty"`

	// DefaultNamespace is the namespace for TenantClusters when not using Teams.
	// Used in Disabled and Optional modes.
	// +kubebuilder:default="butler-tenants"
	// +optional
	DefaultNamespace string `json:"defaultNamespace,omitempty"`

	// DefaultProviderConfigRef references the default ProviderConfig.
	// Used when Teams or TenantClusters don't specify their own.
	// +optional
	DefaultProviderConfigRef *LocalObjectReference `json:"defaultProviderConfigRef,omitempty"`

	// DefaultTeamLimits are the default resource limits for new Teams.
	// Admins can override these when creating individual Teams.
	// +optional
	DefaultTeamLimits *ResourceLimits `json:"defaultTeamLimits,omitempty"`

	// DefaultAddonVersions specifies the default versions for addons.
	// Used when TenantCluster doesn't specify versions.
	// +optional
	DefaultAddonVersions *AddonVersions `json:"defaultAddonVersions,omitempty"`

	// GitProvider configures the default Git provider for GitOps operations.
	// This enables features like exporting clusters to GitOps, enabling Flux
	// on clusters, and managing addons via Git repositories.
	// +optional
	GitProvider *GitProviderConfig `json:"gitProvider,omitempty"`
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
	// Cilium version.
	// +optional
	Cilium string `json:"cilium,omitempty"`

	// MetalLB version.
	// +optional
	MetalLB string `json:"metallb,omitempty"`

	// CertManager version.
	// +optional
	CertManager string `json:"certManager,omitempty"`

	// Longhorn version.
	// +optional
	Longhorn string `json:"longhorn,omitempty"`

	// Traefik version.
	// +optional
	Traefik string `json:"traefik,omitempty"`

	// FluxCD version.
	// +optional
	FluxCD string `json:"fluxcd,omitempty"`
}

// ButlerConfigStatus defines the observed state of ButlerConfig.
type ButlerConfigStatus struct {
	// Conditions represent the latest available observations of the config's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// TeamCount is the current number of Teams.
	// +optional
	TeamCount int32 `json:"teamCount,omitempty"`

	// ClusterCount is the current number of TenantClusters.
	// +optional
	ClusterCount int32 `json:"clusterCount,omitempty"`

	// GitProvider shows the status of the configured Git provider.
	// +optional
	GitProvider *GitProviderStatus `json:"gitProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=bc
// +kubebuilder:printcolumn:name="Mode",type="string",JSONPath=".spec.multiTenancy.mode",description="Multi-tenancy mode"
// +kubebuilder:printcolumn:name="Git",type="string",JSONPath=".spec.gitProvider.type",description="Git provider"
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

// Helper methods

// IsGitProviderConfigured returns true if a Git provider is configured.
func (c *ButlerConfig) IsGitProviderConfigured() bool {
	return c.Spec.GitProvider != nil && c.Spec.GitProvider.SecretRef.Name != ""
}

// GetGitProviderURL returns the Git provider URL with a sensible default.
func (c *ButlerConfig) GetGitProviderURL() string {
	if c.Spec.GitProvider == nil || c.Spec.GitProvider.URL == "" {
		return "https://api.github.com"
	}
	return c.Spec.GitProvider.URL
}
