# Configuration Reference

All Facter components are configured via a YAML file passed with the `--config` flag.

## Components

| Component                                   | Config file              | Default location                   |
| ------------------------------------------- | ------------------------ | ---------------------------------- |
| [facter-oss](facter-oss.md)                 | `config-facter.yml`      | `./configs/config-facter.yml`      |
| [facter-grpc](facter-grpc.md)               | `config-facter-grpc.yml` | `./configs/config-facter-grpc.yml` |
| [facter-api](facter-api.md)                 | `config-facter-api.yml`  | `./configs/config-facter-api.yml`  |
| [facter-rule-engine](facter-rule-engine.md) | `config-facter-re.yml`   | `./configs/config-facter-re.yml`   |

## Common Patterns

### Debug Mode

Enable verbose logging in any component:

```yaml
logs:
  debugMode: true
```

### Performance Profiling

Enable Go pprof:

```yaml
performanceProfiling:
  enabled: true
```

### Connection Pooling (Apache AGE / PostgreSQL)

```yaml
graphDatabase:
  age:
    maxConns: 20
    minConns: 2
    maxConnLifetime: 30m
    maxConnIdleTime: 10m
    healthCheckPeriod: 1m
```
