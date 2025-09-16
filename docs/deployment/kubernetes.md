# Kubernetes Deployment

This guide covers deploying Arches on Kubernetes using our hybrid Kustomize + Helm approach.

## Prerequisites

- Kubernetes cluster (1.24+)
- kubectl configured
- Helm 3.x installed
- Kustomize 4.x installed
- Container registry access

## Architecture Overview

Arches uses a **hybrid Kustomize + Helm deployment system**:

- **Kustomize components**: Plain YAML Kubernetes resources organized by service
- **Helm templating**: Only templates the `kustomization.yaml` file to control component composition
- **Environment-specific values**: Dev/prod configurations control which components are enabled

## Quick Start

### Deploy to Development

```bash
# Preview what will be deployed
make k8s-preview

# Deploy to dev environment
make k8s-deploy-dev

# Or use the script directly
./deployments/scripts/deploy.sh dev
```

### Deploy to Production

```bash
# Deploy to production
make k8s-deploy-prod

# Or with custom namespace
./deployments/scripts/deploy.sh prod archesai-production
```

## Component Architecture

### Kustomize Components

Each service is organized as a Kustomize component with **consistent labeling**:

```yaml
# All components use these standard labels
commonLabels:
  app.kubernetes.io/name: archesai
  app.kubernetes.io/component: <service-name>
```

**Available Components:**

- `api` - REST API server
- `database` - PostgreSQL database
- `platform` - React frontend
- `redis` - Redis cache
- `storage` - MinIO object storage
- `monitoring` - Grafana + Loki
- `ingress` - NGINX ingress rules
- `unstructured` - Document processing service
- `scraper` - Web scraping service
- `migrations` - Database migration job

### Environment Configuration

#### Development (`values-dev.yaml`)

```yaml
# Core services only
api:
  enabled: true
  replicas: 1
  image:
    tag: latest

database:
  enabled: true

platform:
  enabled: true
  replicas: 1

# Disabled in dev
redis:
  enabled: false
monitoring:
  enabled: false
ingress:
  enabled: false
storage:
  enabled: false
```

#### Production (`values-prod.yaml`)

```yaml
# All services enabled
api:
  enabled: true
  replicas: 3
  image:
    tag: "v1.0.0"

database:
  enabled: true

platform:
  enabled: true
  replicas: 2

# Production features
redis:
  enabled: true
monitoring:
  enabled: true
ingress:
  enabled: true
storage:
  enabled: true
unstructured:
  enabled: true
scraper:
  enabled: true
```

## Generated Manifests

The hybrid system generates clean Kubernetes manifests. Here are examples of what gets produced:

### Component Examples

#### API Component (`components/api/`)

```yaml
# deployment.yaml (labels added by Kustomize)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: archesai-api
spec:
  replicas: 1 # Controlled by Helm values
  selector:
    matchLabels:
      # Labels injected by commonLabels
  template:
    spec:
      serviceAccountName: archesai
      containers:
        - name: api
          image: archesai/api:latest # Tag controlled by Helm
          ports:
            - name: http
              containerPort: 3001
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: archesai-database
                  key: DATABASE_URL
```

#### Database Component (`components/database/`)

```yaml
# statefulset.yaml (labels added by Kustomize)
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: archesai-postgres
spec:
  serviceName: archesai-postgres
  replicas: 1
  template:
    spec:
      serviceAccountName: archesai
      containers:
        - name: postgres
          image: pgvector/pgvector:pg16
          env:
            - name: POSTGRES_DB
              value: archesai
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: archesai-database
                  key: POSTGRES_PASSWORD
```

## How It Works

### 1. Helm Templates Kustomization

```bash
# Helm templates the kustomization.yaml with environment values
helm template archesai deployments/helm-minimal \
  -f deployments/helm-minimal/values-dev.yaml > /tmp/kustomization.yaml
```

Generated `kustomization.yaml`:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: archesai-dev

resources:
  - deployments/kustomize/base

components:
  - deployments/kustomize/components/api
  - deployments/kustomize/components/database
  - deployments/kustomize/components/platform
  - deployments/kustomize/components/migrations

replicas:
  - name: archesai-api
    count: 1
  - name: archesai-platform
    count: 1

images:
  - name: archesai/api
    newTag: latest
  - name: archesai/platform
    newTag: latest
```

### 2. Kustomize Builds Final Manifests

```bash
# Kustomize processes components and applies labels
kustomize build /tmp > final-manifests.yaml
kubectl apply -f final-manifests.yaml
```

## Deployment Commands

### Manual Deployment

```bash
# 1. Template the kustomization.yaml
helm template archesai deployments/helm-minimal \
  -f deployments/helm-minimal/values-prod.yaml \
  --set namespace=archesai > /tmp/kustomization.yaml

# 2. Build and apply with Kustomize
kustomize build /tmp | kubectl apply -f -

# 3. Check deployment status
kubectl get all -n archesai
```

### Using Deployment Script

```bash
# Development environment
./deployments/scripts/deploy.sh dev

# Production environment
./deployments/scripts/deploy.sh prod

