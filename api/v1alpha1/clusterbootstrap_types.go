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

// ClusterBootstrapSpec defines the desired state of ClusterBootstrap
type ClusterBootstrapSpec struct {
	// ProviderRef references the provider configuration to use
	ProviderRef ProviderReference `json:"providerRef"`

	// Cluster defines the cluster configuration
	Cluster ClusterBootstrapClusterSpec `json:"cluster"`

	// Network defines the network configuration
	Network ClusterBootstrapNetworkSpec `json:"network"`

	// Talos defines the Talos configuration
	Talos ClusterBootstrapTalosSpec `json:"talos"`

	// Addons defines the addons to install
	// +optional
	Addons ClusterBootstrapAddonsSpec `json:"addons,omitempty"`

	// Paused indicates the reconciliation is paused
	// +optional
	Paused bool `json:"paused,omitempty"`
}

// ClusterBootstrapClusterSpec defines the cluster configuration
type ClusterBootstrapClusterSpec struct {
	// Name is the cluster name
	Name string `json:"name"`

	// ControlPlane defines the control plane node pool
	ControlPlane ClusterBootstrapNodePool `json:"controlPlane"`

	// Workers defines the worker node pool
	// +optional
	Workers *ClusterBootstrapNodePool `json:"workers,omitempty"`
}

