# Script para Setup Completo do Servidor
#!/bin/bash

set -e

DOMAIN="$1"
EMAIL="$2"

if [ -z "$DOMAIN" ] || [ -z "$EMAIL" ]; then
    echo "❌ Uso: $0 <dominio> <email>"
    echo "   Exemplo: $0 production.example.com devops@example.com"
    exit 1
fi

echo "🚀 Configurando servidor de produção..."
echo "   Domínio: $DOMAIN"
echo "   Email: $EMAIL"

# 1. Verificar dependências
echo "🔍 Verificando dependências do servidor..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker não está instalado. Execute:"
    echo "   curl -fsSL https://get.docker.com | sh"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose não está instalado"
    exit 1
fi

echo "✅ Docker instalado: $(docker --version)"
echo "✅ Docker Compose disponível"

# 2. Criar configuração base do servidor
echo "🏗️  Criando configuração base do servidor..."
./harborctl init-server --domain "$DOMAIN" --email "$EMAIL"

if [ ! -f "server-base.yml" ]; then
    echo "❌ Falha ao criar server-base.yml"
    exit 1
fi

echo "✅ Configuração base criada: server-base.yml"

# 3. Deploy da infraestrutura base
echo "🚢 Fazendo deploy da infraestrutura base..."
./harborctl up -f server-base.yml

if [ $? -eq 0 ]; then
    echo "✅ Infraestrutura base deployada com sucesso!"
else
    echo "❌ Falha no deploy da infraestrutura base"
    exit 1
fi

# 4. Aguardar serviços ficarem prontos
echo "⏳ Aguardando serviços ficarem prontos..."
sleep 30

# 5. Verificar se serviços estão rodando
echo "🔍 Verificando status dos serviços..."
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# 6. Testar conectividade
echo "🌐 Testando conectividade..."

# Testar Traefik
if curl -k -f https://$DOMAIN > /dev/null 2>&1; then
    echo "✅ Traefik respondendo em https://$DOMAIN"
else
    echo "⚠️  Traefik ainda não está respondendo (pode levar alguns minutos para certificados)"
fi

# Testar Dozzle
if curl -k -f https://logs.$DOMAIN > /dev/null 2>&1; then
    echo "✅ Dozzle disponível em https://logs.$DOMAIN"
else
    echo "⚠️  Dozzle ainda não está respondendo"
fi

# Testar Beszel
if curl -k -f https://monitor.$DOMAIN > /dev/null 2>&1; then
    echo "✅ Beszel disponível em https://monitor.$DOMAIN"
else
    echo "⚠️  Beszel ainda não está respondendo"
fi

echo ""
echo "🎉 Servidor configurado com sucesso!"
echo ""
echo "📊 Painéis disponíveis:"
echo "   • Logs: https://logs.$DOMAIN"
echo "   • Monitor: https://monitor.$DOMAIN"
echo ""
echo "📦 Para deployar microserviços:"
echo "   harborctl deploy-service --service <nome-servico> --repo <url-repo>"
echo ""
echo "🔧 Configurações criadas:"
echo "   • server-base.yml - Configuração da infraestrutura base"
echo "   • .deploy/ - Arquivos de deploy gerados"
echo ""
echo "⚠️  Notas importantes:"
echo "   1. Certificados TLS podem levar alguns minutos para serem emitidos"
echo "   2. DNS deve estar apontando para este servidor"
echo "   3. Portas 80 e 443 devem estar liberadas no firewall"
echo "   4. Esta configuração base deve permanecer sempre rodando"
