# Concepts

This section explains the core ideas and design decisions behind Facter.

## How Facter Works

Facter follows a **collect → store → query → audit** pipeline:

1. **Collect** — `facter-oss` agents run on each monitored host and collect facts about the system (packages, users, processes, network, SSH keys, containers, compliance, vulnerabilities).
2. **Store** — Facts are serialised as Protobuf messages and pushed to `facter-grpc` via gRPC with mTLS. The server stores them as a property graph in Apache AGE.
3. **Query** — `facter-api` exposes the graph through a GraphQL endpoint. `facter-ui` uses it to display dashboards and per-host drill-downs.
4. **Audit** — `facter-rule-engine` executes Cypher audit rules against the graph via gRPC, stores results in PostgreSQL, and exposes them via REST API. A rule that returns rows produces a finding; a rule that returns no rows is compliant.

## Core Concepts

| Concept             | Description                                                                             |
| ------------------- | --------------------------------------------------------------------------------------- |
| **Inventory**       | A snapshot of all facts collected from a host at a point in time                        |
| **Delta Inventory** | Only the changes since the last snapshot (added/removed packages, users, etc.)          |
| **Graph**           | The property graph in Apache AGE where all inventory data is stored as nodes and edges  |
| **Cypher**          | The query language used to navigate the graph                                           |
| **Rule**            | A YAML document defining an audit check: a Cypher query to detect a condition           |
| **Execution**       | A run of a rule against the graph, producing findings or a compliant result             |
| **Finding**         | A rule execution that returned matched results (potential security or compliance issue) |

## Further Reading

- [Architecture](architecture.md) — how components connect
- [Data Model](data-model.md) — the graph schema
- [Rule Engine](rule-engine.md) — how audit rules work
- [Security](security.md) — authentication and mTLS
