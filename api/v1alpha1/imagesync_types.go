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

// ImageSyncSpec defines the desired state of ImageSync.
type ImageSyncSpec struct {
	// FactoryRef identifies the image in Butler Image Factory.
	// +kubebuilder:validation:Required
	FactoryRef ImageFactoryRef `json:"factoryRef"`

	// ProviderConfigRef references the target ProviderConfig.
	// +kubebuilder:validation:Required
	ProviderConfigRef ProviderReference `json:"providerConfigRef"`

	// Format is the image format to sync (e.g., "qcow2", "raw", "iso").
	// +kubebuilder:default="qcow2"
	// +optional
	Format string `json:"format,omitempty"`

	// TransferMode controls how the image artifact is transferred to the provider.
	// "direct" means the provider downloads directly from the factory URL.
	// "proxy" means the provider controller downloads first, then uploads to the provider.
	// +kubebuilder:validation:Enum=direct;proxy
	// +kubebuilder:default="direct"
	// +optional
	TransferMode string `json:"transferMode,omitempty"`

	// DisplayName is the human-readable name for the image on the provider.
	// If not set, auto-generated from OS type + version + schematic ID prefix.
	// +optional
	DisplayName string `json:"displayName,omitempty"`
}

// ImageFactoryRef identifies an image in the Butler Image Factory.
type ImageFactoryRef struct {
	// SchematicID is the content-addressable schematic identifier (SHA-256 hex).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=8
	SchematicID string `json:"schematicID"`

	// Version is the OS version (e.g., "v1.12.4", "9.5").
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Arch is the CPU architecture.
	// +kubebuilder:validation:Enum=amd64;arm64
	// +kubebuilder:default="amd64"
	// +optional
	Arch string `json:"arch,omitempty"`

	// Platform is the target platform for the image artifact.
	// Maps to the Talos Image Factory platform identifier (e.g., "nocloud", "metal", "vmware").
	// Defaults to "nocloud" which works for KubeVirt/cloud-init environments (Harvester, Nutanix).
	// +kubebuilder:default="nocloud"
	// +optional
	Platform string `json:"platform,omitempty"`
}

// ImageSyncPhase represents the lifecycle phase of an ImageSync.
// +kubebuilder:validation:Enum=Pending;Building;Downloading;Uploading;Ready;Failed
type ImageSyncPhase string

const (
	// ImageSyncPhasePending indicates the sync request has been received but not yet processed.
	ImageSyncPhasePending ImageSyncPhase = "Pending"

	// ImageSyncPhaseBuilding indicates the factory is building the image.
	ImageSyncPhaseBuilding ImageSyncPhase = "Building"

	// ImageSyncPhaseDownloading indicates the image is being downloaded from the factory.
	ImageSyncPhaseDownloading ImageSyncPhase = "Downloading"

	// ImageSyncPhaseUploading indicates the image is being uploaded to the provider.
	ImageSyncPhaseUploading ImageSyncPhase = "Uploading"

	// ImageSyncPhaseReady indicates the image is available on the provider.
	ImageSyncPhaseReady ImageSyncPhase = "Ready"

	// ImageSyncPhaseFailed indicates the image sync failed.
	ImageSyncPhaseFailed ImageSyncPhase = "Failed"
)

// TransferMode constants.
const (
	// TransferModeDirect means the provider downloads directly from the factory URL.
	TransferModeDirect = "direct"

	// TransferModeProxy means the provider controller downloads first, then uploads to the provider.
	TransferModeProxy = "proxy"
)

// ImageSyncStatus defines the observed state of ImageSync.
type ImageSyncStatus struct {
	// Phase represents the current lifecycle phase.
	// +optional
	Phase ImageSyncPhase `json:"phase,omitempty"`

	// Conditions represent the latest available observations.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ProviderImageRef is the provider-specific image reference once synced.
	// Harvester: "namespace/name", Nutanix: UUID, Proxmox: template ID.
	// +optional
	ProviderImageRef string `json:"providerImageRef,omitempty"`

	// ArtifactURL is the factory download URL used.
	// +optional
	ArtifactURL string `json:"artifactURL,omitempty"`

	// ArtifactSHA256 is the verified checksum.
	// +optional
	ArtifactSHA256 string `json:"artifactSHA256,omitempty"`

	// ProviderTaskID tracks a provider-side async task (e.g., Nutanix Prism Central task UUID).
	// Set during Downloading/Uploading phases, cleared on completion or failure.
	// +optional
	ProviderTaskID string `json:"providerTaskID,omitempty"`

	// FailureReason provides a machine-readable failure reason.
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage provides a human-readable failure message.
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// ObservedGeneration is the generation most recently observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastUpdated is the timestamp of the last status update.
	// +optional
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=is
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Sync phase"
// +kubebuilder:printcolumn:name="Schematic",type="string",JSONPath=".spec.factoryRef.schematicID",description="Schematic ID",priority=1
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.factoryRef.version",description="OS version"
// +kubebuilder:printcolumn:name="Provider",type="string",JSONPath=".spec.providerConfigRef.name",description="Target provider"
// +kubebuilder:printcolumn:name="Image Ref",type="string",JSONPath=".status.providerImageRef",description="Provider image reference"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ImageSync is the Schema for the imagesyncs API.
// It represents a request to sync an image from the Butler Image Factory
// to an infrastructure provider.
type ImageSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSyncSpec   `json:"spec,omitempty"`
	Status ImageSyncStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageSyncList contains a list of ImageSync.
type ImageSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageSync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImageSync{}, &ImageSyncList{})
}

// Helper methods for ImageSync

// IsReady returns true if the image sync is complete and the provider image ref is set.
func (is *ImageSync) IsReady() bool {
	return is.Status.Phase == ImageSyncPhaseReady && is.Status.ProviderImageRef != ""
}

// IsFailed returns true if the image sync is in a failed state.
func (is *ImageSync) IsFailed() bool {
	return is.Status.Phase == ImageSyncPhaseFailed
}

// SetPhase updates the phase and last updated timestamp.
func (is *ImageSync) SetPhase(phase ImageSyncPhase) {
	is.Status.Phase = phase
	now := metav1.Now()
	is.Status.LastUpdated = &now
}

// SetFailure sets the failure reason and message.
func (is *ImageSync) SetFailure(reason, message string) {
	is.Status.FailureReason = reason
	is.Status.FailureMessage = message
	is.SetPhase(ImageSyncPhaseFailed)
}
