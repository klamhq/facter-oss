# facter-rule-engine

## Overview

`facter-rule-engine` is the compliance and audit engine. It manages audit rules (stored in PostgreSQL), launches Cypher queries via `facter-grpc`, evaluates results, and exposes findings through a REST API.

---

## REST API Reference

### Rules

| Method   | Path                | Role         | Description                      |
| -------- | ------------------- | ------------ | -------------------------------- |
| `GET`    | `/api/v1/rules`     | viewer/admin | List all enabled rules           |
| `POST`   | `/api/v1/rules`     | admin        | Create one or more rules (array) |
| `PUT`    | `/api/v1/rules/:id` | admin        | Update an existing rule          |
| `DELETE` | `/api/v1/rules/:id` | admin        | Delete a rule                    |

### Executions

| Method | Path                               | Role         | Description             |
| ------ | ---------------------------------- | ------------ | ----------------------- |
| `GET`  | `/api/v1/rule_executions`          | viewer/admin | All execution records   |
| `GET`  | `/api/v1/rule_executions/:rule_id` | viewer/admin | Executions for one rule |

### Audits

| Method | Path                      | Role  | Description                          |
| ------ | ------------------------- | ----- | ------------------------------------ |
| `POST` | `/api/v1/audits`          | admin | Run global audit (all enabled rules) |
| `POST` | `/api/v1/audits/:rule_id` | admin | Run audit for a single rule          |

### System

| Method | Path      | Auth | Description                      |
| ------ | --------- | ---- | -------------------------------- |
| `GET`  | `/health` | None | Service health + DB connectivity |

---

## Rule Validation

When creating or updating rules via the API, the body is validated before insertion:

- `id` must match `R-[A-Z]{3}-[0-9]{4}`
- `cypher_query` is required and must contain both `MATCH` and `RETURN`
- `name` and `description` are required
- Duplicate `id` on creation returns HTTP 409

### Idempotent Updates

Rules are updated using a **checksum** mechanism. If the rule content has not changed since the last update, the API returns `status: unchanged` (HTTP 200) rather than writing a no-op to the database.

---

## System Rules Protection

Rules with IDs starting with `R-SYS-` are **system rules**. They are:

- Loaded from the `rules/` directory on startup
- Read-only via the API (DELETE returns HTTP 403)
- Not overwritable via PUT

---

## gRPC Connection to facter-grpc

`facter-rule-engine` connects to `facter-grpc` using mTLS:

```yaml
facterRuleEngine:
  facterGrpc:
    serverHost: "grpc.example.com"
    serverPort: "56230"
    certificatePath: "./certs/facter_rule_engine_cert.pem"
    certificateKeyPath: "./certs/facter_rule_engine_key.pem"
    caPath: "./certs/ca_cert.pem"
    sslHostname: "grpc.example.com"
    healthCheckInterval: 30s
```

A background goroutine monitors the gRPC connection health every `healthCheckInterval`.

---

## Validate Command

The `validate` sub-command checks YAML rule files without starting the server:

```bash
facter-rule-engine validate --rules-dir ./rules/
```

Useful in CI pipelines to validate rule files before deployment.

---

## Database Migrations

Database schema is managed by SQL migration files in `db/migrations/`. Run migrations before first start:

```bash
./scripts/init_db.sh      # creates DB and runs migrations
```

---

## Build & Run

```bash
cd facter-rule-engine
go build -o bin/facter-rule-engine .
./bin/facter-rule-engine --config configs/config-facter-re.yml
```

### With Docker

```bash
docker build -t facter-rule-engine:latest -f docker/Dockerfile .
docker run -v $(pwd)/configs:/configs -v $(pwd)/certs:/certs \
  facter-rule-engine:latest --config /configs/config-facter-re.yml
```
