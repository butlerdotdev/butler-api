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
	"encoding/binary"
	"fmt"
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterBootstrapPhase represents the current phase of bootstrap
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

// ClusterTopology defines the cluster topology type
type ClusterTopology string

const (
	// ClusterTopologySingleNode is a single-node cluster where the control plane also runs workloads.
	// This is useful for development, testing, and edge deployments where resources are limited.
	// In single-node mode:
	// - Only 1 control plane node is created
	// - Workers section is ignored
	// - Control plane is configured to allow scheduling workloads
	// - etcd runs as a single member (not HA)
	ClusterTopologySingleNode ClusterTopology = "single-node"

	// ClusterTopologyHA is a high-availability cluster with separate control plane and worker nodes.
	// This is the default and recommended mode for production deployments.
	ClusterTopologyHA ClusterTopology = "ha"
)

// ClusterBootstrapSpec defines the desired state of ClusterBootstrap
type ClusterBootstrapSpec struct {
	// Provider is the infrastructure provider type (harvester, nutanix, proxmox)
	// +kubebuilder:validation:Enum=harvester;nutanix;proxmox
	Provider string `json:"provider"`

	// ProviderRef references the ProviderConfig to use for provisioning
	// Reuses existing ProviderReference from common_types.go
	// +kubebuilder:validation:Required
	ProviderRef ProviderReference `json:"providerRef"`

	// Cluster defines the cluster configuration
	// +kubebuilder:validation:Required
	Cluster ClusterBootstrapClusterSpec `json:"cluster"`

	// Network defines network configuration for the cluster
	// +kubebuilder:validation:Required
	Network ClusterBootstrapNetworkSpec `json:"network"`

	// Talos defines Talos-specific configuration
	// +kubebuilder:validation:Required
	Talos ClusterBootstrapTalosSpec `json:"talos"`

	// Addons defines which addons to install
	// +optional
	Addons ClusterBootstrapAddonsSpec `json:"addons,omitempty"`

	// Paused can be set to true to pause reconciliation
	// +optional
	Paused bool `json:"paused,omitempty"`
}

// ClusterBootstrapClusterSpec defines the cluster topology for bootstrap
type ClusterBootstrapClusterSpec struct {
	// Name is the cluster name used for resource naming
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name"`

	// Topology defines the cluster topology
	// - "single-node": Single control plane node that also runs workloads (no workers needed)
	// - "ha": High-availability with separate control plane and worker nodes (default)
	// +kubebuilder:validation:Enum=single-node;ha
	// +kubebuilder:default=ha
	// +optional
	Topology ClusterTopology `json:"topology,omitempty"`

	// ControlPlane defines control plane node configuration
	// +kubebuilder:validation:Required
	ControlPlane ClusterBootstrapNodePool `json:"controlPlane"`

	// Workers defines worker node configuration
	// Ignored when topology is "single-node"
	// +optional
	Workers *ClusterBootstrapNodePool `json:"workers,omitempty"`
}

