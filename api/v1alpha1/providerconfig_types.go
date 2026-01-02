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

// ProviderType defines the type of infrastructure provider
// +kubebuilder:validation:Enum=harvester;nutanix;proxmox
type ProviderType string

const (
	ProviderTypeHarvester ProviderType = "harvester"
	ProviderTypeNutanix   ProviderType = "nutanix"
	ProviderTypeProxmox   ProviderType = "proxmox"
)

// ProviderConfigSpec defines the desired state of ProviderConfig
type ProviderConfigSpec struct {
	// Provider is the type of infrastructure provider
	Provider ProviderType `json:"provider"`

	// Harvester contains Harvester-specific configuration
	// +optional
	Harvester *HarvesterProviderConfig `json:"harvester,omitempty"`

	// Nutanix contains Nutanix-specific configuration
	// +optional
	Nutanix *NutanixProviderConfig `json:"nutanix,omitempty"`

	// Proxmox contains Proxmox-specific configuration
	// +optional
	Proxmox *ProxmoxProviderConfig `json:"proxmox,omitempty"`
}

// HarvesterProviderConfig defines Harvester-specific configuration
type HarvesterProviderConfig struct {
	// KubeconfigSecretRef references a secret containing the Harvester kubeconfig
	KubeconfigSecretRef SecretReference `json:"kubeconfigSecretRef"`

	// Namespace is the Harvester namespace for VMs
	// +kubebuilder:default="default"
	Namespace string `json:"namespace,omitempty"`

	// NetworkName is the Harvester network name
	NetworkName string `json:"networkName"`

	// ImageName is the Talos image name in Harvester
	ImageName string `json:"imageName"`
}

// NutanixProviderConfig defines Nutanix-specific configuration
type NutanixProviderConfig struct {
	// Endpoint is the Prism Central endpoint URL
	Endpoint string `json:"endpoint"`

	// CredentialsSecretRef references a secret containing username/password
	CredentialsSecretRef SecretReference `json:"credentialsSecretRef"`

	// ClusterUUID is the Nutanix cluster UUID
	ClusterUUID string `json:"clusterUUID"`

	// SubnetUUID is the Nutanix subnet UUID
	SubnetUUID string `json:"subnetUUID"`

	// ImageUUID is the Talos image UUID
	ImageUUID string `json:"imageUUID"`

	// Insecure allows insecure TLS connections
	// +optional
	Insecure bool `json:"insecure,omitempty"`
}

// ProxmoxProviderConfig defines Proxmox-specific configuration
type ProxmoxProviderConfig struct {
	// Endpoint is the Proxmox API endpoint URL
	Endpoint string `json:"endpoint"`

	// CredentialsSecretRef references a secret containing username/password or API token
	CredentialsSecretRef SecretReference `json:"credentialsSecretRef"`

	// Node is the Proxmox node name
	Node string `json:"node"`

	// StorageLocation is the storage location for VMs
	StorageLocation string `json:"storageLocation"`

	// VMIDRange defines the range for VM IDs
	VMIDRange VMIDRange `json:"vmidRange"`

	// TemplateID is the VM template to clone
	TemplateID int `json:"templateID"`

	// Insecure allows insecure TLS connections
	// +optional
	Insecure bool `json:"insecure,omitempty"`
}

// VMIDRange defines a range of VM IDs
type VMIDRange struct {
	// Start is the starting VM ID
	Start int `json:"start"`

	// End is the ending VM ID
	End int `json:"end"`
}

// ProviderConfigStatus defines the observed state of ProviderConfig
type ProviderConfigStatus struct {
	// Ready indicates the provider is ready to be used
	Ready bool `json:"ready,omitempty"`

	// Message provides additional status information
	// +optional
	Message string `json:"message,omitempty"`

	// Conditions represent the latest observations
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastUpdated is the last time the status was updated
	// +optional
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Provider",type="string",JSONPath=".spec.provider"
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ProviderConfig is the Schema for the providerconfigs API
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec,omitempty"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProviderConfig{}, &ProviderConfigList{})
}
