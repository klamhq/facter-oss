# facter-oss Configuration

## Full Reference

```yaml
facter:
  enabled: true                     # Master switch

  store:
    path: "/tmp/facter-store.db"    # Local SQLite store for delta computation

  logs:
    debugMode: false                # Verbose logging

  performanceProfiling:
    enabled: false                  # Enable Go pprof

  sink:
    output:
      format: "proto"               # Serialisation format (proto only)
      type: "file"                  # "file" (local) or "remote" (gRPC)

      # Used when type = "file"
      outputDirectory: "/tmp"
      outputFilename: "export.iya"

      # Used when type = "remote"
      facterServer:
        serverHost: "grpc.example.com"
        serverPort: "56230"
        certificatePath: "./certs/facter_cert.pem"
        certificateKeyPath: "./certs/facter_key.pem"
        caPath: "./certs/ca_cert.pem"
        sslHostname: "grpc.example.com"   # TLS SNI hostname

  inventory:

    applications:
      enabled: true
      docker:
        enabled: true               # Collect Docker info

    systemdService:
      enabled: true                 # Collect systemd units

    platform:
      enabled: true
      system:
        initCheckPath: "/proc/1/exe"          # Used to detect init system (Linux)
        # initCheckPath: "/sbin/init"         # Use this on macOS / older systems
        machineID: "/etc/machine-id"
        machineUUID: "/sys/class/dmi/id/product_uuid"
      hardware:
        enabled: true                         # CPU, RAM, Disk
      os:
        enabled: true                         # OS name, version, family
      virtualization:
        enabled: true                         # VM/container detection
      kernel:
        enabled: true                         # Kernel version

    packages:
      enabled: true                           # Installed packages

    process:
      enabled: true                           # Running processes

    ssh:
      enabled: true                           # SSH keys and known hosts

    user:
      enabled: true
      passwdFile: "/etc/passwd"               # File to read users from

    networks:
      enabled: true
      geoIp:
        enabled: false                        # GeoIP lookup (requires API key)
        timeout: 10                           # HTTP timeout in seconds
        googleGeoUrl: "https://www.googleapis.com/geolocation/v1/geolocate"
        googleGeoApikey: ""                   # Google Geolocation API key
      ports:
        enabled: true                         # Open ports
      publicIp:
        enabled: true                         # External IP detection
        publicIpApiUrl: "https://ifconfig.me/"
        timeout: 15
      firewall:
        enabled: true                         # iptables rules (requires root)
      connections:
        enabled: true                         # Established TCP/UDP connections

  compliance:
    enabled: false                            # OpenSCAP compliance scan
    profile: "xccdf_org.ssgproject.content_profile_cis_server_l1"
    resultFile: "/tmp/openscap-results.xml"

  vulnerabilities:
    enabled: false                            # CVE vulnerability matching
```

## Parameter Reference

### `sink.output`

| Key               | Type   | Default      | Description                                   |
| ----------------- | ------ | ------------ | --------------------------------------------- |
| `type`            | string | `file`       | Output mode: `file` or `remote`               |
| `format`          | string | `proto`      | Serialisation format (only `proto` supported) |
| `outputDirectory` | string | `/tmp`       | Directory for local file output               |
| `outputFilename`  | string | `export.iya` | Filename for local file output                |

### `sink.output.facterServer`

| Key                  | Type   | Description                                             |
| -------------------- | ------ | ------------------------------------------------------- |
| `serverHost`         | string | `facter-grpc` hostname or IP                            |
| `serverPort`         | string | `facter-grpc` gRPC port (default: `56230`)              |
| `certificatePath`    | string | Path to client certificate                              |
| `certificateKeyPath` | string | Path to client private key                              |
| `caPath`             | string | Path to CA certificate                                  |
| `sslHostname`        | string | TLS SNI hostname (must match server certificate CN/SAN) |

### `store`

| Key    | Type   | Description                                              |
| ------ | ------ | -------------------------------------------------------- |
| `path` | string | Path to local SQLite database used for delta computation |

!!! info "Delta mode"
    The local store keeps hashes of the previous inventory to compute deltas. Delete the store file to force a full inventory on the next run.