// ClusterBootstrapNodePool defines a pool of nodes for bootstrap
// Uses same units as MachineRequest (MemoryMB, DiskGB) for consistency
type ClusterBootstrapNodePool struct {
	// Replicas is the number of nodes in this pool
	// For single-node topology, controlPlane.replicas is forced to 1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Replicas int32 `json:"replicas"`

	// CPU is the number of CPU cores per node
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=128
	CPU int32 `json:"cpu"`

	// MemoryMB is the memory in MB per node (matches MachineRequest)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=2048
	MemoryMB int32 `json:"memoryMB"`

	// DiskGB is the root disk size in GB per node (matches MachineRequest)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=20
	DiskGB int32 `json:"diskGB"`

	// ExtraDisks defines additional disks to attach to each node
	// Reuses DiskSpec from machinerequest_types.go
	// +optional
	ExtraDisks []DiskSpec `json:"extraDisks,omitempty"`

	// Labels to apply to nodes in this pool
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// ClusterBootstrapNetworkSpec defines cluster networking for bootstrap
type ClusterBootstrapNetworkSpec struct {
	// PodCIDR is the CIDR for pod networking
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}/[0-9]{1,2}$`
	PodCIDR string `json:"podCIDR"`

	// ServiceCIDR is the CIDR for service networking
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}/[0-9]{1,2}$`
	ServiceCIDR string `json:"serviceCIDR"`

	// VIP is the virtual IP for the control plane endpoint (kube-vip)
	// This IP is used ONLY for kube-apiserver HA and must NOT be in LoadBalancerPool
	// For single-node topology, the VIP still provides a stable endpoint for the API server
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}$`
	VIP string `json:"vip"`

	// VIPInterface is the network interface for the VIP (optional, auto-detected)
	// +optional
	VIPInterface string `json:"vipInterface,omitempty"`

	// LoadBalancerPool defines the IP range for MetalLB LoadBalancer services
	// This range must NOT include the VIP address to avoid conflicts between
	// kube-vip (control plane) and MetalLB (services)
	// +optional
	LoadBalancerPool *LoadBalancerPoolSpec `json:"loadBalancerPool,omitempty"`
}

// LoadBalancerPoolSpec defines an IP address range for LoadBalancer services
type LoadBalancerPoolSpec struct {
	// Start is the first IP in the pool (inclusive)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}$`
	Start string `json:"start"`

	// End is the last IP in the pool (inclusive)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}$`
	End string `json:"end"`
}

// Validate validates the LoadBalancerPoolSpec
func (p *LoadBalancerPoolSpec) Validate() error {
	if p == nil {
		return nil
	}

	startIP := net.ParseIP(p.Start)
	if startIP == nil {
		return fmt.Errorf("invalid start IP: %s", p.Start)
	}

	endIP := net.ParseIP(p.End)
	if endIP == nil {
		return fmt.Errorf("invalid end IP: %s", p.End)
	}

	if ipToUint32(startIP) > ipToUint32(endIP) {
		return fmt.Errorf("start IP %s must be <= end IP %s", p.Start, p.End)
	}

	return nil
}

// ContainsIP checks if the given IP is within the pool range
func (p *LoadBalancerPoolSpec) ContainsIP(ip string) bool {
	if p == nil {
		return false
	}

	checkIP := net.ParseIP(ip)
	if checkIP == nil {
		return false
	}

	startIP := net.ParseIP(p.Start)
	endIP := net.ParseIP(p.End)
	if startIP == nil || endIP == nil {
		return false
	}

	checkVal := ipToUint32(checkIP)
	startVal := ipToUint32(startIP)
	endVal := ipToUint32(endIP)

	return checkVal >= startVal && checkVal <= endVal
}

// ToAddressRange returns the pool as "start-end" string for MetalLB
func (p *LoadBalancerPoolSpec) ToAddressRange() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%s-%s", p.Start, p.End)
}

// ipToUint32 converts an IPv4 address to a uint32
func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip)
}

// Validate validates the network configuration
func (n *ClusterBootstrapNetworkSpec) Validate() error {
	vip := net.ParseIP(n.VIP)
	if vip == nil {
		return fmt.Errorf("invalid VIP address: %s", n.VIP)
	}

	if n.LoadBalancerPool != nil {
		if err := n.LoadBalancerPool.Validate(); err != nil {
			return fmt.Errorf("invalid loadBalancerPool: %w", err)
		}

		if n.LoadBalancerPool.ContainsIP(n.VIP) {
			return fmt.Errorf("VIP %s must not be within loadBalancerPool range %s-%s; "+
				"kube-vip and MetalLB will conflict if they share IPs",
				n.VIP, n.LoadBalancerPool.Start, n.LoadBalancerPool.End)
		}
	}

	return nil
}

