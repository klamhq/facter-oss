# Deployments

Facter can be deployed in several ways depending on your environment.

## Deployment Options

| Method                       | Effort    | Best For                 |
| ---------------------------- | --------- | ------------------------ |
| [Quick Start](quickstart.md) | ⭐ Low     | Local testing, demo      |
| [Docker Compose](docker.md)  | ⭐⭐ Medium | Development, small teams |
| [Kubernetes](kubernetes.md)  | ⭐⭐⭐ High  | Production at scale      |

## Architecture Requirements

Regardless of deployment method, you need:

- **Apache AGE** (PostgreSQL ≥ 14 with AGE extension) — stores the inventory graph
- **PostgreSQL** — stores rule definitions and execution results for `facter-rule-engine`
- **Keycloak** — identity provider for authentication
- **TLS certificates** — for mTLS between `facter-oss` and `facter-grpc` (and optionally `facter-rule-engine` → `facter-grpc`)

## Component Ports

| Component            | Default Port             | Protocol   |
| -------------------- | ------------------------ | ---------- |
| `facter-grpc`        | 56230                    | gRPC/TLS   |
| `facter-api`         | 56231                    | HTTP       |
| `facter-rule-engine` | 8081                     | HTTP       |
| `facter-ui`          | 80 (Docker) / 9000 (dev) | HTTP       |
| Apache AGE           | 5433                     | PostgreSQL |
| PostgreSQL           | 5432                     | PostgreSQL |
| Keycloak             | 8080                     | HTTP       |
