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
// This determines the ongoing relationship between Butler and the cluster's addons.
// +kubebuilder:validation:Enum=Active;Observe;GitOps
type ManagementMode string

const (
	// ManagementModeActive means Butler actively manages addons throughout
	// the cluster lifecycle. New addons added to spec are installed automatically.
	// Addon removal must be user-initiated to prevent accidental deletion.
	ManagementModeActive ManagementMode = "Active"

	// ManagementModeObserve means Butler only observes after initial installation.
	// Changes to spec.addons are ignored once the cluster reaches Ready state.
	// Useful for teams that want Butler for provisioning but manual addon control.
	ManagementModeObserve ManagementMode = "Observe"

	// ManagementModeGitOps means Butler bootstraps Flux and hands off addon management.
	// After initial setup, Flux manages the cluster from the configured Git repository.
	// Butler continues to manage infrastructure but not workloads/addons.
	ManagementModeGitOps ManagementMode = "GitOps"
)

// OSType defines the operating system for worker nodes.
// +kubebuilder:validation:Enum=rocky;flatcar
type OSType string

const (
	// OSTypeRocky uses Rocky Linux, a RHEL-compatible enterprise distribution.
	// Recommended for production workloads requiring long-term support.
	OSTypeRocky OSType = "rocky"

	// OSTypeFlatcar uses Flatcar Container Linux, an immutable container-optimized OS.
	// Recommended for security-focused deployments with automatic updates.
	OSTypeFlatcar OSType = "flatcar"
)

// TenantClusterSpec defines the desired state of TenantCluster.
// A TenantCluster represents a complete Kubernetes cluster with hosted control plane
// (via Kamaji) and worker nodes provisioned on the configured infrastructure provider.
type TenantClusterSpec struct {
	// KubernetesVersion is the target Kubernetes version for this cluster.
	// Must be a supported version (check ButlerConfig for allowed versions).
	// Format: vX.Y.Z (e.g., v1.31.4)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^v\d+\.\d+\.\d+$`
	KubernetesVersion string `json:"kubernetesVersion"`

	// TeamRef references the Team this cluster belongs to.
	// Required when ButlerConfig.spec.multiTenancy.mode is Enforced.
	// The Team determines namespace placement, resource quotas, and RBAC.
	// +optional
	TeamRef *LocalObjectReference `json:"teamRef,omitempty"`

	// ProviderConfigRef references the ProviderConfig for infrastructure provisioning.
	// If not specified, the Team's default or platform default ProviderConfig is used.
	// This determines which infrastructure (Harvester, Nutanix, etc.) hosts the workers.
	// +optional
	ProviderConfigRef *LocalObjectReference `json:"providerConfigRef,omitempty"`

	// ControlPlane configures the Kamaji-hosted control plane.
	// Kamaji runs control plane components as pods in the management cluster,
	// providing efficient multi-tenancy without dedicated control plane VMs.
	// +optional
	ControlPlane ControlPlaneSpec `json:"controlPlane,omitempty"`

	// Workers configures the worker nodes that run tenant workloads.
	// Workers are provisioned as VMs on the infrastructure provider.
	// +kubebuilder:validation:Required
	Workers WorkersSpec `json:"workers"`

	// Networking configures cluster networking including pod/service CIDRs
	// and LoadBalancer IP pools for tenant services.
	// +optional
	Networking NetworkingSpec `json:"networking,omitempty"`

	// ManagementPolicy defines how Butler manages this cluster after creation.
	// Controls addon lifecycle management and upgrade behavior.
	// +optional
	ManagementPolicy ManagementPolicySpec `json:"managementPolicy,omitempty"`

	// Addons defines the initial addons to install during cluster creation.
	// These form the base platform capabilities (CNI, storage, ingress, etc.).
	// Additional addons can be added later via TenantAddon resources.
	// +optional
	Addons AddonsSpec `json:"addons,omitempty"`

	// InfrastructureOverride allows overriding provider-specific settings for this cluster.
	// These take precedence over ProviderConfig defaults, useful for cluster-specific
	// requirements like dedicated networks or storage.
	// +optional
	InfrastructureOverride *InfrastructureOverride `json:"infrastructureOverride,omitempty"`
}

