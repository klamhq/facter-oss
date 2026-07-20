# facter-schema

## Overview

`facter-schema` is a shared Go module that contains the [Protocol Buffers](https://protobuf.dev/) definitions for all gRPC communication between `facter-oss` and `facter-grpc`.

**Repository:** `github.com/klamhq/facter-schema`

---

## gRPC Service

```proto
service FactGrpcService {
  rpc Inventory(InventoryRequest) returns (InventoryResponse);
  rpc CheckRules(CheckRulesRequest) returns (CheckRulesResponse);
}
```

### `Inventory(InventoryRequest)`

Sent by `facter-oss` to push a host's collected facts to `facter-grpc`.

The request is a `oneof` that accepts either:

- `HostInventory` — full inventory snapshot
- `HostDeltaInventory` — only the changes since last run

### `CheckRules(CheckRulesRequest)`

Sent by `facter-rule-engine` to execute a Cypher query on the graph via `facter-grpc`.

---

## Key Message Types

### `HostInventory`

| Field                  | Type                | Description                       |
| ---------------------- | ------------------- | --------------------------------- |
| `hostname`             | string              | Host name                         |
| `packages`             | Package[]           | Installed packages                |
| `network`              | Network             | Network info                      |
| `platform`             | Platform            | OS, kernel, hardware              |
| `users`                | User[]              | Local users                       |
| `ssh_key_info`         | SshKeyInfo[]        | SSH keys                          |
| `ssh_key_access`       | SshKeyAccess[]      | authorized_keys entries           |
| `application`          | Application[]       | Docker apps                       |
| `systemd_service`      | SystemdService[]    | Systemd units                     |
| `known_host`           | KnownHost[]         | known_hosts entries               |
| `processes`            | Process[]           | Running processes                 |
| `vulnerability_report` | VulnerabilityReport | CVE findings                      |
| `compliance_report`    | ComplianceReport    | OpenSCAP results                  |
| `identifier`           | Identifier          | machine-id + UUID                 |
| `metadata`             | Metadata            | Run context (user, date, version) |
| `created_at`           | string              | Inventory timestamp               |

### `HostDeltaInventory`

Carries only the _changes_ since the previous full inventory:
`packages_added`, `packages_removed`, `users_added`, `users_removed`, `applications_added`, `applications_removed`, `systemdservices_added`, `systemdservices_removed`, `knownhosts_added`, `knownhosts_removed`, `sshkeyaccess_added`, `sshkeyaccess_removed`, `sshkeyinfo_added`, `sshkeyinfo_removed`, plus updated `platform` and `network`.

---

## Using facter-schema in Go

Add the module to your `go.mod`:

```bash
go get github.com/klamhq/facter-schema
```

Import in Go code:

```go
import facterv1 "github.com/klamhq/facter-schema/proto/rpc/facter/v1"
```

### Local Development Override

If you are working on both `facter-grpc` and `facter-schema` simultaneously, redirect the import to your local copy by adding to `go.mod`:

```go
replace github.com/klamhq/facter-schema => /path/to/local/facter-schema
```

!!! warning
    Remove this `replace` directive before committing and pushing.

---

## Generating Go Code

The generated `.pb.go` and `_grpc.pb.go` files are committed to the repository. To regenerate after modifying the `.proto` file:

```bash
cd facter-schema
make generate       # runs buf generate
```

Requires [buf](https://buf.build/docs/installation) and the Go protoc plugins.
