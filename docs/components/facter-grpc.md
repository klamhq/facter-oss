# facter-grpc

## Overview

`facter-grpc` is the central ingest server. It receives `HostInventory` or `HostDeltaInventory` Protobuf messages from `facter-oss` agents and persists them as a property graph in Apache AGE.

It also serves as the Cypher execution gateway for `facter-rule-engine`.

---

## gRPC Service

The service implements two RPC methods defined in [facter-schema](facter-schema.md):

```proto
service FactGrpcService {
  rpc Inventory(InventoryRequest) returns (InventoryResponse);
  rpc CheckRules(CheckRulesRequest) returns (CheckRulesResponse);
}
```

### `Inventory`

Accepts either a full `HostInventory` or a `HostDeltaInventory`. The server:

1. Extracts the host identifier (machine-id / UUID)
2. If **full inventory**: upserts all nodes and relationships in the graph
3. If **delta inventory**: applies only the added/removed items since the last snapshot

### `CheckRules`

Accepts a Cypher query string (sent by `facter-rule-engine`) and executes it directly against Apache AGE. Returns raw result rows as JSON.

---

## Authentication

Connections are authenticated via **mTLS + SPIFFE**. The server validates client certificates against the configured CA and checks that the SPIFFE identity in the certificate SAN matches the `facterPrincipal` or `facterRuleEnginePrincipal` allowlist.

See [Security](../concepts/security.md) for certificate setup.

---

## Graph Database

`facter-grpc` connects to **Apache AGE**, a PostgreSQL extension that adds openCypher graph support.

Connection is configured via the `graphDatabase.age` block (see [Configuration](../settings/facter-grpc.md)).

---

## Scheduled Tasks

### Snapshots

Creates a point-in-time copy of the graph. Enables trend analysis and historical comparisons.

| Setting                          | Default  | Description               |
| -------------------------------- | -------- | ------------------------- |
| `cron.snapshot.enabled`          | `true`   | Enable snapshots          |
| `cron.snapshot.schedule`         | `@daily` | Cron schedule             |
| `cron.snapshot.purge.enabled`    | `true`   | Auto-delete old snapshots |
| `cron.snapshot.purge.schedule`   | `@daily` | Purge schedule            |
| `cron.snapshot.purge.maxAgeDays` | `7`      | Maximum snapshot age      |

### Hard Deletions

Purges stale nodes from the graph (hosts that have not sent inventory for a while, orphaned nodes).

| Setting                          | Default        | Description                      |
| -------------------------------- | -------------- | -------------------------------- |
| `cron.deletion.enabled`          | `true`         | Enable hard deletions            |
| `cron.deletion.purge.enabled`    | `true`         | Enable purge                     |
| `cron.deletion.purge.schedule`   | `*/10 * * * *` | Purge schedule                   |
| `cron.deletion.purge.maxAgeDays` | `7`            | Maximum node age before deletion |

---

## File Import

`facter-grpc` can import inventory from `.iya` (Protobuf binary) files collected by `facter-oss` in local mode.

```yaml
configuration:
  import:
    enabled: true
    filesToImports:
      - "/path/to/export.iya"
```

This is useful for air-gapped environments where agents cannot reach the server directly.

---

## Build & Run

```bash
cd facter-grpc
go build -o bin/facter-grpc .
./bin/facter-grpc --config configs/config-facter-grpc.yml
```

### With Docker

```bash
docker build -t facter-grpc:latest -f docker/Dockerfile .
docker run -v $(pwd)/configs:/configs -v $(pwd)/certs:/certs \
  facter-grpc:latest --config /configs/config-facter-grpc.yml
```

---

## Health Check

`facter-grpc` also exposes a gRPC health check endpoint (standard [gRPC Health Checking Protocol](https://grpc.io/docs/guides/health-checking/)).

| Setting              | Default |
| -------------------- | ------- |
| `grpc.healthTimeout` | `10s`   |