// ControlPlaneSpec configures the Kamaji-hosted control plane.
// Kamaji provides hosted Kubernetes control planes that run as pods in the
// management cluster, enabling efficient multi-tenancy.
type ControlPlaneSpec struct {
	// Replicas is the number of API server replicas for high availability.
	// Kamaji manages leader election and etcd clustering automatically.
	// Use 1 for development, 3 for production HA.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// DataStoreRef references the Kamaji DataStore for etcd storage.
	// If not specified, the platform default DataStore is used.
	// DataStores can be shared across tenants or dedicated per-tenant.
	// +optional
	DataStoreRef *LocalObjectReference `json:"dataStoreRef,omitempty"`

	// ExposureMode defines how the control plane API server is exposed to clients.
	// This determines how kubectl and other tools connect to the cluster.
	// If not specified, inherits from ButlerConfig.spec.controlPlane.defaultExposureMode.
	//
	// Options:
	// - Gateway: SNI-based routing via Gateway API (most IP-efficient, recommended)
	// - LoadBalancer: Dedicated LoadBalancer IP per cluster (simple but IP-intensive)
	// - NodePort: NodePort service (for edge/restricted environments)
	// +optional
	ExposureMode ControlPlaneExposureMode `json:"exposureMode,omitempty"`

	// Gateway contains Gateway-specific configuration for control plane exposure.
	// Only used when ExposureMode is Gateway.
	// +optional
	Gateway *TenantGatewayConfig `json:"gateway,omitempty"`

	// CertSANs are additional Subject Alternative Names for the API server certificate.
	// Use this to add custom DNS names or IPs for API server access.
	// When using Gateway mode, the generated hostname is automatically added.
	// +optional
	CertSANs []string `json:"certSANs,omitempty"`

	// ExternalCloudProvider enables --cloud-provider=external on apiserver and
	// controller-manager. Required for cloud provider integrations like Harvester
	// CCM, vSphere CPI, etc. that manage LoadBalancers and node metadata.
	// +kubebuilder:default=true
	// +optional
	ExternalCloudProvider *bool `json:"externalCloudProvider,omitempty"`

	// ServiceType for the control plane endpoint.
	// DEPRECATED: Use ExposureMode instead. This field is maintained for backward
	// compatibility and will be removed in v1beta1.
	// +kubebuilder:validation:Enum=LoadBalancer;NodePort;ClusterIP
	// +optional
	ServiceType string `json:"serviceType,omitempty"`
}

// TenantGatewayConfig contains per-tenant Gateway configuration.
// Allows customizing how this specific tenant's control plane is exposed
// via the shared Gateway.
type TenantGatewayConfig struct {
	// Hostname overrides the auto-generated hostname for this cluster's API server.
	// If not specified, hostname is generated as {cluster-name}.{gateway-domain}
	// where gateway-domain comes from ButlerConfig.spec.controlPlane.gateway.domain.
	//
	// Example: Setting this to "prod-api.k8s.example.com" instead of auto-generated
	// "my-cluster.k8s.example.com"
	// +optional
	Hostname string `json:"hostname,omitempty"`
}

// WorkersSpec configures the worker nodes that run tenant workloads.
type WorkersSpec struct {
	// Replicas is the desired number of worker nodes.
	// Butler uses CAPI MachineDeployments to manage worker lifecycle.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas"`

	// MachineTemplate defines the VM specification for worker nodes.
	// All workers in the cluster use the same specification.
	// +optional
	MachineTemplate MachineTemplateSpec `json:"machineTemplate,omitempty"`
}

