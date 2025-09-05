# Development Setup

## GCP Setup

### I forgot what this is but we need it

https://cloud.google.com/sql/docs/mysql/connect-kubernetes-engine#gsa

```
gcloud iam service-accounts add-iam-policy-binding \
--role="roles/iam.workloadIdentityUser" \
--member="serviceAccount:archesai.svc.id.goog[archesai-stage/cloud-sql-service-account]" \
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

REPOSITORIES WITH INCLUDES

- UserRepository (Auth Domain)
- OrganizationRepository (Organizations Domain)
- PipelineRepository (Workflows Domain)
- ArtifactRepository (Content Domain)

## Installing Minikube

sudo apt install -y conntrack
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
minikube version
sudo apt install -y kubectl
kubectl get nodes

### Check for rbac enabled

```
kubectl api-resources | grep 'RoleTypeEnum\|ClusterRoleTypeEnum'

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

---

# BUSINESS

# Arches AI Design and Use Cases

## TODO

- Add pipeline ux for building stuff with connectors
- Implement ability to save "forms" that can be run against content and extract data to fill out the form. this form type can then be used as input to other tools
- Get labels to work better and add faceting
- Notifications with websockets
- [DONE] Fix websocket auth
- Ensure that refresh tokens work as expected and write tests
- Implement react flow to create pipelines with steps

## Introduction

**Arches AI** is a comprehensive data processing platform designed to empower businesses to efficiently manage, analyze, and transform their diverse data assets. Similar to Palantir Foundry, Arches AI enables organizations to upload various types of content—including files, audio, text, images, and websites—and index them for seamless parsing, querying, and transformation. Leveraging advanced embedding models and a suite of transformation tools, Arches AI provides flexible and powerful data processing capabilities tailored to meet the unique needs of different industries.

## Core Features

### Data Upload and Indexing

- **Multi-Format Support:** Seamlessly upload and manage files, audio, text, images, and websites.
- **Automated Indexing:** Efficiently index all uploaded content for quick retrieval and management.

### Transformation Tools

- **Text-to-Speech:** Convert textual data into natural-sounding audio.
- **Text-to-Image:** Generate high-quality images based on textual descriptions.
- **Text-to-Text:** Advanced text manipulation, generation, and transformation capabilities.
- **Random Files to Text:** Extract and convert content from various file types into text format.

### Embedding Models

- **Advanced Embeddings:** Utilize state-of-the-art models to embed text content into vector representations.
- **Semantic Search:** Enable sophisticated querying and semantic search for enhanced data accessibility.

### Data Querying and Transformation

- **Intuitive Query Interface:** User-friendly tools for querying indexed data with ease.
- **Data Transformation Tools:** Flexible tools to transform data to meet specific business requirements.

### Workflow Building

- **Custom Workflows:** Design and implement data processing workflows using individual tools through the workflows domain.
- **Automation:** Automate complex data workflows tailored to organizational needs.
- **Directed Acyclical Graph**: The workflows are DAGs, so you can represent all possible processing chains.
- **Pipeline Runs:** Track and monitor workflow execution with detailed run history and status.

### Support and Consulting

- **Integration Support:** Expert assistance in integrating Arches AI with existing systems.
- **Data Strategy Consulting:** Help businesses optimize their data strategies for maximum impact.

## Design Concepts

### Scalability

- **Modular Architecture:** Easily add or remove components to scale with business growth.
- **Cloud-Native Infrastructure:** Built on scalable cloud platforms to handle increasing data volumes.

### Usability

- **Intuitive Interface:** User-friendly dashboards and interfaces to lower the barrier to entry.
- **Customizable Workflows:** Flexible pipeline creation to suit various business processes.

### Security

- **Data Encryption:** Ensure data is securely stored and transmitted using advanced encryption standards.
- **Access Controls:** Robust authentication and authorization mechanisms to protect sensitive data.

### Integration

- **APIs:** RESTful and GraphQL APIs for seamless integration with other tools and services.
- **Third-Party Integrations:** Support for integrating with popular third-party applications and services.

