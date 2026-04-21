# butler-api

CRD type definitions for the Butler Kubernetes-as-a-Service platform. This repo is the API contract -- every other Butler component (controller, server, console, CLI, bootstrap, provider controllers) imports these types. It produces no running binary; its artifacts are Go types and generated CRD YAML.

**Go module:** `github.com/butlerdotdev/butler-api`
**API group:** `butler.butlerlabs.dev/v1alpha1`
**Kubebuilder:** v4.10.1, layout `go.kubebuilder.io/v4`
**Go version:** 1.24.6

## Directory Structure

```
butler-api/
  api/v1alpha1/
    groupversion_info.go         # SchemeBuilder, GroupVersion, AddToScheme
    common_types.go              # Shared reference types, labels, annotations, finalizers, conditions, reasons
    tenantcluster_types.go       # TenantCluster (tc) - the main resource
    providerconfig_types.go      # ProviderConfig (pc) - infrastructure provider credentials
    team_types.go                # Team (tm) - multi-tenant team with RBAC
    user_types.go                # User (usr) - platform user account
    butlerconfig_types.go        # ButlerConfig (bc) - singleton platform config
    clusterbootstrap_types.go    # ClusterBootstrap (cb) - management cluster bootstrap
    addondefinition_types.go     # AddonDefinition (ad) - reusable addon templates
    tenantaddon_types.go         # TenantAddon (ta) - addon instance on tenant cluster
    managementaddon_types.go     # ManagementAddon (ma) - addon on management cluster
    machinerequest_types.go      # MachineRequest (mr) - VM provisioning request
    identityprovider_types.go    # IdentityProvider (idp) - OIDC provider config
    gitprovider_types.go         # GitProvider config types (embedded in ButlerConfig)
    networkpool_types.go         # NetworkPool (np) - platform-level IP pool for IPAM
    ipallocation_types.go        # IPAllocation (ipa) - individual IP allocation from pool
    workspace_types.go           # Workspace (ws) - cloud dev environment in tenant cluster
    workspacetemplate_types.go   # WorkspaceTemplate (wst) - pre-configured workspace spec
    zz_generated.deepcopy.go     # Auto-generated (do not edit)
  config/crd/bases/              # Generated CRD YAML manifests
  hack/boilerplate.go.txt        # Copyright header for generated files
  Makefile
  PROJECT                        # Kubebuilder project metadata
```

## All CRD Resources

| Kind | Short Name | Scope | File | Purpose |
|------|-----------|-------|------|---------|
| TenantCluster | tc | Namespaced | tenantcluster_types.go | Complete tenant K8s cluster lifecycle. References Team, ProviderConfig. Creates Steward TCP + MachineRequests. |
| ProviderConfig | pc | Namespaced | providerconfig_types.go | Infrastructure provider credentials + defaults (Harvester, Nutanix, Proxmox, AWS, Azure, GCP). |
| Team | tm | Cluster | team_types.go | Multi-tenancy unit. Owns a namespace. Users/groups with roles (admin/operator/viewer). Resource quotas. |
| User | usr | Cluster | user_types.go | Platform user. Internal (email/password with invite flow) or SSO. SSH keys for workspaces. |
| ButlerConfig | bc | Cluster | butlerconfig_types.go | Singleton ("butler"). Multi-tenancy mode, default provider, team limits, git provider, control plane exposure. |
| ClusterBootstrap | cb | Namespaced | clusterbootstrap_types.go | Management cluster bootstrap via Talos. Provisions machines, configures Talos, installs addons. |
| AddonDefinition | ad, adddef | Cluster | addondefinition_types.go | Addon catalog entry. Helm chart, versions, category, dependencies. |
| TenantAddon | ta | Namespaced | tenantaddon_types.go | Per-cluster addon lifecycle. References TenantCluster and optionally AddonDefinition. |
| ManagementAddon | ma, maddon | Cluster | managementaddon_types.go | Addon on the management cluster (Steward, cert-manager, MetalLB, etc.). |
| MachineRequest | mr | Namespaced | machinerequest_types.go | VM provisioning request. Interface contract between bootstrap/controller and provider controllers. |
| IdentityProvider | idp | Cluster | identityprovider_types.go | OIDC provider (Google, Microsoft, Okta). Issuer URL, client credentials, group/email claims. |
| NetworkPool | np | Namespaced | networkpool_types.go | IP address pool for on-prem IPAM. CIDR, reserved ranges, tenant allocation defaults. |
| IPAllocation | ipa | Namespaced | ipallocation_types.go | Individual IP allocation from a NetworkPool for a TenantCluster (nodes or loadbalancer). |
| Workspace | ws | Namespaced | workspace_types.go | Cloud dev environment. Pod with SSH access in tenant cluster. Multi-repo, dotfiles, editor config. |
| WorkspaceTemplate | wst | Namespaced | workspacetemplate_types.go | Pre-configured workspace spec for one-click creation. Data-only, no controller reconciliation. |