// MachineTemplateSpec defines VM specifications for worker nodes.
type MachineTemplateSpec struct {
	// CPU is the number of CPU cores allocated to each worker VM.
	// +kubebuilder:default=4
	// +kubebuilder:validation:Minimum=1
	// +optional
	CPU int32 `json:"cpu,omitempty"`

	// Memory is the amount of RAM allocated to each worker VM.
	// Accepts Kubernetes quantity format (e.g., "16Gi", "32G").
	// +kubebuilder:default="16Gi"
	// +optional
	Memory resource.Quantity `json:"memory,omitempty"`

	// DiskSize is the root disk size for each worker VM.
	// This should be large enough for the OS, container images, and ephemeral storage.
	// Accepts Kubernetes quantity format (e.g., "100Gi", "200G").
	// +kubebuilder:default="100Gi"
	// +optional
	DiskSize resource.Quantity `json:"diskSize,omitempty"`

	// OS configures the operating system for worker nodes.
	// +optional
	OS OSSpec `json:"os,omitempty"`
}

// OSSpec configures the operating system for worker nodes.
type OSSpec struct {
	// Type is the OS distribution to use.
	// +kubebuilder:default="rocky"
	// +optional
	Type OSType `json:"type,omitempty"`

	// Version is the OS version (e.g., "9.5" for Rocky Linux 9.5).
	// +kubebuilder:default="9.5"
	// +optional
	Version string `json:"version,omitempty"`

	// ImageRef references a specific VM image to use, overriding Type and Version.
	// Format is provider-specific (e.g., "default/rocky-9.5-cloud" for Harvester).
	// +optional
	ImageRef string `json:"imageRef,omitempty"`
}

// InfrastructureOverride allows overriding provider-specific settings per-cluster.
// Use this when a cluster needs different infrastructure settings than the
// ProviderConfig defaults (e.g., a dedicated network or specific storage).
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

// HarvesterOverride contains Harvester-specific settings that override ProviderConfig.
type HarvesterOverride struct {
	// Namespace is the Harvester namespace where VMs will be created.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NetworkName is the Harvester VM network to attach workers to.
	// Format: namespace/name (e.g., "default/vlan-100").
	// +optional
	NetworkName string `json:"networkName,omitempty"`

	// ImageName is the VM image to use for workers.
	// Format: namespace/name (e.g., "default/rocky-9.5-cloud").
	// +optional
	ImageName string `json:"imageName,omitempty"`
}

// NutanixOverride contains Nutanix-specific settings that override ProviderConfig.
type NutanixOverride struct {
	// ClusterUUID is the Nutanix cluster UUID to deploy VMs on.
	// +optional
	ClusterUUID string `json:"clusterUUID,omitempty"`

	// SubnetUUID is the Nutanix subnet UUID for VM networking.
	// +optional
	SubnetUUID string `json:"subnetUUID,omitempty"`

	// ImageUUID is the Nutanix image UUID for the worker VM template.
	// +optional
	ImageUUID string `json:"imageUUID,omitempty"`

	// StorageContainerUUID is the Nutanix storage container for VM disks.
	// +optional
	StorageContainerUUID string `json:"storageContainerUUID,omitempty"`
}

// ProxmoxOverride contains Proxmox-specific settings that override ProviderConfig.
type ProxmoxOverride struct {
	// Node is the Proxmox node name to deploy VMs on.
	// +optional
	Node string `json:"node,omitempty"`

	// Storage is the Proxmox storage identifier for VM disks.
	// +optional
	Storage string `json:"storage,omitempty"`

	// TemplateID is the Proxmox VM template ID to clone from.
	// +optional
	TemplateID int `json:"templateID,omitempty"`
}

