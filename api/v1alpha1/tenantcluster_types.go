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

// ManagementMode defines how Butler manages addons after initial installation.
// +kubebuilder:validation:Enum=Active;Observe;GitOps
type ManagementMode string

const (
	// ManagementModeActive means Butler actively manages addons.
	// New addons in spec are installed. Removal is user-initiated.
	ManagementModeActive ManagementMode = "Active"

	// ManagementModeObserve means Butler only observes after initial install.
	// Changes to spec.addons are ignored after cluster is ready.
	ManagementModeObserve ManagementMode = "Observe"

	// ManagementModeGitOps means Butler bootstraps Flux and hands off.
	// Flux manages the cluster from the configured Git repository.
	ManagementModeGitOps ManagementMode = "GitOps"
)

// OSType defines the operating system for worker nodes.
// +kubebuilder:validation:Enum=rocky;flatcar;talos
type OSType string

const (
	// OSTypeRocky is Rocky Linux.
	OSTypeRocky OSType = "rocky"

	// OSTypeFlatcar is Flatcar Container Linux.
	OSTypeFlatcar OSType = "flatcar"

	// OSTypeTalos is Talos Linux (immutable OS).
	OSTypeTalos OSType = "talos"
)

// TenantClusterSpec defines the desired state of TenantCluster.
type TenantClusterSpec struct {
	// KubernetesVersion is the target Kubernetes version.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^v\d+\.\d+\.\d+$`
	KubernetesVersion string `json:"kubernetesVersion"`

	// TeamRef references the Team this cluster belongs to.
	// Required when multi-tenancy mode is Enforced.
	// +optional
	TeamRef *LocalObjectReference `json:"teamRef,omitempty"`

	// ProviderConfigRef references the ProviderConfig for infrastructure.
	// If not specified, defaults are used (Team's or platform's).
	// +optional
	ProviderConfigRef *LocalObjectReference `json:"providerConfigRef,omitempty"`

	// ControlPlane configures the Steward-hosted control plane.
	// +optional
	ControlPlane ControlPlaneSpec `json:"controlPlane,omitempty"`

	// Workers configures the worker nodes.
	// +kubebuilder:validation:Required
	Workers WorkersSpec `json:"workers"`

	// Networking configures cluster networking.
	// +optional
	Networking NetworkingSpec `json:"networking,omitempty"`

	// ManagementPolicy defines how Butler manages this cluster.
	// +optional
	ManagementPolicy ManagementPolicySpec `json:"managementPolicy,omitempty"`

	// Addons defines the initial addons to install.
	// These are installed at cluster creation time.
	// Additional addons can be added via TenantAddon resources.
	// +optional
	Addons AddonsSpec `json:"addons,omitempty"`

	// InfrastructureOverride allows overriding provider-specific settings.
	// These take precedence over ProviderConfig defaults.
	// +optional
	InfrastructureOverride *InfrastructureOverride `json:"infrastructureOverride,omitempty"`
}

// ControlPlaneSpec configures the Steward-hosted control plane.
type ControlPlaneSpec struct {
	// Replicas is the number of API server replicas.
	// Steward manages high availability automatically.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// DataStoreRef references the Steward DataStore to use.
	// If not specified, the default DataStore is used.
	// +optional
	DataStoreRef *LocalObjectReference `json:"dataStoreRef,omitempty"`

	// ServiceType for the control plane endpoint.
	// If not specified, inherits from ButlerConfig.spec.controlPlaneExposure.mode.
	// Only set this to override the platform-level setting for this specific cluster.
	// +kubebuilder:validation:Enum=LoadBalancer;NodePort;ClusterIP
	// +optional
	ServiceType string `json:"serviceType,omitempty"`

	// CertSANs are additional Subject Alternative Names for the API server certificate.
	// Use this to add custom DNS names or IPs for API server access.
	// +optional
	CertSANs []string `json:"certSANs,omitempty"`

	// ExternalCloudProvider enables --cloud-provider=external on apiserver and controller-manager.
	// Required for Harvester, vSphere, and other infrastructure providers.
	// +kubebuilder:default=true
	// +optional
	ExternalCloudProvider *bool `json:"externalCloudProvider,omitempty"`
}

// WorkersSpec configures worker nodes.
type WorkersSpec struct {
	// Replicas is the desired number of worker nodes.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas"`

	// MachineTemplate defines the VM specification for workers.
	// +optional
	MachineTemplate MachineTemplateSpec `json:"machineTemplate,omitempty"`
}

