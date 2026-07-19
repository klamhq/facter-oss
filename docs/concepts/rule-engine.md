# Rule Engine

## Overview

The `facter-rule-engine` evaluates compliance and security audit rules against the infrastructure graph. Each rule is a YAML document that defines:

1. **A Cypher query** to detect a condition in the graph
2. **Metadata** (severity, category, description, references)

Rules are stored in PostgreSQL and executed via the `facter-grpc` `CheckRules` gRPC endpoint.

---

## Rule Structure

```yaml
id: R-SYS-1001                           # Unique ID — format: R-XXX-NNNN
name: "Shared SSH Keys Detection"
description: "Detects SSH keys deployed on multiple hosts"
category: IDENTITY_ACCESS                # See categories below
severity: HIGH                           # CRITICAL | HIGH | MEDIUM | LOW | INFO
status: approved                         # draft | review | approved | deprecated
enabled: true
priority: P1                             # P1 (highest) to P4

# The Cypher query executed against Apache AGE
cypher_query: |
  MATCH (h:Host)-[:HAS_SSH_KEY]->(k:SshKey)
  WITH k, collect(DISTINCT h.hostname) AS hosts, count(DISTINCT h) AS cnt
  WHERE cnt > 1
  RETURN k.fingerprint AS fingerprint, hosts, cnt
```

### Field Reference

| Field          | Required | Type   | Description                                           |
| -------------- | -------- | ------ | ----------------------------------------------------- |
| `id`           | ✓        | string | Unique rule ID. Format: `R-[A-Z]{3}-[0-9]{4}`         |
| `name`         | ✓        | string | Short human-readable name                             |
| `description`  | ✓        | string | Detailed description of what the rule checks          |
| `category`     | ✓        | string | Rule category (see below)                             |
| `severity`     | ✓        | string | `CRITICAL`, `HIGH`, `MEDIUM`, `LOW`, `INFO`           |
| `cypher_query` | ✓        | string | The Cypher query. Must contain `MATCH` and `RETURN`   |
| `status`       | —        | string | `draft` (default), `review`, `approved`, `deprecated` |
| `enabled`      | —        | bool   | `true` by default                                     |
| `priority`     | —        | string | `P1` to `P4`                                          |
| `version`      | —        | int    | Schema version, starts at 1                           |
| `created_by`   | —        | string | Author                                                |

---

## Categories

| Value              | Description                            |
| ------------------ | -------------------------------------- |
| `IDENTITY_ACCESS`  | Users, SSH keys, privileges            |
| `NETWORK_SECURITY` | Connections, firewall, exposed ports   |
| `DATA_PROTECTION`  | Encryption, storage, secrets           |
| `COMPLIANCE`       | Regulatory / benchmark compliance      |
| `INFRASTRUCTURE`   | Services, packages, platform hardening |

---

## Execution Flow

```
┌──────────────────────────────────────────────────────────┐
│                  Audit Trigger                           │
│  POST /api/v1/audits         (global — all rules)        │
│  POST /api/v1/audits/:id     (single rule)               │
└──────────────────────────┬───────────────────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  Load rules from       │
              │  PostgreSQL            │
              └────────────┬───────────┘
                           │ for each rule
                           ▼
              ┌────────────────────────┐
              │  Send cypher_query via │
              │  CheckRules gRPC RPC   │
              │  → facter-grpc         │
              └────────────┬───────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  facter-grpc executes  │
              │  Cypher on Apache AGE  │
              │  returns raw rows      │
              └────────────┬───────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  Rule engine evaluates │
              │  results:              │
              │  matched = rows > 0    │
              └────────────┬───────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  Persist execution     │
              │  record to PostgreSQL  │
              │  (rule_executions)     │
              └────────────────────────┘
```

### Result Semantics

The result of a rule execution is determined by whether the Cypher query returns any rows:

| `matched` | Meaning                                  | UI Status       |
| --------- | ---------------------------------------- | --------------- |
| `true`    | Query returned rows → condition detected | 🟠 **Finding**   |
| `false`   | Query returned no rows → all clear       | 🟢 **Compliant** |

The **evidence** is the full set of rows returned by the Cypher query, stored as JSON. For each row, scalar values (strings, numbers, booleans) are displayed in a table; complex values (lists, maps) are in a collapsible section.

---

## Rule ID Convention

All built-in system rules use the `R-SYS-XXXX` prefix and can only be read (not deleted or modified) via the API. Custom rules use any other prefix following the format `R-[A-Z]{3}-[0-9]{4}`.

Examples:
- `R-SYS-1001` — built-in: Shared SSH Key Detection
- `R-NET-0001` — custom: Non-standard outbound port detection
- `R-PKG-0042` — custom: Outdated OpenSSL check

---

## PostgreSQL Schema

### `rules` table

| Column         | Type      | Description                                      |
| -------------- | --------- | ------------------------------------------------ |
| `id`           | varchar   | Rule ID (PK)                                     |
| `name`         | varchar   | Rule name                                        |
| `description`  | text      | Description                                      |
| `status`       | varchar   | Lifecycle status                                 |
| `category`     | varchar   | Category                                         |
| `severity`     | varchar   | Severity level                                   |
| `cypher_query` | text      | The Cypher query                                 |
| `version`      | int       | Rule version                                     |
| `priority`     | varchar   | Priority                                         |
| `created_by`   | varchar   | Author                                           |
| `enabled`      | bool      | Whether active                                   |
| `created_at`   | timestamp | Creation time                                    |
| `updated_at`   | timestamp | Last update time                                 |
| `checksum`     | varchar   | SHA-256 of rule content (prevents no-op updates) |

### `rule_executions` table

| Column           | Type    | Description                         |
| ---------------- | ------- | ----------------------------------- |
| `id`             | uuid    | Execution ID (PK)                   |
| `rule_id`        | varchar | FK → rules.id                       |
| `executed_by`    | varchar | Component that triggered the run    |
| `execution_type` | varchar | `simulation`, `production`          |
| `cypher_query`   | text    | Snapshot of query at execution time |
| `matched`        | bool    | Whether findings were detected      |
| `evidence`       | jsonb   | Raw Cypher result rows              |
| `duration_ms`    | bigint  | Query execution duration            |
| `created_at`     | bigint  | Unix timestamp of execution         |