Note: GitProvider is not a standalone CRD. `GitProviderConfig` and `GitProviderStatus` are embedded in ButlerConfig.

## Key Struct Relationships

```
ButlerConfig (singleton "butler")
  .spec.defaultProviderConfigRef    -> ProviderConfig
  .spec.defaultTeamLimits           -> ResourceLimits (platform-wide defaults)
  .spec.gitProvider                 -> GitProviderConfig (embedded, not a CRD)
  .spec.controlPlaneExposure        -> ControlPlaneExposureSpec

Team (cluster-scoped, owns namespace team-{name})
  .spec.access.users[]              -> TeamUser (email + role)
  .spec.access.groups[]             -> TeamGroup (name + role + optional IdP)
  .spec.resourceLimits              -> TeamResourceLimits (overrides ButlerConfig defaults)
  .spec.providerConfigRef           -> ProviderConfig (team-level override)
  .spec.clusterDefaults             -> ClusterDefaults
  .status.resourceUsage             -> TeamResourceUsage

TenantCluster (in team namespace)
  .spec.teamRef                     -> Team (via LocalObjectReference)
  .spec.providerConfigRef           -> ProviderConfig (via ProviderReference)
  .spec.controlPlane.dataStoreRef   -> Steward DataStore
  .spec.workers.machineTemplate     -> MachineTemplateSpec (CPU, Memory, DiskSize)
  .spec.networking.loadBalancerPool -> IPPool (start/end)
  .spec.addons                      -> AddonsSpec (CNI, LB, cert-manager, storage, ingress, gitops)
  .spec.workspaces                  -> WorkspacesConfig (enable CDE feature)
  .status.ipAllocationRef           -> IPAllocation (node IPs)
  .status.lbAllocationRef           -> IPAllocation (LB IPs)
  .status.kubeconfigSecretRef       -> Secret (admin kubeconfig)

ProviderConfig (namespaced)
  .spec.provider                    -> ProviderType enum (harvester/nutanix/proxmox/aws/azure/gcp)
  .spec.credentialsRef              -> Secret
  .spec.scope                       -> ProviderConfigScope (platform/team)
  .spec.network                     -> ProviderNetworkConfig (IPAM mode, pool refs, LB config, quotas)
  .spec.network.poolRefs[]          -> NetworkPool (ordered by priority)
  .spec.network.loadBalancer        -> ProviderLBConfig (allocation mode, elastic IPAM)
  .spec.network.quotaPerTenant      -> NetworkQuota

NetworkPool (namespaced)
  .spec.cidr                        -> Network range
  .spec.reserved[]                  -> ReservedRange (excluded from allocation)
  .spec.tenantAllocation            -> TenantAllocationConfig (allocatable sub-range + defaults)

IPAllocation (namespaced)
  .spec.poolRef                     -> NetworkPool
  .spec.tenantClusterRef            -> TenantCluster (NamespacedObjectReference)
  .spec.type                        -> IPAllocationType (nodes/loadbalancer)
  .spec.pinnedRange                 -> PinnedIPRange (for migration/reservation)

TenantAddon (namespaced)
  .spec.clusterRef                  -> TenantCluster
  .spec.addon                       -> AddonDefinition name (built-in addons)
  .spec.helm                        -> HelmChartSpec (custom charts)
  .spec.dependsOn[]                 -> TenantAddon

Workspace (namespaced, in team namespace)
  .spec.clusterRef                  -> TenantCluster
  .spec.owner                       -> User email
  .spec.repositories[]              -> WorkspaceRepository (Git repos to clone)
  .spec.envFrom                     -> WorkspaceEnvSource (copy env from tenant workload)

WorkspaceTemplate (namespaced)
  .spec.template                    -> WorkspaceTemplateBody (image, repos, resources)
  .spec.scope                       -> cluster | team

IdentityProvider (cluster-scoped)
  .spec.oidc.clientSecretRef        -> Secret
  .spec.oidc.googleWorkspace        -> GoogleWorkspaceConfig (Admin SDK for group fetch)
```

