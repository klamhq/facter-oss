# Data Model

## Overview

All host inventory data is stored as a **property graph** in [Apache AGE](https://age.apache.org/), a PostgreSQL extension that implements the openCypher graph query language.

The graph name is configurable (default: `facter`). Every node and edge has typed labels and properties.

The complete list of relationship types is auto-generated in [RELATIONS.md](https://github.com/klamhq/facter-grpc/blob/main/RELATIONS.md) (see `make export-relations` in the `facter-grpc` repository).

---

## Node Types

### Host

The root node representing a monitored machine.

| Property           | Type   | Description                      |
| ------------------ | ------ | -------------------------------- |
| `hostname`         | string | FQDN or short hostname           |
| `operating_system` | string | OS name (e.g. Ubuntu 24.04)      |
| `kernel`           | string | Kernel version string            |
| `uptime`           | int    | Uptime in seconds                |
| `memoryTotal`      | int    | Total RAM in bytes               |
| `memoryUsed`       | int    | Used RAM in bytes                |
| `swapTotal`        | int    | Total swap in bytes              |
| `vzSystem`         | string | Virtualisation system            |
| `vzRole`           | string | Virtualisation role (host/guest) |

---

### Package

An installed software package on a host.

| Property             | Type   | Description                        |
| -------------------- | ------ | ---------------------------------- |
| `name`               | string | Package name                       |
| `version`            | string | Installed version                  |
| `architecture`       | string | amd64, arm64, …                    |
| `description`        | string | Package description                |
| `upgradable_version` | string | Available upgrade version (if any) |
| `is_up_to_date`      | bool   | false when an upgrade is available |
| `active`             | bool   | false when soft-deleted            |

---

### User

A local system user.

| Property        | Type   | Description             |
| --------------- | ------ | ----------------------- |
| `username`      | string | Login name              |
| `uid`           | string | User ID                 |
| `gid`           | string | Primary group ID        |
| `home`          | string | Home directory          |
| `shell`         | string | Login shell             |
| `canbecomeroot` | bool   | Has sudo / wheel group  |
| `active`        | bool   | false when soft-deleted |

---

### SshKey

An SSH key found on disk (private key, public key, or authorized_keys entry).

| Property             | Type   | Description                         |
| -------------------- | ------ | ----------------------------------- |
| `fingerprint`        | string | Key fingerprint (SHA-256)           |
| `type`               | string | rsa, ecdsa, ed25519, dsa            |
| `name`               | string | File name (e.g. authorized_keys)    |
| `path`               | string | Full path on disk                   |
| `owner`              | string | Username who owns the key           |
| `comment`            | string | Key comment field                   |
| `length`             | string | Key size (bits)                     |
| `fromAuthorizedKeys` | bool   | true when read from authorized_keys |
| `active`             | bool   | false when soft-deleted             |

---

### NetworkInterface

A network interface on a host.

| Property        | Type   | Description                    |
| --------------- | ------ | ------------------------------ |
| `name`          | string | Interface name (eth0, ens3, …) |
| `hardware_addr` | string | MAC address                    |

---

### IP

An IP address, linked to a network interface or used as a public address.

| Property       | Type   | Description                    |
| -------------- | ------ | ------------------------------ |
| `addr`         | string | IP address                     |
| `version`      | string | ipv4 / ipv6                    |
| `externalIp`   | bool   | true for the public IP address |
| `forwardedFor` | string | X-Forwarded-For header value   |

---

### Network

A network/subnet, linked to IP addresses.

| Property | Type   | Description   |
| -------- | ------ | ------------- |
| `cidr`   | string | CIDR notation |

---

### GeoLocalisation

Geographic location of the host (derived from public IP).

| Property    | Type  | Description              |
| ----------- | ----- | ------------------------ |
| `latitude`  | float | Latitude                 |
| `longitude` | float | Longitude                |
| `accuracy`  | float | Accuracy radius (meters) |

---

### DnsInfo

DNS configuration of a host.

| Property        | Type         | Description            |
| --------------- | ------------ | ---------------------- |
| `nameservers`   | list(string) | DNS server addresses   |
| `searchdomains` | list(string) | DNS search domains     |
| `port`          | int          | DNS port (default: 53) |

---

### Disk

A physical or virtual disk device.

| Property | Type   | Description |
| -------- | ------ | ----------- |
| `device` | string | Device path |
| `uuid`   | string | Disk UUID   |

---

### Partition

A disk partition.

| Property      | Type   | Description        |
| ------------- | ------ | ------------------ |
| `mountpoint`  | string | Mount point        |
| `device`      | string | Parent device path |
| `fsType`      | string | File system type   |
| `total`       | int    | Total size (bytes) |
| `used`        | int    | Used space (bytes) |
| `free`        | int    | Free space (bytes) |
| `usedPercent` | float  | Usage percentage   |

---

### CPU

CPU information for a host.

| Property | Type   | Description       |
| -------- | ------ | ----------------- |
| `model`  | string | CPU model name    |
| `core`   | int    | Number of cores   |
| `mhz`    | float  | Clock speed (MHz) |

---

### OperatingSystem

The operating system installed on a host.

| Property     | Type   | Description                     |
| ------------ | ------ | ------------------------------- |
| `name`       | string | Distribution name               |
| `version`    | string | Distribution version            |
| `initSystem` | string | Init system (systemd, openrc …) |

---

### SystemFamily

The OS family grouping operating systems.

| Property | Type   | Description                 |
| -------- | ------ | --------------------------- |
| `name`   | string | Family name (debian, rhel…) |

---

### HardwareId

DMI product UUID for hardware identification.

| Property | Type   | Description      |
| -------- | ------ | ---------------- |
| `uuid`   | string | DMI product UUID |

---

### MachineId

The `/etc/machine-id` OS-level identifier.

| Property | Type   | Description      |
| -------- | ------ | ---------------- |
| `id`     | string | machine-id value |

---

### Process

A running process on a host.

| Property     | Type   | Description             |
| ------------ | ------ | ----------------------- |
| `pid`        | string | Process ID              |
| `name`       | string | Process name            |
| `cmdline`    | string | Full command line       |
| `username`   | string | Running user            |
| `exe`        | string | Executable path         |
| `createTime` | int    | Creation timestamp      |
| `parent`     | int    | Parent PID              |
| `cpuPercent` | float  | CPU usage %             |
| `memPercent` | float  | Memory usage %          |
| `active`     | bool   | false when soft-deleted |

---

### SystemdService

A systemd unit template (shared across hosts, represents the service definition).

| Property      | Type   | Description                |
| ------------- | ------ | -------------------------- |
| `name`        | string | Unit name (nginx.service…) |
| `description` | string | Unit description           |

---

### ServiceInstance

A specific running instance of a systemd service on a host.

| Property   | Type   | Description                  |
| ---------- | ------ | ---------------------------- |
| `name`     | string | Unit name                    |
| `hostname` | string | Host where the instance runs |
| `active`   | bool   | true when active             |
| `subState` | string | running, dead, exited …      |
| `enabled`  | bool   | Enabled at boot              |
| `pid`      | string | Main process PID             |
| `memory`   | int    | Memory usage (bytes)         |
| `tasks`    | int    | Number of tasks              |
| `cpuUsage` | string | CPU time used                |
| `loaded`   | string | loaded / not-found           |
| `cgroup`   | string | CGroup path                  |

---

### Application

A top-level application marker. Currently only one node is created (name = "Docker") to indicate Docker is present on the host.

| Property | Type   | Description      |
| -------- | ------ | ---------------- |
| `name`   | string | Application name |

---

### DockerContainer

A Docker container running on a host.

| Property      | Type         | Description              |
| ------------- | ------------ | ------------------------ |
| `name`        | string       | Container name           |
| `image`       | string       | Image name:tag           |
| `imageID`     | string       | Image SHA                |
| `state`       | string       | running, stopped, exited |
| `status`      | string       | Human-readable status    |
| `created`     | string       | Creation timestamp       |
| `networkMode` | string       | Network mode             |
| `privatePort` | list(string) | Internal ports           |
| `publicPort`  | list(string) | Exposed ports            |
| `networkID`   | list(string) | Connected network IDs    |

---

### DockerImage

A Docker image present on a host.

| Property     | Type         | Description        |
| ------------ | ------------ | ------------------ |
| `repoTags`   | list(string) | Repository tags    |
| `size`       | string       | Image size         |
| `created`    | string       | Build timestamp    |
| `parentId`   | string       | Parent image SHA   |
| `repoDigest` | list(string) | Repository digests |

---

### DockerNetwork

A Docker network present on a host.

| Property   | Type   | Description                   |
| ---------- | ------ | ----------------------------- |
| `name`     | string | Network name                  |
| `driver`   | string | bridge, overlay, host …       |
| `internal` | bool   | Internal (no external access) |
| `scope`    | string | local / swarm                 |

---

### EstablishedConnection

An established TCP/UDP connection on a host.

| Property     | Type   | Description |
| ------------ | ------ | ----------- |
| `localPort`  | string | Local port  |
| `remotePort` | string | Remote port |
| `protocol`   | string | tcp / udp   |

---

### ListeningConnection

A socket in the LISTEN state on a host.

| Property    | Type   | Description    |
| ----------- | ------ | -------------- |
| `localPort` | string | Listening port |
| `protocol`  | string | tcp / udp      |

---

### Vulnerability

A CVE vulnerability matched against installed packages.

| Property       | Type   | Description                         |
| -------------- | ------ | ----------------------------------- |
| `id`           | string | CVE identifier (e.g. CVE-2024-1234) |
| `severity`     | string | CRITICAL, HIGH, MEDIUM, LOW         |
| `title`        | string | Short title                         |
| `description`  | string | Full description                    |
| `fixedVersion` | string | Package version that fixes the CVE  |

---

### ComplianceReport

An OpenSCAP benchmark result for a host.

| Property       | Type   | Description            |
| -------------- | ------ | ---------------------- |
| `scoreMaximum` | float  | Maximum possible score |
| `scoreValue`   | float  | Achieved score         |
| `profile`      | string | XCCDF profile applied  |

---

### ComplianceRule

An individual XCCDF rule result within a compliance report.

| Property      | Type   | Description          |
| ------------- | ------ | -------------------- |
| `id`          | string | XCCDF rule ID        |
| `title`       | string | Rule title           |
| `description` | string | Rule description     |
| `severity`    | string | Rule severity        |
| `fix`         | string | Remediation guidance |

The relationship `(:ComplianceReport)-[:HAS_COMPLIANCE_RULE {result}]->(:ComplianceRule)` carries the `result` property (`pass`, `fail`, `notselected` …).

---

### ExternalNode

A host referenced in a known_hosts file but not present in the Facter inventory.

| Property     | Type         | Description                |
| ------------ | ------------ | -------------------------- |
| `hostname`   | string       | Hostname or IP as known    |
| `identities` | list(string) | Hostnames/IPs in the entry |
| `from`       | string       | Origin (`known_host`)      |

---

## Relationship Types

The following relationships are created by the `facter-grpc` inventory pipeline. All are written with `MERGE` (idempotent).

```
# Host → inventory nodes
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

# Host → platform nodes
(:Host)-[:HAS_DISK]->(:Disk)
(:Host)-[:HAS_CPU]->(:CPU)
(:Host)-[:HAS_OS]->(:OperatingSystem)
(:Host)-[:HARDWARE_ID]->(:HardwareId)
(:Host)-[:OS_ID]->(:MachineId)

# Host → network nodes
(:Host)-[:HAS_INTERFACE]->(:NetworkInterface)
(:Host)-[:HAS_LOCATION]->(:GeoLocalisation)
(:Host)-[:HAS_IP]->(:IP)
(:Host)-[:HAS_PUBLIC_IP]->(:IP)
(:Host)-[:HAS_DNS]->(:DnsInfo)

# User relationships
(:User)-[:USER_HAS_SSH_KEY]->(:SshKey)
(:User)-[:AUTHORIZES_KEY]->(:SshKey)
(:User)-[:KNOWN_HOST]->(:Host)
(:User)-[:KNOWN_HOST]->(:ExternalNode)
(:User)-[:CAN_BECOME_ROOT_ON]->(:Host)
(:User)-[:CONNECTED_USER {terminal, host, started}]->(:Host)

# SSH access
(:SshKey)-[:CAN_CONNECT_TO {as_user}]->(:Host)

# Package → vulnerability
(:Package)-[:HAS_VULNERABILITY]->(:Vulnerability)

# Process
(:Process)-[:CHILD_OF]->(:Process)
(:Process)-[:RUN_BY_INSTANCE_OF_SYSTEMD_SERVICE]->(:ServiceInstance)
(:Process)-[:OPENS]->(:EstablishedConnection)
(:Process)-[:OPENS]->(:ListeningConnection)

# Systemd
(:ServiceInstance)-[:INSTANCE_OF_SYSTEMD_SERVICE]->(:SystemdService)
(:SystemdService)-[:REQUIRES_SYSTEMD_SERVICE]->(:SystemdService)
(:SystemdService)-[:WANTS_SYSTEMD_SERVICE]->(:SystemdService)
(:SystemdService)-[:BEFORE_SYSTEMD_SERVICE]->(:SystemdService)
(:SystemdService)-[:AFTER_SYSTEMD_SERVICE]->(:SystemdService)

# Platform
(:Disk)-[:HAS_PARTITION]->(:Partition)
(:OperatingSystem)-[:BELONGS_TO_FAMILY]->(:SystemFamily)

# Network
(:NetworkInterface)-[:HAS_IP]->(:IP)
(:IP)-[:BELONGS_TO_NETWORK]->(:Network)
(:EstablishedConnection)-[:HAS_SOURCE_IP]->(:IP)
(:EstablishedConnection)-[:HAS_DESTINATION_IP]->(:IP)
(:ListeningConnection)-[:HAS_SOURCE_IP]->(:IP)

# Docker
(:DockerContainer)-[:USE_DOCKER_NETWORKS]->(:DockerNetwork)

# Compliance
(:ComplianceReport)-[:HAS_COMPLIANCE_RULE {result}]->(:ComplianceRule)
```

> This list is auto-generated from facter-grpc source code. Run `make export-relations` in the `facter-grpc` repository to regenerate.

---

## Example Cypher Queries

Find all hosts with a specific package installed:

```cypher
MATCH (h:Host)-[:HAS_PACKAGE]->(p:Package {name: 'openssl'})
RETURN h.hostname, p.version
ORDER BY h.hostname
```

Find SSH keys that can connect to a host:

```cypher
MATCH (k:SshKey)-[r:CAN_CONNECT_TO]->(h:Host {hostname: 'prod-server-01'})
MATCH (u:User)-[:AUTHORIZES_KEY]->(k)
RETURN u.username, k.fingerprint, k.type, r.as_user
```

Find users with sudo access on more than 3 hosts:

```cypher
MATCH (u:User {canbecomeroot: true})-[:CAN_BECOME_ROOT_ON]->(h:Host)
WITH u.username AS username, collect(h.hostname) AS hosts, count(h) AS cnt
WHERE cnt > 3
RETURN username, hosts, cnt
ORDER BY cnt DESC
```

Find Docker containers and their networks:

```cypher
MATCH (h:Host)-[:HAS_DOCKER_CONTAINER]->(c:DockerContainer)-[:USE_DOCKER_NETWORKS]->(n:DockerNetwork)
RETURN h.hostname, c.name, c.image, n.name AS network
ORDER BY h.hostname
```

---

## Snapshots

`facter-grpc` can take daily snapshots of the graph. A snapshot creates a copy of the current graph state with a timestamp, enabling:

- Trend analysis ("how many vulnerable packages did we have last week?")
- Diff between two points in time
- Historical compliance reports

Snapshots are purged automatically after a configurable number of days (`maxAgeDays`).
