# Writing Audit Rules

This guide explains how to write effective audit rules for the Facter rule engine.

---

## Rule Format

Rules are stored in the database (via the REST API) or can be loaded from YAML files in the `rules/` directory of `facter-rule-engine`.

### YAML File Format

```yaml
id: R-SSH-0001
name: "Root SSH login must be disabled"
description: |
  Detects hosts where the SSH daemon does not explicitly disable root login.
category: IDENTITY_ACCESS
severity: HIGH
status: approved
enabled: true
version: 1
priority: P1
created_by: "security-team"
cypher_query: |
  MATCH (h:Host)-[:HAS_SSH_CONFIG]-(c:SSHConfig)
  WHERE c.permit_root_login IS NULL OR c.permit_root_login <> 'no'
  RETURN h AS seed, h.hostname AS hostname
tests:
  - name: root_login_enabled
    input_graph: root_login_enabled.graph.json
    expect_findings: 1
```

### JSON / REST API Format

```json
{
  "id": "R-SSH-0001",
  "name": "Root SSH login must be disabled",
  "description": "Detects hosts where the SSH daemon does not explicitly disable root login.",
  "category": "IDENTITY_ACCESS",
  "severity": "HIGH",
  "status": "approved",
  "enabled": true,
  "version": 1,
  "priority": "P1",
  "created_by": "security-team",
  "cypher_query": "MATCH (h:Host)-[:HAS_SSH_CONFIG]-(c:SSHConfig) WHERE ... RETURN h AS seed, h.hostname AS hostname"
}
```

---

## Field Reference

