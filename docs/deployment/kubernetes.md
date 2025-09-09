# Kubernetes Deployment

This guide covers deploying ArchesAI on Kubernetes using Helm charts and raw manifests.

## Prerequisites

- Kubernetes cluster (1.24+)
- kubectl configured
- Helm 3.x installed
- Container registry access

## Quick Start with Helm

### Install Helm Chart

```bash
# Add the ArchesAI repository
helm repo add archesai https://charts.archesai.com
helm repo update

# Install with default values
helm install archesai archesai/archesai

# Install with custom values
helm install archesai archesai/archesai -f values.yaml
```

### Custom Values

Create a `values.yaml` file:

```yaml
# Application
api:
  replicaCount: 3
  image:
    repository: archesai/api
    tag: latest
    pullPolicy: IfNotPresent

  resources:
    requests:
      memory: "256Mi"
      cpu: "250m"
    limits:
      memory: "1Gi"
      cpu: "1000m"

  autoscaling:
    enabled: true
    minReplicas: 3
    maxReplicas: 10
    targetCPUUtilizationPercentage: 70

# Database
postgresql:
  enabled: true
  auth:
    username: archesai
    password: changeme
    database: archesai
  primary:
    persistence:
      enabled: true
      size: 10Gi

# Redis
redis:
  enabled: true
  auth:
    enabled: true
    password: changeme
  master:
    persistence:
      enabled: true
      size: 5Gi

# Ingress
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: api.archesai.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: archesai-tls
      hosts:
        - api.archesai.com
```

## Kubernetes Manifests

### Namespace

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: archesai
```

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: archesai-config
  namespace: archesai
data:
  API_PORT: "8080"
  LOG_LEVEL: "info"
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_NAME: "archesai"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: archesai-secret
  namespace: archesai
type: Opaque
stringData:
  DB_USER: archesai
  DB_PASSWORD: secure-password
  REDIS_PASSWORD: redis-password
  JWT_SECRET: your-jwt-secret
  JWT_REFRESH_SECRET: your-refresh-secret
  OPENAI_API_KEY: your-openai-key
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: archesai-api
  namespace: archesai
  labels:
    app: archesai-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: archesai-api
  template:
    metadata:
      labels:
        app: archesai-api
    spec:
      containers:
        - name: api
          image: archesai/api:latest
          ports:
            - containerPort: 8080
              name: http
          envFrom:
            - configMapRef:
                name: archesai-config
            - secretRef:
                name: archesai-secret
          env:
            - name: DATABASE_URL
              value: "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=require"
            - name: REDIS_URL
              value: "redis://:$(REDIS_PASSWORD)@$(REDIS_HOST):$(REDIS_PORT)"
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "1000m"
          livenessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health/ready
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
          volumeMounts:
            - name: data
              mountPath: /data
      volumes:
        - name: data
          emptyDir: {}
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: archesai-api
  namespace: archesai
spec:
  selector:
    app: archesai-api
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
  type: ClusterIP
```

### Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: archesai-ingress
  namespace: archesai
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
    - hosts:
        - api.archesai.com
      secretName: archesai-tls
  rules:
    - host: api.archesai.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: archesai-api
                port:
                  number: 80
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
    "title": "ArchesAI Metrics",
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
