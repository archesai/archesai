#!/bin/bash

# Install GPU Stuff
sudo nvidia-ctk runtime configure --runtime=docker && sudo systemctl restart docker

# Create the secret for tls
mkcert -install
mkcert -cert-file deploy/kubernetes/overlays/dev/cert.pem -key-file deploy/kubernetes/overlays/dev/key.pem arches-api.test arches-ui.test arches-grafana.test

# Start minikube
minikube start --driver=docker --container-runtime docker --gpus all

# Pull base images
minikube ssh docker pull node:20-alpine
minikube ssh docker pull quay.io/unstructured-io/base-images@sha256:38de7347cad45c069b1fd0c2ab8f52947aaf45e8a5eda553d8d968e7167510e4

# Enable the ingress controller
echo "default/mkcert-tls" | minikube addons configure ingress
minikube addons enable ingress
minikube addons enable ingress-dns

# Write to resolved.conf.d/minikube.com
sudo bash -c "cat > /etc/systemd/resolved.conf.d/minikube.conf <<EOF
[Resolve]
DNSStubListener=no
DNS=127.0.0.1
EOF"

# Restart systemd-resolved
sudo systemctl restart systemd-resolved

# Install dnsmasq
sudo apt update
sudo apt install dnsmasq

# Get Minikube IP
MINIKUBE_IP=$(minikube ip)

# Update dnsmasq configuration
sudo bash -c "cat > /etc/dnsmasq.d/minikube.conf <<EOF
server=8.8.8.8
server=/test/$MINIKUBE_IP
no-resolv
no-poll
EOF"

# Restart services
sudo systemctl restart dnsmasq
sudo systemctl restart systemd-resolved
