# Contributing to Butler API

Thank you for your interest in contributing to Butler API! This document provides guidelines for contributing to the API type definitions.

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you are expected to uphold this code.

## Developer Certificate of Origin

By contributing to this project, you agree to the Developer Certificate of Origin (DCO). Every commit must be signed off:

```bash
git commit -s -m "Your commit message"
```

## Getting Started

### Prerequisites

- Go 1.24+
- make
- kubectl (for testing CRD installation)

### Setting Up Your Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/butler-api.git
   cd butler-api
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/butlerdotdev/butler-api.git
   ```

## Making Changes

### Adding a New CRD Type

1. Create `api/v1alpha1/{typename}_types.go`
2. Follow the existing pattern:
   - Define `{Type}Spec` struct with fields and validation markers
   - Define `{Type}Status` struct
   - Define `{Type}` root type with kubebuilder markers
   - Define `{Type}List` type
   - Register in `init()` function
3. Run code generation:
   ```bash
   make generate
   make manifests
   ```
4. Update this documentation if the CRD is user-facing

### Modifying Existing Types

1. Add/modify fields with appropriate json tags
2. Add kubebuilder validation markers as needed
3. Run `make generate && make manifests`
4. Consider backward compatibility

### Code Style

- Follow standard Go formatting (`gofmt`)
- Use kubebuilder validation markers for field validation
- Add doc comments for exported types and fields
- Include Apache 2.0 license headers on all source files

### Kubebuilder Markers

Common markers used in this project:

```go
// Validation
// +kubebuilder:validation:Required
// +kubebuilder:validation:Enum=option1;option2
// +kubebuilder:validation:Minimum=1
// +kubebuilder:validation:Pattern=`^[a-z]+$`

// Defaults
// +kubebuilder:default=defaultValue

// Resource metadata
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=xyz
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
```

### Commit Messages

Follow conventional commits:

```
type(scope): description

[optional body]

Signed-off-by: Your Name <your.email@example.com>
```

Types:
- `feat`: New CRD type or field
- `fix`: Bug fix in type definitions
- `docs`: Documentation changes
- `chore`: Maintenance tasks
- `refactor`: Code refactoring

### Pull Request Process

1. Create a feature branch:
   ```bash
   git checkout -b feat/your-feature
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -s -m "feat(api): add new field to TenantCluster"
   ```

3. Ensure code generation is up to date:
   ```bash
   make generate
   make manifests
   ```

4. Push to your fork:
   ```bash
   git push origin feat/your-feature
   ```

5. Open a Pull Request against `main`

6. Ensure all checks pass

## Important Guidelines

### Do NOT

- Add controller logic — that belongs in butler-controller
- Add server/API logic — that belongs in butler-server
- Add CLI logic — that belongs in butler-cli
- Break backward compatibility without a migration plan
- Use `interface{}` or `map[string]interface{}` — use typed structs

### Do

- Keep this module lightweight (minimal dependencies)
- Use kubebuilder markers for validation
- Follow Kubernetes API conventions for status conditions
- Document all exported types and fields

## Testing

```bash
# Run tests
make test

# Verify CRD generation
make manifests
kubectl apply -f config/crd/bases/ --dry-run=client
```

## Getting Help

- Open an [issue](https://github.com/butlerdotdev/butler-api/issues) for bugs or feature requests
- Check existing issues before creating new ones

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
