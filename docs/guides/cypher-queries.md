# Cypher Query Guide

The Facter graph database is powered by [Apache AGE](https://age.apache.org/), a PostgreSQL extension that adds openCypher graph query support. This guide covers how to write queries against the Facter inventory graph.

---

## Graph Model Overview

```
(Host) -[:DEPLOYED_ON]-> (Package)
(Host) -[:HAS_USER]-> (User)
(User) -[:HAS_SSH_KEY]-> (SSHKey)
(Host) -[:HAS_INTERFACE]-> (Interface)
(Host) -[:LISTENS_ON]-> (Port)
(Host) -[:HAS_PROCESS]-> (Process)
(Host) -[:HAS_SERVICE]-> (SystemdService)
(Host) -[:RUNS_APP]-> (Application)   -- Docker containers
(Host) -[:HAS_VULN]-> (Vulnerability)
(Host) -[:HAS_OS]-> (OS)
(Host) -[:HAS_HARDWARE]-> (Hardware)
```

See [Data Model](../concepts/data-model.md) for the full node and property reference.

---

## Basic Query Structure

```cypher
MATCH <pattern>
WHERE <filter>
RETURN <columns>
```

All Facter audit queries must have both `MATCH` and `RETURN` clauses.

---

## Node Labels and Key Properties

### Host

```cypher
MATCH (h:Host)
RETURN h.hostname, h.os, h.os_version, h.kernel_version, h.arch
```

| Property         | Type   | Description              |
| ---------------- | ------ | ------------------------ |
| `hostname`       | string | Fully qualified hostname |
| `os`             | string | OS family (e.g. `linux`) |
| `os_version`     | string | Distribution + version   |
| `kernel_version` | string | Kernel version string    |
| `arch`           | string | CPU architecture         |

### Package

```cypher
MATCH (h:Host)-[:DEPLOYED_ON]-(p:Package)
RETURN h.hostname, p.name, p.version, p.manager
LIMIT 20
```

| Property  | Description                                |
| --------- | ------------------------------------------ |
| `name`    | Package name                               |
| `version` | Installed version                          |
| `manager` | Package manager (`apt`, `rpm`, `pip`, ...) |

### User

```cypher
MATCH (h:Host)-[:HAS_USER]-(u:User)
RETURN h.hostname, u.name, u.uid, u.gid, u.shell, u.home
```

| Property | Description      |
| -------- | ---------------- |
| `name`   | Username         |
| `uid`    | User ID          |
| `gid`    | Primary group ID |
| `shell`  | Login shell      |
| `home`   | Home directory   |

### SSHKey

```cypher
MATCH (h:Host)-[:HAS_USER]-(u:User)-[:HAS_SSH_KEY]-(k:SSHKey)
RETURN h.hostname, u.name, k.fingerprint, k.key_type, k.comment
```

| Property      | Description                     |
| ------------- | ------------------------------- |
| `fingerprint` | SHA256 fingerprint              |
| `key_type`    | Algorithm (ed25519, rsa, ecdsa) |
| `comment`     | Key comment / label             |

### Interface

```cypher
MATCH (h:Host)-[:HAS_INTERFACE]-(i:Interface)
RETURN h.hostname, i.name, i.mac_address, i.ip_address, i.cidr
```

### Port

```cypher
MATCH (h:Host)-[:LISTENS_ON]-(p:Port)
RETURN h.hostname, p.port, p.proto, p.service, p.pid
```

### Process

```cypher
MATCH (h:Host)-[:HAS_PROCESS]-(proc:Process)
RETURN h.hostname, proc.pid, proc.name, proc.user, proc.cmdline
```

### SystemdService

```cypher
MATCH (h:Host)-[:HAS_SERVICE]-(s:SystemdService)
RETURN h.hostname, s.name, s.state, s.sub_state, s.enabled
```

### Application (Docker)

```cypher
MATCH (h:Host)-[:RUNS_APP]-(app:Application)
RETURN h.hostname, app.name, app.image, app.status, app.ports
```

### Vulnerability

```cypher
MATCH (h:Host)-[:HAS_VULN]-(v:Vulnerability)
RETURN h.hostname, v.cve_id, v.severity, v.description
ORDER BY v.severity DESC
```

---

## Common Query Patterns

### Count: how many hosts have a given package?

```cypher
MATCH (h:Host)-[:DEPLOYED_ON]-(p:Package {name: "nginx"})
RETURN p.name AS package, count(h) AS host_count
```

### Find hosts WITHOUT a package

```cypher
MATCH (h:Host)
WHERE NOT (h)-[:DEPLOYED_ON]-(:Package {name: "curl"})
RETURN h.hostname AS hostname
```

### Cross-host: find shared SSH keys

```cypher
MATCH (h1:Host)-[:HAS_USER]-(u1:User)-[:HAS_SSH_KEY]-(k:SSHKey)-[:HAS_SSH_KEY]-(u2:User)-[:HAS_USER]-(h2:Host)
WHERE h1.hostname < h2.hostname
RETURN k.fingerprint AS key,
       u1.name + '@' + h1.hostname AS user1,
       u2.name + '@' + h2.hostname AS user2
```

### Find users with interactive shells who have UID 0

```cypher
MATCH (h:Host)-[:HAS_USER]-(u:User)
WHERE u.uid = 0 AND u.shell <> '/sbin/nologin' AND u.shell <> '/bin/false'
RETURN h.hostname, u.name, u.uid, u.shell
```

### List all open ports with service name

```cypher
MATCH (h:Host)-[:LISTENS_ON]-(p:Port)
WHERE p.proto = 'tcp' AND p.state = 'LISTEN'
RETURN h.hostname, p.port, p.service, p.pid
ORDER BY h.hostname, p.port
```

### Docker containers with port 80 exposed

```cypher
MATCH (h:Host)-[:RUNS_APP]-(app:Application)
WHERE app.ports CONTAINS '80'
RETURN h.hostname, app.name, app.image, app.ports
```

### Hosts with critical CVEs

```cypher
MATCH (h:Host)-[:HAS_VULN]-(v:Vulnerability {severity: "CRITICAL"})
RETURN h.hostname, collect(v.cve_id) AS critical_cves, count(v) AS vuln_count
ORDER BY vuln_count DESC
```

### Services in failed state

```cypher
MATCH (h:Host)-[:HAS_SERVICE]-(s:SystemdService)
WHERE s.state = 'failed'
RETURN h.hostname, s.name, s.sub_state
ORDER BY h.hostname
```

### Find all sudo-capable users

```cypher
MATCH (h:Host)-[:HAS_USER]-(u:User)
WHERE u.groups CONTAINS 'sudo' OR u.groups CONTAINS 'wheel'
RETURN h.hostname, u.name, u.uid, u.groups
```

---

## Aggregation

### count()

```cypher
MATCH (h:Host)-[:DEPLOYED_ON]-(p:Package)
RETURN h.hostname, count(p) AS package_count
ORDER BY package_count DESC
```

### collect()

```cypher
MATCH (h:Host)-[:HAS_USER]-(u:User)
WHERE u.uid >= 1000
RETURN h.hostname, collect(u.name) AS regular_users
```

### WITH for multi-step aggregations

```cypher
MATCH (h:Host)-[:DEPLOYED_ON]-(p:Package)
WITH h, count(p) AS pkg_count
WHERE pkg_count > 500
RETURN h.hostname, pkg_count
ORDER BY pkg_count DESC
```

---

## Filtering with WHERE

### String contains

```cypher
WHERE u.shell CONTAINS 'bash'
```

### Property exists

```cypher
WHERE k.fingerprint IS NOT NULL
```

### Property does not exist

```cypher
WHERE u.sudo_entry IS NULL
```

### List / contains (string match)

```cypher
WHERE app.ports CONTAINS '443'
```

### Comparison

```cypher
WHERE u.uid >= 1000 AND u.uid < 65534
```

---

## Tips for Writing Audit Queries

1. **Always alias return columns** — use `RETURN x.prop AS alias`. Templates and conditions use these aliases.
2. **Test queries directly** — connect to Apache AGE using `psql` or pgAdmin and run queries in the `age` database before adding them to rules.
3. **Use `LIMIT` when exploring** — omit `LIMIT` in production rules.
4. **Avoid Cartesian products** — always specify relationships between nodes using `-[r:REL_TYPE]-` syntax.
5. **Optional matches** — use `OPTIONAL MATCH` for properties that may not exist on all nodes.

---

## Running a Query Directly Against Apache AGE

```sql
-- Connect to the age database
SET search_path = ag_catalog, "$user", public;
SELECT * FROM cypher('facter', $$
  MATCH (h:Host)-[:DEPLOYED_ON]-(p:Package {name: 'curl'})
  RETURN h.hostname AS hostname
$$) AS (hostname agtype);
```

Or using the `validate` command of facter-rule-engine:

```bash
facter-rule-engine validate --config configs/config.yml --rule R-PKG-0001
```