### Technology Stack

- **Backend:** Go with Echo framework, Domain-Driven Design architecture with four domains (auth, organizations, workflows, content)
- **Database:** PostgreSQL with vector extensions for efficient storage and querying of embeddings, type-safe queries with sqlc
- **Frontend:** TypeScript/React with TanStack Router, built with Vite in a monorepo structure
- **API:** OpenAPI-first development with automatic code generation for type safety
- **AI Models:** Proprietary and third-party AI models for data transformation and embedding
- **Cloud Infrastructure:** Kubernetes-native with Helm charts for scalable deployment

## Use Cases by Industry

### Finance

- **Fraud Detection:** Analyze transaction data to identify and prevent fraudulent activities.
- **Risk Management:** Assess and manage financial risks through comprehensive data analysis.
- **Customer Insights:** Gain deeper understanding of customer behaviors and preferences to enhance services.

### Healthcare

- **Medical Records Management:** Organize and analyze patient data for improved healthcare delivery.
- **Research and Development:** Facilitate medical research by managing and processing large datasets.
- **Telemedicine:** Enhance telemedicine services through efficient data processing and transformation.

### Retail

- **Inventory Management:** Optimize inventory levels and reduce stockouts through data-driven insights.
- **Personalized Marketing:** Create targeted marketing campaigns based on customer data analysis.
- **Sales Analytics:** Analyze sales data to identify trends and improve sales strategies.

### Technology

- **Product Development:** Streamline product development processes with efficient data management.
- **User Experience Analysis:** Analyze user data to enhance product usability and satisfaction.
- **IT Operations:** Improve IT operations through data-driven monitoring and management.

### Manufacturing

- **Supply Chain Optimization:** Enhance supply chain efficiency through comprehensive data analysis.
- **Quality Control:** Implement data-driven quality control measures to reduce defects.
- **Predictive Maintenance:** Use data to predict and prevent equipment failures, minimizing downtime.

### Education

- **Student Performance Tracking:** Analyze student data to improve educational outcomes.
- **Curriculum Development:** Use data insights to develop and refine educational programs.
- **Administrative Efficiency:** Streamline administrative tasks through effective data management.

### Media and Entertainment

- **Content Management:** Organize and manage large volumes of media content efficiently.
- **Audience Analytics:** Gain insights into audience preferences and behaviors to tailor content.
- **Content Personalization:** Deliver personalized content experiences based on data analysis.

### Logistics

- **Route Optimization:** Improve delivery routes through data-driven insights, reducing costs and increasing efficiency.
- **Fleet Management:** Manage and monitor fleet operations efficiently using real-time data.
- **Demand Forecasting:** Predict demand to optimize logistics and reduce operational costs.

### Legal

- **Document Management:** Organize and search through large volumes of legal documents with ease.
- **Case Analysis:** Analyze case data to identify patterns and support legal strategies.
- **Compliance Monitoring:** Ensure compliance with regulations through continuous data monitoring and reporting.

### Energy

- **Resource Management:** Optimize the use of resources through detailed data analysis.
- **Predictive Maintenance:** Predict equipment failures and schedule maintenance to prevent downtime.
- **Energy Consumption Analysis:** Analyze energy usage patterns to improve efficiency and reduce costs.

### Real Estate

- **Property Management:** Manage property data efficiently, including documents, images, and tenant information.
- **Market Analysis:** Analyze market trends to inform investment and development strategies.
- **Customer Relationship Management:** Enhance client interactions through detailed data insights.

## Conclusion

Arches AI offers a versatile and powerful data processing platform designed to meet the diverse needs of businesses across various industries. By providing a comprehensive suite of tools for data management, transformation, and analysis, Arches AI empowers organizations to unlock the full potential of their data, drive innovation, and achieve strategic goals. Whether it's enhancing operational efficiency, enabling advanced analytics, or fostering data-driven decision-making, Arches AI is positioned to be an indispensable partner for businesses seeking to thrive in the data-centric landscape.

# Project golang

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
