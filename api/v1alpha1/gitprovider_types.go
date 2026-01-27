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

// GitProviderType defines supported Git providers.
// +kubebuilder:validation:Enum=github;gitlab;bitbucket
type GitProviderType string

const (
	// GitProviderTypeGitHub is GitHub.com or GitHub Enterprise.
	GitProviderTypeGitHub GitProviderType = "github"

	// GitProviderTypeGitLab is GitLab.com or self-hosted GitLab.
	GitProviderTypeGitLab GitProviderType = "gitlab"

	// GitProviderTypeBitbucket is Bitbucket Cloud or Server.
	GitProviderTypeBitbucket GitProviderType = "bitbucket"
)

// GitProviderConfig configures a Git provider for GitOps operations.
// This is stored in ButlerConfig and used as the default for all GitOps operations.
type GitProviderConfig struct {
	// Type is the Git provider type.
	// +kubebuilder:validation:Required
	Type GitProviderType `json:"type"`

	// URL is the Git provider API URL.
	// For GitHub: https://api.github.com (or https://github.example.com/api/v3 for enterprise)
	// For GitLab: https://gitlab.com (or self-hosted URL)
	// +kubebuilder:default="https://api.github.com"
	// +optional
	URL string `json:"url,omitempty"`

	// Organization is the default organization/group for repositories.
	// When set, repository listings will be scoped to this org.
	// +optional
	Organization string `json:"organization,omitempty"`

	// SecretRef references the Secret containing credentials.
	// Required keys depend on provider type:
	// - GitHub: "token" (Personal Access Token with repo scope)
	// - GitLab: "token" (Personal Access Token with api scope)
	// - Bitbucket: "username" and "app-password"
	// +kubebuilder:validation:Required
	SecretRef LocalObjectReference `json:"secretRef"`
}

// GitProviderStatus shows the status of the Git provider configuration.
type GitProviderStatus struct {
	// Connected indicates whether the provider credentials are valid.
	// +optional
	Connected bool `json:"connected,omitempty"`

	// Username is the authenticated username (from token validation).
	// +optional
	Username string `json:"username,omitempty"`

	// LastValidated is when the credentials were last validated.
	// +optional
	LastValidated *metav1.Time `json:"lastValidated,omitempty"`

	// Message provides additional status information.
	// +optional
	Message string `json:"message,omitempty"`
}

// GitOpsExportFormat defines the output format for GitOps exports.
// +kubebuilder:validation:Enum=flux;argocd;raw;kustomize
type GitOpsExportFormat string

const (
	// GitOpsExportFormatFlux generates Flux HelmRelease and Kustomization resources.
	GitOpsExportFormatFlux GitOpsExportFormat = "flux"

	// GitOpsExportFormatArgoCD generates ArgoCD Application resources.
	GitOpsExportFormatArgoCD GitOpsExportFormat = "argocd"

	// GitOpsExportFormatRaw generates plain Kubernetes manifests.
	GitOpsExportFormatRaw GitOpsExportFormat = "raw"

	// GitOpsExportFormatKustomize generates Kustomization structure.
	GitOpsExportFormatKustomize GitOpsExportFormat = "kustomize"
)

// GitOpsDirectoryLayout defines the standard directory structure for GitOps repositories.
type GitOpsDirectoryLayout struct {
	// ClustersPath is the path for cluster-specific configurations.
	// +kubebuilder:default="clusters"
	// +optional
	ClustersPath string `json:"clustersPath,omitempty"`

	// InfrastructurePath is the path for infrastructure components.
	// +kubebuilder:default="infrastructure"
	// +optional
	InfrastructurePath string `json:"infrastructurePath,omitempty"`

	// AppsPath is the path for application workloads.
	// +kubebuilder:default="apps"
	// +optional
	AppsPath string `json:"appsPath,omitempty"`

	// PlatformPath is the path for platform components (observability, security, etc).
	// +kubebuilder:default="platform"
	// +optional
	PlatformPath string `json:"platformPath,omitempty"`
}

// DefaultGitOpsDirectoryLayout returns the default directory layout.
func DefaultGitOpsDirectoryLayout() GitOpsDirectoryLayout {
	return GitOpsDirectoryLayout{
		ClustersPath:       "clusters",
		InfrastructurePath: "infrastructure",
		AppsPath:           "apps",
		PlatformPath:       "platform",
	}
}
