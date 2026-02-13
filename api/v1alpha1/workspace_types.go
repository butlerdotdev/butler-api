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

// WorkspacePhase represents the current lifecycle phase of a Workspace.
// +kubebuilder:validation:Enum=Pending;Creating;Running;Starting;Stopped;Failed
type WorkspacePhase string

const (
	// WorkspacePhasePending indicates the workspace is awaiting reconciliation.
	WorkspacePhasePending WorkspacePhase = "Pending"

	// WorkspacePhaseCreating indicates tenant resources (PVC, Pod) are being provisioned.
	WorkspacePhaseCreating WorkspacePhase = "Creating"

	// WorkspacePhaseRunning indicates the workspace pod is running and SSH server is ready.
	WorkspacePhaseRunning WorkspacePhase = "Running"

	// WorkspacePhaseStarting indicates a stopped workspace is resuming.
	WorkspacePhaseStarting WorkspacePhase = "Starting"

	// WorkspacePhaseStopped indicates the pod was deleted after inactivity. PVC persists.
	WorkspacePhaseStopped WorkspacePhase = "Stopped"

	// WorkspacePhaseFailed indicates a terminal error occurred.
	WorkspacePhaseFailed WorkspacePhase = "Failed"
)

// Workspace condition types.
const (
	// WorkspaceConditionPVCReady indicates the PVC is created and bound.
	WorkspaceConditionPVCReady = "PVCReady"

	// WorkspaceConditionPodReady indicates the pod is created and running.
	WorkspaceConditionPodReady = "PodReady"

	// WorkspaceConditionRepositoryCloned indicates the Git repository was cloned.
	WorkspaceConditionRepositoryCloned = "RepositoryCloned"

	// WorkspaceConditionDotfilesInstalled indicates dotfiles were installed.
	WorkspaceConditionDotfilesInstalled = "DotfilesInstalled"

	// WorkspaceConditionSSHReady indicates the SSH server is accessible.
	WorkspaceConditionSSHReady = "SSHReady"

	// WorkspaceConditionReady indicates the workspace is fully operational.
	WorkspaceConditionReady = "Ready"
)

// WorkspaceSpec defines the desired state of a Workspace.
type WorkspaceSpec struct {
	// ClusterRef references the TenantCluster this workspace runs in.
	// The workspace pod is created in the tenant cluster's "workspaces" namespace.
	// +kubebuilder:validation:Required
	ClusterRef LocalObjectReference `json:"clusterRef"`

	// Owner is the email of the Butler user who owns this workspace.
	// Set by the server from the authenticated user's JWT.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Owner string `json:"owner"`

	// Image is the workspace container image.
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Repository to clone into the workspace on creation.
	// +optional
	Repository *WorkspaceRepository `json:"repository,omitempty"`

	// EnvFrom copies environment variables from an existing workload
	// running in the tenant cluster.
	// +optional
	EnvFrom *WorkspaceEnvSource `json:"envFrom,omitempty"`

	// Resources for the workspace pod.
	// +optional
	Resources *WorkspaceResources `json:"resources,omitempty"`

	// Dotfiles repo to clone and run install script on first creation.
	// +optional
	Dotfiles *DotfilesSpec `json:"dotfiles,omitempty"`

	// IdleTimeout after which the SSH service is torn down.
	// Pod and PVC persist. Only the network path is removed.
	// Measured from connect time (no activity tracking).
	// +kubebuilder:default="4h"
	// +optional
	IdleTimeout *metav1.Duration `json:"idleTimeout,omitempty"`

	// AutoStopAfter stops the workspace pod after this duration of no SSH connections.
	// Measured from last disconnect time. PVC persists.
	// Set to 0 to disable auto-stop.
	// +kubebuilder:default="8h"
	// +optional
	AutoStopAfter *metav1.Duration `json:"autoStopAfter,omitempty"`

	// StorageSize for the workspace PVC.
	// +kubebuilder:default="50Gi"
	// +optional
	StorageSize *resource.Quantity `json:"storageSize,omitempty"`

	// SSHPublicKeys for authorized access. If empty, keys are resolved
	// from the owner's User profile (spec.sshKeys).
	// +optional
	SSHPublicKeys []string `json:"sshPublicKeys,omitempty"`
}

