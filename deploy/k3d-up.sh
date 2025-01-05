#!/bin/bash

set -euo pipefail # Exit on error, undefined variable, and fail on pipe errors

# Constants and Configurations
CLUSTER_NAME="tower"
K3D_CONFIG_FILE="deploy/k3d-config.yaml"
NAMESPACE="default"
TLS_SECRET="archesai-tls"
INGRESS_NGINX_NAMESPACE="ingress-nginx"
INGRESS_NGINX_REPO="https://kubernetes.github.io/ingress-nginx"
INGRESS_NGINX_NAME="ingress-nginx"

# Ensure required commands are available
source ./deploy/commands.sh

# Create k3d cluster if it doesn't exist
if k3d cluster list | grep -qw "$CLUSTER_NAME"; then
    echo "✅ Cluster '$CLUSTER_NAME' already exists. Skipping creation."
else
    echo "🔄 Creating cluster '$CLUSTER_NAME'..."
    k3d cluster create "$CLUSTER_NAME" --config "$K3D_CONFIG_FILE"
    echo "✅ Cluster '$CLUSTER_NAME' created successfully."
fi

# Extract TLS certificate and key
source ./deploy/extract-tls.sh

# Install Ingress Nginx using Helm
if helm ls -n "$INGRESS_NGINX_NAMESPACE" | grep -qw "$INGRESS_NGINX_NAME"; then
    echo "✅ Ingress Nginx is already installed in namespace '$INGRESS_NGINX_NAMESPACE'."
else
    echo "🔄 Installing Ingress Nginx..."
    helm upgrade --install $INGRESS_NGINX_NAME ingress-nginx \
        --repo $INGRESS_NGINX_REPO \
        --namespace $INGRESS_NGINX_NAMESPACE --create-namespace
    # --set controller.extraArgs.default-ssl-certificate=$NAMESPACE/$TLS_SECRET
    echo "✅ Ingress Nginx installed successfully."
fi

helm repo add jetstack https://charts.jetstack.io --force-update
helm upgrade --install cert-manager jetstack/cert-manager \
    --namespace cert-manager --create-namespace \
    --version v1.16.2 \
    --set crds.enabled=true

echo "🎉 Setup complete."