// MachineTemplateSpec defines VM specifications.
type MachineTemplateSpec struct {
	// CPU is the number of CPU cores.
	// +kubebuilder:default=4
	// +kubebuilder:validation:Minimum=1
	// +optional
	CPU int32 `json:"cpu,omitempty"`

	// Memory is the amount of RAM.
	// +kubebuilder:default="16Gi"
	// +optional
	Memory resource.Quantity `json:"memory,omitempty"`

	// DiskSize is the root disk size.
	// +kubebuilder:default="100Gi"
	// +optional
	DiskSize resource.Quantity `json:"diskSize,omitempty"`

	// OS configures the operating system.
	// +optional
	OS OSSpec `json:"os,omitempty"`
}

// OSSpec configures the operating system.
type OSSpec struct {
	// Type is the OS type.
	// +kubebuilder:default="rocky"
	// +optional
	Type OSType `json:"type,omitempty"`

	// Version is the OS version.
	// +kubebuilder:default="9.5"
	// +optional
	Version string `json:"version,omitempty"`

	// ImageRef references a specific image to use.
	// Overrides Type and Version if specified.
	// +optional
	ImageRef string `json:"imageRef,omitempty"`

	// Talos provides Talos-specific worker node configuration.
	// Required when type is "talos".
	// +optional
	Talos *TalosConfig `json:"talos,omitempty"`
}

// TalosConfig provides Talos-specific worker node configuration.
type TalosConfig struct {
	// InstallDisk is the disk where Talos will be installed.
	// +kubebuilder:default="/dev/vda"
	// +optional
	InstallDisk string `json:"installDisk,omitempty"`

	// InstallerImage is the Talos installer image
	// (e.g., factory.talos.dev/installer/<schematic>:v1.9.3).
	// +optional
	InstallerImage string `json:"installerImage,omitempty"`

	// Version is the Talos version.
	// +kubebuilder:default="v1.9.3"
	// +optional
	Version string `json:"version,omitempty"`
}

// InfrastructureOverride allows overriding provider-specific settings per-cluster.
type InfrastructureOverride struct {
	// Harvester contains Harvester-specific overrides.
	// +optional
	Harvester *HarvesterOverride `json:"harvester,omitempty"`

	// Nutanix contains Nutanix-specific overrides.
	// +optional
	Nutanix *NutanixOverride `json:"nutanix,omitempty"`

	// Proxmox contains Proxmox-specific overrides.
	// +optional
	Proxmox *ProxmoxOverride `json:"proxmox,omitempty"`
}

// HarvesterOverride contains Harvester-specific settings.
type HarvesterOverride struct {
	// Namespace is the Harvester namespace for VMs.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NetworkName is the Harvester network to use (format: namespace/name).
	// +optional
	NetworkName string `json:"networkName,omitempty"`

	// ImageName is the VM image to use (format: namespace/name).
	// +optional
	ImageName string `json:"imageName,omitempty"`
}

// NutanixOverride contains Nutanix-specific settings.
type NutanixOverride struct {
	// ClusterUUID is the Nutanix cluster UUID.
	// +optional
	ClusterUUID string `json:"clusterUUID,omitempty"`

	// SubnetUUID is the Nutanix subnet UUID.
	// +optional
	SubnetUUID string `json:"subnetUUID,omitempty"`

	// ImageUUID is the Nutanix image UUID.
	// +optional
	ImageUUID string `json:"imageUUID,omitempty"`

	// StorageContainerUUID is the Nutanix storage container UUID.
	// +optional
	StorageContainerUUID string `json:"storageContainerUUID,omitempty"`
}

// ProxmoxOverride contains Proxmox-specific settings.
type ProxmoxOverride struct {
	// Node is the Proxmox node to deploy VMs on.
	// +optional
	Node string `json:"node,omitempty"`

	// Storage is the Proxmox storage to use.
	// +optional
	Storage string `json:"storage,omitempty"`

	// TemplateID is the VM template ID.
	// +optional
	TemplateID int `json:"templateID,omitempty"`
}

// NetworkingSpec configures cluster networking.
type NetworkingSpec struct {
	// PodCIDR is the CIDR for pod IPs.
	// +kubebuilder:default="10.244.0.0/16"
	// +optional
	PodCIDR string `json:"podCIDR,omitempty"`

	// ServiceCIDR is the CIDR for service IPs.
	// +kubebuilder:default="10.96.0.0/12"
	// +optional
	ServiceCIDR string `json:"serviceCIDR,omitempty"`

	// LoadBalancerPool defines the IP pool for LoadBalancer services.
	// +optional
	LoadBalancerPool *IPPool `json:"loadBalancerPool,omitempty"`
}

