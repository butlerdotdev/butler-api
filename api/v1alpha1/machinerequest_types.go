/*
Copyright 2025 The Butler Authors.

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

// MachineRole defines the role of a machine in a cluster.
// +kubebuilder:validation:Enum=control-plane;worker
type MachineRole string

const (
	// MachineRoleControlPlane is a control plane node.
	MachineRoleControlPlane MachineRole = "control-plane"

	// MachineRoleWorker is a worker node.
	MachineRoleWorker MachineRole = "worker"
)

// MachinePhase represents the lifecycle phase of a MachineRequest.
// +kubebuilder:validation:Enum=Pending;Creating;Running;Failed;Deleting;Deleted;Unknown
type MachinePhase string

const (
	// MachinePhasePending indicates the request has been received but not yet processed.
	MachinePhasePending MachinePhase = "Pending"

	// MachinePhaseCreating indicates the machine is being created.
	MachinePhaseCreating MachinePhase = "Creating"

	// MachinePhaseRunning indicates the machine is running and has an IP address.
	MachinePhaseRunning MachinePhase = "Running"

	// MachinePhaseFailed indicates the machine creation failed.
	MachinePhaseFailed MachinePhase = "Failed"

	// MachinePhaseDeleting indicates the machine is being deleted.
	MachinePhaseDeleting MachinePhase = "Deleting"

	// MachinePhaseDeleted indicates the machine has been deleted.
	MachinePhaseDeleted MachinePhase = "Deleted"

	// MachinePhaseUnknown indicates the machine state cannot be determined.
	MachinePhaseUnknown MachinePhase = "Unknown"
)

// MachineRequestSpec defines the desired state of MachineRequest.
// This is the interface contract between the bootstrap controller and
// infrastructure provider controllers.
type MachineRequestSpec struct {
	// ProviderRef references the ProviderConfig to use for this machine.
	// +kubebuilder:validation:Required
	ProviderRef ProviderReference `json:"providerRef"`

	// MachineName is the desired name for the virtual machine.
	// Must be unique within the provider's namespace/project.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	MachineName string `json:"machineName"`

	// Role indicates the intended role of this machine in the cluster.
	// +kubebuilder:validation:Required
	Role MachineRole `json:"role"`

	// CPU is the number of virtual CPU cores.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=128
	CPU int32 `json:"cpu"`

	// MemoryMB is the amount of memory in megabytes.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1024
	MemoryMB int32 `json:"memoryMB"`

	// DiskGB is the root disk size in gigabytes.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=10
	DiskGB int32 `json:"diskGB"`

	// ExtraDisks defines additional disks to attach to the machine.
	// +optional
	ExtraDisks []DiskSpec `json:"extraDisks,omitempty"`

	// Image overrides the default OS image from ProviderConfig.
	// Format is provider-specific:
	// - harvester: "namespace/image-name"
	// - nutanix: UUID
	// - proxmox: template ID or image name
	// +optional
	Image string `json:"image,omitempty"`

	// UserData is cloud-init user data to configure the machine.
	// This typically contains the Talos machine configuration.
	// +optional
	UserData string `json:"userData,omitempty"`

	// NetworkData is cloud-init network configuration.
	// +optional
	NetworkData string `json:"networkData,omitempty"`

	// Labels are key-value pairs to apply to the VM in the provider.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// DiskSpec defines an additional disk to attach to a machine.
type DiskSpec struct {
	// SizeGB is the disk size in gigabytes.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	SizeGB int32 `json:"sizeGB"`

	// StorageClass is the provider-specific storage class or tier.
	// +optional
	StorageClass string `json:"storageClass,omitempty"`
}

// MachineRequestStatus defines the observed state of MachineRequest.
type MachineRequestStatus struct {
	// Phase represents the current lifecycle phase of the machine.
	// +optional
	Phase MachinePhase `json:"phase,omitempty"`

	// ProviderID is the provider-specific identifier for the machine.
	// Format is provider-specific (e.g., Harvester VM UID, Nutanix VM UUID).
	// +optional
	ProviderID string `json:"providerID,omitempty"`

	// IPAddress is the primary IP address of the machine.
	// This is set when the machine reaches the Running phase.
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`

	// IPAddresses contains all IP addresses assigned to the machine.
	// +optional
	IPAddresses []string `json:"ipAddresses,omitempty"`

	// MACAddress is the primary MAC address of the machine.
	// +optional
	MACAddress string `json:"macAddress,omitempty"`

	// FailureReason provides a machine-readable failure reason.
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage provides a human-readable failure message.
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// Conditions represent the latest available observations of the
	// MachineRequest's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastUpdated is the timestamp of the last status update.
	// +optional
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`

	// ObservedGeneration is the generation most recently observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=mr
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".spec.machineName",description="VM name"
// +kubebuilder:printcolumn:name="Role",type="string",JSONPath=".spec.role",description="Machine role"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current phase"
// +kubebuilder:printcolumn:name="IP",type="string",JSONPath=".status.ipAddress",description="IP address"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// MachineRequest is the Schema for the machinerequests API.
// It represents a request to create a virtual machine on an infrastructure provider.
// This resource serves as the interface contract between the Butler bootstrap
// controller and provider-specific controllers.
type MachineRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineRequestSpec   `json:"spec,omitempty"`
	Status MachineRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineRequestList contains a list of MachineRequest.
type MachineRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachineRequest{}, &MachineRequestList{})
}

// Helper methods for MachineRequest

// IsReady returns true if the machine is in the Running phase with an IP address.
func (mr *MachineRequest) IsReady() bool {
	return mr.Status.Phase == MachinePhaseRunning && mr.Status.IPAddress != ""
}

// IsFailed returns true if the machine is in a failed state.
func (mr *MachineRequest) IsFailed() bool {
	return mr.Status.Phase == MachinePhaseFailed
}

// IsTerminating returns true if the machine is being deleted.
func (mr *MachineRequest) IsTerminating() bool {
	return mr.Status.Phase == MachinePhaseDeleting || mr.Status.Phase == MachinePhaseDeleted
}

// SetPhase updates the phase and last updated timestamp.
func (mr *MachineRequest) SetPhase(phase MachinePhase) {
	mr.Status.Phase = phase
	now := metav1.Now()
	mr.Status.LastUpdated = &now
}

// SetFailure sets the failure reason and message.
func (mr *MachineRequest) SetFailure(reason, message string) {
	mr.Status.FailureReason = reason
	mr.Status.FailureMessage = message
	mr.SetPhase(MachinePhaseFailed)
}