## Important Constants (common_types.go)

### Labels

```go
// Kubernetes standard
LabelManagedBy = "app.kubernetes.io/managed-by"

// Butler-specific
LabelTeam            = "butler.butlerlabs.dev/team"
LabelTenant          = "butler.butlerlabs.dev/tenant"
LabelSourceNamespace = "butler.butlerlabs.dev/source-namespace"
LabelSourceName      = "butler.butlerlabs.dev/source-name"
LabelNetworkPool     = "butler.butlerlabs.dev/network-pool"
LabelProviderConfig  = "butler.butlerlabs.dev/provider-config"
LabelWorkspaceOwner  = "butler.butlerlabs.dev/workspace-owner"
LabelAllocationType  = "butler.butlerlabs.dev/allocation-type"
```

### Annotations

```go
AnnotationDescription = "butler.butlerlabs.dev/description"
AnnotationCreatedBy   = "butler.butlerlabs.dev/created-by"
AnnotationConnect     = "butler.butlerlabs.dev/connect"
AnnotationConnectTime = "butler.butlerlabs.dev/connect-time"
```

### Finalizers

```go
FinalizerTeam           = "butler.butlerlabs.dev/team"
FinalizerTenantCluster  = "butler.butlerlabs.dev/tenantcluster"
FinalizerTenantAddon    = "butler.butlerlabs.dev/tenantaddon"
FinalizerUser           = "butler.butlerlabs.dev/user"
FinalizerNetworkPool    = "butler.butlerlabs.dev/networkpool"
FinalizerIPAllocation   = "butler.butlerlabs.dev/ipallocation"
FinalizerProviderConfig = "butler.butlerlabs.dev/providerconfig"
FinalizerWorkspace      = "butler.butlerlabs.dev/workspace"
```

### Generic Condition Types

```go
ConditionTypeReady       = "Ready"
ConditionTypeProgressing = "Progressing"
ConditionTypeDegraded    = "Degraded"
```

### Condition Reasons

The full set of condition reasons is defined in common_types.go lines 291-357. Key ones:
- `ReasonPending`, `ReasonCreating`, `ReasonCreated`, `ReasonRunning`
- `ReasonFailed`, `ReasonDeleting`, `ReasonDeleted`
- `ReasonProviderError`, `ReasonInvalidConfiguration`
- `ReasonReady`, `ReasonReconciling`, `ReasonValidationFailed`
- `ReasonQuotaExceeded`, `ReasonPoolExhausted`, `ReasonAllocationFailed`
- `ReasonProviderAccessDenied`, `ReasonNetworkNotReady`, `ReasonCredentialsInvalid`
- `ReasonNetworkReachable`, `ReasonPoolAvailable`

## Phase Enums Per Resource

| Resource | Phases |
|----------|--------|
| TenantCluster | Pending, Provisioning, Installing, Ready, Updating, Deleting, Failed |
| Team | Pending, Ready, Terminating, Failed |
| User | Pending, Active, Disabled, Locked |
| IdentityProvider | Pending, Ready, Failed |
| TenantAddon | Pending, Installing, Installed, Upgrading, Degraded, Failed, Deleting |
| ManagementAddon | Pending, Installing, Installed, Upgrading, Failed, Uninstalling |
| MachineRequest | Pending, Creating, Running, Failed, Deleting, Deleted, Unknown |
| ClusterBootstrap | Pending, ProvisioningMachines, ConfiguringTalos, BootstrappingCluster, InstallingAddons, Pivoting, Ready, Failed |
| IPAllocation | Pending, Allocated, Released, Failed |
| Workspace | Pending, Creating, Running, Starting, Stopped, Failed |

## Resource-Specific Condition Types

**TenantCluster:** InfrastructureReady, ControlPlaneReady, WorkersReady, AddonsReady, Ready, NetworkReady, ProviderAccessGranted

**Team:** NamespaceReady, RBACReady, Ready, QuotaExceeded

**User:** Ready, InvitePending, InviteExpired

**TenantAddon:** ClusterReady, DependenciesMet, Installed, Healthy, Ready

**IdentityProvider:** Discovered, SecretValid, Ready

**Workspace:** PVCReady, PodReady, RepositoryCloned, DotfilesInstalled, SSHReady, Ready

## Shared Reference Types (common_types.go)

