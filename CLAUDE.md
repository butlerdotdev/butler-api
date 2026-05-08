# Butler API — Claude Code Notes

## Tooling

- controller-gen pinned at v0.19.0 (Makefile `CONTROLLER_TOOLS_VERSION`). Downstream repos that regenerate CRDs from butler-api types must match this version. Drift produces subtly different schema markers and breaks schema consistency.

## Conventions

- Data-only API types without methods do not get unit test files. Tests exist only for types with behavioral methods (e.g., `GetEffectiveTier`, `GetEffectivePlatformRole`) or serialization edge cases (e.g., `IconData` omitempty). See `addondefinition_types_test.go` and `user_types_test.go` for the pattern.
