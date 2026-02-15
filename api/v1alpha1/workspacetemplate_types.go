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

// WorkspaceTemplateCategory groups templates in the UI picker.
// +kubebuilder:validation:Enum=backend;frontend;data;devops;custom
type WorkspaceTemplateCategory string

const (
	WorkspaceTemplateCategoryBackend  WorkspaceTemplateCategory = "backend"
	WorkspaceTemplateCategoryFrontend WorkspaceTemplateCategory = "frontend"
	WorkspaceTemplateCategoryData     WorkspaceTemplateCategory = "data"
	WorkspaceTemplateCategoryDevOps   WorkspaceTemplateCategory = "devops"
	WorkspaceTemplateCategoryCustom   WorkspaceTemplateCategory = "custom"
)

// WorkspaceTemplateScope determines template visibility.
// +kubebuilder:validation:Enum=cluster;team
type WorkspaceTemplateScope string

const (
	// WorkspaceTemplateScopeCluster makes the template visible to all teams.
	// Created by platform admins in the butler-system namespace.
	WorkspaceTemplateScopeCluster WorkspaceTemplateScope = "cluster"

	// WorkspaceTemplateScopeTeam makes the template visible only to the owning team.
	WorkspaceTemplateScopeTeam WorkspaceTemplateScope = "team"
)

// WorkspaceTemplateSpec defines the desired state of a WorkspaceTemplate.
type WorkspaceTemplateSpec struct {
	// DisplayName shown in the template picker.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	DisplayName string `json:"displayName"`

	// Description of what this template provides.
	// +optional
	Description string `json:"description,omitempty"`

	// Icon is a Material UI icon name or URL for the template card.
	// +optional
	Icon string `json:"icon,omitempty"`

	// Category groups templates in the picker UI.
	// +kubebuilder:default="custom"
	// +optional
	Category WorkspaceTemplateCategory `json:"category,omitempty"`

	// Scope determines visibility.
	// "cluster" templates are visible to all teams (created by platform admins).
	// "team" templates are visible only to the team that owns them.
	// +kubebuilder:default="team"
	// +optional
	Scope WorkspaceTemplateScope `json:"scope,omitempty"`

	// Template is the workspace spec that gets applied when using this template.
	// Owner and ClusterRef are set at creation time by the server.
	// +kubebuilder:validation:Required
	Template WorkspaceTemplateBody `json:"template"`
}

// WorkspaceTemplateBody defines the workspace configuration within a template.
type WorkspaceTemplateBody struct {
	// Image for the workspace container.
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Repository to clone into the workspace.
	// Deprecated: Use Repositories for multi-repo support.
	// +optional
	Repository *WorkspaceRepository `json:"repository,omitempty"`

	// Repositories is a list of Git repositories to clone into the workspace.
	// Each repository is cloned as a sibling directory under /workspace/{repo-name}/.
	// +optional
	Repositories []WorkspaceRepository `json:"repositories,omitempty"`

	// EnvFrom workload to copy environment variables from.
	// +optional
	EnvFrom *WorkspaceEnvSource `json:"envFrom,omitempty"`

	// Dotfiles repo to clone and install.
	// +optional
	Dotfiles *DotfilesSpec `json:"dotfiles,omitempty"`

	// Resources for the workspace pod.
	// +optional
	Resources *WorkspaceResources `json:"resources,omitempty"`

	// StorageSize for the workspace PVC.
	// +optional
	StorageSize *resource.Quantity `json:"storageSize,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=wst
// +kubebuilder:printcolumn:name="Display Name",type="string",JSONPath=".spec.displayName",description="Template display name"
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.template.image",description="Default workspace image"
// +kubebuilder:printcolumn:name="Category",type="string",JSONPath=".spec.category",description="Template category"
// +kubebuilder:printcolumn:name="Scope",type="string",JSONPath=".spec.scope",description="Visibility scope"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// WorkspaceTemplate is a pre-configured workspace specification for one-click creation.
// Templates are data-only resources — no controller reconciliation is needed.
// Cluster-scoped templates live in butler-system and are visible to all teams.
// Team-scoped templates live in the team namespace and are visible only to that team.
type WorkspaceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec WorkspaceTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// WorkspaceTemplateList contains a list of WorkspaceTemplates.
type WorkspaceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkspaceTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkspaceTemplate{}, &WorkspaceTemplateList{})
}
