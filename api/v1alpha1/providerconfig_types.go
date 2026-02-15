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
// +kubebuilder:validation:Enum=harvester;nutanix;proxmox;azure;aws;gcp
type ProviderType string

const (
	// ProviderTypeHarvester is the Harvester HCI provider.
	ProviderTypeHarvester ProviderType = "harvester"

	// ProviderTypeNutanix is the Nutanix AHV provider.
	ProviderTypeNutanix ProviderType = "nutanix"

	// ProviderTypeProxmox is the Proxmox VE provider.
	ProviderTypeProxmox ProviderType = "proxmox"

	// ProviderTypeAzure is the Microsoft Azure provider.
	ProviderTypeAzure ProviderType = "azure"

	// ProviderTypeAWS is the Amazon Web Services provider.
	ProviderTypeAWS ProviderType = "aws"

	// ProviderTypeGCP is the Google Cloud Platform provider.
	ProviderTypeGCP ProviderType = "gcp"
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

	// Azure contains Azure-specific configuration.
	// Required when provider is "azure".
	// +optional
	Azure *AzureProviderConfig `json:"azure,omitempty"`

	// AWS contains AWS-specific configuration.
	// Required when provider is "aws".
	// +optional
	AWS *AWSProviderConfig `json:"aws,omitempty"`

	// GCP contains GCP-specific configuration.
	// Required when provider is "gcp".
	// +optional
	GCP *GCPProviderConfig `json:"gcp,omitempty"`

	// Scope defines the visibility of this ProviderConfig.
	// Platform-scoped providers are available to all teams.
	// Team-scoped providers are restricted to a specific team.
	// +optional
	Scope *ProviderConfigScope `json:"scope,omitempty"`

	// Network configures IPAM and network settings for this provider.
	// +optional
	Network *ProviderNetworkConfig `json:"network,omitempty"`

	// Limits defines resource limits enforced per-team on this provider.
	// +optional
	Limits *ProviderLimits `json:"limits,omitempty"`
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

	// TemplateID is the VM template ID to clone.
	// +optional
	TemplateID int32 `json:"templateID,omitempty"`

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

// AzureProviderConfig contains Azure-specific configuration.
type AzureProviderConfig struct {
	// SubscriptionID is the Azure subscription ID.
	// +kubebuilder:validation:Required
	SubscriptionID string `json:"subscriptionID"`

	// ResourceGroup is the Azure resource group.
	// +kubebuilder:validation:Required
	ResourceGroup string `json:"resourceGroup"`

	// Location is the Azure region.
	// +optional
	Location string `json:"location,omitempty"`

	// VNetName is the Azure Virtual Network name.
	// +optional
	VNetName string `json:"vnetName,omitempty"`

	// SubnetName is the subnet within the VNet.
	// +optional
	SubnetName string `json:"subnetName,omitempty"`
}

// AWSProviderConfig contains AWS-specific configuration.
type AWSProviderConfig struct {
	// Region is the AWS region.
	// +kubebuilder:validation:Required
	Region string `json:"region"`

	// VPCID is the VPC identifier.
	// +optional
	VPCID string `json:"vpcID,omitempty"`

	// SubnetIDs are the subnet identifiers for VM placement.
	// +optional
	SubnetIDs []string `json:"subnetIDs,omitempty"`

	// SecurityGroupIDs are the security group identifiers.
	// +optional
	SecurityGroupIDs []string `json:"securityGroupIDs,omitempty"`
}

// GCPProviderConfig contains GCP-specific configuration.
type GCPProviderConfig struct {
	// ProjectID is the GCP project identifier.
	// +kubebuilder:validation:Required
	ProjectID string `json:"projectID"`

	// Region is the GCP region.
	// +kubebuilder:validation:Required
	Region string `json:"region"`

	// Network is the VPC network name.
	// +optional
	Network string `json:"network,omitempty"`

	// Subnetwork is the subnetwork name.
	// +optional
	Subnetwork string `json:"subnetwork,omitempty"`
}

// ProviderConfigScopeType defines the visibility scope.
// +kubebuilder:validation:Enum=platform;team
type ProviderConfigScopeType string

const (
	// ProviderConfigScopePlatform means the provider is available to all teams.
	ProviderConfigScopePlatform ProviderConfigScopeType = "platform"

	// ProviderConfigScopeTeam means the provider is restricted to a specific team.
	ProviderConfigScopeTeam ProviderConfigScopeType = "team"
)

// ProviderConfigScope defines the visibility of a ProviderConfig.
type ProviderConfigScope struct {
	// Type is the scope type.
	// +kubebuilder:default="platform"
	// +optional
	Type ProviderConfigScopeType `json:"type,omitempty"`

	// TeamRef references the Team when type is "team".
	// Required when type is "team".
	// +optional
	TeamRef *LocalObjectReference `json:"teamRef,omitempty"`
}

// ProviderNetworkConfig configures IPAM and network settings.
type ProviderNetworkConfig struct {
	// Mode determines how IP addresses are managed.
	// "ipam" uses NetworkPool-based automated allocation.
	// "cloud" relies on the cloud provider's native networking.
	// +kubebuilder:validation:Enum=ipam;cloud
	// +kubebuilder:default="cloud"
	// +optional
	Mode string `json:"mode,omitempty"`

	// PoolRefs references NetworkPools for IPAM allocation, ordered by priority.
	// Required when mode is "ipam". Allocator tries first pool, falls back to next if exhausted.
	// +optional
	PoolRefs []PoolReference `json:"poolRefs,omitempty"`

	// Subnet is the network name for VM placement (e.g., "VM Network - VLAN 40").
	// +optional
	Subnet string `json:"subnet,omitempty"`

	// Gateway is the network gateway address.
	// +optional
	Gateway string `json:"gateway,omitempty"`

	// DNSServers are the DNS server addresses.
	// +optional
	DNSServers []string `json:"dnsServers,omitempty"`

	// LoadBalancer configures load balancer IP allocation defaults.
	// +optional
	LoadBalancer *ProviderLBConfig `json:"loadBalancer,omitempty"`

	// QuotaPerTenant defines per-tenant network resource quotas.
	// +optional
	QuotaPerTenant *NetworkQuota `json:"quotaPerTenant,omitempty"`
}

// PoolReference references a NetworkPool with a priority.
type PoolReference struct {
	// Name is the name of the NetworkPool.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Priority determines allocation order (lower = higher priority).
	// Pools at the same priority are tried in list order.
	// +kubebuilder:default=0
	// +optional
	Priority *int32 `json:"priority,omitempty"`
}

// ProviderLBConfig configures load balancer defaults.
type ProviderLBConfig struct {
	// DefaultPoolSize is the default number of LB IPs per tenant in static mode.
	// +kubebuilder:default=8
	// +kubebuilder:validation:Minimum=1
	// +optional
	DefaultPoolSize *int32 `json:"defaultPoolSize,omitempty"`

	// AllocationMode controls how LB IPs are allocated to tenants.
	// "static" pre-allocates a fixed block (DefaultPoolSize IPs).
	// "elastic" starts small and grows/shrinks based on usage.
	// +kubebuilder:validation:Enum=static;elastic
	// +kubebuilder:default="static"
	// +optional
	AllocationMode string `json:"allocationMode,omitempty"`

	// InitialPoolSize is the number of LB IPs initially allocated in elastic mode.
	// +kubebuilder:default=2
	// +kubebuilder:validation:Minimum=1
	// +optional
	InitialPoolSize *int32 `json:"initialPoolSize,omitempty"`

	// GrowthIncrement is the number of IPs added per expansion in elastic mode.
	// +kubebuilder:default=2
	// +kubebuilder:validation:Minimum=1
	// +optional
	GrowthIncrement *int32 `json:"growthIncrement,omitempty"`
}

// NetworkQuota defines per-tenant network resource quotas.
type NetworkQuota struct {
	// MaxNodeIPs limits the number of node IPs per tenant.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxNodeIPs *int32 `json:"maxNodeIPs,omitempty"`

	// MaxLoadBalancerIPs limits the number of LB IPs per tenant.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxLoadBalancerIPs *int32 `json:"maxLoadBalancerIPs,omitempty"`
}

// ProviderLimits defines per-team resource limits on a provider.
type ProviderLimits struct {
	// MaxClustersPerTeam limits the number of clusters per team.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxClustersPerTeam *int32 `json:"maxClustersPerTeam,omitempty"`

	// MaxNodesPerTeam limits the total nodes per team.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxNodesPerTeam *int32 `json:"maxNodesPerTeam,omitempty"`
}

// ProviderCapacity reports the available capacity of a provider.
type ProviderCapacity struct {
	// AvailableIPs is the number of available IPs across all pools.
	// +optional
	AvailableIPs int32 `json:"availableIPs,omitempty"`

	// EstimatedTenants is the estimated number of tenants that can be provisioned.
	// +optional
	EstimatedTenants int32 `json:"estimatedTenants,omitempty"`
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

	// Ready indicates overall readiness of the provider.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// LastProbeTime is the timestamp of the last health probe.
	// +optional
	LastProbeTime *metav1.Time `json:"lastProbeTime,omitempty"`

	// Capacity reports the available capacity of this provider.
	// +optional
	Capacity *ProviderCapacity `json:"capacity,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=pc
// +kubebuilder:printcolumn:name="Provider",type="string",JSONPath=".spec.provider",description="Infrastructure provider type"
// +kubebuilder:printcolumn:name="Scope",type="string",JSONPath=".spec.scope.type",description="Visibility scope"
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready",description="Provider ready"
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
