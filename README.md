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

<!-- name: Cut Docs
on:
  workflow_dispatch:

jobs:
  rdme-openapi:
    runs-on: ubuntu-latest
    steps:
      - name: Check out repo 📚
        uses: actions/checkout@v3

      - name: Run `openapi` command 🚀
        uses: readmeio/rdme@v8
        with:
          rdme: openapi https://api.archesai.com/-json --key=${{ secrets.README_SECRET }} --id=64837ab02aa53c002a2ceccd -->

REPOSITORES WITH INCLUDES

- UserRepository
- PipelineRepository