// ClusterBootstrapTalosSpec defines Talos configuration for bootstrap
type ClusterBootstrapTalosSpec struct {
	// Version is the Talos version to use
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^v[0-9]+\.[0-9]+\.[0-9]+$`
	Version string `json:"version"`

	// Schematic is the Talos factory schematic ID for the image
	// +kubebuilder:validation:Required
	Schematic string `json:"schematic"`

	// ConfigPatches allows inline Talos config patches
	// +optional
	ConfigPatches []TalosConfigPatch `json:"configPatches,omitempty"`

	// InstallDisk overrides the default install disk
	// +optional
	// +kubebuilder:default="/dev/vda"
	InstallDisk string `json:"installDisk,omitempty"`
}

// TalosConfigPatch defines a Talos config patch
type TalosConfigPatch struct {
	// Op is the patch operation (add, remove, replace)
	// +kubebuilder:validation:Enum=add;remove;replace
	Op string `json:"op"`

	// Path is the JSON path to patch
	Path string `json:"path"`

	// Value is the value to set (for add/replace)
	// +optional
	Value string `json:"value,omitempty"`
}

// ClusterBootstrapAddonsSpec defines which addons to install during bootstrap
type ClusterBootstrapAddonsSpec struct {
	// CNI defines the CNI configuration
	// +optional
	CNI *CNIAddonSpec `json:"cni,omitempty"`

	// Storage defines storage configuration
	// +optional
	Storage *StorageAddonSpec `json:"storage,omitempty"`

	// LoadBalancer defines load balancer configuration
	// +optional
	LoadBalancer *LoadBalancerAddonSpec `json:"loadBalancer,omitempty"`

	// GitOps defines GitOps configuration
	// +optional
	GitOps *GitOpsAddonSpec `json:"gitOps,omitempty"`

	// ControlPlaneHA defines control plane HA configuration
	// +optional
	ControlPlaneHA *ControlPlaneHAAddonSpec `json:"controlPlaneHA,omitempty"`

	// CertManager defines cert-manager configuration
	// +optional
	CertManager *CertManagerAddonSpec `json:"certManager,omitempty"`

	// Ingress defines ingress controller configuration
	// +optional
	Ingress *IngressAddonSpec `json:"ingress,omitempty"`

	// ControlPlaneProvider defines hosted control plane provider (Kamaji)
	// +optional
	ControlPlaneProvider *ControlPlaneProviderAddonSpec `json:"controlPlaneProvider,omitempty"`

	// CAPI defines Cluster API configuration
	// +optional
	CAPI *CAPIAddonSpec `json:"capi,omitempty"`

	// ButlerController defines butler-controller configuration
	// +optional
	ButlerController *ButlerControllerAddonSpec `json:"butlerController,omitempty"`

	// Console defines Butler Console configuration
	// +optional
	Console *ConsoleAddonSpec `json:"console,omitempty"`
}

// CNIAddonSpec defines CNI configuration
type CNIAddonSpec struct {
	// Type is the CNI type
	// +kubebuilder:validation:Enum=cilium;none
	// +kubebuilder:default=cilium
	Type string `json:"type,omitempty"`

	// Version is the addon version
	// +optional
	Version string `json:"version,omitempty"`

	// HubbleEnabled enables Hubble observability (Cilium only)
	// +optional
	// +kubebuilder:default=true
	HubbleEnabled bool `json:"hubbleEnabled,omitempty"`
}

// StorageAddonSpec defines storage configuration
type StorageAddonSpec struct {
	// Type is the storage type
	// +kubebuilder:validation:Enum=longhorn;none
	// +kubebuilder:default=longhorn
	Type string `json:"type,omitempty"`

	// Version is the addon version
	// +optional
	Version string `json:"version,omitempty"`

	// ReplicaCount is the default replica count for Longhorn volumes
	// For single-node topology, this is automatically set to 1
	// +optional
	// +kubebuilder:default=3
	ReplicaCount *int32 `json:"replicaCount,omitempty"`
}

// LoadBalancerAddonSpec defines load balancer configuration
type LoadBalancerAddonSpec struct {
	// Type is the load balancer type
	// +kubebuilder:validation:Enum=metallb;none
	// +kubebuilder:default=metallb
	Type string `json:"type,omitempty"`

	// AddressPool is the IP address range for MetalLB
	// DEPRECATED: Use network.loadBalancerPool instead for proper validation
	// +optional
	AddressPool string `json:"addressPool,omitempty"`
}

// GitOpsAddonSpec defines GitOps configuration
type GitOpsAddonSpec struct {
	// Type is the GitOps type
	// +kubebuilder:validation:Enum=flux;none
	// +kubebuilder:default=flux
	Type string `json:"type,omitempty"`

	// Enabled controls whether GitOps is installed
	// +optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`
}

