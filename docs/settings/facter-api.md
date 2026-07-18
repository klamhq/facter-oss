# facter-api Configuration

## Full Reference

```yaml
facterApi:
  performanceProfiling:
    enabled: false

  logs:
    debugMode: false

  configuration:
    api:
      ssl: false                              # Enable TLS on HTTP server
      address: "0.0.0.0"
      port: "56231"

      # TLS (leave empty if ssl: false)
      caPath: ""
      certificatePath: ""
      certificateKeyPath: ""
      sslVerify: false

      # HTTP timeouts
      readTimeout: 15s
      writeTimeout: 15s
      idleTimeout: 60s

      cors:
        allowedOrigins:
          - "http://localhost:9000"           # In production: your UI domain
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
        url: "http://localhost:8080"          # Keycloak base URL
        realm: "facter-audit"                 # Keycloak realm name
        clientId: "facter-ui"                 # Client ID in Keycloak

  graphDatabase:
    age:
      enabled: true
      host: "localhost"
      port: 5433
      username: "ageUser"
      password: "age"
      sslmode: "disable"
      dbName: "age"
      graphName: "facter"
      maxConns: 20
      minConns: 2
      maxConnLifetime: 30m
      maxConnIdleTime: 10m
      healthCheckPeriod: 1m
```

## Parameter Reference

### `configuration.api`

| Key                  | Type     | Default   | Description                           |
| -------------------- | -------- | --------- | ------------------------------------- |
| `ssl`                | bool     | `false`   | Enable TLS                            |
| `address`            | string   | `0.0.0.0` | Bind address                          |
| `port`               | string   | `56231`   | Listen port                           |
| `certificatePath`    | string   | —         | Server certificate (when `ssl: true`) |
| `certificateKeyPath` | string   | —         | Server private key (when `ssl: true`) |
| `caPath`             | string   | —         | CA certificate                        |
| `readTimeout`        | duration | `15s`     | HTTP read timeout                     |
| `writeTimeout`       | duration | `15s`     | HTTP write timeout                    |
| `idleTimeout`        | duration | `60s`     | Keep-alive idle timeout               |

### `configuration.api.keycloak`

| Key        | Type   | Description                                         |
| ---------- | ------ | --------------------------------------------------- |
| `url`      | string | Base URL of your Keycloak instance                  |
| `realm`    | string | Name of the Keycloak realm                          |
| `clientId` | string | Client ID registered in Keycloak for JWT validation |

### `configuration.api.cors`

| Key              | Type     | Description                                             |
| ---------------- | -------- | ------------------------------------------------------- |
| `allowedOrigins` | string[] | Allowed CORS origins (use your UI domain in production) |
| `allowedMethods` | string[] | Allowed HTTP methods                                    |
| `allowedHeaders` | string[] | Allowed request headers                                 |
