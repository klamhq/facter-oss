# Data Model

## Overview

All host inventory data is stored as a **property graph** in [Apache AGE](https://age.apache.org/), a PostgreSQL extension that implements the openCypher graph query language.

The graph name is configurable (default: `facter`). Every node and edge has typed labels and properties.

---

## Node Types

### Host

The root node representing a monitored machine.

| Property     | Type      | Description                   |
| ------------ | --------- | ----------------------------- |
| `hostname`   | string    | FQDN or short hostname        |
| `machine_id` | string    | `/etc/machine-id` value       |
| `uuid`       | string    | DMI product UUID              |
| `created_at` | timestamp | First time this host was seen |
| `updated_at` | timestamp | Last inventory update         |

---

### Package

An installed software package on a host.

| Property  | Type   | Description                    |
| --------- | ------ | ------------------------------ |
| `name`    | string | Package name                   |
| `version` | string | Installed version              |
| `arch`    | string | Architecture (amd64, arm64, …) |
| `manager` | string | apt, rpm, pacman, brew         |

---

### User

A local system user.

| Property          | Type   | Description            |
| ----------------- | ------ | ---------------------- |
| `username`        | string | Login name             |
| `uid`             | string | User ID                |
| `gid`             | string | Primary group ID       |
| `home`            | string | Home directory         |
| `shell`           | string | Login shell            |
| `can_become_root` | bool   | Has sudo / wheel group |

---

### SSHKey

An SSH key found on disk.

| Property      | Type   | Description               |
| ------------- | ------ | ------------------------- |
| `fingerprint` | string | Key fingerprint (SHA-256) |
| `type`        | string | rsa, ecdsa, ed25519, dsa  |
| `bits`        | int    | Key size in bits          |
| `comment`     | string | Key comment field         |

---

### Interface

A network interface on a host.

| Property        | Type   | Description                    |
| --------------- | ------ | ------------------------------ |
| `name`          | string | Interface name (eth0, ens3, …) |
| `hardware_addr` | string | MAC address                    |

---

### Process

A running process.

| Property | Type   | Description       |
| -------- | ------ | ----------------- |
| `pid`    | int    | Process ID        |
| `name`   | string | Process name      |
| `cmd`    | string | Full command line |
| `user`   | string | Running user      |

---

### Application

A Docker application stack or container group.

| Property | Type   | Description                        |
| -------- | ------ | ---------------------------------- |
| `name`   | string | Application / compose project name |
| `type`   | string | docker                             |

### Container

| Property | Type   | Description              |
| -------- | ------ | ------------------------ |
| `id`     | string | Container ID             |
| `name`   | string | Container name           |
| `image`  | string | Image name:tag           |
| `state`  | string | running, stopped, exited |
| `status` | string | Up 2 hours, Exited (0) … |

---

### SystemdService

| Property    | Type   | Description                    |
| ----------- | ------ | ------------------------------ |
| `name`      | string | Unit name (e.g. nginx.service) |
| `active`    | string | active / inactive              |
| `loaded`    | string | loaded / not-found             |
| `enabled`   | bool   | Enabled at boot                |
| `sub_state` | string | running, dead, exited …        |

---

### VulnerabilityReport / ComplianceReport

Aggregated results from OpenSCAP scans. Stored on the Host node as embedded properties or linked nodes depending on configuration.

---

## Relationship Types

```
(:Host)-[:HAS_PACKAGE]->(:Package)
(:Host)-[:HAS_USER]->(:User)
(:Host)-[:HAS_INTERFACE]->(:Interface)
(:Host)-[:HAS_PROCESS]->(:Process)
(:Host)-[:HAS_SERVICE]->(:SystemdService)
(:Host)-[:HAS_APPLICATION]->(:Application)
(:Application)-[:HAS_CONTAINER]->(:Container)
(:SSHKey)-[:DEPLOYED_ON]->(:Host)
(:User)-[:HAS_SSH_KEY]->(:SSHKey)
(:User)-[:AUTHORIZED_KEY]->(:SSHKey)
(:Host)-[:KNOWS]->(:KnownHost)
(:Interface)-[:HAS_IP]->(:Ip)
(:Host)-[:HAS_CONNECTION]->(:Connection)
(:Connection)-[:TO]->(:Ip)
```

---

## Example Cypher Queries

Find all hosts with a specific package installed:

```cypher
MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package {name: 'openssl'})
RETURN h.hostname, p.version
ORDER BY h.hostname
```

Find SSH keys shared across multiple hosts (shared key detection):

```cypher
MATCH (k:SSHKey)-[:DEPLOYED_ON]->(h:Host)
WITH k, collect(DISTINCT h.hostname) AS hosts, count(DISTINCT h) AS cnt
WHERE cnt > 1
RETURN k.fingerprint, k.type, hosts, cnt
ORDER BY cnt DESC
```

Find users with sudo capability on more than 3 hosts:

```cypher
MATCH (u:User {can_become_root: true})--(h:Host)
WITH u.username AS username, collect(h.hostname) AS hosts, count(h) AS cnt
WHERE cnt > 3
RETURN username, hosts, cnt
ORDER BY cnt DESC
```

Find all Docker containers using the `latest` tag:

```cypher
MATCH (h:Host)-[:HAS_APPLICATION]->(a:Application)-[:HAS_CONTAINER]->(c:Container)
WHERE c.image ENDS WITH ':latest'
RETURN h.hostname, a.name, c.name, c.image
```

---

## Snapshots

`facter-grpc` can take daily snapshots of the graph. A snapshot creates a copy of the current graph state with a timestamp, enabling:

- Trend analysis ("how many vulnerable packages did we have last week?")
- Diff between two points in time
- Historical compliance reports

Snapshots are purged automatically after a configurable number of days (`maxAgeDays`).
