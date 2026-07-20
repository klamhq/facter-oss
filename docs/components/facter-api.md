# facter-api

## Overview

`facter-api` is the read-only query gateway that exposes the infrastructure graph to `facter-ui` via GraphQL.

---

## Endpoints

| Method | Path       | Auth         | Description   |
| ------ | ---------- | ------------ | ------------- |
| `GET`  | `/health`  | None         | Health check  |
| `GET`  | `/graphql` | JWT required | GraphQL query |
| `POST` | `/graphql` | JWT required | GraphQL query |

All GraphQL endpoints are **read-only**. No mutations are exposed.

---

## Authentication

`facter-api` validates JWT tokens issued by Keycloak.

1. On startup, it fetches the JWKS from Keycloak's well-known endpoint
2. Each request must carry `Authorization: Bearer <access_token>`
3. The token is validated (signature, expiry, issuer, audience)
4. Roles are extracted from the `resource_access.<clientId>.roles` claim

**Required roles:** `facter-admin` or `facter-viewer`

---

## GraphQL Schema

The GraphQL schema mirrors the property graph structure. Typical query capabilities:

- List all hosts with their properties
- Fetch a specific host's full inventory (packages, users, processes, services, SSH keys, Docker, network, compliance, vulnerabilities)
- Cross-host queries (e.g. all hosts running a specific package version)
- Vulnerability and compliance summaries

---

## Connection to Apache AGE

`facter-api` connects directly to Apache AGE (PostgreSQL) as a graph query engine. It uses connection pooling configured via `graphDatabase.age`.

!!! tip
    For production, configure a **read-only** PostgreSQL user for `facter-api` to limit blast radius.

---

## Build & Run

```bash
cd facter-api
go build -o bin/facter-api .
./bin/facter-api --config configs/config-facter-api.yml
```

### With Docker

```bash
docker build -t facter-api:latest -f docker/Dockerfile .
docker run -v $(pwd)/configs:/configs \
  facter-api:latest --config /configs/config-facter-api.yml
```

---

## CORS

Configure allowed origins to match your `facter-ui` deployment URL:

```yaml
facterApi:
  configuration:
    api:
      cors:
        allowedOrigins:
          - "https://facter.example.com"
        allowedMethods:
          - "GET"
          - "POST"
          - "OPTIONS"
        allowedHeaders:
          - "Authorization"
          - "Content-Type"
```
