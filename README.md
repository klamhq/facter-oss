# facter

## Presentation

## Components

### Facter

// TODO doc on facter

## Performance Profiling 

Enabled performance profiling in `config.yaml`:

```yaml
facter:
  performanceProfiling:
    enabled: true
```

```shell
apt install graphviz
make profile
```

## Release

```
make compress
```

Binary is written is `bin` folder and compressed by UPX.

