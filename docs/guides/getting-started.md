# Getting Started

Ce guide vous accompagne lors de votre première utilisation de la plateforme Facter, du démarrage de la stack à la création et l'exécution de votre première règle d'audit.

## Prerequisites

- Docker Engine ≥ 24 et Docker Compose ≥ 2.20 installés
- `facter-oss` binaire compilé (ou image Docker)
- Accès aux dépôts `facter-grpc`, `facter-api`, `facter-rule-engine` et `facter-ui`

---

## Step 1 — Start the stack

Follow the [Quick Start](../deployments/quickstart.md) guide to bring up all services with Docker Compose.

Once all containers are healthy:

```bash
docker compose ps
```

You should see all services in the `running` state:

```
facter-grpc          running
facter-api           running
facter-rule-engine   running
facter-ui            running
keycloak             running
apache-age           running
postgresql           running
```

---

## Step 2 — Send your first inventory

### Option A — Run facter-oss on the local machine

Create a minimal configuration file `config-local.yml`:

```yaml
agent:
  hostname: "my-first-host"
  remote: true
  grpcServerAddress: "localhost:56230"
  tlsEnabled: true
  caCertFile: "/path/to/certs/ca_cert.pem"
  certFile: "/path/to/certs/facter_cert.pem"
  keyFile: "/path/to/certs/facter_key.pem"

collectors:
  packages:
    enabled: true
  users:
    enabled: true
  network:
    enabled: true
  processes:
    enabled: false
```

Run the agent:

```bash
facter-oss --config config-local.yml
```

### Option B — Use the local output mode (no gRPC)

```yaml
agent:
  hostname: "my-first-host"
  remote: false  # write to local file

collectors:
  packages:
    enabled: true
```

```bash
facter-oss --config config-local.yml
# Inventory saved to: my-first-host.iya
```

You can then import the file through the gRPC server using the file import endpoint (see [facter-grpc component](../components/facter-grpc.md)).

---

## Step 3 — Verify the inventory in the UI

1. Open the Facter UI at `http://localhost:5173`
2. Log in with your Keycloak credentials
3. Navigate to **Hosts**

Your host should appear with all collected inventory items: packages, users, network interfaces, etc.

Navigate to a host and explore the tabs:
- **Overview** — system summary
- **Packages** — installed software
- **Access** — users and SSH keys
- **Connections** — network interfaces and open ports

---

## Step 4 — Create your first audit rule

In the UI, navigate to **Rule Engine > Rules** and click **New Rule**.

Or call the REST API directly:

```bash
curl -X POST http://localhost:8081/api/v1/rules \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "R-PKG-0001",
    "name": "Curl must be installed",
    "description": "Ensures that curl is present on all monitored hosts.",
    "category": "INFRASTRUCTURE",
    "severity": "MEDIUM",
    "status": "approved",
    "enabled": true,
    "version": 1,
    "priority": "P3",
    "cypher_query": "MATCH (h:Host)-[:HAS_PACKAGE]-(p:Package {name: \"curl\"}) RETURN p AS seed, h.hostname AS hostname, p.name AS package, p.version AS version"
  }'
```

**Compliance logic**: the query returns rows only when `curl` is found. If it returns 0 rows, the rule produces a finding (curl is absent). If it returns rows, it is compliant.

### Key rule fields

| Field          | Description                                                                                     |
| -------------- | ----------------------------------------------------------------------------------------------- |
| `id`           | Unique rule identifier, format `R-[A-Z]{3}-[0-9]{4}`                                            |
| `cypher_query` | openCypher query — rows returned = findings, no rows = compliant                                |
| `category`     | One of `IDENTITY_ACCESS`, `NETWORK_SECURITY`, `DATA_PROTECTION`, `COMPLIANCE`, `INFRASTRUCTURE` |
| `severity`     | `CRITICAL`, `HIGH`, `MEDIUM`, or `LOW`                                                          |
| `status`       | `draft`, `review`, `approved`, or `deprecated`                                                  |

See the [Writing Audit Rules](writing-rules.md) guide for all fields.

---

## Step 5 — Run your first audit

### Run a single rule

```bash
curl -X POST http://localhost:8081/api/v1/audits/R-PKG-0001 \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Run all rules (global audit)

```bash
curl -X POST http://localhost:8081/api/v1/audits \
  -H "Authorization: Bearer <your-jwt-token>"
```

---

## Step 6 — Review results in the UI

1. Navigate to **Rule Engine > Executions** in the UI
2. Click on a rule execution to see the detailed results:
   - **Matched** results → findings (non-compliant)
   - **Compliant** results → hosts where the rule passed
3. Use the **Global Report** view for a cross-rule summary with charts

---

## Next Steps

- [Writing Audit Rules](writing-rules.md) — learn the full YAML rule format
- [Cypher Query Guide](cypher-queries.md) — write more sophisticated queries
- [Configuration Reference](../settings/index.md) — tune each component
