#!/usr/bin/env bash
set -e 

## THIS FILE NEEDS TO BE EXECUTED IN THE CI/CD PIPELINE WITH $VERSION
## $nlp_sha, $api_sha, $ui_sha are read from a source file

source /workspace/tags.sh

VERSION=$(echo "${VERSION}" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9-]/-/g')
TIMESTAMP=$(date +'%Y%m%d%H%M%S')
RELEASE_NAME="release-${VERSION}-${TIMESTAMP}"
RELEASE_NAME="${RELEASE_NAME:0:63}"
gcloud components install gke-gcloud-auth-plugin
gcloud container clusters get-credentials archesai-cluster --zone us-east1-c --project archesai
gcloud deploy releases create "$RELEASE_NAME" \
--project=archesai \
--region=us-east1 \
--delivery-pipeline=arches-deployment \
--images=us-east4-docker.pkg.dev/archesai/images/arches-api=us-east4-docker.pkg.dev/archesai/images/arches-api:${api_sha},us-east4-docker.pkg.dev/archesai/images/arches-ui=us-east4-docker.pkg.dev/archesai/images/arches-ui:${ui_sha},us-east4-docker.pkg.dev/archesai/images/arches-nlp=us-east4-docker.pkg.dev/archesai/images/arches-nlp:${nlp_sha}
