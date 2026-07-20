# facter-oss

[![Go Report Card](https://goreportcard.com/badge/github.com/klamhq/facter-oss)](https://goreportcard.com/report/github.com/klamhq/facter-oss)
[![CI](https://github.com/klamhq/facter-oss/actions/workflows/go.yml/badge.svg)](https://github.com/klamhq/facter-oss/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Lightweight Go agent that collects system inventory facts (packages, processes, users, networks, SSH keys, Docker containers, etc.) and exports them as protobuf messages — either to a local file or to [facter-grpc](https://github.com/klamhq/facter-grpc) via mTLS gRPC.

📖 **[Full documentation](https://klamhq.github.io/facter)**

---

## Architecture

```
facter-oss (agent)
    ├── local file output  →  export.iya  (protobuf)
    └── remote output      →  [gRPC / mTLS]  →  facter-grpc
```

## Tech stack

| Component   | Technology                                                                 |
| ----------- | -------------------------------------------------------------------------- |
| Language    | Go 1.25                                                                    |
| Transport   | gRPC + protobuf ([facter-schema](https://github.com/klamhq/facter-schema)) |
| Local store | bbolt                                                                      |
| Config      | YAML (Viper) + Cobra CLI                                                   |

---

## Prerequisites

- Go >= 1.25
- `make`
- (optional) Docker — for integration tests

---

## Quick start

```shell
# 1. Build (macOS)
make buildMac

# 2. Run
./bin/facter-oss --config configs/config-facter.yml
```

### Remote output (send to facter-grpc)

Set `sink.output.type: remote` in config and provide the gRPC server address and certificates:

```yaml
facter:
  sink:
    output:
      type: remote
      facterServer:
        serverHost: localhost
        serverPort: "56230"
        caPath: "./certs/facter_ca_cert.pem"
        certificatePath: "./certs/client_cert.pem"
        certificateKeyPath: "./certs/client_key.pem"
        sslHostname: "test.facter.fr"
        insecureSkipTlsVerify: false
```

---

## Configuration

Reference config: [`configs/config-facter.yml`](configs/config-facter.yml)

Key parameters:

| Key                       | Description                       |
| ------------------------- | --------------------------------- |
| `facter.sink.output.type` | `file` (local) or `remote` (gRPC) |
| `facter.inventory.*`      | Enable/disable each collector     |
| `facter.store.path`       | Path to local bbolt database      |

---

## Tests

```shell
make test
make integration-test
```

---

## Build

```shell
make buildMac   # macOS
make build      # Linux amd64
make release    # Linux + UPX compression
```

---

## Docker

```shell
make dockerBuild
```

---

## Security (mTLS / SPIFFE)

When using remote output, the gRPC connection uses mutual TLS. Authorization is based on SPIFFE identities — not CN. See [Security (docs)](https://klamhq.github.io/facter-oss/concepts/security/) for details.

The `insecureSkipTlsVerify` flag disables certificate verification. Use only in controlled dev/lab environments.

---

## Local development with facter-schema

To use a local `facter-schema` checkout, add to `go.mod`:

```
replace github.com/klamhq/facter-schema => /path/to/facter-schema
```

> Remove before committing.

