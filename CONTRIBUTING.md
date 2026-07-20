# Contributing to Facter

Thank you for your interest in contributing. This document applies to all repositories in the `klamhq` organisation that are part of the Facter platform.

---

## Repositories

| Repository                                                         | Description                                         |
| ------------------------------------------------------------------ | --------------------------------------------------- |
| [facter-oss](https://github.com/klamhq/facter-oss)                 | Go agent — system inventory collection              |
| [facter-grpc](https://github.com/klamhq/facter-grpc)               | Go gRPC server — inventory ingest and graph storage |
| [facter-api](https://github.com/klamhq/facter-api)                 | Go GraphQL API — query interface                    |
| [facter-rule-engine](https://github.com/klamhq/facter-rule-engine) | Go compliance rule engine                           |
| [facter-ui](https://github.com/klamhq/facter-ui)                   | Vue 3 / Quasar web dashboard                        |
| [facter-schema](https://github.com/klamhq/facter-schema)           | Shared protobuf schema                              |

---

## Commit convention

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>

[optional body]
[optional footer]
```

### Types

| Type       | When to use                        |
| ---------- | ---------------------------------- |
| `feat`     | New feature                        |
| `fix`      | Bug fix                            |
| `chore`    | Tooling, dependencies, config      |
| `docs`     | Documentation only                 |
| `test`     | Adding or fixing tests             |
| `refactor` | Code change without feature or fix |
| `perf`     | Performance improvement            |
| `ci`       | CI/CD pipeline changes             |

### Version bump rules

| Commit type                            | Version bump  |
| -------------------------------------- | ------------- |
| `fix`                                  | patch (1.0.x) |
| `feat`                                 | minor (1.x.0) |
| `feat!` or `BREAKING CHANGE` in footer | major (x.0.0) |

### Examples

```
feat(collector): add systemd drop-in unit support
fix(grpc): handle empty hostname in MERGE query
docs(readme): update prerequisites section
chore(deps): bump google.golang.org/grpc to v1.83
```

---

## Branch naming

```
feat/<short-description>
fix/<issue-or-description>
chore/<description>
docs/<description>
```

Examples:
- `feat/add-docker-network-collector`
- `fix/grpc-mtls-cert-reload`
- `docs/update-writing-rules-guide`

---

## Workflow

1. Fork the repository (or create a branch if you have access)
2. Create a branch from `main` following the naming convention
3. Make your changes
4. Ensure tests pass: `make test`
5. Open a pull request targeting `main`
6. Fill in the pull request template
7. Wait for review

---

## Adding a collector to facter-oss

Each collector lives in `pkg/inventory/<collector-name>/`. To add a new one:

1. Create a new package under `pkg/inventory/`
2. Implement the `Collector` interface
3. Add the collector to the config struct in `pkg/options/`
4. Register it in `internal/inventory/runner.go`
5. Add protobuf fields in `facter-schema` if new data types are needed
6. Update `RELATIONS.md` in `facter-grpc` after adding the corresponding graph writer

---

## Writing a new audit rule

See [Writing Audit Rules](https://klamhq.github.io/facter-oss/guides/writing-rules/) in the documentation.

Rules live in `facter-rule-engine/rules/<category>/`. IDs must follow the format `R-[A-Z]{3}-[0-9]{4}`.

System rules (prefix `R-SYS-`) are reserved. Use a different prefix for custom rules.

---

## Code style

- **Go**: run `golangci-lint run` before opening a PR (`make golint`)
- **Vue/JS**: run `npm run lint` and `npm run format`
- **Proto**: run `make lint` and `make format` in `facter-schema`

---

## Updating the graph schema documentation

When adding a new collector that creates new node types or relationships:

1. Run `make export-relations` in `facter-grpc` to regenerate `RELATIONS.md`
2. Update `facter-oss/docs/concepts/data-model.md` to reflect the new node types and properties
3. Commit both changes together

---

## Security vulnerabilities

Do **not** open a public issue for security vulnerabilities. Contact the maintainers directly.