# Custom namespace
./deployments/scripts/deploy.sh prod my-namespace

# Dry run (see generated manifests)
./deployments/scripts/deploy.sh dev default true
```

### Using Makefile

```bash
# Preview deployment
make k8s-preview

# Deploy development
make k8s-deploy-dev

# Deploy production
make k8s-deploy-prod

# Dry run
make k8s-dry-run
```

## Database Deployment

### PostgreSQL StatefulSet

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: archesai
spec:
  serviceName: postgres-service
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:16-alpine
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: archesai-secret
                  key: DB_USER
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: archesai-secret
                  key: DB_PASSWORD
            - name: POSTGRES_DB
              value: archesai
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - metadata:
        name: postgres-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi
```

### Redis StatefulSet

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
  namespace: archesai
spec:
  serviceName: redis-service
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          ports:
            - containerPort: 6379
          command:
            - redis-server
            - --requirepass
            - $(REDIS_PASSWORD)
          env:
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: archesai-secret
                  key: REDIS_PASSWORD
          volumeMounts:
            - name: redis-storage
              mountPath: /data
  volumeClaimTemplates:
    - metadata:
        name: redis-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 5Gi
```

## Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: archesai-api-hpa
  namespace: archesai
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: archesai-api
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

## Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: archesai-api-policy
  namespace: archesai
spec:
  podSelector:
    matchLabels:
      app: archesai-api
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: postgres
      ports:
        - protocol: TCP
          port: 5432
    - to:
        - podSelector:
            matchLabels:
              app: redis
      ports:
        - protocol: TCP
          port: 6379
    - to:
        - namespaceSelector: {}
          podSelector:
            matchLabels:
              k8s-app: kube-dns
      ports:
        - protocol: UDP
          port: 53
```

## Monitoring with Prometheus

### ServiceMonitor

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: archesai-api
  namespace: archesai
spec:
  selector:
    matchLabels:
      app: archesai-api
  endpoints:
    - port: http
      path: /metrics
      interval: 30s
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "Arches Metrics",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{job=\"archesai-api\"}[5m])"
          }
        ]
      },
      {
        "title": "Response Time",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))"
          }
        ]
      }
    ]
  }
}
```

## GitOps with ArgoCD

### Application Manifest

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: archesai
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/archesai/archesai
    targetRevision: HEAD
    path: deployments/kubernetes
  destination:
    server: https://kubernetes.default.svc
    namespace: archesai
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

## Backup and Restore

### Velero Backup

```bash
# Install Velero
velero install \
  --provider aws \
  --bucket archesai-backups \
  --secret-file ./credentials-velero

# Create backup
velero backup create archesai-backup \
  --include-namespaces archesai \
  --ttl 720h

# Restore from backup
velero restore create \
  --from-backup archesai-backup
```

## Security

### Pod Security Policy

```yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: archesai-psp
spec:
  privileged: false
  allowPrivilegeEscalation: false
  requiredDropCapabilities:
    - ALL
  volumes:
    - "configMap"
    - "emptyDir"
    - "projected"
    - "secret"
    - "persistentVolumeClaim"
  runAsUser:
    rule: "MustRunAsNonRoot"
  seLinux:
    rule: "RunAsAny"
  fsGroup:
    rule: "RunAsAny"
```

### RBAC

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: archesai-role
  namespace: archesai
rules:
  - apiGroups: [""]
    resources: ["pods", "services"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: archesai-rolebinding
  namespace: archesai
subjects:
  - kind: ServiceAccount
    name: archesai-sa
    namespace: archesai
roleRef:
  kind: Role
  name: archesai-role
  apiGroup: rbac.authorization.k8s.io
```

## Troubleshooting

### Debug Commands

```bash
# Check pod status
kubectl get pods -n archesai

# Describe pod
kubectl describe pod archesai-api-xxx -n archesai

# View logs
kubectl logs -f archesai-api-xxx -n archesai

# Execute into pod
kubectl exec -it archesai-api-xxx -n archesai -- /bin/sh

# Check events
kubectl get events -n archesai --sort-by='.lastTimestamp'

# Port forward for debugging
kubectl port-forward -n archesai svc/archesai-api 8080:80
```

### Common Issues

1. **ImagePullBackOff**

   ```bash
   kubectl create secret docker-registry regcred \
     --docker-server=registry.example.com \
     --docker-username=user \
     --docker-password=pass
   ```

2. **CrashLoopBackOff**

   ```bash
   kubectl logs archesai-api-xxx -n archesai --previous
   ```

3. **Pending PVCs**

   ```bash
   kubectl get pvc -n archesai
   kubectl describe pvc postgres-storage-postgres-0 -n archesai
   ```

## Local Development with k3d

```bash
# Create cluster
k3d cluster create archesai \
  --servers 1 \
  --agents 3 \
  --port 8080:80@loadbalancer

# Deploy
kubectl apply -k deployments/kubernetes/

# Access application
curl http://localhost:8080
```

## Next Steps

- [Production Guide](production.md) for production best practices
- [Docker Guide](docker.md) for container configuration
- [Monitoring Setup](../monitoring/overview.md) for observability
