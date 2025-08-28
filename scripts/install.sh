#!/bin/bash
# HarborCtl Installer

set -e

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Download URLs
GITHUB_REPO="leandrodaf/harborctl"
BINARY_NAME="harborctl"

# Get latest release tag
LATEST_TAG=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Build download URL with correct asset naming
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_TAG}/${BINARY_NAME}_${LATEST_TAG}_linux_${ARCH}"

echo "🚀 Installing HarborCtl..."
echo "   OS: $OS"
echo "   Arch: $ARCH"
echo "   Version: $LATEST_TAG"

# Download
echo "📥 Downloading $BINARY_NAME..."
curl -sSLf "$DOWNLOAD_URL" -o "$BINARY_NAME"

# Make executable
chmod +x "$BINARY_NAME"

# Install
if [ -w "/usr/local/bin" ]; then
    mv "$BINARY_NAME" "/usr/local/bin/"
    echo "✅ HarborCtl installed in /usr/local/bin/$BINARY_NAME"
elif sudo -n true 2>/dev/null; then
    sudo mv "$BINARY_NAME" "/usr/local/bin/"
    echo "✅ HarborCtl installed in /usr/local/bin/$BINARY_NAME (with sudo)"
else
    echo "⚠️  Install manually:"
    echo "   sudo mv $BINARY_NAME /usr/local/bin/"
    exit 1
fi

# Verify installation
if command -v $BINARY_NAME >/dev/null 2>&1; then
    echo "🎉 Installation completed!"
    echo ""
    echo "📋 Next steps:"
    echo "   • Server setup: $BINARY_NAME init-server --domain yourdomain.com --email admin@yourdomain.com"
    echo "   • Deploy service: $BINARY_NAME deploy-service --service my-service --repo https://github.com/company/my-service.git"
    echo "   • View help: $BINARY_NAME docs"
    echo ""
    $BINARY_NAME --version
else
    echo "❌ Installation failed"
    exit 1
fi