// ControlPlaneHAAddonSpec defines control plane HA configuration
type ControlPlaneHAAddonSpec struct {
	// Type is the control plane HA type
	// +kubebuilder:validation:Enum=kube-vip;none
	// +kubebuilder:default=kube-vip
	Type string `json:"type,omitempty"`

	// Version is the addon version
	// +optional
	Version string `json:"version,omitempty"`
}

// CertManagerAddonSpec defines cert-manager configuration
type CertManagerAddonSpec struct {
	// Enabled controls whether cert-manager is installed
	// +optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the addon version
	// +optional
	Version string `json:"version,omitempty"`
}

// IngressAddonSpec defines ingress controller configuration
type IngressAddonSpec struct {
	// Type is the ingress controller type
	// +kubebuilder:validation:Enum=traefik;nginx;none
	// +kubebuilder:default=traefik
	Type string `json:"type,omitempty"`

	// Enabled controls whether the ingress controller is installed
	// +optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the addon version
	// +optional
	Version string `json:"version,omitempty"`
}

// ControlPlaneProviderAddonSpec defines hosted control plane provider configuration
type ControlPlaneProviderAddonSpec struct {
	// Type is the control plane provider type
	// +kubebuilder:validation:Enum=kamaji;none
	// +kubebuilder:default=kamaji
	Type string `json:"type,omitempty"`

	// Enabled controls whether Kamaji is installed
	// +optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the addon version
	// +optional
	Version string `json:"version,omitempty"`
}

// CAPIAddonSpec defines Cluster API configuration
type CAPIAddonSpec struct {
	// Enabled controls whether CAPI is installed
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the CAPI core version
	// +kubebuilder:default="v1.9.4"
	// +optional
	Version string `json:"version,omitempty"`

	// InfrastructureProviders lists additional infrastructure providers to install
	// The management cluster's provider is ALWAYS included automatically
	// +optional
	InfrastructureProviders []CAPIInfraProviderSpec `json:"infrastructureProviders,omitempty"`
}

// CAPIInfraProviderSpec defines an infrastructure provider configuration
type CAPIInfraProviderSpec struct {
	// Name is the provider name
	// +kubebuilder:validation:Enum=harvester;nutanix;proxmox
	Name string `json:"name"`

	// Version overrides the default provider version
	// +optional
	Version string `json:"version,omitempty"`

	// CredentialsSecretRef points to provider credentials
	// Required for providers other than the management cluster's provider
	// +optional
	CredentialsSecretRef *SecretReference `json:"credentialsSecretRef,omitempty"`
}

// ButlerControllerAddonSpec defines Butler controller configuration
type ButlerControllerAddonSpec struct {
	// Enabled controls whether butler-controller is installed
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the butler-controller version (image tag)
	// +kubebuilder:default="latest"
	// +optional
	Version string `json:"version,omitempty"`

	// Image is the full image reference (overrides default)
	// +optional
	// +kubebuilder:default="ghcr.io/butlerdotdev/butler-controller"
	Image string `json:"image,omitempty"`
}

// ConsoleAddonSpec defines Butler Console configuration
type ConsoleAddonSpec struct {
	// Enabled controls whether butler-console is installed
	// +kubebuilder:default=false
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Version is the console version (image tag)
	// +kubebuilder:default="latest"
	// +optional
	Version string `json:"version,omitempty"`

	// Ingress defines ingress configuration for the console
	// +optional
	Ingress *ConsoleIngressSpec `json:"ingress,omitempty"`
}

// ConsoleIngressSpec defines ingress configuration for the Butler Console
type ConsoleIngressSpec struct {
	// Enabled controls whether to create an Ingress resource
	// +kubebuilder:default=false
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// Host is the hostname for the console (e.g., "butler.example.com")
	// If not set and ingress is enabled, uses "butler.<cluster-name>.local"
	// +optional
	Host string `json:"host,omitempty"`

	// ClassName is the ingress class (e.g., "traefik", "nginx")
	// +optional
	ClassName string `json:"className,omitempty"`

	// TLS enables TLS termination
	// +kubebuilder:default=false
	// +optional
	TLS bool `json:"tls,omitempty"`

	// TLSSecretName is the name of the TLS secret
	// +optional
	TLSSecretName string `json:"tlsSecretName,omitempty"`
}

