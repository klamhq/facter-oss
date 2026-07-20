# facter-rule-engine Configuration

## Full Reference

```yaml
facterRuleEngine:
  performanceProfiling:
    enabled: false

  logs:
    debugMode: false

  postgresql:
    uri: "postgres://facter_rule_engine:password@localhost:5432/rule_engine"

  facterApi:
    serverHost: "localhost"
    serverPort: "8081"                        # Port this service listens on

    cors:
      allowedOrigins:
        - "http://localhost:9000"             # In production: your UI domain
      allowedMethods:
        - "GET"
        - "POST"
        - "PUT"
        - "DELETE"
        - "OPTIONS"
      allowedHeaders:
        - "Authorization"
        - "Content-Type"

    keycloak:
      url: "http://localhost:8080"
      realm: "facter-audit"
      clientId: "facter-ui"

  facterGrpc:
    serverHost: "localhost"
    serverPort: "56230"
    certificatePath: "./certs/facter_rule_engine_cert.pem"
    certificateKeyPath: "./certs/facter_rule_engine_key.pem"
    caPath: "./certs/ca_cert.pem"
    sslHostname: "grpc.example.com"           # TLS SNI hostname for facter-grpc (optional, defaults to serverHost)
    insecureSkipTlsVerify: false              # Disable TLS verification — never use in production
    healthCheckInterval: 30s                  # gRPC health monitoring interval
```

## Parameter Reference

### `postgresql`

| Key   | Type   | Description                                                            |
| ----- | ------ | ---------------------------------------------------------------------- |
| `uri` | string | Full PostgreSQL connection URI including credentials and database name |

Format: `postgres://user:password@host:port/dbname`

!!! warning "Credentials"
    Do not commit production credentials. Use environment variable substitution or a secrets manager.

### `facterApi` (HTTP server settings)

| Key          | Type   | Default     | Description  |
| ------------ | ------ | ----------- | ------------ |
| `serverHost` | string | `localhost` | Bind address |
| `serverPort` | string | `8081`      | Listen port  |

### `facterApi.cors`

| Key              | Type     | Description          |
| ---------------- | -------- | -------------------- |
| `allowedOrigins` | string[] | CORS allowed origins |
| `allowedMethods` | string[] | Allowed HTTP methods |
| `allowedHeaders` | string[] | Allowed headers      |

### `facterApi.keycloak`

| Key        | Type   | Description                            |
| ---------- | ------ | -------------------------------------- |
| `url`      | string | Keycloak base URL                      |
| `realm`    | string | Keycloak realm                         |
| `clientId` | string | Client ID for JWKS/audience validation |

### `facterGrpc`

| Key                     | Type     | Default        | Description                                                                     |
| ----------------------- | -------- | -------------- | ------------------------------------------------------------------------------- |
| `serverHost`            | string   |                | `facter-grpc` host                                                              |
| `serverPort`            | string   | `56230`        | `facter-grpc` port                                                              |
| `certificatePath`       | string   |                | Client certificate for mTLS                                                     |
| `certificateKeyPath`    | string   |                | Client private key for mTLS                                                     |
| `caPath`                | string   |                | CA certificate                                                                  |
| `sslHostname`           | string   | *(serverHost)* | TLS SNI hostname (must match server certificate CN/SAN). Defaults to serverHost |
| `insecureSkipTlsVerify` | bool     | `false`        | Skip TLS certificate verification. **Never enable in production.**              |
| `healthCheckInterval`   | duration |                | How often to check gRPC connection health                                       |