// NetworkingSpec configures cluster networking.
type NetworkingSpec struct {
	// PodCIDR is the CIDR block for pod IP addresses.
	// Must not overlap with ServiceCIDR or infrastructure networks.
	// +kubebuilder:default="10.244.0.0/16"
	// +optional
	PodCIDR string `json:"podCIDR,omitempty"`

	// ServiceCIDR is the CIDR block for Kubernetes service IPs.
	// Must not overlap with PodCIDR or infrastructure networks.
	// +kubebuilder:default="10.96.0.0/12"
	// +optional
	ServiceCIDR string `json:"serviceCIDR,omitempty"`

	// LoadBalancerPool defines the IP pool for LoadBalancer services in the tenant cluster.
	// These IPs are announced via MetalLB for services of type LoadBalancer.
	//
	// DEPRECATED: Use Addons.LoadBalancer.AddressPool instead.
	// If both are specified, Addons.LoadBalancer.AddressPool takes precedence.
	// This field is maintained for backward compatibility and will be removed in v1beta1.
	// +optional
	LoadBalancerPool *IPPool `json:"loadBalancerPool,omitempty"`
}

// IPPool defines a range of IP addresses for MetalLB.
type IPPool struct {
	// Start is the first IP address in the pool (inclusive).
	// +kubebuilder:validation:Required
	Start string `json:"start"`

	// End is the last IP address in the pool (inclusive).
	// +kubebuilder:validation:Required
	End string `json:"end"`
}

// ManagementPolicySpec defines how Butler manages the cluster after creation.
type ManagementPolicySpec struct {
	// Mode determines how Butler manages addons after initial installation.
	// See ManagementMode constants for detailed behavior descriptions.
	// +kubebuilder:default="Active"
	// +optional
	Mode ManagementMode `json:"mode,omitempty"`

	// AutoUpgrade enables automatic Kubernetes version upgrades.
	// When enabled, Butler upgrades clusters to newer patch versions automatically
	// during maintenance windows. Minor version upgrades require manual approval.
	// +kubebuilder:default=false
	// +optional
	AutoUpgrade bool `json:"autoUpgrade,omitempty"`

	// MaintenanceWindow defines when automatic upgrades and maintenance can occur.
	// If not specified, maintenance can occur at any time (not recommended for production).
	// +optional
	MaintenanceWindow *MaintenanceWindowSpec `json:"maintenanceWindow,omitempty"`
}

// MaintenanceWindowSpec defines when automated maintenance operations can occur.
type MaintenanceWindowSpec struct {
	// DaysOfWeek specifies which days maintenance is allowed.
	// Values: "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"
	// +optional
	DaysOfWeek []string `json:"daysOfWeek,omitempty"`

	// StartTime is the start of the maintenance window in HH:MM format (24-hour, UTC).
	// Example: "02:00" for 2 AM UTC.
	// +optional
	StartTime string `json:"startTime,omitempty"`

	// Duration is how long the maintenance window lasts.
	// Format: Go duration string (e.g., "4h", "2h30m").
	// +optional
	Duration string `json:"duration,omitempty"`
}

// AddonsSpec defines the platform addons to install on the tenant cluster.
// These addons provide essential cluster capabilities like networking, storage,
// and ingress. Butler installs these during cluster creation.
type AddonsSpec struct {
	// CNI configures the Container Network Interface plugin.
	// Required for pod networking. Cilium is the supported CNI.
	// +optional
	CNI *CNISpec `json:"cni,omitempty"`

	// LoadBalancer configures the load balancer for services of type LoadBalancer.
	// MetalLB is used to announce service IPs on the network.
	// +optional
	LoadBalancer *LoadBalancerSpec `json:"loadBalancer,omitempty"`

	// CertManager configures cert-manager for automatic certificate management.
	// Provides TLS certificates for ingress and internal services.
	// +optional
	CertManager *CertManagerSpec `json:"certManager,omitempty"`

	// Storage configures persistent storage provisioning.
	// Longhorn or LINSTOR provide replicated block storage for stateful workloads.
	// +optional
	Storage *StorageSpec `json:"storage,omitempty"`

	// Ingress configures the ingress controller for HTTP/HTTPS routing.
	// Traefik or NGINX can be used to expose services externally.
	// +optional
	Ingress *IngressSpec `json:"ingress,omitempty"`

	// GitOps configures GitOps tooling for declarative cluster management.
	// Flux CD or Argo CD can be bootstrapped for continuous deployment.
	// +optional
	GitOps *GitOpsSpec `json:"gitOps,omitempty"`
}