| Field          | Type             | Required | Description                                                                            |
| -------------- | ---------------- | -------- | -------------------------------------------------------------------------------------- |
| `id`           | string           | ✅        | Unique rule ID, format `R-[A-Z]{3}-[0-9]{4}`. IDs starting with `R-SYS-` are reserved. |
| `name`         | string           | ✅        | Short, human-readable name (5–100 characters)                                          |
| `description`  | string           | ✅        | Detailed description of what the rule checks (10–500 characters)                       |
| `cypher_query` | string           | ✅        | openCypher query to execute. Must contain `MATCH` and `RETURN` clauses.                |
| `category`     | string           | —        | One of the [allowed categories](#categories)                                           |
| `severity`     | string           | —        | `CRITICAL`, `HIGH`, `MEDIUM`, or `LOW`                                                 |
| `status`       | string           | —        | `draft` (default), `review`, `approved`, `deprecated`                                  |
| `enabled`      | boolean          | —        | Whether the rule runs during audit. Default: `true`                                    |
| `version`      | integer          | —        | Rule version, starts at 1                                                              |
| `priority`     | string           | —        | `P1` (highest) to `P4`                                                                 |
| `created_by`   | string           | —        | Author identifier                                                                      |
| `context`      | object           | —        | Arbitrary key/value metadata (e.g. `require_host_tag`)                                 |
| `tests`        | array of objects | —        | Optional unit tests referencing fixture graph files                                    |

---

## Categories

The allowed category values are:

| Value              | Example checks                         |
| ------------------ | -------------------------------------- |
| `IDENTITY_ACCESS`  | Users, SSH keys, privileges            |
| `NETWORK_SECURITY` | Open ports, firewall rules, interfaces |
| `DATA_PROTECTION`  | Encryption, secrets exposure           |
| `COMPLIANCE`       | OpenSCAP results, benchmark compliance |
| `INFRASTRUCTURE`   | Services, packages, platform hardening |

---

## Compliance Logic

**No `condition` field exists.** Compliance is determined directly from the Cypher query results:

| Query returns | `matched` | Meaning                                         |
| ------------- | --------- | ----------------------------------------------- |
| ≥ 1 rows      | `true`    | 🟠 **Finding** — each row is a piece of evidence |
| 0 rows        | `false`   | 🟢 **Compliant** — nothing was found             |
| `nodata=true` | `false`   | ⚪ **No data** — host not applicable             |

**Design principle**: write your query so that it returns rows only when something is wrong. Each returned row becomes a finding stored in `evidence`.

---

## Writing Cypher Queries

The `cypher_query` field must be a valid openCypher query compatible with Apache AGE.

### Required clauses

All queries must contain `MATCH` and `RETURN`.

### `seed` column convention

Return the main node as `seed` so the UI can link the finding back to a specific object:

```cypher
MATCH (h:Host)-[:HAS_USER]-(u:User)
WHERE u.uid = 0 AND u.name <> 'root'
RETURN h AS seed, h.hostname AS hostname, u.name AS username, u.uid AS uid
```

### No data sentinel

If the query applies only to hosts that have a certain collector enabled, return `nodata: true` in the `WITH … WHERE n = 0` branch to mark the rule as not applicable rather than as compliant:

```cypher
MATCH (h:Host)
OPTIONAL MATCH (h)-[:HAS_SSH_CONFIG]-(c:SSHConfig)
WITH h, count(c) AS config_count
WHERE config_count = 0
RETURN true AS nodata, h.hostname AS hostname
```

### Filter in query, not outside

All filtering should happen inside the Cypher query. There is no post-query evaluation step.

---

## Example Rules

### 1. Detect shared SSH keys

```yaml
id: R-SYS-1001
name: "Shared SSH key across multiple hosts"
description: |
  Detects SSH keys (fingerprint) that are present on more than one host.
category: IDENTITY_ACCESS
severity: HIGH
status: approved
enabled: true
version: 1
priority: P2
created_by: "facter-rule-engine"
cypher_query: |
  MATCH (k:SSHKey)-[:DEPLOYED_ON]->(h:Host)
  WITH k, collect(DISTINCT h.hostname) AS hosts, count(DISTINCT h) AS cnt
  WHERE cnt > 1
  RETURN k AS seed, hosts, cnt, k.fingerprint AS fingerprint
```

### 2. Hosts with outdated packages

```yaml
id: R-SYS-2001
name: "Hosts with outdated packages"
description: |
  Detects hosts that have packages that are outdated.
category: COMPLIANCE
severity: HIGH
status: approved
enabled: true
version: 1
priority: P2
created_by: "facter-rule-engine"
cypher_query: |
  MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package {is_up_to_date: false})
  WITH h, collect(DISTINCT h.hostname) AS hosts, count(DISTINCT h) AS cnt, count(DISTINCT p) AS outdated_package_count
  WHERE cnt > 0
  RETURN h AS seed, hosts, cnt, outdated_package_count, h.hostname AS hostname
```

### 3. Users with UID 0 other than root

```yaml
id: R-IAM-0001
name: "Non-root users with UID 0"
description: |
  Detects local user accounts with UID 0 that are not named root.
category: IDENTITY_ACCESS
severity: CRITICAL
status: approved
enabled: true
version: 1
priority: P1
created_by: "security-team"
cypher_query: |
  MATCH (h:Host)-[:HAS_USER]-(u:User)
  WHERE u.uid = 0 AND u.name <> 'root'
  RETURN u AS seed, h.hostname AS hostname, u.name AS username, u.uid AS uid
```

### 4. SSH listening on a non-standard port

```yaml
id: R-NET-0001
name: "SSH on non-standard port"
description: |
  Detects hosts where SSH is bound to a port other than 22.
category: NETWORK_SECURITY
severity: MEDIUM
status: draft
enabled: true
version: 1
priority: P3
cypher_query: |
  MATCH (h:Host)-[:LISTENS_ON]-(p:Port)
  WHERE p.service = 'ssh' AND p.port <> 22
  RETURN p AS seed, h.hostname AS hostname, p.port AS actual_port
```

---

## Naming Conventions

Rule IDs follow the format `R-[A-Z]{3}-[0-9]{4}`:

| Prefix   | Suggested use                      |
| -------- | ---------------------------------- |
| `R-SYS-` | Reserved for built-in system rules |
| `R-IAM-` | Identity & access management       |
| `R-NET-` | Network security                   |
| `R-PKG-` | Package / software                 |
| `R-OS-`  | Operating system                   |
| `R-SVC-` | Systemd services / processes       |
| `R-DKR-` | Docker / containers                |
| `R-VUL-` | Vulnerabilities                    |
| `R-CMP-` | Compliance benchmarks              |

---

## Validating Rules

Use the `validate` sub-command of `facter-rule-engine` to check rule YAML syntax and schema without persisting:

```bash
facter-rule-engine validate --config configs/config-facter-re.yml --rule rules/identity_access/R-IAM-0001.yaml
```

---

## Loading Rules from Files

Place individual YAML files in the `rules/` directory (one rule per file). They are loaded and upserted into PostgreSQL at startup.

Built-in system rules (those with IDs starting with `R-SYS-`) are read-only and cannot be deleted via the API.