// IPPool defines a range of IP addresses.
type IPPool struct {
	// Start is the first IP in the pool.
	// +kubebuilder:validation:Required
	Start string `json:"start"`

	// End is the last IP in the pool.
	// +kubebuilder:validation:Required
	End string `json:"end"`
}

// ManagementPolicySpec defines how Butler manages the cluster.
type ManagementPolicySpec struct {
	// Mode determines how Butler manages addons.
	// +kubebuilder:default="Active"
	// +optional
	Mode ManagementMode `json:"mode,omitempty"`
}

// AddonsSpec defines the addons to install.
type AddonsSpec struct {
	// CNI configures the Container Network Interface.
	// +optional
	CNI *CNISpec `json:"cni,omitempty"`

	// LoadBalancer configures the load balancer.
	// +optional
	LoadBalancer *LoadBalancerSpec `json:"loadBalancer,omitempty"`

	// CertManager configures cert-manager.
	// +optional
	CertManager *CertManagerSpec `json:"certManager,omitempty"`

	// Storage configures persistent storage.
	// +optional
	Storage *StorageSpec `json:"storage,omitempty"`

	// Ingress configures the ingress controller.
	// +optional
	Ingress *IngressSpec `json:"ingress,omitempty"`

	// GitOps configures GitOps (Flux or ArgoCD).
	// +optional
	GitOps *GitOpsSpec `json:"gitops,omitempty"`
}

