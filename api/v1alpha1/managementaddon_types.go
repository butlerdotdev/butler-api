/*
Copyright 2024 The Butler Authors.

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
	"k8s.io/apimachinery/pkg/runtime"
)

// ManagementAddonPhase represents the current phase of a management addon
// +kubebuilder:validation:Enum=Pending;Installing;Installed;Upgrading;Failed;Uninstalling
type ManagementAddonPhase string

const (
	ManagementAddonPhasePending      ManagementAddonPhase = "Pending"
	ManagementAddonPhaseInstalling   ManagementAddonPhase = "Installing"
	ManagementAddonPhaseInstalled    ManagementAddonPhase = "Installed"
	ManagementAddonPhaseUpgrading    ManagementAddonPhase = "Upgrading"
	ManagementAddonPhaseFailed       ManagementAddonPhase = "Failed"
	ManagementAddonPhaseUninstalling ManagementAddonPhase = "Uninstalling"
)

// ManagementAddonSpec defines the desired state of ManagementAddon
type ManagementAddonSpec struct {
	// Addon is the name of the addon to install (must match an AddonDefinition)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Addon string `json:"addon"`

	// Version is the specific version to install. If not specified, uses the
	// default version from the AddonDefinition.
	// +optional
	Version string `json:"version,omitempty"`

	// Values are the Helm values to use for installation.
	// These are merged with any default values from the AddonDefinition.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *runtime.RawExtension `json:"values,omitempty"`

	// Paused indicates whether reconciliation of this addon is paused.
	// When paused, the controller will not make any changes to the addon.
	// +optional
	Paused bool `json:"paused,omitempty"`
}

// ManagementAddonStatus defines the observed state of ManagementAddon
type ManagementAddonStatus struct {
	// Phase represents the current phase of the addon installation
	// +optional
	Phase ManagementAddonPhase `json:"phase,omitempty"`

	// InstalledVersion is the currently installed version of the addon
	// +optional
	InstalledVersion string `json:"installedVersion,omitempty"`

	// HelmRelease contains information about the Helm release
	// +optional
	HelmRelease *HelmReleaseStatus `json:"helmRelease,omitempty"`

	// Message provides additional information about the current state
	// +optional
	Message string `json:"message,omitempty"`

	// LastAttemptedVersion is the version that was last attempted to install
	// +optional
	LastAttemptedVersion string `json:"lastAttemptedVersion,omitempty"`

	// Conditions represent the latest available observations of the addon's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ObservedGeneration is the last observed generation of the ManagementAddon
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=ma;maddon
// +kubebuilder:printcolumn:name="Addon",type="string",JSONPath=".spec.addon",description="Addon name"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="Requested version"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current phase"
// +kubebuilder:printcolumn:name="Installed",type="string",JSONPath=".status.installedVersion",description="Installed version"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ManagementAddon is the Schema for the managementaddons API.
// It represents an addon to be installed on the management cluster.
type ManagementAddon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagementAddonSpec   `json:"spec,omitempty"`
	Status ManagementAddonStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ManagementAddonList contains a list of ManagementAddon
type ManagementAddonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagementAddon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagementAddon{}, &ManagementAddonList{})
}