// WorkspaceRepository configures a Git repository to clone into the workspace.
type WorkspaceRepository struct {
	// URL is the Git repository URL.
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Branch to checkout.
	// +kubebuilder:default="main"
	// +optional
	Branch string `json:"branch,omitempty"`

	// SecretRef for Git credentials (SSH key or token).
	// References a Secret in the team namespace on the management cluster.
	// The controller copies it to the tenant cluster's workspaces namespace.
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty"`
}

// WorkspaceEnvSource configures environment variable copying from an existing workload.
type WorkspaceEnvSource struct {
	// Kind of the source workload.
	// +kubebuilder:validation:Enum=Deployment;StatefulSet
	// +kubebuilder:default="Deployment"
	Kind string `json:"kind,omitempty"`

	// Name of the workload to copy env from.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Namespace of the workload in the tenant cluster. Defaults to "default".
	// +kubebuilder:default="default"
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Container name within the workload. Defaults to first container.
	// +optional
	Container string `json:"container,omitempty"`
}

// WorkspaceResources configures compute resources for the workspace pod.
type WorkspaceResources struct {
	// CPU request and limit for the workspace.
	// +kubebuilder:default="2"
	// +optional
	CPU string `json:"cpu,omitempty"`

	// Memory request and limit for the workspace.
	// +kubebuilder:default="4Gi"
	// +optional
	Memory string `json:"memory,omitempty"`
}

// DotfilesSpec configures a dotfiles repository to clone and install on workspace creation.
type DotfilesSpec struct {
	// URL is the Git repository URL for dotfiles.
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// InstallCommand to run after cloning. If not specified, the controller
	// looks for install.sh, install, bootstrap.sh, bootstrap, setup.sh, setup, or Makefile.
	// +optional
	InstallCommand string `json:"installCommand,omitempty"`
}

// WorkspaceStatus defines the observed state of Workspace.
type WorkspaceStatus struct {
	// Conditions represent the latest available observations of the workspace's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase of the workspace lifecycle.
	// +optional
	Phase WorkspacePhase `json:"phase,omitempty"`

	// PodName is the name of the workspace pod in the tenant cluster.
	// +optional
	PodName string `json:"podName,omitempty"`

	// PVCName is the name of the workspace PVC in the tenant cluster.
	// +optional
	PVCName string `json:"pvcName,omitempty"`

	// ServiceName is the SSH service name when connected.
	// +optional
	ServiceName string `json:"serviceName,omitempty"`

	// SSHEndpoint is the IP:port for SSH access when connected.
	// +optional
	SSHEndpoint string `json:"sshEndpoint,omitempty"`

	// Connected indicates whether the SSH service is currently active.
	// +optional
	Connected bool `json:"connected,omitempty"`

	// LastActivityTime tracks the last SSH connect time.
	// Used for idle timeout calculation.
	// +optional
	LastActivityTime *metav1.Time `json:"lastActivityTime,omitempty"`

	// LastDisconnectTime tracks when the SSH service was last removed.
	// Used for auto-stop calculation.
	// +optional
	LastDisconnectTime *metav1.Time `json:"lastDisconnectTime,omitempty"`

	// ObservedGeneration is the last observed generation of the workspace spec.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ws
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current lifecycle phase"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.clusterRef.name",description="Target tenant cluster"
// +kubebuilder:printcolumn:name="Owner",type="string",JSONPath=".spec.owner",description="Workspace owner email"
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image",description="Container image"
// +kubebuilder:printcolumn:name="Connected",type="boolean",JSONPath=".status.connected",description="SSH service active"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Workspace represents a cloud development environment running inside a tenant cluster.
// The Workspace CRD lives on the management cluster (in the team namespace),
// while the actual pod, PVC, and SSH service are created in the tenant cluster's
// "workspaces" namespace by the workspace controller.
type Workspace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkspaceSpec   `json:"spec,omitempty"`
	Status WorkspaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WorkspaceList contains a list of Workspaces.
type WorkspaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workspace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Workspace{}, &WorkspaceList{})
}

// IsRunning returns true if the workspace pod is running.
func (w *Workspace) IsRunning() bool {
	return w.Status.Phase == WorkspacePhaseRunning
}

// IsStopped returns true if the workspace pod has been stopped.
func (w *Workspace) IsStopped() bool {
	return w.Status.Phase == WorkspacePhaseStopped
}

// IsConnected returns true if the SSH service is active.
func (w *Workspace) IsConnected() bool {
	return w.Status.Connected
}
