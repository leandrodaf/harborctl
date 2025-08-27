#!/bin/bash
# Instalador do Harbor CLI

set -e

# Detectar OS e arquitetura
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "❌ Arquitetura não suportada: $ARCH"; exit 1 ;;
esac

# URLs de download
GITHUB_REPO="company/harborctlr"
BINARY_NAME="harborctl"
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"

echo "🚀 Instalando Harbor CLI..."
echo "   OS: $OS"
echo "   Arch: $ARCH"

# Download
echo "📥 Baixando $BINARY_NAME..."
curl -sSL "$DOWNLOAD_URL" -o "$BINARY_NAME"

# Tornar executável
chmod +x "$BINARY_NAME"

# Instalar
if [ -w "/usr/local/bin" ]; then
    mv "$BINARY_NAME" "/usr/local/bin/"
    echo "✅ Harbor CLI instalado em /usr/local/bin/$BINARY_NAME"
elif sudo -n true 2>/dev/null; then
    sudo mv "$BINARY_NAME" "/usr/local/bin/"
    echo "✅ Harbor CLI instalado em /usr/local/bin/$BINARY_NAME (com sudo)"
else
    echo "⚠️  Instale manualmente:"
    echo "   sudo mv $BINARY_NAME /usr/local/bin/"
    exit 1
fi

# Verificar instalação
if command -v $BINARY_NAME >/dev/null 2>&1; then
    echo "🎉 Instalação concluída!"
    echo ""
    echo "📋 Próximos passos:"
    echo "   • Setup servidor: $BINARY_NAME init-server --domain yourdomain.com --email admin@yourdomain.com"
    echo "   • Deploy serviço: $BINARY_NAME deploy-service --service my-service --repo https://github.com/company/my-service.git"
    echo "   • Ver ajuda: $BINARY_NAME docs"
    echo ""
    $BINARY_NAME --version
else
    echo "❌ Falha na instalação"
    exit 1
fi
