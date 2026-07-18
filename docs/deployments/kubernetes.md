# Kubernetes Deployment

!!! note "Work in progress"
    This page covers Kubernetes deployment guidelines. Helm charts are planned for a future release.

## Overview

All Facter components are stateless HTTP/gRPC services and are well-suited for Kubernetes deployment.

## Prerequisites

- Kubernetes cluster ≥ 1.28
- `kubectl` and `helm` (optional) configured
- A PostgreSQL operator or managed PostgreSQL service (for Apache AGE and rule-engine DB)
- A TLS certificate manager (e.g. cert-manager)
- A Keycloak instance

---

## Storage

### Apache AGE

Apache AGE requires PostgreSQL with the AGE extension. Options:
- Deploy Apache AGE as a StatefulSet with a PVC
- Use a managed PostgreSQL service with the AGE extension (check compatibility)

### PostgreSQL (rule-engine)

Standard PostgreSQL StatefulSet or managed service.

---

## Secrets

Create secrets for database credentials and TLS certificates:

```bash
# TLS certificates for facter-grpc
kubectl create secret generic facter-grpc-tls \
  --from-file=ca_cert.pem=./certs/ca_cert.pem \
  --from-file=server_cert.pem=./certs/server_cert.pem \
  --from-file=server_key.pem=./certs/server_key.pem

# TLS certificates for facter-rule-engine client
kubectl create secret generic facter-rule-engine-tls \
  --from-file=ca_cert.pem=./certs/ca_cert.pem \
  --from-file=facter_rule_engine_cert.pem=./certs/facter_rule_engine_cert.pem \
  --from-file=facter_rule_engine_key.pem=./certs/facter_rule_engine_key.pem

# Database credentials
kubectl create secret generic facter-grpc-db \
  --from-literal=password=<age-password>

kubectl create secret generic facter-rule-engine-db \
  --from-literal=uri=postgres://user:pass@host:5432/rule_engine
```

---

## ConfigMaps

Store configuration files as ConfigMaps:

```bash
kubectl create configmap facter-grpc-config \
  --from-file=config.yml=./configs/config-facter-grpc.yml

kubectl create configmap facter-api-config \
  --from-file=config.yml=./configs/config-facter-api.yml

kubectl create configmap facter-rule-engine-config \
  --from-file=config.yml=./configs/config-facter-re.yml
```

---

## Sample Deployment: facter-grpc

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: facter-grpc
spec:
  replicas: 1
  selector:
    matchLabels:
      app: facter-grpc
  template:
    metadata:
      labels:
        app: facter-grpc
    spec:
      containers:
        - name: facter-grpc
          image: ghcr.io/klamhq/facter-grpc:latest
          args: ["--config", "/configs/config.yml"]
          ports:
            - containerPort: 56230
              name: grpc
          volumeMounts:
            - name: config
              mountPath: /configs
            - name: tls
              mountPath: /certs
              readOnly: true
          livenessProbe:
            grpc:
              port: 56230
            initialDelaySeconds: 10
            periodSeconds: 30
      volumes:
        - name: config
          configMap:
            name: facter-grpc-config
        - name: tls
          secret:
            secretName: facter-grpc-tls
---
apiVersion: v1
kind: Service
metadata:
  name: facter-grpc
spec:
  selector:
    app: facter-grpc
  ports:
    - port: 56230
      targetPort: 56230
      name: grpc
  type: ClusterIP
```

---

## Ingress

Expose `facter-api` and `facter-rule-engine` through an Ingress controller. For gRPC (`facter-grpc`), use an L4 LoadBalancer or a gRPC-capable ingress (nginx with `nginx.ingress.kubernetes.io/backend-protocol: GRPC`).

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: facter-api-ingress
  annotations:
    nginx.ingress.kubernetes.io/backend-protocol: HTTP
spec:
  rules:
    - host: api.facter.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: facter-api
                port:
                  number: 56231
```

---

## facter-oss as a DaemonSet

To deploy the agent on every cluster node:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: facter-oss
  namespace: facter-system
spec:
  selector:
    matchLabels:
      app: facter-oss
  template:
    metadata:
      labels:
        app: facter-oss
    spec:
      hostPID: true
      hostNetwork: true
      tolerations:
        - operator: Exists
      containers:
        - name: facter-oss
          image: ghcr.io/klamhq/facter-oss:latest
          args: ["--config", "/configs/config.yml"]
          securityContext:
            privileged: true          # Required for firewall/process collectors
          volumeMounts:
            - name: config
              mountPath: /configs
            - name: tls
              mountPath: /certs
              readOnly: true
            - name: host-proc
              mountPath: /proc
              readOnly: true
              hostPath:
                path: /proc
      volumes:
        - name: config
          configMap:
            name: facter-oss-config
        - name: tls
          secret:
            secretName: facter-oss-tls
```

!!! warning
    Running the agent as a DaemonSet with `privileged: true` grants it broad host access. Review your security policy before enabling this in production.