// CNISpec configures the Container Network Interface addon.
type CNISpec struct {
	// Provider is the CNI implementation to use.
	// Currently only Cilium is supported, providing eBPF-based networking
	// with advanced features like network policies and observability.
	// +kubebuilder:validation:Enum=cilium
	// +kubebuilder:default="cilium"
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the Cilium version to install.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for Cilium customization.
	// See Cilium Helm chart documentation for available options.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// LoadBalancerSpec configures the load balancer addon (MetalLB).
type LoadBalancerSpec struct {
	// Provider is the load balancer implementation to use.
	// MetalLB provides LoadBalancer services in bare-metal environments.
	// +kubebuilder:validation:Enum=metallb
	// +kubebuilder:default="metallb"
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the MetalLB version to install.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// AddressPool defines the IP address pool for LoadBalancer services.
	// MetalLB will assign IPs from this range to services of type LoadBalancer.
	// This is the preferred location for address pool configuration.
	// +optional
	AddressPool *AddressPoolSpec `json:"addressPool,omitempty"`

	// Values are Helm values for MetalLB customization.
	// See MetalLB Helm chart documentation for available options.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// AddressPoolSpec defines an IP address pool for MetalLB LoadBalancer services.
type AddressPoolSpec struct {
	// Start is the first IP address in the pool (inclusive).
	// +optional
	Start string `json:"start,omitempty"`

	// End is the last IP address in the pool (inclusive).
	// +optional
	End string `json:"end,omitempty"`
}

// CertManagerSpec configures the cert-manager addon.
type CertManagerSpec struct {
	// Enabled indicates whether cert-manager should be installed.
	// Set to false to skip cert-manager installation if managing certificates externally.
	// +kubebuilder:default=true
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// Version is the cert-manager version to install.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for cert-manager customization.
	// See cert-manager Helm chart documentation for available options.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// StorageSpec configures the persistent storage addon.
type StorageSpec struct {
	// Provider is the storage solution to use.
	// - longhorn: Distributed block storage with replication and snapshots
	// - linstor: DRBD-based storage with synchronous replication
	// +kubebuilder:validation:Enum=longhorn;linstor
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the storage provider version to install.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for storage provider customization.
	// See the respective Helm chart documentation for available options.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// IngressSpec configures the ingress controller addon.
type IngressSpec struct {
	// Provider is the ingress controller to use.
	// - traefik: Modern cloud-native ingress with automatic HTTPS
	// - nginx: Battle-tested ingress with extensive configuration options
	// +kubebuilder:validation:Enum=traefik;nginx
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the ingress controller version to install.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Values are Helm values for ingress controller customization.
	// See the respective Helm chart documentation for available options.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *ExtensionValues `json:"values,omitempty"`
}

// GitOpsSpec configures GitOps tooling for declarative cluster management.
type GitOpsSpec struct {
	// Provider is the GitOps tool to use.
	// - fluxcd: CNCF graduated GitOps toolkit
	// - argocd: Declarative GitOps CD for Kubernetes
	// +kubebuilder:validation:Enum=fluxcd;argocd
	// +optional
	Provider string `json:"provider,omitempty"`

	// Version is the GitOps tool version to install.
	// +optional
	Version string `json:"version,omitempty"`

	// Repository configures the Git repository for GitOps synchronization.
	// The cluster's desired state is pulled from this repository.
	// +optional
	Repository *GitRepositorySpec `json:"repository,omitempty"`
}

// GitRepositorySpec configures a Git repository for GitOps synchronization.
type GitRepositorySpec struct {
	// URL is the Git repository URL.
	// Supports HTTPS and SSH URLs (e.g., "https://github.com/org/repo.git").
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Branch is the Git branch to synchronize from.
	// +kubebuilder:default="main"
	// +optional
	Branch string `json:"branch,omitempty"`

	// Path is the directory path within the repository for this cluster's manifests.
	// Useful for monorepos with multiple clusters (e.g., "clusters/production/my-cluster").
	// +optional
	Path string `json:"path,omitempty"`

	// SecretRef references a Secret containing Git credentials.
	// Required for private repositories. The Secret should contain
	// 'username' and 'password' keys for HTTPS, or 'identity' for SSH.
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty"`
}

// ExtensionValues holds arbitrary Helm values for addon customization.
// Values are passed directly to the Helm chart during installation.
// +kubebuilder:pruning:PreserveUnknownFields
type ExtensionValues struct {
	// Raw is the raw JSON/YAML values passed to Helm.
	// +optional
	Raw []byte `json:"-"`
}

// TenantClusterPhase represents the current lifecycle phase of a TenantCluster.
// +kubebuilder:validation:Enum=Pending;Provisioning;Installing;Ready;Updating;Deleting;Failed
type TenantClusterPhase string

const (
	// TenantClusterPhasePending indicates the cluster is pending creation.
	// Butler has accepted the spec but hasn't started provisioning yet.
	TenantClusterPhasePending TenantClusterPhase = "Pending"

	// TenantClusterPhaseProvisioning indicates infrastructure is being provisioned.
	// Control plane and worker VMs are being created.
	TenantClusterPhaseProvisioning TenantClusterPhase = "Provisioning"

	// TenantClusterPhaseInstalling indicates addons are being installed.
	// Infrastructure is ready and Butler is installing CNI, storage, etc.
	TenantClusterPhaseInstalling TenantClusterPhase = "Installing"

	// TenantClusterPhaseReady indicates the cluster is ready for use.
	// All components are healthy and the cluster can accept workloads.
	TenantClusterPhaseReady TenantClusterPhase = "Ready"

	// TenantClusterPhaseUpdating indicates the cluster is being updated.
	// A spec change triggered an update operation (scaling, version upgrade, etc.).
	TenantClusterPhaseUpdating TenantClusterPhase = "Updating"

	// TenantClusterPhaseDeleting indicates the cluster is being deleted.
	// Butler is cleaning up all resources associated with the cluster.
	TenantClusterPhaseDeleting TenantClusterPhase = "Deleting"

	// TenantClusterPhaseFailed indicates a failure occurred.
	// Check conditions for detailed error information.
	TenantClusterPhaseFailed TenantClusterPhase = "Failed"
)

// TenantClusterStatus defines the observed state of TenantCluster.
type TenantClusterStatus struct {
	// Conditions represent the latest available observations of the cluster's state.
	// Standard conditions: Ready, InfrastructureReady, ControlPlaneReady, WorkersReady, AddonsReady
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current lifecycle phase of the cluster.
	// +optional
	Phase TenantClusterPhase `json:"phase,omitempty"`

	// TenantNamespace is the namespace containing CAPI and Kamaji resources for this cluster.
	// Format: tc-{cluster-name} or team-{team-name} depending on multi-tenancy mode.
	// +optional
	TenantNamespace string `json:"tenantNamespace,omitempty"`

	// ControlPlaneEndpoint is the API server endpoint URL.
	// DEPRECATED: Use ControlPlane.Endpoint instead for detailed exposure information.
	// +optional
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint,omitempty"`

	// ControlPlane contains detailed control plane exposure status.
	// Includes endpoint, exposure mode, and mode-specific details.
	// +optional
	ControlPlane *ControlPlaneStatus `json:"controlPlane,omitempty"`

	// KubeconfigSecretRef references the Secret containing the admin kubeconfig.
	// Use this kubeconfig to access the tenant cluster as cluster-admin.
	// +optional
	KubeconfigSecretRef *LocalObjectReference `json:"kubeconfigSecretRef,omitempty"`

	// ObservedGeneration is the spec generation that was last reconciled.
	// Used to detect when the spec has changed and reconciliation is needed.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastTransitionTime is when the Phase last changed.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// ObservedState captures the current observed state of the cluster.
	// Useful for comparing desired vs actual state.
	// +optional
	ObservedState *ObservedClusterState `json:"observedState,omitempty"`

	// WorkerNodesReady is the count of worker nodes in Ready state.
	// +optional
	WorkerNodesReady int32 `json:"workerNodesReady,omitempty"`

	// WorkerNodesDesired is the desired count of worker nodes from spec.
	// +optional
	WorkerNodesDesired int32 `json:"workerNodesDesired,omitempty"`
}

// ControlPlaneStatus contains detailed control plane exposure status.
// Provides information about how to connect to the tenant cluster API server.
type ControlPlaneStatus struct {
	// ExposureMode is the active exposure mode for this control plane.
	// +optional
	ExposureMode ControlPlaneExposureMode `json:"exposureMode,omitempty"`

	// Endpoint is the control plane endpoint URL for kubectl/API access.
	// Format depends on exposure mode:
	// - Gateway: https://{hostname}:6443
	// - LoadBalancer: https://{ip}:6443
	// - NodePort: https://{node-ip}:{port}
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// Hostname is the DNS hostname for API server access (Gateway mode only).
	// Clients connect using this hostname, which routes through the Gateway.
	// +optional
	Hostname string `json:"hostname,omitempty"`

	// GatewayReady indicates the Gateway TLSRoute is configured and ready (Gateway mode only).
	// When true, traffic can flow from the Gateway to this control plane.
	// +optional
	GatewayReady bool `json:"gatewayReady,omitempty"`

	// LoadBalancerIP is the allocated external IP (LoadBalancer mode only).
	// This IP is assigned by MetalLB or the cloud provider's load balancer.
	// +optional
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`

	// Ready indicates the control plane endpoint is accessible and healthy.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Message provides additional status information or error details.
	// +optional
	Message string `json:"message,omitempty"`
}

// ObservedClusterState captures the current observed state of the cluster.
type ObservedClusterState struct {
	// KubernetesVersion is the actual Kubernetes version running on the cluster.
	// May differ from spec during upgrades.
	// +optional
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// Workers shows the current state of worker nodes.
	// +optional
	Workers *WorkerStatus `json:"workers,omitempty"`

	// Addons shows the status of installed addons.
	// +optional
	Addons []AddonStatus `json:"addons,omitempty"`
}

// WorkerStatus shows the current state of worker nodes.
type WorkerStatus struct {
	// Desired is the desired number of workers from spec.
	Desired int32 `json:"desired"`

	// Ready is the number of workers in Ready state.
	Ready int32 `json:"ready"`

	// Nodes lists the names of worker nodes in the cluster.
	// +optional
	Nodes []string `json:"nodes,omitempty"`
}

// AddonStatus shows the status of an installed addon.
type AddonStatus struct {
	// Name is the addon name (e.g., "cilium", "metallb", "longhorn").
	Name string `json:"name"`

	// Version is the installed addon version.
	Version string `json:"version"`

	// Status is the addon health status.
	// +kubebuilder:validation:Enum=Pending;Installing;Healthy;Degraded;Failed
	Status string `json:"status"`

	// ManagedBy indicates who manages this addon's lifecycle.
	// +kubebuilder:validation:Enum=butler;flux;argocd;manual
	ManagedBy string `json:"managedBy"`
}

// TenantCluster condition types following Kubernetes API conventions.
const (
	// TenantClusterConditionInfrastructureReady indicates CAPI infrastructure resources are ready.
	// True when VMs are provisioned and networking is configured.
	TenantClusterConditionInfrastructureReady = "InfrastructureReady"

	// TenantClusterConditionControlPlaneReady indicates the Kamaji control plane is ready.
	// True when API server is accessible and responding to requests.
	TenantClusterConditionControlPlaneReady = "ControlPlaneReady"

	// TenantClusterConditionWorkersReady indicates worker nodes are ready.
	// True when the desired number of workers have joined and are Ready.
	TenantClusterConditionWorkersReady = "WorkersReady"

	// TenantClusterConditionAddonsReady indicates all addons are installed and healthy.
	// True when CNI, storage, and other configured addons report healthy status.
	TenantClusterConditionAddonsReady = "AddonsReady"

	// TenantClusterConditionReady indicates the cluster is fully ready for use.
	// True when all other conditions are true.
	TenantClusterConditionReady = "Ready"

	// TenantClusterConditionGatewayReady indicates Gateway routing is ready (Gateway mode only).
	// True when the TLSRoute is accepted by the Gateway and traffic can flow.
	TenantClusterConditionGatewayReady = "GatewayReady"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=tc
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current lifecycle phase"
// +kubebuilder:printcolumn:name="K8s Version",type="string",JSONPath=".spec.kubernetesVersion",description="Target Kubernetes version"
// +kubebuilder:printcolumn:name="Workers",type="string",JSONPath=".status.observedState.workers.ready",description="Ready worker count"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.controlPlane.endpoint",description="API server endpoint"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// TenantCluster is the Schema for the tenantclusters API.
// It represents a complete Kubernetes cluster managed by Butler, including
// a Kamaji-hosted control plane and worker nodes on the configured infrastructure.
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

// GetExposureMode returns the configured exposure mode for this cluster.
// Returns empty string if not explicitly set (caller should check ButlerConfig for default).
func (tc *TenantCluster) GetExposureMode() ControlPlaneExposureMode {
	return tc.Spec.ControlPlane.ExposureMode
}

// GetGatewayHostname returns the configured gateway hostname override.
// Returns empty string if using auto-generated hostname.
func (tc *TenantCluster) GetGatewayHostname() string {
	if tc.Spec.ControlPlane.Gateway != nil && tc.Spec.ControlPlane.Gateway.Hostname != "" {
		return tc.Spec.ControlPlane.Gateway.Hostname
	}
	return ""
}

// GenerateGatewayHostname generates the default hostname for this cluster given a domain.
// Format: {cluster-name}.{domain}
func (tc *TenantCluster) GenerateGatewayHostname(domain string) string {
	if domain == "" {
		return ""
	}
	return tc.Name + "." + domain
}

// GetControlPlaneEndpoint returns the control plane endpoint from status.
// Prefers the new ControlPlane.Endpoint field, falls back to deprecated ControlPlaneEndpoint.
func (tc *TenantCluster) GetControlPlaneEndpoint() string {
	if tc.Status.ControlPlane != nil && tc.Status.ControlPlane.Endpoint != "" {
		return tc.Status.ControlPlane.Endpoint
	}
	return tc.Status.ControlPlaneEndpoint
}

// GetLoadBalancerPool returns the effective LoadBalancer address pool configuration.
// Checks the preferred location (Addons.LoadBalancer.AddressPool) first,
// then falls back to the deprecated location (Networking.LoadBalancerPool).
// Returns empty strings if no pool is configured.
func (tc *TenantCluster) GetLoadBalancerPool() (start, end string) {
	// Preferred location: Addons.LoadBalancer.AddressPool
	if tc.Spec.Addons.LoadBalancer != nil && tc.Spec.Addons.LoadBalancer.AddressPool != nil {
		return tc.Spec.Addons.LoadBalancer.AddressPool.Start, tc.Spec.Addons.LoadBalancer.AddressPool.End
	}
	// Deprecated location: Networking.LoadBalancerPool (for backward compatibility)
	if tc.Spec.Networking.LoadBalancerPool != nil {
		return tc.Spec.Networking.LoadBalancerPool.Start, tc.Spec.Networking.LoadBalancerPool.End
	}
	return "", ""
}