// ClusterBootstrapStatus defines the observed state of ClusterBootstrap
type ClusterBootstrapStatus struct {
	// Phase is the current phase of bootstrap
	// +optional
	Phase ClusterBootstrapPhase `json:"phase,omitempty"`

	// ControlPlaneEndpoint is the endpoint for the control plane
	// +optional
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint,omitempty"`

	// Kubeconfig contains the base64-encoded kubeconfig for the cluster
	// +optional
	Kubeconfig string `json:"kubeconfig,omitempty"`

	// TalosConfig contains the base64-encoded talosconfig for the cluster
	// +optional
	TalosConfig string `json:"talosconfig,omitempty"`

	// ConsoleURL is the URL to access the Butler Console
	// +optional
	ConsoleURL string `json:"consoleURL,omitempty"`

	// Machines contains the status of each machine
	// +optional
	Machines []ClusterBootstrapMachineStatus `json:"machines,omitempty"`

	// FailureReason indicates why bootstrap failed
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage provides details about the failure
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// Conditions represents the current conditions of the ClusterBootstrap
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastUpdated is the timestamp of the last status update
	// +optional
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`

	// ObservedGeneration is the last observed generation
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// AddonsInstalled tracks which addons have been installed
	// +optional
	AddonsInstalled map[string]bool `json:"addonsInstalled,omitempty"`
}

// ClusterBootstrapMachineStatus tracks the status of a machine in the cluster
type ClusterBootstrapMachineStatus struct {
	// Name is the MachineRequest name
	Name string `json:"name"`

	// Role is the machine role (control-plane or worker)
	Role string `json:"role"`

	// Phase is the MachineRequest phase
	Phase string `json:"phase"`

	// IPAddress is the machine's IP address
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`

	// TalosConfigured indicates if Talos config has been applied
	// +optional
	TalosConfigured bool `json:"talosConfigured,omitempty"`

	// Ready indicates if the node has joined the cluster
	// +optional
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=cb
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.cluster.name"
// +kubebuilder:printcolumn:name="Topology",type="string",JSONPath=".spec.cluster.topology"
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

// Helper methods

// IsReady returns true if the cluster bootstrap is complete
func (c *ClusterBootstrap) IsReady() bool {
	return c.Status.Phase == ClusterBootstrapPhaseReady
}

// IsFailed returns true if the cluster bootstrap has failed
func (c *ClusterBootstrap) IsFailed() bool {
	return c.Status.Phase == ClusterBootstrapPhaseFailed
}

// IsSingleNode returns true if this is a single-node topology
func (c *ClusterBootstrap) IsSingleNode() bool {
	return c.Spec.Cluster.Topology == ClusterTopologySingleNode
}

// GetExpectedMachineCount returns the expected number of machines based on topology
func (c *ClusterBootstrap) GetExpectedMachineCount() int {
	if c.IsSingleNode() {
		// Single-node: only 1 control plane, ignore workers
		return 1
	}
	// HA: control plane replicas + worker replicas
	count := int(c.Spec.Cluster.ControlPlane.Replicas)
	if c.Spec.Cluster.Workers != nil {
		count += int(c.Spec.Cluster.Workers.Replicas)
	}
	return count
}

// GetControlPlaneReplicas returns the effective control plane replicas based on topology
func (c *ClusterBootstrap) GetControlPlaneReplicas() int32 {
	if c.IsSingleNode() {
		return 1 // Single-node always has exactly 1 control plane
	}
	return c.Spec.Cluster.ControlPlane.Replicas
}

// GetControlPlaneIPs returns the IP addresses of control plane nodes
func (c *ClusterBootstrap) GetControlPlaneIPs() []string {
	var ips []string
	for _, m := range c.Status.Machines {
		if m.Role == string(MachineRoleControlPlane) && m.IPAddress != "" {
			ips = append(ips, m.IPAddress)
		}
	}
	return ips
}

// GetWorkerIPs returns the IP addresses of worker nodes
func (c *ClusterBootstrap) GetWorkerIPs() []string {
	var ips []string
	for _, m := range c.Status.Machines {
		if m.Role == string(MachineRoleWorker) && m.IPAddress != "" {
			ips = append(ips, m.IPAddress)
		}
	}
	return ips
}

// AllMachinesRunning returns true if all machines are in Running phase with IPs
func (c *ClusterBootstrap) AllMachinesRunning() bool {
	expectedCount := c.GetExpectedMachineCount()

	if len(c.Status.Machines) != expectedCount {
		return false
	}

	for _, m := range c.Status.Machines {
		if m.Phase != string(MachinePhaseRunning) || m.IPAddress == "" {
			return false
		}
	}
	return true
}

// IsCAPIEnabled returns whether CAPI should be installed
func (s *ClusterBootstrapAddonsSpec) IsCAPIEnabled() bool {
	if s == nil || s.CAPI == nil || s.CAPI.Enabled == nil {
		return true // Default enabled
	}
	return *s.CAPI.Enabled
}

// GetCAPIVersion returns the CAPI version to install
func (s *ClusterBootstrapAddonsSpec) GetCAPIVersion() string {
	if s == nil || s.CAPI == nil || s.CAPI.Version == "" {
		return "v1.9.4"
	}
	return s.CAPI.Version
}

// IsButlerControllerEnabled returns whether butler-controller should be installed
func (s *ClusterBootstrapAddonsSpec) IsButlerControllerEnabled() bool {
	if s == nil || s.ButlerController == nil || s.ButlerController.Enabled == nil {
		return true // Default enabled
	}
	return *s.ButlerController.Enabled
}

// GetButlerControllerImage returns the full butler-controller image reference
func (s *ClusterBootstrapAddonsSpec) GetButlerControllerImage() string {
	image := "ghcr.io/butlerdotdev/butler-controller"
	version := "latest"

	if s != nil && s.ButlerController != nil {
		if s.ButlerController.Image != "" {
			image = s.ButlerController.Image
		}
		if s.ButlerController.Version != "" {
			version = s.ButlerController.Version
		}
	}

	return image + ":" + version
}

// GetLoadBalancerAddressPool returns the address pool string for MetalLB
// Prefers network.loadBalancerPool (validated), falls back to addons.loadBalancer.addressPool (legacy)
func (c *ClusterBootstrap) GetLoadBalancerAddressPool() string {
	// Prefer network.loadBalancerPool (new way with validation)
	if c.Spec.Network.LoadBalancerPool != nil {
		return c.Spec.Network.LoadBalancerPool.ToAddressRange()
	}

	// Fall back to addons.loadBalancer.addressPool (legacy, deprecated)
	if c.Spec.Addons.LoadBalancer != nil && c.Spec.Addons.LoadBalancer.AddressPool != "" {
		return c.Spec.Addons.LoadBalancer.AddressPool
	}

	return ""
}

// IsConsoleEnabled returns whether butler-console should be installed
func (s *ClusterBootstrapAddonsSpec) IsConsoleEnabled() bool {
	if s == nil || s.Console == nil || s.Console.Enabled == nil {
		return false // Default disabled - user must opt-in
	}
	return *s.Console.Enabled
}

// GetConsoleVersion returns the console version to install
func (s *ClusterBootstrapAddonsSpec) GetConsoleVersion() string {
	if s == nil || s.Console == nil || s.Console.Version == "" {
		return "latest"
	}
	return s.Console.Version
}

// GetConsoleIngressHost returns the ingress host, with fallback to cluster name
func (s *ClusterBootstrapAddonsSpec) GetConsoleIngressHost(clusterName string) string {
	if s == nil || s.Console == nil || s.Console.Ingress == nil || s.Console.Ingress.Host == "" {
		return fmt.Sprintf("butler.%s.local", clusterName)
	}
	return s.Console.Ingress.Host
}

// GetStorageReplicaCount returns the effective storage replica count based on topology
func (c *ClusterBootstrap) GetStorageReplicaCount() int32 {
	if c.IsSingleNode() {
		return 1 // Single-node can only have 1 replica
	}
	if c.Spec.Addons.Storage != nil && c.Spec.Addons.Storage.ReplicaCount != nil {
		return *c.Spec.Addons.Storage.ReplicaCount
	}
	return 3 // Default for HA
}
