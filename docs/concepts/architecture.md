# Architecture

## Overview

The Facter platform is composed of six independent services that communicate over defined protocols. The diagram below shows the full architecture.

> **Diagram**: [architecture-overview.drawio](../assets/diagrams/architecture-overview.drawio) — open with [draw.io](https://app.diagrams.net/) or the VS Code draw.io extension.

---

## Components

### facter-oss (Agent)

The agent runs directly on each Linux host to be monitored. It collects system facts and sends them to `facter-grpc`.

**Key characteristics:**

- Written in Go, compiles to a single static binary (~10 MB)
- Runs as a one-shot command or in a scheduled job (cron, systemd timer)
- Two output modes: `local` (writes `.iya` file) or `remote` (streams via gRPC)
- Supports **full** and **delta** inventory modes
- Requires root or specific capabilities for some collectors (firewall, processes)

**Collectors:**

| Collector       | Gathers                                                                         |
| --------------- | ------------------------------------------------------------------------------- |
| Platform        | OS, kernel, hardware (CPU, RAM, disk), virtualisation, init system              |
| Packages        | APT, RPM, Pacman, Homebrew packages                                             |
| Users           | Local users, groups, sudo capabilities, active sessions                         |
| Networks        | Interfaces, IPs, established connections, firewall rules, DNS, public IP, GeoIP |
| SSH             | SSH keys (RSA, DSA, ECDSA, Ed25519), authorized_keys, known_hosts               |
| Processes       | Running processes with PID, owner, command                                      |
| Applications    | Docker containers, images, networks                                             |
| Systemd         | Unit states (active, loaded, enabled), dependencies                             |
| Compliance      | OpenSCAP XCCDF benchmark results                                                |
| Vulnerabilities | CVE matching against installed packages                                         |

---

### facter-grpc

The central ingest server. Receives `HostInventory` or `HostDeltaInventory` Protobuf messages and persists them as a property graph in Apache AGE.

**Key characteristics:**

- gRPC server (port `56230`) with **mTLS** (SPIFFE-based)
- Principal allowlist: only SPIFFE identities in `facterPrincipal` or `facterRuleEnginePrincipal` are accepted
- Handles both full and delta inventory processing
- Scheduled tasks:
    - **Snapshot** (daily by default): point-in-time graph copies for trend analysis
    - **Hard deletion** (configurable): purges old nodes and snapshots
- Also exposes `CheckRules` RPC used by `facter-rule-engine` to run Cypher audit queries

---

### facter-api

The read-only query gateway for the UI.

**Key characteristics:**

- HTTP server (port `56231`) exposing `GET /graphql` and `POST /graphql`
- Connects directly to Apache AGE with a read-only connection pool
- Authentication via Keycloak JWT (roles: `facter-admin`, `facter-viewer`)
- CORS configurable per environment

---

### facter-rule-engine

The compliance and audit engine.

**Key characteristics:**

- HTTP REST API (port `8081`) with full CRUD over rules and audit execution
- Stores rules and execution results in PostgreSQL
- Sends `CheckRules` gRPC requests to `facter-grpc` to execute Cypher queries
- Compliance determined directly from Cypher results: rows returned = finding, no rows = compliant
- Authentication via Keycloak JWT (roles: `facter-admin`, `facter-viewer`)

**REST endpoints:**

| Method   | Path                               | Role         | Description                            |
| -------- | ---------------------------------- | ------------ | -------------------------------------- |
| `GET`    | `/api/v1/rules`                    | viewer/admin | List all enabled rules                 |
| `POST`   | `/api/v1/rules`                    | admin        | Create one or more rules               |
| `PUT`    | `/api/v1/rules/:id`                | admin        | Update a rule                          |
| `DELETE` | `/api/v1/rules/:id`                | admin        | Delete a rule                          |
| `GET`    | `/api/v1/rule_executions`          | viewer/admin | All execution records                  |
| `GET`    | `/api/v1/rule_executions/:rule_id` | viewer/admin | Executions for one rule                |
| `POST`   | `/api/v1/audits`                   | admin        | Run a global audit (all enabled rules) |
| `POST`   | `/api/v1/audits/:rule_id`          | admin        | Run audit for one rule                 |

---

### facter-ui

The web dashboard built with Vue.js 3 and Quasar Framework.

**Views:**

| Route                           | View                                                                                                                                     |
| ------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| `/dashboard`                    | Global metrics: host count, packages, vulnerabilities                                                                                    |
| `/hosts`                        | List of all inventoried hosts                                                                                                            |
| `/hosts/:hostname`              | Per-host drill-down (overview, packages, users, SSH, Docker, hardware, services, processes, vulnerabilities, connections, network graph) |
| `/networks`                     | Network topology and connections                                                                                                         |
| `/packages`                     | Packages across all hosts                                                                                                                |
| `/vulnerabilities`              | CVE dashboard                                                                                                                            |
| `/systemd`                      | Systemd services across all hosts                                                                                                        |
| `/rule-engine`                  | Audit rules management                                                                                                                   |
| `/rule-engine/reports`          | Global compliance report with charts                                                                                                     |
| `/rule_engine/reports/:rule_id` | Per-rule execution history                                                                                                               |

**Authentication:** Keycloak OIDC (Authorization Code Flow with PKCE)

---

### facter-schema

A shared Go module containing the Protobuf definitions used by `facter-oss` (client) and `facter-grpc` (server).

**gRPC service:**

```proto
service FactGrpcService {
  rpc Inventory(InventoryRequest) returns (InventoryResponse);
  rpc CheckRules(CheckRulesRequest) returns (CheckRulesResponse);
}
```

- `Inventory`: accepts a full `HostInventory` or a `HostDeltaInventory`
- `CheckRules`: accepts a Cypher query string, returns raw graph results

---

## Communication Matrix

```
┌─────────────────────┬──────────────────────┬─────────────┬────────────┐
│ From                │ To                   │ Protocol    │ Auth       │
├─────────────────────┼──────────────────────┼─────────────┼────────────┤
│ facter-oss          │ facter-grpc          │ gRPC / mTLS │ SPIFFE     │
│ facter-grpc         │ Apache AGE           │ SQL/Cypher  │ DB creds   │
│ facter-api          │ Apache AGE           │ SQL/Cypher  │ DB creds   │
│ facter-rule-engine  │ facter-grpc          │ gRPC / mTLS │ SPIFFE     │
│ facter-rule-engine  │ PostgreSQL           │ SQL         │ DB creds   │
│ facter-ui           │ facter-api           │ HTTPS/WS    │ JWT        │
│ facter-ui           │ facter-rule-engine   │ HTTPS       │ JWT        │
│ facter-ui           │ Keycloak             │ HTTPS/OIDC  │ —          │
│ facter-api          │ Keycloak JWKS        │ HTTPS       │ —          │
│ facter-rule-engine  │ Keycloak JWKS        │ HTTPS       │ —          │
└─────────────────────┴──────────────────────┴─────────────┴────────────┘
```

---

## Port Reference

| Service                  | Port  | Protocol   |
| ------------------------ | ----- | ---------- |
| facter-grpc              | 56230 | gRPC/TLS   |
| facter-api               | 56231 | HTTP       |
| facter-rule-engine       | 8081  | HTTP       |
| facter-ui                | 9000  | HTTP (dev) |
| Apache AGE               | 5433  | PostgreSQL |
| PostgreSQL (rule-engine) | 5432  | PostgreSQL |
| Keycloak                 | 8080  | HTTP       |
