# Butler API

<p align="left">
  <a href="https://github.com/butlerdotdev/butler-api/blob/main/LICENSE"><img src="https://img.shields.io/github/license/butlerdotdev/butler-api" alt="License"></a>
  <img src="https://img.shields.io/github/go-mod/go-version/butlerdotdev/butler-api" alt="Go Version">
  <a href="https://goreportcard.com/report/github.com/butlerdotdev/butler-api"><img src="https://goreportcard.com/badge/github.com/butlerdotdev/butler-api" alt="Go Report Card"></a>
  <a href="https://pkg.go.dev/github.com/butlerdotdev/butler-api"><img src="https://pkg.go.dev/badge/github.com/butlerdotdev/butler-api.svg" alt="Go Reference"></a>
</p>

API type definitions for the [Butler](https://github.com/butlerdotdev/butler) Kubernetes-as-a-Service platform. This repository defines all Butler Custom Resource Definitions (CRDs) and serves as the API contract for the entire platform.

## Overview

- **Go module**: `github.com/butlerdotdev/butler-api`
- **API group**: `butler.butlerlabs.dev`
- **API version**: `v1alpha1`
- **Types path**: `api/v1alpha1/`

Every Butler component imports this module to ensure type consistency across the platform.

## CRD Types

| CRD | Scope | Short Name | Purpose |
|-----|-------|------------|---------|
| `ButlerConfig` | Cluster | `bc` | Platform configuration singleton. Multi-tenancy mode, default provider, team limits, Git provider config. |
| `Team` | Cluster | `tm` | Multi-tenancy unit. Owns a namespace, has users/groups with roles (admin/operator/viewer), resource quotas. |
| `User` | Cluster | — | Internal user accounts. Email, password hash in Secret, invite flow, SSO metadata. |
| `IdentityProvider` | Cluster | `idp` | OIDC/SSO configuration. Issuer URL, client ID/secret, scopes, group/email claims. |
| `ProviderConfig` | Namespaced | `pc` | Infrastructure provider credentials + defaults (Harvester/Nutanix/Proxmox). |
| `TenantCluster` | Namespaced | `tc` | Main resource. Declares a complete tenant K8s cluster with control plane, workers, networking, addons. |
| `TenantAddon` | Namespaced | — | Per-cluster addon lifecycle (install/upgrade/remove Helm charts on tenant). |
| `AddonDefinition` | Cluster | — | Addon catalog entry. Describes available addons with versions, categories, Helm repo info. |
| `ManagementAddon` | Cluster | — | Addon installed on the management cluster itself (Steward, cert-manager, MetalLB, etc.). |
| `MachineRequest` | Namespaced | — | VM lifecycle request dispatched to provider controllers. |
| `ClusterBootstrap` | Namespaced | — | Talos-based bootstrap workflow for management cluster initial setup. |

## Installation

### As a Go Dependency

```go
import butlerv1alpha1 "github.com/butlerdotdev/butler-api/api/v1alpha1"
```

Add to your `go.mod`:

```bash
go get github.com/butlerdotdev/butler-api@latest
```

### Installing CRDs to a Cluster

```bash
make install
```

Or use the [butler-crds Helm chart](https://github.com/butlerdotdev/butler-charts/tree/main/charts/butler-crds):

```bash
helm install butler-crds oci://ghcr.io/butlerdotdev/charts/butler-crds
```

## Usage

### Registering Types with a Scheme

```go
import (
    "k8s.io/apimachinery/pkg/runtime"
    butlerv1alpha1 "github.com/butlerdotdev/butler-api/api/v1alpha1"
)

func init() {
    scheme := runtime.NewScheme()
    _ = butlerv1alpha1.AddToScheme(scheme)
}
```

### Creating a TenantCluster

```go
cluster := &butlerv1alpha1.TenantCluster{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-cluster",
        Namespace: "team-dev",
    },
    Spec: butlerv1alpha1.TenantClusterSpec{
        TeamRef: butlerv1alpha1.LocalObjectReference{Name: "dev"},
        ProviderConfigRef: butlerv1alpha1.ProviderReference{
            Name: "harvester-prod",
        },
        ControlPlane: butlerv1alpha1.ControlPlaneSpec{
            Replicas: 2,
        },
        Workers: butlerv1alpha1.WorkersSpec{
            Replicas: 3,
        },
    },
}
```

## Development

### Prerequisites

- Go 1.24+
- make
- controller-gen (installed automatically via make)

### Building

```bash
make build
```

### Generating Code

```bash
# Generate DeepCopy methods
make generate

# Generate CRD YAML manifests
make manifests
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

## Adding New CRDs

1. Create `{typename}_types.go` in `api/v1alpha1/`
2. Define structs: `{Type}Spec`, `{Type}Status`, `{Type}` (root), `{Type}List`
3. Add kubebuilder markers:
   ```go
   //+kubebuilder:object:root=true
   //+kubebuilder:subresource:status
   //+kubebuilder:resource:scope=Namespaced,shortName=xyz
   //+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
   ```
4. Register in `init()`:
   ```go
   func init() {
       SchemeBuilder.Register(&{Type}{}, &{Type}List{})
   }
   ```
5. Run `make generate && make manifests`
6. Sync CRD YAML to butler-charts via `hack/sync-crds.sh`

## Adding Fields to Existing Types

1. Add field to spec/status struct with json tag and kubebuilder validation markers:
   ```go
   // +kubebuilder:validation:Minimum=1
   // +kubebuilder:default=3
   Replicas int32 `json:"replicas,omitempty"`
   ```
2. Run `make generate && make manifests`
3. Sync CRD YAML to butler-charts

## Project Structure

```
butler-api/
├── api/
│   └── v1alpha1/
│       ├── butlerconfig_types.go      # ButlerConfig CRD
│       ├── team_types.go              # Team CRD
│       ├── user_types.go              # User CRD
│       ├── identityprovider_types.go  # IdentityProvider CRD
│       ├── providerconfig_types.go    # ProviderConfig CRD
│       ├── tenantcluster_types.go     # TenantCluster CRD
│       ├── tenantaddon_types.go       # TenantAddon CRD
│       ├── addondefinition_types.go   # AddonDefinition CRD
│       ├── managementaddon_types.go   # ManagementAddon CRD
│       ├── machinerequest_types.go    # MachineRequest CRD
│       ├── clusterbootstrap_types.go  # ClusterBootstrap CRD
│       ├── common_types.go            # Shared types, labels, conditions
│       ├── gitprovider_types.go       # Git provider configuration
│       ├── groupversion_info.go       # API group registration
│       └── zz_generated.deepcopy.go   # Auto-generated (do not edit)
├── config/
│   └── crd/
│       └── bases/                     # Generated CRD YAML
├── Makefile
└── go.mod
```

## Labels and Annotations

Standard labels used across Butler-managed resources:

```yaml
labels:
  app.kubernetes.io/managed-by: butler
  butler.butlerlabs.dev/team: <team-name>
  butler.butlerlabs.dev/tenant: <cluster-name>
```

## Documentation

- [Butler Documentation](https://docs.butlerlabs.dev/butler/)
- [CRD Reference](https://docs.butlerlabs.dev/butler/reference/crds/)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