```go
ProviderReference          // Name + optional Namespace (for ProviderConfig refs)
SecretReference            // Name + optional Namespace + optional Key
LocalObjectReference       // Name only (same namespace)
NamespacedObjectReference  // Name + Namespace (both required)
```

## TeamResourceLimits and TeamResourceUsage

`TeamResourceLimits` (common_types.go) defines team-level quotas:
- Cluster limits: MaxClusters, MaxNodesPerCluster, MaxTotalNodes
- Compute limits: MaxCPUCores, MaxMemory, MaxStorage (all `*resource.Quantity`)
- Per-cluster defaults: DefaultNodeCount, DefaultCPUPerNode, DefaultMemoryPerNode
- Feature restrictions: AllowedKubernetesVersions, AllowedProviders, AllowedAddons, DeniedAddons

`TeamResourceUsage` (common_types.go) tracks consumption:
- Clusters, TotalNodes, TotalCPU, TotalMemory, TotalStorage
- Utilization percentages: ClusterUtilization, NodeUtilization, CPUUtilization, MemoryUtilization (0-100)

These are separate from `ResourceLimits` in butlerconfig_types.go which defines platform-wide defaults.

## ProviderConfig Scope (Multi-Tenancy)

```go
type ProviderConfigScopeType string  // "platform" | "team"

type ProviderConfigScope struct {
    Type    ProviderConfigScopeType    // default: "platform"
    TeamRef *LocalObjectReference      // required when type is "team"
}
```

Platform-scoped providers are available to all teams. Team-scoped providers are restricted.

## ProviderConfig Network / IPAM

`ProviderNetworkConfig` configures IPAM:
- `Mode`: "ipam" (NetworkPool-based) or "cloud" (native cloud networking)
- `PoolRefs`: ordered list of NetworkPool references with priority
- `Subnet`, `Gateway`, `DNSServers`: static network config
- `LoadBalancer`: `ProviderLBConfig` with:
  - `AllocationMode`: "static" (fixed block) or "elastic" (grow/shrink)
  - `DefaultPoolSize`: IPs per tenant in static mode (default: 8)
  - `InitialPoolSize`: starting IPs in elastic mode (default: 2)
  - `GrowthIncrement`: IPs added per expansion in elastic mode (default: 2)
- `QuotaPerTenant`: `NetworkQuota` with MaxNodeIPs and MaxLoadBalancerIPs

## Cloud Provider Configs

ProviderConfig supports six providers:

| Provider | Config Struct | Key Required Fields |
|----------|--------------|---------------------|
| harvester | HarvesterProviderConfig | Endpoint, Namespace, NetworkName (ns/name format), ImageName, StorageClassName |
| nutanix | NutanixProviderConfig | Endpoint, Port, ClusterUUID, SubnetUUID, ImageUUID, StorageContainerUUID |
| proxmox | ProxmoxProviderConfig | Endpoint, Nodes[], Storage, TemplateID, VMIDRange |
| azure | AzureProviderConfig | SubscriptionID, ResourceGroup, Location, VNetName, SubnetName |
| aws | AWSProviderConfig | Region, VPCID, SubnetIDs[], SecurityGroupIDs[] |
| gcp | GCPProviderConfig | ProjectID, Region, Network, Subnetwork |

Credentials secret keys vary by provider (see `CredentialsRef` comment in providerconfig_types.go).

## NetworkPool and IPAllocation (IPAM)

`NetworkPool` defines a platform-level IP pool:
- `spec.cidr`: network range (e.g., "10.40.0.0/24")
- `spec.reserved[]`: ranges excluded from allocation (ReservedRange with CIDR + description)
- `spec.tenantAllocation`: sub-range + defaults (NodesPerTenant, LBPoolPerTenant)
- `status`: TotalIPs, AllocatedIPs, AvailableIPs, AllocationCount, FragmentationPercent, LargestFreeBlock

`IPAllocation` represents an individual allocation:
- `spec.poolRef`: the NetworkPool to allocate from
- `spec.tenantClusterRef`: the TenantCluster (NamespacedObjectReference)
- `spec.type`: "nodes" or "loadbalancer"
- `spec.count`: number of IPs (uses pool defaults if unset)
- `spec.pinnedRange`: request a specific range (for migration)
- `status`: Phase, CIDR, StartAddress, EndAddress, Addresses[], AllocatedCount, AllocatedBy

## Workspace and WorkspaceTemplate

