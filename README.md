# Development Setup

## Hosts File

You must add the following to your hosts file. This is so that the frontend can communicate with internal docker services.
This is only needed for the offline version

```
127.0.0.1 arches-auth-keycloak
127.0.0.1 arches-minio
```

## GCP Setup

### I forgot what this is but we need it

https://cloud.google.com/sql/docs/mysql/connect-kubernetes-engine#gsa

```
gcloud iam service-accounts add-iam-policy-binding \
--role="roles/iam.workloadIdentityUser" \
--member="serviceAccount:archesai.svc.id.goog[archesai-stage/arches-api-service-account]" \
cloud-sql-proxy@archesai.iam.gserviceaccount.com
```

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

## Installing Minikube

sudo apt install -y conntrack
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
minikube version
sudo apt install -y kubectl
kubectl get nodes

### Check for rbac enabled

```
kubectl api-resources | grep 'Role\|ClusterRole'

minikube dashboard

sudo apt update
sudo apt install -y nvidia-driver-<your-driver-version>
```

### Install nvidia docker

Add the package repository

```
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
```

### Install the nvidia-docker package

```
sudo apt update
sudo apt install -y nvidia-docker2
```

### Restart Docker

sudo systemctl restart docker

put this in /etc/docker/daemon.json

```
{
    "default-runtime": "nvidia",
    "runtimes": {
        "nvidia": {
            "path": "nvidia-container-runtime",
            "runtimeArgs": []
        }
    }
}


minikube start --driver=docker --gpu

```

### Install k8s gpu device plugin

kubectl apply -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/v0.13.0/nvidia-device-plugin.yml

```
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test
spec:
  containers:
  - name: gpu-container
    image: nvidia/cuda:11.0-base
    resources:
      limits:
        nvidia.com/gpu: 1 # Request 1 GPU
    command: ["nvidia-smi"]
```

kubectl apply -f gpu-test.yaml

kubectl logs gpu-test

### Setting image pull secret

```
kubectl create secret docker-registry artifact-registry-key \
 --docker-server=us-east4-docker.pkg.dev \
 --docker-username=\_json_key \
 --docker-password="$(gcloud auth print-access-token)" \
 --docker-email=jonathan@archesai.com

minikube mount $HOME:/host

gcloud auth print-access-token | docker login -u oauth2accesstoken --password-stdin https://us-east4-docker.pkg.dev

cat /home/jonathan/.docker/config.json

kubectl delete secret artifact-registry-key

kubectl create secret generic artifact-registry-key \  INT ✘  minikube ⎈
--from-file=.dockerconfigjson=$HOME/.docker/config.json \
 --type=kubernetes.io/dockerconfigjson
```