// CNISpec configures the CNI addon.
type CNISpec struct {
	// Provider is the CNI provider.
	// +kubebuilder:validation:Enum=cilium
	// +kubebuilder:default="cilium"
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the addon version.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for customization.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// LoadBalancerSpec configures the load balancer addon.
type LoadBalancerSpec struct {
	// Provider is the load balancer provider.
	// +kubebuilder:validation:Enum=metallb
	// +kubebuilder:default="metallb"
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the addon version.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for customization.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// CertManagerSpec configures cert-manager.
type CertManagerSpec struct {
	// Enabled indicates whether cert-manager should be installed.
	// +kubebuilder:default=true
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// Version is the addon version.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for customization.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// StorageSpec configures persistent storage.
type StorageSpec struct {
	// Provider is the storage provider.
	// +kubebuilder:validation:Enum=longhorn;linstor
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the addon version.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for customization.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// IngressSpec configures the ingress controller.
type IngressSpec struct {
	// Provider is the ingress provider.
	// +kubebuilder:validation:Enum=traefik;nginx
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the addon version.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for customization.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// GitOpsSpec configures GitOps tooling.
type GitOpsSpec struct {
	// Provider is the GitOps provider.
	// +kubebuilder:validation:Enum=fluxcd;argocd
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the addon version.
	// +optional
	Version string `json:"version,omitempty"`

	// Repository configures the Git repository for GitOps.
	// +optional
	Repository *GitRepositorySpec `json:"repository,omitempty"`
}

// GitRepositorySpec configures a Git repository for GitOps.
type GitRepositorySpec struct {
	// URL is the Git repository URL.
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Branch is the branch to use.
	// +kubebuilder:default="main"
	// +optional
	Branch string `json:"branch,omitempty"`

	// Path is the path within the repository for this cluster's manifests.
	// +optional
	Path string `json:"path,omitempty"`

	// SecretRef references the Secret containing Git credentials.
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty"`
}

// ExtensionValues holds arbitrary Helm values.
// +kubebuilder:pruning:PreserveUnknownFields
type ExtensionValues struct {
	// Raw is the raw JSON/YAML values.
	// +optional
	Raw []byte `json:"-"`
}

// TenantClusterPhase represents the current phase of a TenantCluster.
// +kubebuilder:validation:Enum=Pending;Provisioning;Installing;Ready;Updating;Deleting;Failed
type TenantClusterPhase string

const (
	// TenantClusterPhasePending indicates the cluster is pending creation.
	TenantClusterPhasePending TenantClusterPhase = "Pending"

	// TenantClusterPhaseProvisioning indicates infrastructure is being provisioned.
	TenantClusterPhaseProvisioning TenantClusterPhase = "Provisioning"

	// TenantClusterPhaseInstalling indicates addons are being installed.
	TenantClusterPhaseInstalling TenantClusterPhase = "Installing"

	// TenantClusterPhaseReady indicates the cluster is ready for use.
	TenantClusterPhaseReady TenantClusterPhase = "Ready"

	// TenantClusterPhaseUpdating indicates the cluster is being updated.
	TenantClusterPhaseUpdating TenantClusterPhase = "Updating"

	// TenantClusterPhaseDeleting indicates the cluster is being deleted.
	TenantClusterPhaseDeleting TenantClusterPhase = "Deleting"

	// TenantClusterPhaseFailed indicates a failure occurred.
	TenantClusterPhaseFailed TenantClusterPhase = "Failed"
)

// TenantClusterStatus defines the observed state of TenantCluster.
type TenantClusterStatus struct {
	// Conditions represent the latest available observations.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current phase of the cluster.
	// +optional
	Phase TenantClusterPhase `json:"phase,omitempty"`

	// TenantNamespace is the namespace containing CAPI/Steward resources.
	// +optional
	TenantNamespace string `json:"tenantNamespace,omitempty"`

	// ControlPlaneEndpoint is the API server endpoint.
	// +optional
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint,omitempty"`

	// KubeconfigSecretRef references the Secret containing the kubeconfig.
	// +optional
	KubeconfigSecretRef *LocalObjectReference `json:"kubeconfigSecretRef,omitempty"`

	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastTransitionTime is when the phase last changed.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// ObservedState is the observed state of the cluster.
	// +optional
	ObservedState *ObservedClusterState `json:"observedState,omitempty"`

	// WorkerNodesReady is the count of ready worker nodes
	// +optional
	WorkerNodesReady int32 `json:"workerNodesReady,omitempty"`

	// WorkerNodesDesired is the desired count of worker nodes
	// +optional
	WorkerNodesDesired int32 `json:"workerNodesDesired,omitempty"`
}

// ObservedClusterState captures the current state of the cluster.
type ObservedClusterState struct {
	// KubernetesVersion is the actual Kubernetes version running.
	// +optional
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// Workers shows worker node status.
	// +optional
	Workers *WorkerStatus `json:"workers,omitempty"`

	// Addons shows installed addon status.
	// +optional
	Addons []AddonStatus `json:"addons,omitempty"`
}

// WorkerStatus shows worker node status.
type WorkerStatus struct {
	// Desired is the desired number of workers.
	Desired int32 `json:"desired"`

	// Ready is the number of ready workers.
	Ready int32 `json:"ready"`

	// Nodes lists the worker nodes.
	// +optional
	Nodes []string `json:"nodes,omitempty"`
}

// AddonStatus shows the status of an installed addon.
type AddonStatus struct {
	// Name is the addon name.
	Name string `json:"name"`

	// Version is the installed version.
	Version string `json:"version"`

	// Status is the addon health status.
	// +kubebuilder:validation:Enum=Pending;Installing;Healthy;Degraded;Failed
	Status string `json:"status"`

	// ManagedBy indicates who manages this addon.
	// +kubebuilder:validation:Enum=butler;flux;argocd;manual
	ManagedBy string `json:"managedBy"`
}

// TenantCluster condition types.
const (
	// TenantClusterConditionInfrastructureReady indicates CAPI resources are ready.
	TenantClusterConditionInfrastructureReady = "InfrastructureReady"

	// TenantClusterConditionControlPlaneReady indicates the control plane is ready.
	TenantClusterConditionControlPlaneReady = "ControlPlaneReady"

	// TenantClusterConditionWorkersReady indicates workers are ready.
	TenantClusterConditionWorkersReady = "WorkersReady"

	// TenantClusterConditionAddonsReady indicates addons are installed.
	TenantClusterConditionAddonsReady = "AddonsReady"

	// TenantClusterConditionReady indicates the cluster is fully ready.
	TenantClusterConditionReady = "Ready"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=tc
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Cluster phase"
// +kubebuilder:printcolumn:name="K8s Version",type="string",JSONPath=".spec.kubernetesVersion",description="Kubernetes version"
// +kubebuilder:printcolumn:name="Workers",type="string",JSONPath=".status.observedState.workers.ready",description="Ready workers"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.controlPlaneEndpoint",description="API endpoint"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// TenantCluster is the Schema for the tenantclusters API.
// It represents a complete Kubernetes cluster managed by Butler.
type TenantCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantClusterSpec   `json:"spec,omitempty"`
	Status TenantClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantClusterList contains a list of TenantCluster.
type TenantClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenantCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TenantCluster{}, &TenantClusterList{})
}
