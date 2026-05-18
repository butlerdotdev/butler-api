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

// OptionType names a curated option list. Closed set in v1; matches the
// four list endpoints butler-server exposes today (image, network,
// cluster, storageContainer). See ADR-018 Decision section 4.
// +kubebuilder:validation:Enum=image;network;cluster;storageContainer
type OptionType string

const (
	// OptionTypeImage curates the per-provider VM image list.
	OptionTypeImage OptionType = "image"

	// OptionTypeNetwork curates the per-provider network list (Nutanix
	// subnets, Harvester NetworkAttachmentDefinitions).
	OptionTypeNetwork OptionType = "network"

	// OptionTypeCluster curates the Nutanix Prism Element cluster list.
	// No effect on providers without a cluster-listing endpoint.
	OptionTypeCluster OptionType = "cluster"

	// OptionTypeStorageContainer curates the Nutanix storage container
	// list. No effect on providers without a storage-container-listing
	// endpoint.
	OptionTypeStorageContainer OptionType = "storageContainer"
)

// OptionMode controls how a rule constrains or decorates an option list.
// See ADR-018 Decision section 5 for full semantics.
// +kubebuilder:validation:Enum=pin;allowList;default;recommended
type OptionMode string

const (
	// OptionModePin restricts the dropdown to the listed values. The
	// admission webhook rejects any TenantCluster spec value outside the
	// pinned set.
	OptionModePin OptionMode = "pin"

	// OptionModeAllowList restricts the dropdown to the listed values.
	// The admission webhook rejects any TenantCluster spec value outside
	// the allow-listed set. Default (when set) is pre-selected.
	OptionModeAllowList OptionMode = "allowList"

	// OptionModeDefault pre-selects Default in the dropdown. The full
	// provider list remains visible. Admission does not enforce.
	OptionModeDefault OptionMode = "default"

	// OptionModeRecommended badges the listed values and sorts them
	// first. The full provider list remains visible. Admission does not
	// enforce.
	OptionModeRecommended OptionMode = "recommended"
)

// OptionRule defines the curation behavior for one option type. See
// ADR-018 Decision section 5 for mode semantics. Cross-field validation
// (pin and allowList require non-empty Values; default requires
// non-empty Default; recommended requires non-empty Values) is enforced
// by the policy admission webhook in butler-controller, not by CRD
// schema; the webhook produces better operator-facing error messages.
type OptionRule struct {
	// Mode selects the curation behavior.
	// +kubebuilder:validation:Required
	Mode OptionMode `json:"mode"`

	// Values lists the affected provider entry IDs. Semantics depend on
	// Mode (see ADR-018 Decision section 5). Validated by the policy
	// admission webhook against the mode-specific requirements.
	// +optional
	Values []string `json:"values,omitempty"`

	// Default is the pre-selected entry ID. Required when Mode is
	// default; optional for allowList (pre-selects within the allowed
	// set); ignored for recommended; auto-resolves to the single Values
	// entry for pin when Values has length 1.
	// +optional
	Default string `json:"default,omitempty"`

	// RecommendedReason is the short string shown alongside the
	// recommended badge in the console. Used only when Mode is
	// recommended.
	// +optional
	// +kubebuilder:validation:MaxLength=120
	RecommendedReason string `json:"recommendedReason,omitempty"`
}

// ClusterWideScope marks a policy as applying to every team. Use for
// org-wide rules such as deprecation lists.
type ClusterWideScope struct{}

// TeamScope marks a policy as applying to one team across all of its
// environments.
type TeamScope struct {
	// TeamRef names the Team this policy applies to.
	// +kubebuilder:validation:Required
	TeamRef LocalObjectReference `json:"teamRef"`
}

// TeamEnvironmentScope marks a policy as applying to one team's named
// environment (see ADR-009 Team.spec.environments).
type TeamEnvironmentScope struct {
	// TeamRef names the Team this policy applies to.
	// +kubebuilder:validation:Required
	TeamRef LocalObjectReference `json:"teamRef"`

	// EnvironmentName must match a Team.spec.environments[].name on the
	// referenced Team.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	EnvironmentName string `json:"environmentName"`
}

