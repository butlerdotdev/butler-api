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

// TeamRole defines the role a user or group has within a Team.
// +kubebuilder:validation:Enum=admin;operator;viewer
type TeamRole string

const (
	// TeamRoleAdmin has full access to manage the team and all its resources.
	// Can: create/delete clusters, manage team members, change settings
	TeamRoleAdmin TeamRole = "admin"

	// TeamRoleOperator can create and manage clusters but cannot manage team settings.
	// Can: create/delete clusters, scale, install addons
	// Cannot: manage team members, change team settings
	TeamRoleOperator TeamRole = "operator"

	// TeamRoleViewer has read-only access to team resources.
	// Can: view clusters, view kubeconfigs, view logs
	// Cannot: create/modify/delete anything
	TeamRoleViewer TeamRole = "viewer"
)

// TeamSpec defines the desired state of Team.
type TeamSpec struct {
	// DisplayName is a human-readable name for the Team.
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// Description provides additional context about the Team.
	// +optional
	Description string `json:"description,omitempty"`

	// Access defines who can access this Team's resources.
	// +optional
	Access TeamAccess `json:"access,omitempty"`

	// ResourceLimits defines the resource quotas and restrictions for this Team.
	// If not specified, defaults from ButlerConfig are used.
	// If ButlerConfig has no defaults, no limits are enforced.
	// +optional
	ResourceLimits *TeamResourceLimits `json:"resourceLimits,omitempty"`

	// ProviderConfigRef references a Team-specific ProviderConfig.
	// If not specified, the platform default is used.
	// +optional
	ProviderConfigRef *LocalObjectReference `json:"providerConfigRef,omitempty"`

	// ClusterDefaults defines default values for new clusters in this team.
	// +optional
	ClusterDefaults *ClusterDefaults `json:"clusterDefaults,omitempty"`
}

// ClusterDefaults defines default values for new TenantClusters.
type ClusterDefaults struct {
	// KubernetesVersion is the default K8s version for new clusters.
	// +optional
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// WorkerCount is the default number of worker nodes.
	// +optional
	// +kubebuilder:validation:Minimum=0
	WorkerCount *int32 `json:"workerCount,omitempty"`

	// WorkerCPU is the default CPU cores per worker.
	// +optional
	// +kubebuilder:validation:Minimum=1
	WorkerCPU *int32 `json:"workerCPU,omitempty"`

	// WorkerMemoryGi is the default memory per worker in Gi.
	// +optional
	// +kubebuilder:validation:Minimum=1
	WorkerMemoryGi *int32 `json:"workerMemoryGi,omitempty"`

	// WorkerDiskGi is the default disk size per worker in Gi.
	// +optional
	// +kubebuilder:validation:Minimum=10
	WorkerDiskGi *int32 `json:"workerDiskGi,omitempty"`

	// DefaultAddons are addons automatically installed on new clusters.
	// +optional
	DefaultAddons []string `json:"defaultAddons,omitempty"`
}

// TeamAccess defines users and groups that have access to the Team.
type TeamAccess struct {
	// Users is a list of users with access to this Team.
	// Users are identified by their email address.
	// +optional
	Users []TeamUser `json:"users,omitempty"`

	// Groups is a list of groups with access to this Team.
	// Groups are matched against OIDC groups or AD groups.
	// +optional
	Groups []TeamGroup `json:"groups,omitempty"`
}

// TeamUser represents a user with access to a Team.
type TeamUser struct {
	// Name is the user identifier (email address).
	// For internal users, this is the email from User.spec.email.
	// For SSO users, this is the email from the OIDC token.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Role is the user's role within the Team.
	// +kubebuilder:default="viewer"
	// +optional
	Role TeamRole `json:"role,omitempty"`
}

// TeamGroup represents a group with access to a Team.
type TeamGroup struct {
	// Name is the group identifier (OIDC group, AD group DN, etc.).
	// This can be the full DN for AD groups or simple names for OIDC.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Role is the group's role within the Team.
	// All members of the group inherit this role.
	// +kubebuilder:default="viewer"
	// +optional
	Role TeamRole `json:"role,omitempty"`

	// IdentityProvider is the name of the IdentityProvider CRD this group comes from.
	// If specified, only users authenticating through this IdP will be matched.
	// If not specified, the group name will be matched against groups from any IdP.
	// +optional
	IdentityProvider string `json:"identityProvider,omitempty"`
}

// TeamPhase represents the current phase of a Team.
// +kubebuilder:validation:Enum=Pending;Ready;Terminating;Failed
type TeamPhase string

const (
	// TeamPhasePending indicates the Team is being set up.
	TeamPhasePending TeamPhase = "Pending"

	// TeamPhaseReady indicates the Team is ready for use.
	TeamPhaseReady TeamPhase = "Ready"

	// TeamPhaseTerminating indicates the Team is being deleted.
	TeamPhaseTerminating TeamPhase = "Terminating"

	// TeamPhaseFailed indicates the Team setup failed.
	TeamPhaseFailed TeamPhase = "Failed"
)

// TeamStatus defines the observed state of Team.
type TeamStatus struct {
	// Conditions represent the latest available observations of the Team's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current phase of the Team.
	// +optional
	Phase TeamPhase `json:"phase,omitempty"`

	// Namespace is the namespace created for this Team.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// ClusterCount is the number of TenantClusters in this Team.
	// +optional
	ClusterCount int32 `json:"clusterCount,omitempty"`

	// MemberCount is the total number of users with access to this Team.
	// +optional
	MemberCount int32 `json:"memberCount,omitempty"`

	// ResourceUsage shows the current resource usage for this Team.
	// +optional
	ResourceUsage *TeamResourceUsage `json:"resourceUsage,omitempty"`

	// QuotaStatus indicates whether the team is within quota.
	// +optional
	// +kubebuilder:validation:Enum=OK;Warning;Exceeded
	QuotaStatus string `json:"quotaStatus,omitempty"`

	// QuotaMessage provides details about quota status.
	// +optional
	QuotaMessage string `json:"quotaMessage,omitempty"`
}

// Team condition types.
const (
	// TeamConditionNamespaceReady indicates the Team namespace exists.
	TeamConditionNamespaceReady = "NamespaceReady"

	// TeamConditionRBACReady indicates RBAC is configured.
	TeamConditionRBACReady = "RBACReady"

	// TeamConditionReady indicates the Team is fully ready.
	TeamConditionReady = "Ready"

	// TeamConditionQuotaExceeded indicates the Team has exceeded quota.
	TeamConditionQuotaExceeded = "QuotaExceeded"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=tm
// +kubebuilder:printcolumn:name="Display Name",type="string",JSONPath=".spec.displayName",description="Human-readable name"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current phase"
// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".status.namespace",description="Team namespace"
// +kubebuilder:printcolumn:name="Clusters",type="integer",JSONPath=".status.clusterCount",description="Number of clusters"
// +kubebuilder:printcolumn:name="Quota",type="string",JSONPath=".status.quotaStatus",description="Quota status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Team is the Schema for the teams API.
// A Team represents a group of users who share access to TenantClusters.
// Each Team gets its own namespace where TenantClusters are created.
type Team struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamSpec   `json:"spec,omitempty"`
	Status TeamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TeamList contains a list of Team.
type TeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Team `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Team{}, &TeamList{})
}
