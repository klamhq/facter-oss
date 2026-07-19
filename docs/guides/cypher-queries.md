# Cypher Query Guide

The Facter graph database is powered by [Apache AGE](https://age.apache.org/), a PostgreSQL extension that adds openCypher graph query support. This guide covers how to write queries against the Facter inventory graph.

---

## Graph Model Overview

```
(:Host)-[:HAS_PACKAGE]->(:Package)
(:Host)-[:IS_NOT_UP_TO_DATE]->(:Package)
(:Host)-[:HAS_USER]->(:User)
(:Host)-[:HAS_SSH_KEY]->(:SshKey)
(:Host)-[:HAS_INTERFACE]->(:NetworkInterface)
(:Host)-[:HAS_PROCESS]->(:Process)
(:Host)-[:RUNS_SYSTEMD_SERVICE]->(:ServiceInstance)
(:Host)-[:HAS_APPLICATION]->(:Application)
(:Host)-[:HAS_DOCKER_CONTAINER]->(:DockerContainer)
(:Host)-[:HAS_DOCKER_IMAGE]->(:DockerImage)
(:Host)-[:HAS_DOCKER_NETWORK]->(:DockerNetwork)
(:Host)-[:IS_VULNERABLE_TO]->(:Vulnerability)
(:Host)-[:HAS_COMPLIANCE_REPORT]->(:ComplianceReport)
(:Host)-[:HAS_DISK]->(:Disk)
(:Host)-[:HAS_CPU]->(:CPU)
(:Host)-[:HAS_OS]->(:OperatingSystem)
(:Host)-[:HARDWARE_ID]->(:HardwareId)
(:Host)-[:OS_ID]->(:MachineId)
(:Host)-[:HAS_LOCATION]->(:GeoLocalisation)
(:Host)-[:HAS_IP]->(:IP)
(:Host)-[:HAS_PUBLIC_IP]->(:IP)
(:Host)-[:HAS_DNS]->(:DnsInfo)
(:User)-[:USER_HAS_SSH_KEY]->(:SshKey)
(:User)-[:AUTHORIZES_KEY]->(:SshKey)
(:User)-[:KNOWN_HOST]->(:Host)
(:User)-[:KNOWN_HOST]->(:ExternalNode)
(:User)-[:CAN_BECOME_ROOT_ON]->(:Host)
(:User)-[:CONNECTED_USER]->(:Host)
(:SshKey)-[:CAN_CONNECT_TO]->(:Host)
(:Package)-[:HAS_VULNERABILITY]->(:Vulnerability)
(:Process)-[:CHILD_OF]->(:Process)
(:Process)-[:RUN_BY_INSTANCE_OF_SYSTEMD_SERVICE]->(:ServiceInstance)
(:Process)-[:OPENS]->(:EstablishedConnection)
(:Process)-[:OPENS]->(:ListeningConnection)
(:ServiceInstance)-[:INSTANCE_OF_SYSTEMD_SERVICE]->(:SystemdService)
(:SystemdService)-[:REQUIRES_SYSTEMD_SERVICE]->(:SystemdService)
(:SystemdService)-[:WANTS_SYSTEMD_SERVICE]->(:SystemdService)
(:SystemdService)-[:BEFORE_SYSTEMD_SERVICE]->(:SystemdService)
(:SystemdService)-[:AFTER_SYSTEMD_SERVICE]->(:SystemdService)
(:Disk)-[:HAS_PARTITION]->(:Partition)
(:OperatingSystem)-[:BELONGS_TO_FAMILY]->(:SystemFamily)
(:NetworkInterface)-[:HAS_IP]->(:IP)
(:IP)-[:BELONGS_TO_NETWORK]->(:Network)
(:EstablishedConnection)-[:HAS_SOURCE_IP]->(:IP)
(:EstablishedConnection)-[:HAS_DESTINATION_IP]->(:IP)
(:ListeningConnection)-[:HAS_SOURCE_IP]->(:IP)
(:DockerContainer)-[:USE_DOCKER_NETWORKS]->(:DockerNetwork)
(:ComplianceReport)-[:HAS_COMPLIANCE_RULE]->(:ComplianceRule)
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
RETURN h.hostname, h.operating_system, h.kernel, h.uptime
```

| Property           | Type   | Description                      |
| ------------------ | ------ | -------------------------------- |
| `hostname`         | string | Fully qualified hostname         |
| `operating_system` | string | OS name (e.g. `Ubuntu 24.04`)    |
| `kernel`           | string | Kernel version string            |
| `uptime`           | int    | Uptime in seconds                |
| `memoryTotal`      | int    | Total RAM in bytes               |
| `memoryUsed`       | int    | Used RAM in bytes                |
| `vzSystem`         | string | Virtualisation system            |
| `vzRole`           | string | Virtualisation role (host/guest) |

### Package

```cypher
MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package)
RETURN h.hostname, p.name, p.version, p.architecture
LIMIT 20
```

| Property             | Description                        |
| -------------------- | ---------------------------------- |
| `name`               | Package name                       |
| `version`            | Installed version                  |
| `architecture`       | amd64, arm64, …                    |
| `upgradable_version` | Available upgrade version (if any) |
| `is_up_to_date`      | false when an upgrade is available |

### User

```cypher
MATCH (h:Host)-[:HAS_USER]->(u:User)
RETURN h.hostname, u.username, u.uid, u.gid, u.shell, u.home
```

| Property        | Description            |
| --------------- | ---------------------- |
| `username`      | Login name             |
| `uid`           | User ID                |
| `gid`           | Primary group ID       |
| `shell`         | Login shell            |
| `home`          | Home directory         |
| `canbecomeroot` | Has sudo / wheel group |

### SshKey

```cypher
MATCH (h:Host)-[:HAS_SSH_KEY]->(k:SshKey)
RETURN h.hostname, k.fingerprint, k.type, k.owner, k.path
```

| Property             | Description                         |
| -------------------- | ----------------------------------- |
| `fingerprint`        | Key fingerprint (SHA-256)           |
| `type`               | rsa, ecdsa, ed25519, dsa            |
| `owner`              | Username who owns the key           |
| `path`               | File path of the key on disk        |
| `fromAuthorizedKeys` | true when read from authorized_keys |

### NetworkInterface

```cypher
MATCH (h:Host)-[:HAS_INTERFACE]->(ni:NetworkInterface)-[:HAS_IP]->(ip:IP)
RETURN h.hostname, ni.name, ni.hardware_addr, ip.addr, ip.version
```

| Node               | Property        | Description             |
| ------------------ | --------------- | ----------------------- |
| `NetworkInterface` | `name`          | Interface name (eth0 …) |
| `NetworkInterface` | `hardware_addr` | MAC address             |
| `IP`               | `addr`          | IP address              |
| `IP`               | `version`       | ipv4 / ipv6             |
| `Network`          | `cidr`          | CIDR range              |

### Process

```cypher
MATCH (h:Host)-[:HAS_PROCESS]->(p:Process)
RETURN h.hostname, p.pid, p.name, p.username, p.cmdline
```

| Property     | Description       |
| ------------ | ----------------- |
| `pid`        | Process ID        |
| `name`       | Process name      |
| `username`   | Running user      |
| `cmdline`    | Full command line |
| `cpuPercent` | CPU usage %       |
| `memPercent` | Memory usage %    |

### SystemdService / ServiceInstance

Services are modelled as two nodes:

- `SystemdService` — the unit template (name and description, shared across hosts)
- `ServiceInstance` — the running instance on a specific host (state, PID, resource usage)

```cypher
MATCH (h:Host)-[:RUNS_SYSTEMD_SERVICE]->(i:ServiceInstance)-[:INSTANCE_OF_SYSTEMD_SERVICE]->(svc:SystemdService)
RETURN h.hostname, svc.name, i.active, i.subState, i.enabled
```

| Node              | Property      | Description                    |
| ----------------- | ------------- | ------------------------------ |
| `SystemdService`  | `name`        | Unit name (e.g. nginx.service) |
| `SystemdService`  | `description` | Unit description               |
| `ServiceInstance` | `active`      | true / false                   |
| `ServiceInstance` | `subState`    | running, dead, exited …        |
| `ServiceInstance` | `enabled`     | Enabled at boot                |
| `ServiceInstance` | `memory`      | Memory usage                   |

### Docker

```cypher
-- Containers
MATCH (h:Host)-[:HAS_DOCKER_CONTAINER]->(c:DockerContainer)
RETURN h.hostname, c.name, c.image, c.state, c.status

-- Images
MATCH (h:Host)-[:HAS_DOCKER_IMAGE]->(i:DockerImage)
RETURN h.hostname, i.repoTags, i.size

-- Networks
MATCH (h:Host)-[:HAS_DOCKER_NETWORK]->(n:DockerNetwork)
RETURN h.hostname, n.name, n.driver, n.internal
```

| Node              | Property   | Description              |
| ----------------- | ---------- | ------------------------ |
| `DockerContainer` | `name`     | Container name           |
| `DockerContainer` | `image`    | Image name used          |
| `DockerContainer` | `state`    | running, stopped, exited |
| `DockerContainer` | `status`   | Human-readable status    |
| `DockerImage`     | `repoTags` | List of repo:tag labels  |
| `DockerNetwork`   | `name`     | Network name             |
| `DockerNetwork`   | `driver`   | bridge, overlay, host …  |

### Vulnerability

```cypher
MATCH (h:Host)-[:IS_VULNERABLE_TO]->(v:Vulnerability)
RETURN h.hostname, v.id, v.severity, v.title
ORDER BY v.severity DESC
```

| Property       | Description                         |
| -------------- | ----------------------------------- |
| `id`           | CVE identifier (e.g. CVE-2024-1234) |
| `severity`     | CRITICAL, HIGH, MEDIUM, LOW         |
| `title`        | Short title                         |
| `description`  | Full description                    |
| `fixedVersion` | Package version that fixes the CVE  |

### Network Connections

```cypher
-- Established connections
MATCH (p:Process)-[:OPENS]->(c:EstablishedConnection)-[:HAS_DESTINATION_IP]->(dst:IP)
MATCH (c)-[:HAS_SOURCE_IP]->(src:IP)
RETURN p.name, src.addr, c.localPort, dst.addr, c.remotePort, c.protocol

-- Listening connections
MATCH (p:Process)-[:OPENS]->(l:ListeningConnection)-[:HAS_SOURCE_IP]->(ip:IP)
RETURN p.name, ip.addr, l.localPort, l.protocol
```

---

## Common Query Patterns

### Count: how many hosts have a given package?

```cypher
MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package {name: "nginx"})
RETURN p.name AS package, count(h) AS host_count
```

### Find hosts WITHOUT a package

```cypher
MATCH (h:Host)
WHERE NOT (h)-[:HAS_PACKAGE]->(:Package {name: "curl"})
RETURN h.hostname AS hostname
```

### Find outdated packages across all hosts

```cypher
MATCH (h:Host)-[:IS_NOT_UP_TO_DATE]->(p:Package)
RETURN h.hostname, p.name, p.version, p.upgradable_version
ORDER BY h.hostname, p.name
```

### Find users with sudo on more than 3 hosts

```cypher
MATCH (u:User {canbecomeroot: true})-[:CAN_BECOME_ROOT_ON]->(h:Host)
WITH u.username AS username, collect(h.hostname) AS hosts, count(h) AS cnt
WHERE cnt > 3
RETURN username, hosts, cnt
ORDER BY cnt DESC
```

### Find users with interactive shells who have UID 0

```cypher
MATCH (h:Host)-[:HAS_USER]->(u:User)
WHERE u.uid = "0" AND u.shell <> '/sbin/nologin' AND u.shell <> '/bin/false'
RETURN h.hostname, u.username, u.uid, u.shell
```

### Cross-host: find shared SSH keys (same fingerprint)

```cypher
MATCH (h1:Host)-[:HAS_SSH_KEY]->(k:SshKey)<-[:HAS_SSH_KEY]-(h2:Host)
WHERE h1.hostname < h2.hostname
RETURN k.fingerprint AS key,
       k.type AS key_type,
       h1.hostname AS host1,
       h2.hostname AS host2
```

### Who can SSH into a host?

```cypher
MATCH (k:SshKey)-[r:CAN_CONNECT_TO]->(h:Host {hostname: "my-server.example.com"})
MATCH (u:User)-[:AUTHORIZES_KEY]->(k)
RETURN h.hostname, u.username, k.fingerprint, r.as_user
```

### Known hosts: which external nodes does a user know?

```cypher
MATCH (h:Host)-[:HAS_USER]->(u:User)-[:KNOWN_HOST]->(target)
RETURN h.hostname, u.username, target.hostname AS known_host, labels(target) AS target_type
ORDER BY u.username
```

### Hosts with critical CVEs

```cypher
MATCH (h:Host)-[:IS_VULNERABLE_TO]->(v:Vulnerability {severity: "CRITICAL"})
RETURN h.hostname, collect(v.id) AS critical_cves, count(v) AS vuln_count
ORDER BY vuln_count DESC
```

### Package → Vulnerability cross-reference

```cypher
MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package)-[:HAS_VULNERABILITY]->(v:Vulnerability)
RETURN h.hostname, p.name, p.version, v.id, v.severity
ORDER BY v.severity, h.hostname
```

### List all listening ports with owning process

```cypher
MATCH (p:Process)-[:OPENS]->(l:ListeningConnection)-[:HAS_SOURCE_IP]->(ip:IP)
MATCH (h:Host)-[:HAS_PROCESS]->(p)
RETURN h.hostname, ip.addr, l.localPort, l.protocol, p.name
ORDER BY h.hostname, l.localPort
```

### Docker containers using a specific network

```cypher
MATCH (c:DockerContainer)-[:USE_DOCKER_NETWORKS]->(n:DockerNetwork {name: "my-net"})
MATCH (h:Host)-[:HAS_DOCKER_CONTAINER]->(c)
RETURN h.hostname, c.name, c.image, c.state
```

### Systemd service dependencies

```cypher
MATCH (svc:SystemdService {name: "nginx.service"})-[:REQUIRES_SYSTEMD_SERVICE|WANTS_SYSTEMD_SERVICE]->(dep:SystemdService)
RETURN svc.name, dep.name
```

### Services in failed state

```cypher
MATCH (h:Host)-[:RUNS_SYSTEMD_SERVICE]->(i:ServiceInstance)
WHERE i.active = false AND i.subState = 'failed'
RETURN h.hostname, i.name, i.subState
ORDER BY h.hostname
```

### Compliance: hosts failing a specific rule

```cypher
MATCH (h:Host)-[:HAS_COMPLIANCE_REPORT]->(r:ComplianceReport)-[rel:HAS_COMPLIANCE_RULE {result: "fail"}]->(rule:ComplianceRule)
WHERE rule.id = "xccdf_org.ssgproject.content_rule_sshd_disable_root_login"
RETURN h.hostname, r.profile, rule.severity, rel.result
ORDER BY h.hostname
```

---

## Aggregation

### count()

```cypher
MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package)
RETURN h.hostname, count(p) AS package_count
ORDER BY package_count DESC
```

### collect()

```cypher
MATCH (h:Host)-[:HAS_USER]->(u:User)
RETURN h.hostname, collect(u.username) AS users
```

### WITH for multi-step aggregations

```cypher
MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package)
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

### Comparison

```cypher
WHERE u.uid >= '1000'
```

---

## Tips for Writing Audit Queries

1. **Always alias return columns** — use `RETURN x.prop AS alias`. Templates and conditions use these aliases.
2. **Test queries directly** — connect to Apache AGE using `psql` or pgAdmin and run queries in the `age` database before adding them to rules.
3. **Use `LIMIT` when exploring** — omit `LIMIT` in production rules.
4. **Avoid Cartesian products** — always specify relationships between nodes using `-[:REL_TYPE]->` syntax.
5. **Optional matches** — use `OPTIONAL MATCH` for properties that may not exist on all nodes.
6. **String UIDs** — `uid` and `gid` are stored as strings in the graph; use string comparisons (`u.uid = "0"` not `= 0`).

---

## Running a Query Directly Against Apache AGE

```sql
-- Connect to the age database
SET search_path = ag_catalog, "$user", public;
SELECT * FROM cypher('facter', $$
  MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package {name: 'curl'})
  RETURN h.hostname AS hostname
$$) AS (hostname agtype);
```

Or using the `validate` command of facter-rule-engine:

```bash
facter-rule-engine validate --config configs/config.yml --rule R-PKG-0001
```
