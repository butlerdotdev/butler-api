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

// AddonCategory defines the category of an addon for UI grouping.
// +kubebuilder:validation:Enum=cni;loadbalancer;storage;certmanager;ingress;observability;backup;gitops;security;other
type AddonCategory string

const (
	AddonCategoryCNI           AddonCategory = "cni"
	AddonCategoryLoadBalancer  AddonCategory = "loadbalancer"
	AddonCategoryStorage       AddonCategory = "storage"
	AddonCategoryCertManager   AddonCategory = "certmanager"
	AddonCategoryIngress       AddonCategory = "ingress"
	AddonCategoryObservability AddonCategory = "observability"
	AddonCategoryBackup        AddonCategory = "backup"
	AddonCategoryGitOps        AddonCategory = "gitops"
	AddonCategorySecurity      AddonCategory = "security"
	AddonCategoryOther         AddonCategory = "other"
)

// AddonDefinitionSpec defines the desired state of AddonDefinition.
// An AddonDefinition is a cluster-scoped resource that defines an addon
// available for installation in tenant clusters.
type AddonDefinitionSpec struct {
	// DisplayName is the human-readable name shown in the Butler UI.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	DisplayName string `json:"displayName"`

	// Description explains what this addon provides.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	Description string `json:"description"`

	// Category groups addons in the UI for easier discovery.
	// +kubebuilder:validation:Required
	Category AddonCategory `json:"category"`

	// Icon is an emoji or icon identifier for UI display.
	// +kubebuilder:validation:MaxLength=8
	// +optional
	Icon string `json:"icon,omitempty"`

	// Chart specifies the Helm chart to install.
	// +kubebuilder:validation:Required
	Chart AddonChartSpec `json:"chart"`

	// Defaults provides installation defaults.
	// These can be overridden in TenantAddon.
	// +optional
	Defaults *AddonDefaults `json:"defaults,omitempty"`

	// Platform marks this as a core platform addon.
	// Platform addons are installed during cluster bootstrap and cannot
	// be uninstalled via the UI. They appear in a separate section.
	// +kubebuilder:default=false
	// +optional
	Platform bool `json:"platform,omitempty"`

	// DependsOn lists addon names that must be installed first.
	// The TenantAddon controller will wait for these dependencies
	// to be in Installed phase before proceeding.
	// +optional
	DependsOn []string `json:"dependsOn,omitempty"`

	// Maintainer identifies who maintains this addon definition.
	// +optional
	Maintainer *AddonMaintainer `json:"maintainer,omitempty"`

	// Links provides URLs for documentation, source, etc.
	// +optional
	Links *AddonLinks `json:"links,omitempty"`
}

// AddonChartSpec specifies the Helm chart to install.
type AddonChartSpec struct {
	// Repository is the Helm repository URL.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^https?://`
	Repository string `json:"repository"`

	// Name is the chart name within the repository.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// DefaultVersion is the chart version used when TenantAddon
	// doesn't specify a version. Should be a stable, tested version.
	// +kubebuilder:validation:Required
	DefaultVersion string `json:"defaultVersion"`

	// AvailableVersions lists other versions known to work with Butler.
	// Used for version dropdown in UI. If empty, only defaultVersion shown.
	// +optional
	AvailableVersions []string `json:"availableVersions,omitempty"`
}

// AddonDefaults provides default installation settings.
type AddonDefaults struct {
	// Namespace is the target namespace for installation.
	// If not specified, defaults to the addon name.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// ReleaseName is the Helm release name.
	// If not specified, defaults to the addon name.
	// +optional
	ReleaseName string `json:"releaseName,omitempty"`

	// CreateNamespace indicates whether to create the namespace.
	// +kubebuilder:default=true
	// +optional
	CreateNamespace bool `json:"createNamespace,omitempty"`

	// Values are default Helm values applied during installation.
	// These can be overridden in TenantAddon.spec.values.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`

	// Timeout for Helm operations.
	// +kubebuilder:default="10m"
	// +optional
	Timeout string `json:"timeout,omitempty"`
}

// AddonMaintainer identifies the maintainer of an addon definition.
type AddonMaintainer struct {
	// Name of the maintainer.
	// +optional
	Name string `json:"name,omitempty"`

	// Email of the maintainer.
	// +optional
	Email string `json:"email,omitempty"`
}

// AddonLinks provides URLs related to the addon.
type AddonLinks struct {
	// Documentation URL.
	// +optional
	Documentation string `json:"documentation,omitempty"`

	// Source code URL.
	// +optional
	Source string `json:"source,omitempty"`

	// Project homepage URL.
	// +optional
	Homepage string `json:"homepage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=ad;adddef
// +kubebuilder:printcolumn:name="Display Name",type="string",JSONPath=".spec.displayName",description="Human-readable name"
// +kubebuilder:printcolumn:name="Category",type="string",JSONPath=".spec.category",description="Addon category"
// +kubebuilder:printcolumn:name="Chart",type="string",JSONPath=".spec.chart.name",description="Helm chart name"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.chart.defaultVersion",description="Default version"
// +kubebuilder:printcolumn:name="Platform",type="boolean",JSONPath=".spec.platform",description="Is platform addon"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// AddonDefinition defines an addon available for installation in tenant clusters.
// AddonDefinitions are cluster-scoped resources that describe Helm charts
// that can be installed via TenantAddon resources.
//
// Butler ships with built-in AddonDefinitions for common CNCF tools.
// Organizations can add custom AddonDefinitions for internal charts.
type AddonDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AddonDefinitionSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AddonDefinitionList contains a list of AddonDefinition.
type AddonDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AddonDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AddonDefinition{}, &AddonDefinitionList{})
}

// Helper methods

// GetNamespace returns the target namespace, defaulting to addon name.
func (a *AddonDefinition) GetNamespace() string {
	if a.Spec.Defaults != nil && a.Spec.Defaults.Namespace != "" {
		return a.Spec.Defaults.Namespace
	}
	return a.Name
}

// GetReleaseName returns the release name, defaulting to addon name.
func (a *AddonDefinition) GetReleaseName() string {
	if a.Spec.Defaults != nil && a.Spec.Defaults.ReleaseName != "" {
		return a.Spec.Defaults.ReleaseName
	}
	return a.Name
}

// IsBuiltIn returns true if this is a Butler-maintained addon.
func (a *AddonDefinition) IsBuiltIn() bool {
	if a.Labels == nil {
		return false
	}
	return a.Labels["butler.butlerlabs.dev/source"] == "builtin"
}
