# facter
![Coverage](https://img.shields.io/badge/Coverage-67.3%25-yellow)

## Presentation

## Components

### Facter

[Doc Website](https://klamhq.github.io/facter-oss)

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

