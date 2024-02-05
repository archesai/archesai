## Hosts File

You must add the following to your hosts file. This is so that the frontend can communicate with internal docker services.

```
127.0.0.1 arches-auth-keycloak
127.0.0.1 arches-minio
```

Version 0.1

## Steps to Set Up Cluster

## This is for Cloud SQL

https://cloud.google.com/sql/docs/mysql/connect-kubernetes-engine#gsa

gcloud iam service-accounts add-iam-policy-binding \
--role="roles/iam.workloadIdentityUser" \
--member="serviceAccount:archesai.svc.id.goog[archesai-stage/arches-api-service-account]" \
cloud-sql-proxy@archesai.iam.gserviceaccount.com
