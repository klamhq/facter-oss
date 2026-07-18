# Docker Deployment

This page covers a production-grade Docker Compose deployment of the full Facter stack.

## Prerequisites

- Docker ≥ 24 and Docker Compose V2
- TLS certificates (see [Security — Generating Certificates](../concepts/security.md#generating-test-certificates))
- A Keycloak instance (can be included in the compose file)

---

## Directory Structure

```
facter/
├── certs/
│   ├── ca_cert.pem
│   ├── server_cert.pem
│   ├── server_key.pem
│   ├── facter_cert.pem
│   ├── facter_key.pem
│   ├── facter_rule_engine_cert.pem
│   └── facter_rule_engine_key.pem
├── configs/
│   ├── config-facter-grpc.yml
│   ├── config-facter-api.yml
│   └── config-facter-re.yml
└── docker-compose.yml
```

---

## docker-compose.yml

```yaml
services:

  # ── Databases ──────────────────────────────────────────────────────────────

  age:
    image: apache/age:release_PG18_1.7.0
    restart: unless-stopped
    environment:
      POSTGRES_PASSWORD: ${AGE_PASSWORD:-changeMe}
      POSTGRES_USER: ageUser
      POSTGRES_DB: age
    ports:
      - "127.0.0.1:5433:5432"
    volumes:
      - age-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ageUser -d age"]
      interval: 10s
      timeout: 5s
      retries: 5

  postgresql:
    image: postgres:18
    restart: unless-stopped
    environment:
      POSTGRES_USER: facter_rule_engine
      POSTGRES_PASSWORD: ${PG_PASSWORD:-changeMe}
      POSTGRES_DB: rule_engine
    ports:
      - "127.0.0.1:5432:5432"
    volumes:
      - pg-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U facter_rule_engine -d rule_engine"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ── Identity Provider ──────────────────────────────────────────────────────

  keycloak:
    image: quay.io/keycloak/keycloak:26
    restart: unless-stopped
    command: start-dev --import-realm
    environment:
      KC_BOOTSTRAP_ADMIN_USERNAME: admin
      KC_BOOTSTRAP_ADMIN_PASSWORD: ${KC_ADMIN_PASSWORD:-admin}
    ports:
      - "127.0.0.1:8080:8080"
    volumes:
      - ./keycloak/realm:/opt/keycloak/data/import

  # ── Backend Services ───────────────────────────────────────────────────────

  facter-grpc:
    image: ghcr.io/klamhq/facter-grpc:latest
    restart: unless-stopped
    depends_on:
      age:
        condition: service_healthy
    ports:
      - "56230:56230"
    volumes:
      - ./configs/config-facter-grpc.yml:/configs/config.yml:ro
      - ./certs:/certs:ro
    command: ["--config", "/configs/config.yml"]

  facter-api:
    image: ghcr.io/klamhq/facter-api:latest
    restart: unless-stopped
    depends_on:
      age:
        condition: service_healthy
      keycloak:
        condition: service_started
    ports:
      - "127.0.0.1:56231:56231"
    volumes:
      - ./configs/config-facter-api.yml:/configs/config.yml:ro
    command: ["--config", "/configs/config.yml"]

  facter-rule-engine:
    image: ghcr.io/klamhq/facter-rule-engine:latest
    restart: unless-stopped
    depends_on:
      postgresql:
        condition: service_healthy
      facter-grpc:
        condition: service_started
      keycloak:
        condition: service_started
    ports:
      - "127.0.0.1:8081:8081"
    volumes:
      - ./configs/config-facter-re.yml:/configs/config.yml:ro
      - ./certs:/certs:ro
    command: ["serve", "--config", "/configs/config.yml"]

  facter-ui:
    image: ghcr.io/klamhq/facter-ui:latest
    restart: unless-stopped
    depends_on:
      - facter-api
      - facter-rule-engine
    ports:
      - "80:80"

volumes:
  age-data:
  pg-data:
```

---

## Configuration Files

### `configs/config-facter-grpc.yml`

```yaml
facterGrpc:
  logs:
    debugMode: false
  configuration:
    policy:
      facterPrincipal:
        - "spiffe://facter.fr/ns/backend/sa/facter"
      facterRuleEnginePrincipal:
        - "spiffe://facter.fr/ns/backend/sa/facter-rule-engine"
    grpc:
      enabled: true
      address: "0.0.0.0"
      port: "56230"
      caPath: "/certs/ca_cert.pem"
      certificatePath: "/certs/server_cert.pem"
      certificateKeyPath: "/certs/server_key.pem"
    cron:
      enabled: true
      snapshot:
        enabled: true
        schedule: "@daily"
        purge:
          enabled: true
          schedule: "@daily"
          maxAgeDays: 30
  graphDatabase:
    age:
      enabled: true
      host: "age"
      port: 5432
      username: "ageUser"
      password: "changeMe"
      dbName: "age"
      graphName: "facter"
      maxConns: 20
      minConns: 2
```

### `configs/config-facter-api.yml`

```yaml
facterApi:
  configuration:
    api:
      address: "0.0.0.0"
      port: "56231"
      cors:
        allowedOrigins:
          - "https://facter.example.com"
        allowedMethods: ["GET", "POST", "OPTIONS"]
        allowedHeaders: ["Authorization", "Content-Type"]
      keycloak:
        url: "http://keycloak:8080"
        realm: "facter-audit"
        clientId: "facter-ui"
  graphDatabase:
    age:
      enabled: true
      host: "age"
      port: 5432
      username: "ageUser"
      password: "changeMe"
      dbName: "age"
      graphName: "facter"
```

### `configs/config-facter-re.yml`

```yaml
facterRuleEngine:
  postgresql:
    uri: "postgres://facter_rule_engine:changeMe@postgresql:5432/rule_engine"
  facterApi:
    serverHost: "0.0.0.0"
    serverPort: "8081"
    cors:
      allowedOrigins:
        - "https://facter.example.com"
      allowedMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
      allowedHeaders: ["Authorization", "Content-Type"]
    keycloak:
      url: "http://keycloak:8080"
      realm: "facter-audit"
      clientId: "facter-ui"
  facterGrpc:
    serverHost: "facter-grpc"
    serverPort: "56230"
    certificatePath: "/certs/facter_rule_engine_cert.pem"
    certificateKeyPath: "/certs/facter_rule_engine_key.pem"
    caPath: "/certs/ca_cert.pem"
    sslHostname: "grpc.example.com"
    healthCheckInterval: 30s
```

---

## Starting the Stack

```bash
# Start infrastructure
docker compose up -d age postgresql keycloak

# Wait for DBs to be ready
docker compose ps

# Start services
docker compose up -d facter-grpc facter-api facter-rule-engine facter-ui

# View logs
docker compose logs -f
```

## Deploying facter-oss on Target Hosts

On each host to monitor:

```bash
curl -L https://github.com/klamhq/facter-oss/releases/latest/download/facter-linux-amd64 \
  -o /usr/local/bin/facter && chmod +x /usr/local/bin/facter

cat > /etc/facter/config.yml <<EOF
facter:
  sink:
    output:
      type: "remote"
      facterServer:
        serverHost: "YOUR_FACTER_GRPC_HOST"
        serverPort: "56230"
        certificatePath: "/etc/facter/certs/facter_cert.pem"
        certificateKeyPath: "/etc/facter/certs/facter_key.pem"
        caPath: "/etc/facter/certs/ca_cert.pem"
        sslHostname: "grpc.example.com"
  inventory:
    platform:
      enabled: true
    packages:
      enabled: true
    user:
      enabled: true
    networks:
      enabled: true
    ssh:
      enabled: true
    process:
      enabled: true
    applications:
      enabled: true
EOF

# Run once
facter --config /etc/facter/config.yml

# Or schedule with systemd timer
```
