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

// ProviderType defines the supported infrastructure providers.
// +kubebuilder:validation:Enum=harvester;nutanix;proxmox
type ProviderType string

const (
	// ProviderTypeHarvester is the Harvester HCI provider.
	ProviderTypeHarvester ProviderType = "harvester"

	// ProviderTypeNutanix is the Nutanix AHV provider.
	ProviderTypeNutanix ProviderType = "nutanix"

	// ProviderTypeProxmox is the Proxmox VE provider.
	ProviderTypeProxmox ProviderType = "proxmox"
)

// ProviderConfigSpec defines the desired state of ProviderConfig.
type ProviderConfigSpec struct {
	// Provider specifies the infrastructure provider type.
	// +kubebuilder:validation:Required
	Provider ProviderType `json:"provider"`

	// CredentialsRef references the Secret containing provider credentials.
	// The Secret must contain the appropriate keys for the provider type:
	// - harvester: "kubeconfig" (Harvester kubeconfig)
	// - nutanix: "username", "password"
	// - proxmox: "username", "password" or "token"
	// +kubebuilder:validation:Required
	CredentialsRef SecretReference `json:"credentialsRef"`

	// Harvester contains Harvester-specific configuration.
	// Required when provider is "harvester".
	// +optional
	Harvester *HarvesterProviderConfig `json:"harvester,omitempty"`

	// Nutanix contains Nutanix-specific configuration.
	// Required when provider is "nutanix".
	// +optional
	Nutanix *NutanixProviderConfig `json:"nutanix,omitempty"`

	// Proxmox contains Proxmox-specific configuration.
	// Required when provider is "proxmox".
	// +optional
	Proxmox *ProxmoxProviderConfig `json:"proxmox,omitempty"`
}

// HarvesterProviderConfig contains Harvester-specific configuration.
type HarvesterProviderConfig struct {
	// Endpoint is the Harvester API server URL.
	// If not specified, extracted from the kubeconfig.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// Namespace is the Harvester namespace for VM resources.
	// +kubebuilder:default="default"
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NetworkName is the VM network in "namespace/name" format.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[a-z0-9-]+/[a-z0-9-]+$`
	NetworkName string `json:"networkName"`

	// ImageName is the default OS image in "namespace/name" format.
	// Used when MachineRequest doesn't specify an image.
	// +optional
	ImageName string `json:"imageName,omitempty"`

	// StorageClassName is the default storage class for VM disks.
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`
}

// NutanixProviderConfig contains Nutanix-specific configuration.
type NutanixProviderConfig struct {
	// Endpoint is the Prism Central API URL.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^https?://`
	Endpoint string `json:"endpoint"`

	// Port is the Prism Central API port.
	// +kubebuilder:default=9440
	// +optional
	Port int32 `json:"port,omitempty"`

	// Insecure allows insecure TLS connections.
	// +kubebuilder:default=false
	// +optional
	Insecure bool `json:"insecure,omitempty"`

	// ClusterUUID is the target Nutanix cluster UUID.
	// +kubebuilder:validation:Required
	ClusterUUID string `json:"clusterUUID"`

	// SubnetUUID is the network subnet UUID for VMs.
	// +kubebuilder:validation:Required
	SubnetUUID string `json:"subnetUUID"`

	// ImageUUID is the default OS image UUID.
	// Used when MachineRequest doesn't specify an image.
	// +optional
	ImageUUID string `json:"imageUUID,omitempty"`

	// StorageContainerUUID is the storage container for VM disks.
	// +optional
	StorageContainerUUID string `json:"storageContainerUUID,omitempty"`
}

// ProxmoxProviderConfig contains Proxmox-specific configuration.
type ProxmoxProviderConfig struct {
	// Endpoint is the Proxmox API URL.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^https?://`
	Endpoint string `json:"endpoint"`

	// Insecure allows insecure TLS connections.
	// +kubebuilder:default=false
	// +optional
	Insecure bool `json:"insecure,omitempty"`

	// Nodes is the list of Proxmox nodes available for VM placement.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Nodes []string `json:"nodes"`

	// Storage is the storage location for VM disks.
	// +kubebuilder:validation:Required
	Storage string `json:"storage"`

	// TemplateID is the location of the ISO to use for VM creation
	// i.e. local:iso/ubuntu-22.04.iso
	// +optional
	TemplateID string `json:"templateID,omitempty"`

	// VMIDRange defines the range of VM IDs to use.
	// +optional
	VMIDRange *VMIDRange `json:"vmidRange,omitempty"`
}

// VMIDRange defines a range of VM IDs.
type VMIDRange struct {
	// Start is the first VM ID in the range.
	// +kubebuilder:validation:Minimum=100
	Start int32 `json:"start"`

	// End is the last VM ID in the range.
	// +kubebuilder:validation:Minimum=100
	End int32 `json:"end"`
}

// ProviderConfigStatus defines the observed state of ProviderConfig.
type ProviderConfigStatus struct {
	// Conditions represent the latest available observations of the ProviderConfig's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Validated indicates whether the provider configuration has been validated.
	// +optional
	Validated bool `json:"validated,omitempty"`

	// LastValidationTime is the timestamp of the last successful validation.
	// +optional
	LastValidationTime *metav1.Time `json:"lastValidationTime,omitempty"`

	// ProviderVersion is the detected version of the infrastructure provider.
	// +optional
	ProviderVersion string `json:"providerVersion,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=pc
// +kubebuilder:printcolumn:name="Provider",type="string",JSONPath=".spec.provider",description="Infrastructure provider type"
// +kubebuilder:printcolumn:name="Validated",type="boolean",JSONPath=".status.validated",description="Configuration validated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ProviderConfig defines the configuration for an infrastructure provider.
// It contains credentials and provider-specific settings needed to create
// and manage virtual machines.
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec,omitempty"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig.
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProviderConfig{}, &ProviderConfigList{})
}