// ClusterBootstrapNodePool defines a node pool configuration
type ClusterBootstrapNodePool struct {
	// Replicas is the number of nodes
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas"`

	// CPU is the number of CPUs per node
	// +kubebuilder:validation:Minimum=1
	CPU int32 `json:"cpu"`

	// MemoryMB is the memory in MB per node
	// +kubebuilder:validation:Minimum=1024
	MemoryMB int32 `json:"memoryMB"`

	// DiskGB is the disk size in GB per node
	// +kubebuilder:validation:Minimum=10
	DiskGB int32 `json:"diskGB"`

	// ExtraDisks defines additional disks to attach
	// +optional
	ExtraDisks []ClusterBootstrapDisk `json:"extraDisks,omitempty"`

	// Labels to apply to nodes
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// ClusterBootstrapDisk defines an extra disk
type ClusterBootstrapDisk struct {
	// SizeGB is the disk size in GB
	// +kubebuilder:validation:Minimum=1
	SizeGB int32 `json:"sizeGB"`
}

// ClusterBootstrapNetworkSpec defines the network configuration
type ClusterBootstrapNetworkSpec struct {
	// PodCIDR is the CIDR for pod networking
	// +kubebuilder:default="10.244.0.0/16"
	PodCIDR string `json:"podCIDR,omitempty"`

	// ServiceCIDR is the CIDR for service networking
	// +kubebuilder:default="10.96.0.0/12"
	ServiceCIDR string `json:"serviceCIDR,omitempty"`

	// VIP is the virtual IP for the control plane
	VIP string `json:"vip"`

	// VIPInterface is the network interface for the VIP
	// +kubebuilder:default="enp1s0"
	// +optional
	VIPInterface string `json:"vipInterface,omitempty"`
}

// ClusterBootstrapTalosSpec defines the Talos configuration
type ClusterBootstrapTalosSpec struct {
	// Version is the Talos version
	// +kubebuilder:default="v1.9.2"
	Version string `json:"version,omitempty"`

	// Schematic is the Talos schematic ID for custom images
	// +optional
	Schematic string `json:"schematic,omitempty"`

	// InstallDisk is the disk to install Talos on
	// +kubebuilder:default="/dev/vda"
	InstallDisk string `json:"installDisk,omitempty"`

	// ConfigPatches are additional config patches to apply
	// +optional
	ConfigPatches []ClusterBootstrapConfigPatch `json:"configPatches,omitempty"`
}

// ClusterBootstrapConfigPatch defines a config patch
type ClusterBootstrapConfigPatch struct {
	// Op is the operation (add, replace, remove)
	Op string `json:"op"`

	// Path is the JSON path
	Path string `json:"path"`

	// Value is the value to set
	// +optional
	Value string `json:"value,omitempty"`
}

// ClusterBootstrapAddonsSpec defines the addons to install
type ClusterBootstrapAddonsSpec struct {
	// ControlPlaneHA defines the control plane HA solution (kube-vip)
	// +optional
	ControlPlaneHA *ControlPlaneHASpec `json:"controlPlaneHA,omitempty"`

	// CNI defines the CNI plugin (cilium)
	// +optional
	CNI *CNISpec `json:"cni,omitempty"`

	// CertManager defines the cert-manager configuration
	// +optional
	CertManager *CertManagerSpec `json:"certManager,omitempty"`

	// Storage defines the storage solution (longhorn)
	// +optional
	Storage *StorageSpec `json:"storage,omitempty"`

	// LoadBalancer defines the load balancer (metallb)
	// +optional
	LoadBalancer *LoadBalancerSpec `json:"loadBalancer,omitempty"`

	// Ingress defines the ingress controller (traefik)
	// +optional
	Ingress *IngressSpec `json:"ingress,omitempty"`

	// ControlPlaneProvider defines the hosted control plane provider (kamaji)
	// +optional
	ControlPlaneProvider *ControlPlaneProviderSpec `json:"controlPlaneProvider,omitempty"`

	// GitOps defines the GitOps solution (flux)
	// +optional
	GitOps *GitOpsSpec `json:"gitOps,omitempty"`
}

// ControlPlaneHASpec defines the control plane HA configuration
type ControlPlaneHASpec struct {
	// Type is the HA solution type
	// +kubebuilder:validation:Enum=kube-vip
	// +kubebuilder:default="kube-vip"
	Type string `json:"type,omitempty"`

	// Version is the kube-vip version
	// +optional
	Version string `json:"version,omitempty"`
}

// CNISpec defines the CNI configuration
type CNISpec struct {
	// Type is the CNI type
	// +kubebuilder:validation:Enum=cilium
	// +kubebuilder:default="cilium"
	Type string `json:"type,omitempty"`

	// Version is the CNI version
	// +optional
	Version string `json:"version,omitempty"`

	// HubbleEnabled enables Hubble observability
	// +kubebuilder:default=true
	// +optional
	HubbleEnabled bool `json:"hubbleEnabled,omitempty"`
}

// CertManagerSpec defines the cert-manager configuration
type CertManagerSpec struct {
	// Enabled controls whether cert-manager is installed
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the cert-manager version
	// +optional
	Version string `json:"version,omitempty"`
}

// StorageSpec defines the storage configuration
type StorageSpec struct {
	// Type is the storage type
	// +kubebuilder:validation:Enum=longhorn
	// +kubebuilder:default="longhorn"
	Type string `json:"type,omitempty"`

	// Version is the storage version
	// +optional
	Version string `json:"version,omitempty"`

	// DefaultReplicaCount is the default replica count for volumes
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=2
	// +optional
	DefaultReplicaCount int32 `json:"defaultReplicaCount,omitempty"`
}

// LoadBalancerSpec defines the load balancer configuration
type LoadBalancerSpec struct {
	// Type is the load balancer type
	// +kubebuilder:validation:Enum=metallb
	// +kubebuilder:default="metallb"
	Type string `json:"type,omitempty"`

	// Version is the load balancer version
	// +optional
	Version string `json:"version,omitempty"`

	// AddressPool is the IP address pool for MetalLB
	// Example: "10.40.0.200-10.40.0.250"
	AddressPool string `json:"addressPool,omitempty"`
}

// IngressSpec defines the ingress controller configuration
type IngressSpec struct {
	// Type is the ingress controller type
	// +kubebuilder:validation:Enum=traefik
	// +kubebuilder:default="traefik"
	Type string `json:"type,omitempty"`

	// Enabled controls whether the ingress controller is installed
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the ingress controller version
	// +optional
	Version string `json:"version,omitempty"`
}

// ControlPlaneProviderSpec defines the hosted control plane provider configuration
type ControlPlaneProviderSpec struct {
	// Type is the control plane provider type
	// +kubebuilder:validation:Enum=kamaji
	// +kubebuilder:default="kamaji"
	Type string `json:"type,omitempty"`

	// Enabled controls whether Kamaji is installed
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the Kamaji version
	// +optional
	Version string `json:"version,omitempty"`
}

// GitOpsSpec defines the GitOps configuration
type GitOpsSpec struct {
	// Type is the GitOps type
	// +kubebuilder:validation:Enum=flux
	// +kubebuilder:default="flux"
	Type string `json:"type,omitempty"`

	// Enabled controls whether GitOps is installed
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the GitOps version
	// +optional
	Version string `json:"version,omitempty"`

	// Repository is the Git repository URL
	// +optional
	Repository string `json:"repository,omitempty"`

	// Branch is the Git branch
	// +optional
	Branch string `json:"branch,omitempty"`

	// Path is the path within the repository
	// +optional
	Path string `json:"path,omitempty"`
}

// ClusterBootstrapPhase represents the current phase of the bootstrap
type ClusterBootstrapPhase string

const (
	ClusterBootstrapPhasePending              ClusterBootstrapPhase = "Pending"
	ClusterBootstrapPhaseProvisioningMachines ClusterBootstrapPhase = "ProvisioningMachines"
	ClusterBootstrapPhaseConfiguringTalos     ClusterBootstrapPhase = "ConfiguringTalos"
	ClusterBootstrapPhaseBootstrappingCluster ClusterBootstrapPhase = "BootstrappingCluster"
	ClusterBootstrapPhaseInstallingAddons     ClusterBootstrapPhase = "InstallingAddons"
	ClusterBootstrapPhasePivoting             ClusterBootstrapPhase = "Pivoting"
	ClusterBootstrapPhaseReady                ClusterBootstrapPhase = "Ready"
	ClusterBootstrapPhaseFailed               ClusterBootstrapPhase = "Failed"
)

// ClusterBootstrapStatus defines the observed state of ClusterBootstrap
type ClusterBootstrapStatus struct {
	// Phase is the current phase of the bootstrap
	// +optional
	Phase ClusterBootstrapPhase `json:"phase,omitempty"`

	// ControlPlaneEndpoint is the endpoint for the control plane
	// +optional
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint,omitempty"`

	// Kubeconfig is the base64-encoded kubeconfig for the cluster
	// +optional
	Kubeconfig string `json:"kubeconfig,omitempty"`

	// TalosConfig is the base64-encoded talosconfig for the cluster
	// +optional
	TalosConfig string `json:"talosconfig,omitempty"`

	// Machines contains the status of each machine
	// +optional
	Machines []ClusterBootstrapMachineStatus `json:"machines,omitempty"`

	// AddonsInstalled tracks which addons have been installed
	// +optional
	AddonsInstalled map[string]bool `json:"addonsInstalled,omitempty"`

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

	// ObservedGeneration is the last observed generation
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// ClusterBootstrapMachineStatus defines the status of a machine
type ClusterBootstrapMachineStatus struct {
	// Name is the machine name
	Name string `json:"name"`

	// Role is the machine role (control-plane or worker)
	Role string `json:"role"`

	// Phase is the machine phase
	Phase string `json:"phase"`

	// IPAddress is the machine IP address
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`

	// TalosConfigured indicates if Talos config has been applied
	// +optional
	TalosConfigured bool `json:"talosConfigured,omitempty"`

	// Ready indicates if the machine is ready
	// +optional
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.controlPlaneEndpoint"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterBootstrap is the Schema for the clusterbootstraps API
type ClusterBootstrap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterBootstrapSpec   `json:"spec,omitempty"`
	Status ClusterBootstrapStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterBootstrapList contains a list of ClusterBootstrap
type ClusterBootstrapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterBootstrap `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterBootstrap{}, &ClusterBootstrapList{})
}
