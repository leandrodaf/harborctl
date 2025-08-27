#!/bin/bash

# Script para criar um novo microserviço com Harbor CLI
# Uso: ./create-microservice.sh <nome-do-servico> <tipo> [dominio]

set -e

SERVICE_NAME="$1"
SERVICE_TYPE="$2"  # api, frontend, worker
DOMAIN="$3"        # Opcional

if [ -z "$SERVICE_NAME" ]; then
    echo "❌ Uso: $0 <nome-do-servico> <tipo> [dominio]"
    echo "   Tipos: api, frontend, worker"
    echo "   Exemplo: $0 auth-service api production.example.com"
    exit 1
fi

if [ -z "$SERVICE_TYPE" ]; then
    SERVICE_TYPE="api"
fi

if [ -z "$DOMAIN" ]; then
    DOMAIN="example.com"
fi

# Diretório do harborctlr (assumindo que o script está em harborctlr/scripts/)
DEPLOYER_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TEMPLATE_DIR="$DEPLOYER_DIR/templates/microservice/$SERVICE_TYPE"

if [ ! -d "$TEMPLATE_DIR" ]; then
    echo "❌ Template não encontrado: $SERVICE_TYPE"
    echo "   Templates disponíveis: api, frontend, worker"
    exit 1
fi

echo "🚀 Criando microserviço: $SERVICE_NAME (tipo: $SERVICE_TYPE)"
echo "   Template: $TEMPLATE_DIR"
echo "   Domínio: $DOMAIN"

# Criar diretório do microserviço
mkdir -p "$SERVICE_NAME"
cd "$SERVICE_NAME"

# Copiar template
cp -r "$TEMPLATE_DIR"/* . 2>/dev/null || true

# Criar estruturas necessárias
mkdir -p deploy/{secrets,environments} src tests .github/workflows


# Copiar GitHub Action template
if [ -f "$DEPLOYER_DIR/templates/github-actions/deploy.yml" ]; then
    cp "$DEPLOYER_DIR/templates/github-actions/deploy.yml" .github/workflows/deploy.yml
    sed -i.bak "s/{{SERVICE_NAME}}/$SERVICE_NAME/g" .github/workflows/deploy.yml
    rm -f .github/workflows/deploy.yml.bak
fi

# Criar secrets templates b??sicos
mkdir -p deploy/secrets deploy/environments

cat > deploy/secrets/.gitignore << 'EOFGIT'
# Ignorar todos os secrets reais
*
!.gitignore
!*.example
!README.md
EOFGIT

cat > deploy/secrets/database_password.txt.example << 'EOFDB'
your_secure_database_password_here
EOFDB

echo "??? Microservi??o $SERVICE_NAME criado com sucesso!"
echo ""
echo "???? Estrutura criada em: $(pwd)"
echo "???? Pr??ximos passos:"
echo "   1. Configure secrets no GitHub"
echo "   2. Implemente seu c??digo em src/"
echo "   3. Commit e push para deploy autom??tico"
