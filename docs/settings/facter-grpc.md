# facter-grpc Configuration

## Full Reference

```yaml
facterGrpc:
  performanceProfiling:
    enabled: false

  logs:
    debugMode: false

  configuration:

    policy:
      # SPIFFE identities allowed for facter-oss agents
      facterPrincipal:
        - "spiffe://facter.fr/ns/backend/sa/facter"
      # SPIFFE identities allowed for facter-rule-engine
      facterRuleEnginePrincipal:
        - "spiffe://facter.fr/ns/backend/sa/facter-rule-engine"

    import:
      enabled: false                          # Import from .iya files on startup
      filesToImports:
        - "./export.iya"

    grpc:
      enabled: true
      address: "0.0.0.0"
      port: "56230"
      caPath: "./scripts/certs/facter_ca_cert.pem"
      certificatePath: "./scripts/certs/server_cert.pem"
      certificateKeyPath: "./scripts/certs/server_key.pem"
      healthTimeout: 10s                      # gRPC health check timeout

    cron:
      enabled: true

      snapshot:
        enabled: true
        schedule: "@daily"                    # Cron expression or @daily/@hourly etc.
        purge:
          enabled: true
          schedule: "@daily"
          maxAgeDays: 7                       # Keep snapshots for 7 days

      deletion:
        enabled: true
        purge:
          enabled: true
          schedule: "*/10 * * * *"            # Run every 10 minutes
          maxAgeDays: 7                       # Delete nodes older than 7 days

  graphDatabase:
    age:
      enabled: true
      host: "localhost"
      port: 5433                              # Apache AGE default port
      username: "ageUser"
      password: "age"
      sslmode: "disable"                      # disable | require | verify-full
      dbName: "age"
      graphName: "facter"                     # Name of the graph in AGE

      # Connection pool settings
      maxConns: 20
      minConns: 2
      maxConnLifetime: 30m
      maxConnIdleTime: 10m
      healthCheckPeriod: 1m
```

## Parameter Reference

### `configuration.grpc`

| Key                  | Type     | Default   | Description                       |
| -------------------- | -------- | --------- | --------------------------------- |
| `enabled`            | bool     | `true`    | Enable gRPC server                |
| `address`            | string   | `0.0.0.0` | Bind address                      |
| `port`               | string   | `56230`   | Listen port                       |
| `caPath`             | string   | —         | Path to CA certificate (for mTLS) |
| `certificatePath`    | string   | —         | Path to server certificate        |
| `certificateKeyPath` | string   | —         | Path to server private key        |
| `healthTimeout`      | duration | `10s`     | Health check request timeout      |

### `configuration.policy`

| Key                         | Type     | Description                                              |
| --------------------------- | -------- | -------------------------------------------------------- |
| `facterPrincipal`           | string[] | SPIFFE URIs allowed for `facter-oss` connections         |
| `facterRuleEnginePrincipal` | string[] | SPIFFE URIs allowed for `facter-rule-engine` connections |

### `configuration.cron.snapshot`

| Key                | Type   | Default  | Description                     |
| ------------------ | ------ | -------- | ------------------------------- |
| `enabled`          | bool   | `true`   | Enable scheduled snapshots      |
| `schedule`         | string | `@daily` | Cron schedule                   |
| `purge.enabled`    | bool   | `true`   | Enable automatic snapshot purge |
| `purge.schedule`   | string | `@daily` | Purge schedule                  |
| `purge.maxAgeDays` | int    | `7`      | Days to retain snapshots        |

### `graphDatabase.age`

| Key                 | Type     | Description                            |
| ------------------- | -------- | -------------------------------------- |
| `host`              | string   | Apache AGE host                        |
| `port`              | int      | Apache AGE port (typically `5433`)     |
| `username`          | string   | PostgreSQL username                    |
| `password`          | string   | PostgreSQL password                    |
| `sslmode`           | string   | `disable`, `require`, or `verify-full` |
| `dbName`            | string   | PostgreSQL database name               |
| `graphName`         | string   | Name of the AGE graph to use           |
| `maxConns`          | int      | Maximum connections in pool            |
| `minConns`          | int      | Minimum idle connections               |
| `maxConnLifetime`   | duration | Maximum connection lifetime            |
| `maxConnIdleTime`   | duration | Maximum connection idle time           |
| `healthCheckPeriod` | duration | Pool health check interval             |
