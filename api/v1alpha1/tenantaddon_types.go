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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantAddonSpec defines the desired state of TenantAddon.
type TenantAddonSpec struct {
	// ClusterRef references the TenantCluster to install this addon into.
	// +kubebuilder:validation:Required
	ClusterRef LocalObjectReference `json:"clusterRef"`

	// Addon specifies a known Butler addon name.
	// Use this for built-in addons like cilium, metallb, etc.
	// Mutually exclusive with Helm.
	// +optional
	Addon string `json:"addon,omitempty"`

	// Version is the addon version to install.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Helm specifies a custom Helm chart to install.
	// Use this for arbitrary Helm charts not built into Butler.
	// Mutually exclusive with Addon.
	// +optional
	Helm *HelmChartSpec `json:"helm,omitempty"`

	// Values are Helm values for customization.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`

	// DependsOn specifies other TenantAddons that must be ready first.
	// +optional
	DependsOn []LocalObjectReference `json:"dependsOn,omitempty"`
}

// HelmChartSpec defines a custom Helm chart to install.
type HelmChartSpec struct {
	// Repository is the Helm repository URL.
	// +kubebuilder:validation:Required
	Repository string `json:"repository"`

	// Chart is the chart name within the repository.
	// +kubebuilder:validation:Required
	Chart string `json:"chart"`

	// ReleaseName is the Helm release name.
	// If not specified, defaults to the TenantAddon name.
	// +optional
	ReleaseName string `json:"releaseName,omitempty"`

	// Namespace is the target namespace for the Helm release.
	// If not specified, a namespace is chosen based on the chart.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// CreateNamespace creates the namespace if it doesn't exist.
	// +kubebuilder:default=true
	// +optional
	CreateNamespace bool `json:"createNamespace,omitempty"`
}

// TenantAddonPhase represents the current phase of a TenantAddon.
// +kubebuilder:validation:Enum=Pending;Installing;Installed;Upgrading;Degraded;Failed;Deleting
type TenantAddonPhase string

const (
	// TenantAddonPhasePending indicates the addon is waiting to be installed.
	TenantAddonPhasePending TenantAddonPhase = "Pending"

	// TenantAddonPhaseInstalling indicates the addon is being installed.
	TenantAddonPhaseInstalling TenantAddonPhase = "Installing"

	// TenantAddonPhaseInstalled indicates the addon is installed and healthy.
	TenantAddonPhaseInstalled TenantAddonPhase = "Installed"

	// TenantAddonPhaseUpgrading indicates the addon is being upgraded.
	TenantAddonPhaseUpgrading TenantAddonPhase = "Upgrading"

	// TenantAddonPhaseDegraded indicates the addon is installed but unhealthy.
	TenantAddonPhaseDegraded TenantAddonPhase = "Degraded"

	// TenantAddonPhaseFailed indicates addon installation failed.
	TenantAddonPhaseFailed TenantAddonPhase = "Failed"

	// TenantAddonPhaseDeleting indicates the addon is being removed.
	TenantAddonPhaseDeleting TenantAddonPhase = "Deleting"
)

// TenantAddonStatus defines the observed state of TenantAddon.
type TenantAddonStatus struct {
	// Conditions represent the latest available observations.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current phase of the addon.
	// +optional
	Phase TenantAddonPhase `json:"phase,omitempty"`

	// InstalledVersion is the currently installed version.
	// +optional
	InstalledVersion string `json:"installedVersion,omitempty"`

	// HelmRelease contains Helm release information.
	// +optional
	HelmRelease *HelmReleaseStatus `json:"helmRelease,omitempty"`

	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastTransitionTime is when the phase last changed.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// Message provides human-readable status information.
	// +optional
	Message string `json:"message,omitempty"`
}

// HelmReleaseStatus contains Helm release information.
type HelmReleaseStatus struct {
	// Name is the Helm release name.
	Name string `json:"name"`

	// Namespace is the release namespace.
	Namespace string `json:"namespace"`

	// Revision is the Helm revision number.
	Revision int32 `json:"revision"`

	// Status is the Helm release status.
	Status string `json:"status"`
}

// TenantAddon condition types.
const (
	// TenantAddonConditionClusterReady indicates the target cluster is ready.
	TenantAddonConditionClusterReady = "ClusterReady"

	// TenantAddonConditionDependenciesMet indicates dependencies are satisfied.
	TenantAddonConditionDependenciesMet = "DependenciesMet"

	// TenantAddonConditionInstalled indicates the addon is installed.
	TenantAddonConditionInstalled = "Installed"

	// TenantAddonConditionHealthy indicates the addon is healthy.
	TenantAddonConditionHealthy = "Healthy"

	// TenantAddonConditionReady indicates the addon is fully ready.
	TenantAddonConditionReady = "Ready"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ta
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.clusterRef.name",description="Target cluster"
// +kubebuilder:printcolumn:name="Addon",type="string",JSONPath=".spec.addon",description="Addon name"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="Desired version"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current phase"
// +kubebuilder:printcolumn:name="Installed",type="string",JSONPath=".status.installedVersion",description="Installed version"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// TenantAddon is the Schema for the tenantaddons API.
// It represents an addon to be installed in a TenantCluster.
// Unlike addons defined in TenantCluster.spec.addons (which are monotonic),
// TenantAddons can be deleted to remove the addon from the cluster.
type TenantAddon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantAddonSpec   `json:"spec,omitempty"`
	Status TenantAddonStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantAddonList contains a list of TenantAddon.
type TenantAddonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenantAddon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TenantAddon{}, &TenantAddonList{})
}
