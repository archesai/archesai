#!/bin/bash

set -euo pipefail # Exit on error, undefined variable, and fail on pipe errors

# Constants and Configurations
CERT_DIR="/etc/letsencrypt/live/archesai.dev"
CERT_FILE_SOURCE="$CERT_DIR/fullchain.pem"
KEY_FILE_SOURCE="$CERT_DIR/privkey.pem"
CERT_FILE_DEST="/tmp/fullchain.pem"
KEY_FILE_DEST="/tmp/privkey.pem"
NAMESPACE="default"
TLS_SECRET="archesai-tls"

# Verify certificate files exist with sudo access
if ! sudo bash -c "test -f '$CERT_FILE_SOURCE' && test -f '$KEY_FILE_SOURCE'"; then
    echo "❌ Certificate files not found in $CERT_DIR."
    exit 1
fi

# Copy cert files to /tmp and set permissions
echo "🔄 Copying certificate files to /tmp..."
sudo cp "$CERT_FILE_SOURCE" "$CERT_FILE_DEST"
sudo cp "$KEY_FILE_SOURCE" "$KEY_FILE_DEST"
sudo chmod 644 "$CERT_FILE_DEST" "$KEY_FILE_DEST"
sudo chown "$USER":"$USER" "$CERT_FILE_DEST" "$KEY_FILE_DEST"
echo "✅ Certificate files copied and permissions set."

# Create TLS secret if it doesn't exist
if kubectl get secret "$TLS_SECRET" --namespace "$NAMESPACE" &>/dev/null; then
    echo "✅ Secret '$TLS_SECRET' already exists in namespace '$NAMESPACE'. Skipping creation."
else
    echo "🔄 Creating TLS secret '$TLS_SECRET' in namespace '$NAMESPACE'..."
    kubectl create secret tls "$TLS_SECRET" \
        --cert="$CERT_FILE_DEST" \
        --key="$KEY_FILE_DEST" \
        --namespace "$NAMESPACE"
    echo "✅ TLS secret '$TLS_SECRET' created successfully."
fi