`Workspace` defines a cloud development environment:
- CRD lives on management cluster (team namespace)
- Pod, PVC, SSH service created in tenant cluster's "workspaces" namespace
- Multi-repo support via `spec.repositories[]`
- Environment copying from existing tenant workloads via `spec.envFrom`
- Dotfiles support via `spec.dotfiles`
- Editor config via `spec.editorConfig` (Neovim config repo/inline init.lua)
- SSH access via `spec.sshPublicKeys` (falls back to User.spec.sshKeys)
- Connect/disconnect lifecycle via `AnnotationConnect` annotation
- Auto-stop via `spec.autoStopAfter`, idle timeout via `spec.idleTimeout`

`WorkspaceTemplate` is data-only (no controller):
- `spec.scope`: "cluster" (visible to all) or "team" (team-only)
- `spec.category`: backend, frontend, data, devops, custom
- `spec.template`: WorkspaceTemplateBody (image, repos, envFrom, dotfiles, resources, storageSize)

## ControlPlaneExposure

Defined in clusterbootstrap_types.go, used in both ClusterBootstrap and ButlerConfig:

```go
type ControlPlaneExposureMode string  // LoadBalancer | Ingress | Gateway

type ControlPlaneExposureSpec struct {
    Mode             ControlPlaneExposureMode  // default: LoadBalancer
    Hostname         string                     // required for Ingress/Gateway ("*.k8s.example.com")
    IngressClassName string                     // for Ingress mode
    ControllerType   string                     // haproxy | nginx | traefik | generic
    GatewayRef       string                     // for Gateway mode ("namespace/name")
}
```

CEL validation on ClusterBootstrap enforces hostname is required for Ingress/Gateway modes, and gatewayRef is required for Gateway mode.

## ManagementMode (TenantCluster Addon Management)

```go
type ManagementMode string  // Active | Observe | GitOps
```

- **Active**: Butler actively manages addons. New addons in spec are installed.
- **Observe**: Butler only observes after initial install. Spec changes ignored.
- **GitOps**: Butler bootstraps Flux and hands off to GitOps.

## Helper Methods

Several types have helper methods defined in their files:

- **ButlerConfig**: `IsGitProviderConfigured()`, `GetGitProviderURL()`, `GetControlPlaneExposureMode()`, `IsTCPProxyRequired()`, `GetControlPlaneExposureHostname()`, `GetControlPlaneExposureGatewayRef()`, `GetControlPlaneExposureIngressClassName()`, `GetControlPlaneExposureControllerType()`
- **ClusterBootstrap**: `IsReady()`, `IsFailed()`, `IsSingleNode()`, `GetExpectedMachineCount()`, `GetControlPlaneReplicas()`, `GetControlPlaneIPs()`, `GetWorkerIPs()`, `AllMachinesRunning()`, `GetLoadBalancerAddressPool()`, `GetStorageReplicaCount()`, `GetControlPlaneExposureMode()`, `IsTCPProxyRequired()`
- **ClusterBootstrapAddonsSpec**: `IsCAPIEnabled()`, `GetCAPIVersion()`, `IsButlerControllerEnabled()`, `GetButlerControllerImage()`, `IsConsoleEnabled()`, `GetConsoleVersion()`, `GetConsoleIngressHost()`
- **LoadBalancerPoolSpec**: `Validate()`, `ContainsIP()`, `ToAddressRange()`
- **MachineRequest**: `IsReady()`, `IsFailed()`, `IsTerminating()`, `SetPhase()`, `SetFailure()`
- **AddonDefinition**: `GetNamespace()`, `GetReleaseName()`, `IsBuiltIn()`
- **User**: `IsSSO()`, `IsInternal()`, `IsActive()`, `IsDisabled()`, `IsPlatformAdmin()`
- **IdentityProvider**: `GetScopes()`, `GetGroupsClaim()`, `GetEmailClaim()`, `GetDisplayName()`, `IsReady()`
- **Workspace**: `IsRunning()`, `IsStopped()`, `IsConnected()`

## Make Targets

| Target | Description |
|--------|-------------|
| `make generate` | Run controller-gen to generate DeepCopy methods (`zz_generated.deepcopy.go`) |
| `make manifests` | Generate CRD YAML into `config/crd/bases/` |
| `make fmt` | Run `go fmt` |
| `make vet` | Run `go vet` |
| `make lint` | Run golangci-lint (v2.5.0) |
| `make test` | Run unit tests (runs generate + manifests first) |
| `make build` | Build binary (runs generate + manifests first) |
| `make install` | Apply CRDs to cluster via kustomize |
| `make uninstall` | Remove CRDs from cluster |

