#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install jq if not present (optional but recommended)
if ! command_exists jq; then
    echo "🔄 jq not found. Installing jq..."
    sudo apt update
    sudo apt install -y jq
else
    echo "✅ jq is already installed."
fi

# Install GPU Stuff
DOCKER_DAEMON_CONFIG="/etc/docker/daemon.json"
if [ -f "$DOCKER_DAEMON_CONFIG" ]; then
    # Check if 'nvidia' runtime is already configured
    if jq '.runtimes.nvidia' "$DOCKER_DAEMON_CONFIG" >/dev/null 2>&1; then
        echo "✅ NVIDIA runtime for Docker is already configured."
    else
        echo "🔄 Configuring NVIDIA runtime for Docker..."
        sudo nvidia-ctk runtime configure --runtime=docker
        sudo systemctl restart docker
    fi
else
    echo "🔄 Docker daemon.json not found. Configuring NVIDIA runtime for Docker..."
    sudo nvidia-ctk runtime configure --runtime=docker
    sudo systemctl restart docker
fi

# Create the secret for TLS
CERT_DIR="deploy/kubernetes/overlays/development"
CERT_FILE="$CERT_DIR/cert.pem"
KEY_FILE="$CERT_DIR/key.pem"
if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    echo "🔄 Generating TLS certificates with mkcert..."
    mkcert -install
    mkdir -p "$CERT_DIR"
    mkcert -cert-file "$CERT_FILE" -key-file "$KEY_FILE" *.archesai.test
else
    echo "✅ TLS certificates already exist."
fi

# Start minikube
if ! minikube status >/dev/null 2>&1; then
    echo "🔄 Starting Minikube..."
    minikube start --driver=docker --container-runtime docker --gpus all
else
    echo "✅ Minikube is already running."
fi

# Pull base images
BASE_IMAGES=(
    "node:20-alpine"
)

for image in "${BASE_IMAGES[@]}"; do
    IMAGE_NAME="${image%%:*}" # Extract the image name (e.g., node)
    IMAGE_TAG="${image##*:}"  # Extract the image tag (e.g., 20-alpine)

    # Check if the image is already present
    IMAGE_FOUND=$(minikube ssh -- docker images "$IMAGE_NAME" --format '{{.Repository}}:{{.Tag}}' | grep -F "$IMAGE_TAG")

    if [ -z "$IMAGE_FOUND" ]; then
        echo "🔄 Pulling Docker image: $image..."
        minikube ssh -- "docker pull $image"
    else
        echo "✅ Docker image $image already exists in Minikube."
    fi
done

# Enable the ingress controller
if minikube addons list | grep -E 'ingress\s+\|\s+minikube\s+\|\s+disabled' >/dev/null; then
    echo "Enabling ingress addon..."
    minikube addons enable ingress
else
    echo "✅ Ingress addon is already enabled."
fi

# Enable the ingress controller
if minikube addons list | grep -E 'ingress-dns\s+\|\s+minikube\s+\|\s+disabled' >/dev/null; then
    echo "🔄 Enabling ingress addon..."
    minikube addons enable ingress-dns
else
    echo "✅ Ingress addon is already enabled."
fi

# Configure ingress TLS
if ! kubectl -n ingress-nginx get deployment ingress-nginx-controller -o yaml | grep -q "default/mkcert-tls"; then
    echo "🔄 Configuring ingress TLS... and refreshing the ingress addon"
    kubectl create secret tls mkcert-tls \
        --cert=$CERT_FILE \
        --key=$KEY_FILE \
        --namespace=default
    echo "default/mkcert-tls" | minikube addons configure ingress
    minikube addons enable ingress --refresh
else
    echo "✅ Ingress TLS is already configured."
fi

# Write to resolved.conf.d/minikube.com if not already present
RESOLVED_CONF="/etc/systemd/resolved.conf.d/minikube.conf"
RESOLVED_CONTENT="[Resolve]
DNSStubListener=no
DNS=127.0.0.1"

if [ ! -f "$RESOLVED_CONF" ] || ! grep -Fxq "DNSStubListener=no" "$RESOLVED_CONF"; then
    echo "🔄 Writing Minikube DNS configuration..."
    sudo mkdir -p /etc/systemd/resolved.conf.d
    sudo bash -c "cat > $RESOLVED_CONF <<EOF
$RESOLVED_CONTENT
EOF"
    # Restart systemd-resolved only if configuration changed
    sudo systemctl restart systemd-resolved
else
    echo "✅ Minikube DNS configuration already exists."
fi

# Install dnsmasq if not installed
if ! command_exists dnsmasq; then
    echo "🔄 Installing dnsmasq..."
    sudo apt update
    sudo apt install -y dnsmasq
else
    echo "✅ dnsmasq is already installed."
fi

# Get Minikube IP
MINIKUBE_IP=$(minikube ip)
if [ -z "$MINIKUBE_IP" ]; then
    echo "🚫 Failed to retrieve Minikube IP."
    exit 1
fi
echo "✅ Minikube IP: $MINIKUBE_IP"

# Update dnsmasq configuration
DNSMASQ_CONF="/etc/dnsmasq.d/minikube.conf"
DESIRED_DNSMASQ_CONTENT=$(
    cat <<EOF
server=192.168.1.178
server=192.168.1.1
server=/test/$MINIKUBE_IP
listen-address=127.0.0.1
no-resolv
no-poll
EOF
)

if [ ! -f "$DNSMASQ_CONF" ] || ! grep -Fxq "server=/test/$MINIKUBE_IP" "$DNSMASQ_CONF"; then
    echo "🔄 Updating dnsmasq configuration..."
    sudo bash -c "cat > $DNSMASQ_CONF <<EOF
$DESIRED_DNSMASQ_CONTENT
EOF"
    # Restart dnsmasq only if configuration changed
    sudo systemctl restart dnsmasq
else
    echo "✅ dnsmasq configuration is already up-to-date."
fi

echo "✅ Setup complete."
