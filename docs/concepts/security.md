# Security

## Transport Security (mTLS)

All gRPC communication between agents and `facter-grpc`, and between `facter-rule-engine` and `facter-grpc`, uses **mutual TLS (mTLS)**.

Both client and server present certificates signed by a shared CA. The server validates the client's certificate against the CA, and the client validates the server's certificate against the same CA.

### Certificate Requirements

| Component                     | Certificate                   | Key                          | CA            |
| ----------------------------- | ----------------------------- | ---------------------------- | ------------- |
| `facter-grpc` (server)        | `server_cert.pem`             | `server_key.pem`             | `ca_cert.pem` |
| `facter-oss` (client)         | `facter_cert.pem`             | `facter_key.pem`             | `ca_cert.pem` |
| `facter-rule-engine` (client) | `facter_rule_engine_cert.pem` | `facter_rule_engine_key.pem` | `ca_cert.pem` |

### SPIFFE Identity

`facter-grpc` validates client connections based on SPIFFE identities embedded in the certificates' Subject Alternative Name (SAN) field. Only principals listed in the configuration are authorised.

```yaml
configuration:
  policy:
    facterPrincipal:
      - "spiffe://facter.fr/ns/backend/sa/facter"
    facterRuleEnginePrincipal:
      - "spiffe://facter.fr/ns/backend/sa/facter-rule-engine"
```

!!! warning "Production certificates"
    The `scripts/certs/create.sh` helper generates self-signed test certificates. **Never use these in production.** Use a proper PKI (e.g. cert-manager, Vault, or your organisation's CA).

### Generating Test Certificates

```bash
cd facter-grpc/scripts/certs
./create.sh
```

This produces:
- `ca_cert.pem` / `ca_key.pem` — Certificate Authority
- `server_cert.pem` / `server_key.pem` — facter-grpc server
- `facter_cert.pem` / `facter_key.pem` — facter-oss client
- `facter_rule_engine_cert.pem` / `facter_rule_engine_key.pem` — rule engine client

---

## Identity & Access Management (Keycloak)

`facter-api`, `facter-rule-engine`, and `facter-ui` use [Keycloak](https://www.keycloak.org/) as an identity provider.

### OIDC Flow

1. `facter-ui` redirects the user to Keycloak for login (Authorization Code Flow with PKCE)
2. After login, Keycloak issues a JWT access token
3. The UI attaches the JWT as `Authorization: Bearer <token>` on every API request
4. `facter-api` and `facter-rule-engine` validate the JWT against Keycloak's JWKS endpoint
5. Role claims are extracted from the token to enforce RBAC

### Realm & Client

Default configuration:

| Setting   | Value          |
| --------- | -------------- |
| Realm     | `facter-audit` |
| Client ID | `facter-ui`    |

### Roles

| Role            | Access                                                                    |
| --------------- | ------------------------------------------------------------------------- |
| `facter-viewer` | Read-only access to all data and rule execution results                   |
| `facter-admin`  | Full access including rule creation, update, deletion and audit execution |

!!! note
    System rules (prefixed `R-SYS-`) cannot be modified or deleted even by `facter-admin`.

---

## API Security

### facter-api (GraphQL)

- All endpoints require a valid JWT
- No write operations — read-only GraphQL resolvers only
- CORS origin allowlist configured per environment

### facter-rule-engine (REST)

- All `/api/v1/*` endpoints require a valid JWT
- Write endpoints (`POST`, `PUT`, `DELETE`) are restricted to `facter-admin`
- Read endpoints support both `facter-admin` and `facter-viewer`
- Public endpoint: `GET /health` (no auth)
- System rule protection middleware prevents deletion of `R-SYS-*` rules

---

## Security Checklist for Production

- [ ] Replace test certificates with production-grade PKI certificates
- [ ] Configure SPIFFE identities to match your organisation's trust domain
- [ ] Enable TLS on `facter-api` and `facter-rule-engine` HTTP servers (set `ssl: true` and provide `certificatePath` / `certificateKeyPath`)
- [ ] Restrict CORS `allowedOrigins` to your actual frontend domain
- [ ] Set secure PostgreSQL credentials (not default `password`)
- [ ] Rotate Keycloak client secrets periodically
- [ ] Pin the Apache AGE connection to a read-only user for `facter-api`
- [ ] Run `facter-oss` with the minimum required OS capabilities