Tool versions:
- controller-gen: v0.19.0
- kustomize: v5.7.1
- golangci-lint: v2.5.0

## Change Propagation Flow

1. Edit type files in `api/v1alpha1/`
2. `make generate` -- regenerates `zz_generated.deepcopy.go`
3. `make manifests` -- regenerates CRD YAML in `config/crd/bases/`
4. Commit and push to main
5. GitHub Action `sync-crds-to-charts.yaml` auto-creates PR in `butler-charts` repo
6. Renovate auto-creates PRs in dependent repos:
   - butler-controller
   - butler-bootstrap
   - butler-server
   - butler-cli
   - butler-provider-harvester
   - butler-provider-nutanix

## Conventions

### Kubebuilder Markers

Every CRD type has these root markers:
```go
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=xx,scope=Cluster  // or scope=Namespaced (default)
// +kubebuilder:printcolumn:name="...",type="...",JSONPath="..."
```

Validation markers:
```go
// +kubebuilder:validation:Required
// +kubebuilder:validation:Enum=value1;value2
// +kubebuilder:validation:Minimum=0
// +kubebuilder:validation:Maximum=100
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=63
// +kubebuilder:validation:Pattern=`^regex$`
// +kubebuilder:validation:MinItems=1
// +kubebuilder:default="value"
// +optional
```

CEL validation (used on ClusterBootstrap):
```go
// +kubebuilder:validation:XValidation:rule="...",message="..."
```

### resource.Quantity Usage

CPU, memory, and storage fields use `resource.Quantity` from `k8s.io/apimachinery/pkg/api/resource`:
- `TeamResourceLimits`: MaxCPUCores, MaxMemory, MaxStorage, DefaultCPUPerNode, DefaultMemoryPerNode
- `TeamResourceUsage`: TotalCPU, TotalMemory, TotalStorage
- `ResourceLimits` (ButlerConfig): MaxTotalCPU, MaxTotalMemory, MaxTotalStorage
- `MachineTemplateSpec` (TenantCluster workers): Memory, DiskSize
- `WorkspaceSpec`: StorageSize
- `WorkspaceTemplateBody`: StorageSize

MachineRequest uses plain int32 fields (MemoryMB, DiskGB, CPU) for the provider interface contract.

### Conditions Pattern

All status types use `[]metav1.Condition` with list map key `type`:
```go
// +optional
// +listType=map
// +listMapKey=type
Conditions []metav1.Condition `json:"conditions,omitempty"`
```

### JSON Serialization

- Required fields omit `omitempty`: `json:"name"`
- Optional fields use `omitempty`: `json:"name,omitempty"`
- Inline embedding: `json:",inline"`
- Raw extension: `+kubebuilder:pruning:PreserveUnknownFields`

### API Evolution Rules

1. Additive changes only in v1alpha1 during development
2. Never remove fields -- deprecate first (see `WorkspaceSpec.Repository` marked deprecated)
3. New required fields must have defaults
4. Breaking changes require a new API version

### ExtensionValues Pattern

Arbitrary Helm values use the `ExtensionValues` type:
```go
// +kubebuilder:pruning:PreserveUnknownFields
type ExtensionValues struct {
    Raw []byte `json:"-"`
}
```

ManagementAddon uses `*runtime.RawExtension` instead for the same purpose.

### init() Registration

Every type file has an `init()` function registering its types with the SchemeBuilder:
```go
func init() {
    SchemeBuilder.Register(&TypeName{}, &TypeNameList{})
}
```

## Important Notes

- ButlerConfig is a singleton -- only one named "butler" should exist
- GitProvider types are NOT a standalone CRD; they are embedded in ButlerConfig
- WorkspaceTemplate has NO status subresource (data-only, no controller)
- AddonDefinition has NO status subresource (catalog entry, no reconciliation)
- The `AnnotationConnect`/`AnnotationConnectTime` annotations drive the workspace SSH lifecycle
- ProviderConfig supports both on-prem (harvester, nutanix, proxmox) and cloud (aws, azure, gcp) providers
- Network IPAM is opt-in via `ProviderConfig.spec.network.mode=ipam`
- Elastic LB allocation (`AllocationMode=elastic`) is configured per-provider, not per-tenant