// PolicyScope discriminator. Exactly one of clusterWide, team, or
// teamAndEnvironment must be set. Enforced via CEL XValidation here at
// schema time; the policy admission webhook in butler-controller layers
// additional checks (referenced team and environment existence,
// intra-tier conflict detection) per ADR-018 Decision section 7.
// +kubebuilder:validation:XValidation:rule="(has(self.clusterWide)?1:0) + (has(self.team)?1:0) + (has(self.teamAndEnvironment)?1:0) == 1",message="exactly one of clusterWide, team, or teamAndEnvironment must be set"
type PolicyScope struct {
	// ClusterWide applies the policy to every team.
	// +optional
	ClusterWide *ClusterWideScope `json:"clusterWide,omitempty"`

	// Team applies the policy to one team across all of its environments.
	// +optional
	Team *TeamScope `json:"team,omitempty"`

	// TeamAndEnvironment applies the policy to one team's named
	// environment.
	// +optional
	TeamAndEnvironment *TeamEnvironmentScope `json:"teamAndEnvironment,omitempty"`
}

// ClusterCreationPolicySpec declares admin-curated constraints over the
// create-cluster option lists. Scalar soft defaults remain on
// Team.spec.clusterDefaults (see ADR-009); this CRD owns option-list
// curation only. See ADR-018 Decision section 1.
type ClusterCreationPolicySpec struct {
	// Scope selects which TenantCluster create requests this policy
	// applies to.
	// +kubebuilder:validation:Required
	Scope PolicyScope `json:"scope"`

	// TargetProviders names the ProviderType values this policy filters.
	// Empty list means all providers.
	// +optional
	TargetProviders []ProviderType `json:"targetProviders,omitempty"`

	// Options maps option-type name to its curation rule. Each
	// ClusterCreationPolicy contains zero or more OptionRule entries
	// within Options, one rule per option type. A policy is the unit of
	// authorship and RBAC; a rule is the unit of resolution.
	// +optional
	Options map[OptionType]OptionRule `json:"options,omitempty"`
}

// ClusterCreationPolicyStatus captures observed conditions on the
// policy. v1 admission webhooks do not populate any fields here; the
// fields are reserved for the deferred status reconciler that surfaces
// stale references and policy health (per ADR-018 Deferred).
type ClusterCreationPolicyStatus struct {
	// Conditions records observed conditions on the policy. Reserved
	// for the deferred status reconciler; not populated by v1 admission
	// webhooks.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// StaleReferences names policy entry IDs that no longer exist in
	// the current provider option list. Reserved for the deferred
	// status reconciler.
	// +optional
	StaleReferences []string `json:"staleReferences,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=ccp
// +kubebuilder:printcolumn:name="Team",type="string",JSONPath=".spec.scope.team.teamRef.name",description="Team this policy targets (team scope)"
// +kubebuilder:printcolumn:name="Env",type="string",JSONPath=".spec.scope.teamAndEnvironment.environmentName",description="Environment this policy targets (teamAndEnvironment scope)"
// +kubebuilder:printcolumn:name="Providers",type="string",JSONPath=".spec.targetProviders",description="Providers this policy filters (empty means all)"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterCreationPolicy is the Schema for admin-curated defaults, pins,
// allow-lists, and recommendations applied to the create-cluster
// option lists (subnets, images, storage containers, clusters).
// Enforced at both butler-server option-list time and at TenantCluster
// admission. See ADR-018 for the full design.
type ClusterCreationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterCreationPolicySpec   `json:"spec,omitempty"`
	Status ClusterCreationPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterCreationPolicyList contains a list of ClusterCreationPolicy.
type ClusterCreationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterCreationPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterCreationPolicy{}, &ClusterCreationPolicyList{})
}
