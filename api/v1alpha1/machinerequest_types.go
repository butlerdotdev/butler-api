/*
Copyright 2026 Butler Labs.
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

// MachineRequestSpec defines the desired state of MachineRequest
type MachineRequestSpec struct {
	// ProviderRef references the provider configuration to use
	ProviderRef ProviderReference `json:"providerRef"`

	// MachineName is the name for the machine
	MachineName string `json:"machineName"`

	// Role is the machine role (control-plane or worker)
	Role MachineRole `json:"role"`

	// CPU is the number of CPUs
	// +kubebuilder:validation:Minimum=1
	CPU int32 `json:"cpu"`

	// MemoryMB is the memory in MB
	// +kubebuilder:validation:Minimum=1024
	MemoryMB int32 `json:"memoryMB"`

	// DiskGB is the disk size in GB
	// +kubebuilder:validation:Minimum=10
	DiskGB int32 `json:"diskGB"`

	// ExtraDisks defines additional disks to attach
	// +optional
	ExtraDisks []MachineRequestDisk `json:"extraDisks,omitempty"`

	// Labels to apply to the machine
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// UserData is the cloud-init or machine configuration
	// +optional
	UserData string `json:"userData,omitempty"`
}

// MachineRequestDisk defines an extra disk
type MachineRequestDisk struct {
	// SizeGB is the disk size in GB
	// +kubebuilder:validation:Minimum=1
	SizeGB int32 `json:"sizeGB"`
}

// MachineRequestStatus defines the observed state of MachineRequest
type MachineRequestStatus struct {
	// Phase is the current phase of the machine
	Phase MachinePhase `json:"phase,omitempty"`

	// ProviderID is the provider-specific ID for the machine
	// +optional
	ProviderID string `json:"providerID,omitempty"`

	// IPAddress is the primary IP address of the machine
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`

	// MACAddress is the MAC address of the machine
	// +optional
	MACAddress string `json:"macAddress,omitempty"`

	// FailureReason is the reason for failure if Phase is Failed
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage is a human-readable message for the failure
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// Conditions represent the latest observations
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastUpdated is the last time the status was updated
	// +optional
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="IP",type="string",JSONPath=".status.ipAddress"
// +kubebuilder:printcolumn:name="Role",type="string",JSONPath=".spec.role"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// MachineRequest is the Schema for the machinerequests API
type MachineRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineRequestSpec   `json:"spec,omitempty"`
	Status MachineRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineRequestList contains a list of MachineRequest
type MachineRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachineRequest{}, &MachineRequestList{})
}
