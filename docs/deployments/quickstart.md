# Quick Start

This guide deploys the entire Facter stack locally using Docker Compose in under 10 minutes.

## Prerequisites

- Docker ≥ 24 and Docker Compose V2
- Git
- 4 GB RAM available for Docker

## Step 1 — Clone the repositories

```bash
git clone https://github.com/klamhq/facter-grpc
git clone https://github.com/klamhq/facter-rule-engine
git clone https://github.com/klamhq/facter-api
git clone https://github.com/klamhq/facter-ui
git clone https://github.com/klamhq/facter-oss
```

## Step 2 — Generate TLS Certificates

```bash
cd facter-grpc/scripts/certs
./create.sh
cd ../../..
```

Copy the certificates to the rule-engine certs directory:

```bash
cp facter-grpc/scripts/certs/ca_cert.pem facter-rule-engine/certs/
cp facter-grpc/scripts/certs/facter_rule_engine_cert.pem facter-rule-engine/certs/
cp facter-grpc/scripts/certs/facter_rule_engine_key.pem facter-rule-engine/certs/
```

## Step 3 — Start the Infrastructure

Start the databases and Keycloak:

```bash
cd facter-grpc
docker compose up -d age keycloak
cd ../facter-rule-engine
docker compose up -d postgresql
```

Wait ~30 seconds for services to start, then verify:

```bash
docker compose ps
```

## Step 4 — Initialise the Rule Engine Database

```bash
cd facter-rule-engine
./scripts/init_db.sh
```

## Step 5 — Start the Services

In separate terminals (or as background services):

```bash
# facter-grpc
cd facter-grpc
go run main.go --config configs/config-facter-grpc.yml
```

```bash
# facter-api
cd facter-api
go run main.go --config configs/config-facter-api.yml
```

```bash
# facter-rule-engine
cd facter-rule-engine
go run main.go serve --config configs/config-facter-re.yml
```

```bash
# facter-ui (dev mode)
cd facter-ui
npm install && npm run dev
```

## Step 6 — Run the Agent

On this machine (or any Linux host):

```bash
cd facter-oss
go run main.go --config configs/config-facter-local-mac.yml
```

!!! note "Local mode"
    `config-facter-local-mac.yml` uses `type: remote` pointing to `host.docker.internal:56230`. Adjust `serverHost` to the actual IP if running on a different machine.

## Step 7 — Open the UI

Navigate to [http://localhost:9000](http://localhost:9000) and log in with your Keycloak credentials.

!!! tip
    The default Keycloak realm `facter-audit` is automatically imported from `facter-ui/keycloak/realm/`. Default credentials depend on your Keycloak setup.

---

## Verify Everything Works

```bash
# Check facter-grpc is up
curl -s http://localhost:56230/health    # (or use grpcurl)

# Check facter-api is up
curl -s http://localhost:56231/health

# Check facter-rule-engine is up
curl -s http://localhost:8081/health
```

---

## Troubleshooting

**Agent can't connect to facter-grpc**  
→ Check certificate paths and `sslHostname` match the server certificate CN/SAN.

**UI shows "Unauthorized"**  
→ Ensure Keycloak is running and the realm/client config matches `facter-api` and `facter-rule-engine`.

**No hosts appear in the UI**  
→ Verify the agent ran successfully (`facter-oss` exit code 0) and Apache AGE is populated:
```sql
SELECT * FROM ag_catalog.cypher('facter', $$ MATCH (h:Host) RETURN h.hostname $$) AS (hostname agtype);
```
