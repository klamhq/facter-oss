# facter-oss (Agent)

## Overview

`facter-oss` is the lightweight agent that runs on each Linux host to be monitored. It collects system facts and either writes them to a local file or streams them to `facter-grpc` via gRPC with mTLS.

## Installation

### Pre-built Binary

Download the latest binary from the [releases page](https://github.com/klamhq/facter-oss/releases):

```bash
curl -L https://github.com/klamhq/facter-oss/releases/latest/download/facter-linux-amd64 -o /usr/local/bin/facter
chmod +x /usr/local/bin/facter
```

### Build from Source

```bash
git clone https://github.com/klamhq/facter-oss
cd facter-oss
make build
# Binary written to bin/facter-oss
```

For a compressed production binary (requires UPX):

```bash
make compress
```

---

## Usage

### Basic Command

```bash
facter --config /etc/facter/config.yml
```

### Flags

| Flag       | Default | Description                                    |
| ---------- | ------- | ---------------------------------------------- |
| `--config` | —       | Path to the YAML configuration file (required) |

---

## Output Modes

### Local Mode (`type: file`)

Facts are serialised as a Protobuf binary and written to a `.iya` file. Useful for offline analysis or debugging.

```yaml
facter:
  sink:
    output:
      type: "file"
      format: "proto"
      outputDirectory: "/tmp"
      outputFilename: "export.iya"
```

### Remote Mode (`type: remote`)

Facts are streamed directly to `facter-grpc` over gRPC with mTLS.

```yaml
facter:
  sink:
    output:
      type: "remote"
      facterServer:
        serverHost: "grpc.example.com"
        serverPort: "56230"
        certificatePath: "/etc/facter/certs/facter_cert.pem"
        certificateKeyPath: "/etc/facter/certs/facter_key.pem"
        caPath: "/etc/facter/certs/ca_cert.pem"
        sslHostname: "grpc.example.com"
```

---

## Collectors

### Platform

Collects OS, kernel, hardware, and virtualisation information.

```yaml
facter:
  inventory:
    platform:
      enabled: true
      system:
        initCheckPath: "/proc/1/exe"       # Path to detect init system
        machineID: "/etc/machine-id"        # Machine unique ID file
        machineUUID: "/sys/class/dmi/id/product_uuid"
      hardware:
        enabled: true
      os:
        enabled: true
      virtualization:
        enabled: true
      kernel:
        enabled: true
```

### Packages

Detects the package manager automatically and lists all installed packages.

Supported managers: `apt` (Debian/Ubuntu), `rpm` (RHEL/Rocky/AlmaLinux/Fedora), `pacman` (Arch), `brew` (macOS).

```yaml
facter:
  inventory:
    packages:
      enabled: true
```

### Users

Reads `/etc/passwd` to list local users and detects sudo capabilities.

```yaml
facter:
  inventory:
    user:
      enabled: true
      passwdFile: "/etc/passwd"
```

### Networks

Collects network interfaces, IP addresses, active connections, firewall rules, DNS settings, public IP, and optionally GeoIP.

```yaml
facter:
  inventory:
    networks:
      enabled: true
      geoIp:
        enabled: true
        timeout: 10
        googleGeoUrl: "https://www.googleapis.com/geolocation/v1/geolocate"
        googleGeoApikey: "<your-api-key>"
      ports:
        enabled: true
      publicIp:
        enabled: true
        publicIpApiUrl: "https://ifconfig.me/"
        timeout: 15
      firewall:
        enabled: true
      connections:
        enabled: true
```

!!! note "Root required"
    Firewall rules (`iptables`) and some connection states require root privileges.

### SSH

Scans for SSH keys in standard locations (`/root/.ssh`, `/home/*/.ssh`).

```yaml
facter:
  inventory:
    ssh:
      enabled: true
```

Collects per key:

- Fingerprint (SHA-256)
- Type (rsa, ecdsa, ed25519, dsa)
- Key size (bits)
- Comment
- Whether it appears in `authorized_keys`
- Associated known hosts

### Processes

Reads `/proc` to list running processes with PID, user, and command line.

```yaml
facter:
  inventory:
    process:
      enabled: true
```

### Applications (Docker)

Connects to the Docker socket to list containers, images and networks.

```yaml
facter:
  inventory:
    applications:
      enabled: true
      docker:
        enabled: true
```

### Systemd Services

Uses D-Bus / `systemctl` to collect unit state information.

```yaml
facter:
  inventory:
    systemdService:
      enabled: true
```

### Compliance (OpenSCAP)

Runs an OpenSCAP XCCDF evaluation and parses the result file.

```yaml
facter:
  compliance:
    enabled: true
    profile: "xccdf_org.ssgproject.content_profile_cis_server_l1"
    resultFile: "/tmp/openscap-results.xml"
```

!!! note
    Requires `oscap` and appropriate SCAP content (`scap-security-guide`) installed on the host.

### Vulnerabilities

Matches installed packages against CVE databases to detect vulnerable versions.

```yaml
facter:
  vulnerabilities:
    enabled: true
```

---

## Scheduling

`facter-oss` is a one-shot command. Use your system's scheduler to run it periodically:

### systemd Timer (recommended for Linux)

```ini
# /etc/systemd/system/facter.timer
[Unit]
Description=Facter inventory collection

[Timer]
OnCalendar=hourly
Persistent=true

[Install]
WantedBy=timers.target
```

```ini
# /etc/systemd/system/facter.service
[Unit]
Description=Facter inventory collection

[Service]
Type=oneshot
ExecStart=/usr/local/bin/facter --config /etc/facter/config.yml
```

```bash
systemctl enable --now facter.timer
```

### Cron

```bash
0 * * * * /usr/local/bin/facter --config /etc/facter/config.yml >> /var/log/facter.log 2>&1
```

---

## Performance Profiling

Enable Go pprof profiling:

```yaml
facter:
  performanceProfiling:
    enabled: true
```

Then generate a flame graph:

```bash
apt install graphviz
make profile
```
