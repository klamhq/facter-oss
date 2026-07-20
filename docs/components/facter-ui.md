# facter-ui

## Overview

`facter-ui` is the web interface for the Facter platform, built with [Vue.js 3](https://vuejs.org/) and the [Quasar Framework](https://quasar.dev/). It provides dashboards, per-host drill-downs, compliance reports and audit rule management.

---

## Views

### Main Dashboard (`/dashboard`)

Global KPI cards showing:

- Connected hosts count
- Packages across all hosts
- Vulnerabilities summary

### Hosts (`/hosts`)

A searchable, filterable table of all inventoried hosts with quick navigation to each host's detail page.

### Per-Host View (`/hosts/:hostname`)

Deep-dive into a single host via tabbed navigation:

| Tab                 | Content                                 |
| ------------------- | --------------------------------------- |
| **Overview**        | OS, kernel, hardware summary, uptime    |
| **Compliance**      | OpenSCAP benchmark results              |
| **Applications**    | Docker containers and images            |
| **Packages**        | All installed packages with versions    |
| **Services**        | Systemd units and their states          |
| **Hardware**        | CPU, memory, disk, network interfaces   |
| **Vulnerabilities** | CVEs matched against installed packages |
| **Connections**     | Established network connections         |
| **Processes**       | Running processes                       |
| **Access**          | Local users, SSH keys, known hosts      |
| **Graph**           | Interactive network topology graph      |

### Networks (`/networks`)

Network topology view and cross-host connection analysis.

### Packages (`/packages`)

Cross-host package dashboard with outdated package detection and CVE correlation.

### Vulnerabilities (`/vulnerabilities`)

Global vulnerability dashboard with severity charts and CVE details.

### Systemd (`/systemd`)

Fleet-wide systemd service status overview.

### Audit Rule Engine (`/rule-engine`)

- List, create, edit and delete audit rules (admin only)
- Launch per-rule or global audits (admin only)
- Filter rules by category, severity, status

### Global Audit Report (`/rule-engine/reports`)

Compliance overview across all rules:

- KPI cards: total rules, findings count, compliance rate, last audit time
- Donut chart: compliant vs findings vs no data
- Bar chart: findings by severity
- Bar chart: findings by category
- Summary table: per-rule status with link to detailed execution history

### Rule Execution Reports (`/rule_engine/reports/:rule_id`)

Detailed execution history for a single rule:

- Status per execution (Finding / Compliant)
- Evidence: scalar properties in a table, nested structures in collapsible sections
- Duration and execution metadata

---

## Authentication

`facter-ui` uses Keycloak OIDC (**Authorization Code Flow with PKCE**).

On visit, unauthenticated users are redirected to the Keycloak login page. After login, the access token is stored in memory (not localStorage) and attached to every API request.

Configure Keycloak in `src/boot/keycloak.js` or via environment variables at build time.

The following roles are enforced:
- `facter-viewer`: read-only access to all views
- `facter-admin`: full access including write operations and audit launches

---

## API Clients

`facter-ui` uses two API clients:

| Client             | File                               | Target                    |
| ------------------ | ---------------------------------- | ------------------------- |
| `gqlClient`        | `src/graphql/client.js`            | `facter-api` GraphQL      |
| `ruleEngineClient` | `src/services/ruleEngineClient.js` | `facter-rule-engine` REST |

Both clients automatically attach the Keycloak access token as a Bearer token.

---

## Technology Stack

| Library                | Version | Use                                            |
| ---------------------- | ------- | ---------------------------------------------- |
| Vue.js                 | 3       | Reactive UI framework                          |
| Quasar                 | 2       | Component library (q-table, q-chip, q-card, …) |
| Vue Router             | 4       | Client-side routing                            |
| Axios                  | —       | HTTP requests (via ruleEngineClient)           |
| graphql-request        | —       | GraphQL client                                 |
| chart.js + vue-chartjs | 4 / 5   | Charts (bar, doughnut)                         |
| d3-force               | 3       | Network graph force simulation                 |
| Tailwind CSS           | —       | Utility classes                                |
| Keycloak-js            | —       | OIDC authentication                            |
| Vitest                 | —       | Unit testing                                   |

---

## Development

```bash
cd facter-ui
npm install
npm run dev          # Start dev server on http://localhost:9000
npm run build        # Production build to dist/
npm test             # Run unit tests
```

### Environment

The dev server proxies API requests based on the configuration in `quasar.config.js`. Update the proxy targets to point to your local backend services.

---

## Build & Deploy

### Docker

```bash
docker build -t facter-ui:latest -f docker/Dockerfile .
docker run -p 80:80 facter-ui:latest
```

The Docker image uses nginx to serve the static build. The nginx configuration is in `docker/nginx.conf`.
