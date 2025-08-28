# Complete Guide - HarborCtl

Detailed documentation of the HarborCtl system for microservice deployment.

## 📖 Index

1. [Basic Concepts](#basic-concepts)
2. [Installation](#installation)
3. [Server Commands](#server-commands)
4. [Remote Commands](#remote-commands)
5. [Configuration](#configuration)
6. [Practical Examples](#practical-examples)
7. [Troubleshooting](#troubleshooting)

## 🎯 Basic Concepts

### Architecture

HarborCtl separates responsibilities into two layers:

- **🏗️ Base Server**: Centralized infrastructure (Traefik, observability)
- **🚀 Microservices**: Isolated applications deployed via Git

### Workflow

1. **Admin** configures base server once
2. **Developers** deploy microservices independently
3. **CI/CD** automates deployments via GitHub Actions

## 💻 Installation

### Super Quick Installation (Direct Binary)
```bash
sudo curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64 -o /usr/local/bin/harborctl && sudo chmod +x /usr/local/bin/harborctl
```

### Auto-detect Architecture
```bash
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

curl -sSLf "https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_${ARCH}.tar.gz" | sudo tar -xzC /usr/local/bin harborctl
```

### Verify Installation
```bash
harborctl --version
```

### Build from Source
```bash
git clone https://github.com/leandrodaf/harborctl.git
cd harborctlr
go build -o harborctl ./cmd/harborctl
```

## 🏗️ Comandos do Servidor

Execute estes comandos **no servidor de produção**.

### Inicialização
```bash
# Setup completo da infraestrutura
harborctl init-server --domain production.example.com --email admin@example.com

# Validar configuração antes de aplicar
harborctl validate -f server-base.yml

# Aplicar infraestrutura
harborctl up -f server-base.yml
```

### Gerenciamento
```bash
# Ver status de todos os serviços
harborctl status

# Parar todos os serviços
harborctl down

# Escalar serviços específicos
harborctl scale traefik --replicas 2
harborctl scale dozzle --replicas 1
```

### Utilitários
```bash
# Gerar senha para autenticação básica
harborctl hash-password --password "mypassword"

# Auditoria de segurança
harborctl security-audit

# Ver logs de serviço específico
harborctl logs traefik --tail 50

# Documentação
harborctl docs
```

## 🚀 Comandos Remotos

Execute estes comandos **remotamente** (local ou CI/CD).

### Deploy de Microserviços
```bash
# Deploy básico
harborctl deploy-service --service auth-service --repo https://github.com/company/auth-service.git

# Deploy com branch específica
harborctl deploy-service --service auth-service --repo https://github.com/company/auth-service.git --branch develop

# Deploy local (código já clonado)
harborctl deploy-service --service auth-service

# Deploy com scaling
harborctl deploy-service --service auth-service --replicas 5
```

### Desenvolvimento Local
```bash
# Inicializar novo microserviço
harborctl init --project my-service --domain localhost

# Validar configuração local
harborctl validate -f deploy/stack.yml

# Testar localmente
harborctl up -f deploy/stack.yml
```

## ⚙️ Configuração

### Variáveis de Ambiente

#### Servidor
```bash
# /etc/environment no servidor
DOMAIN=production.example.com
ACME_EMAIL=admin@example.com
LOG_LEVEL=info
```

#### Microserviço
```yaml
# deploy/stack.yml
version: 1
project: auth-service

services:
  - name: auth-api
    subdomain: auth
    build:
      context: .
      dockerfile: Dockerfile
    expose: 8080
    replicas: 2
    
    env:
      APP_ENV: production
      LOG_LEVEL: info
      DATABASE_URL: ${DATABASE_URL}
    
    secrets:
      - name: database_password
        file: secrets/database_password.txt
    
    traefik: true
```

### GitHub Secrets

Configure no repositório do microserviço:

```bash
# Secrets sensíveis
DATABASE_PASSWORD=secret_password
JWT_SECRET=secret_key_32_chars_minimum
API_KEY=external_api_key
ENCRYPTION_KEY=base64_encoded_key

# Deploy
DEPLOY_TOKEN=github_token
HARBOR_SERVER_HOST=production.example.com
HARBOR_SERVER_USER=harbor
HARBOR_SSH_KEY=private_ssh_key
```

### GitHub Variables

```bash
# URLs e configurações
DATABASE_URL=postgresql://user:${DATABASE_PASSWORD}@postgres:5432/db
API_BASE_URL=https://api.example.com
LOG_LEVEL=info
MONITORING_ENABLED=true
```

## 📋 Exemplos Práticos

### 1. Setup Inicial Completo

```bash
# No servidor de produção
sudo useradd -m -s /bin/bash harbor
sudo usermod -aG docker harbor
sudo su - harbor

# Clone o harborctlr
git clone https://github.com/leandrodaf/harborctl.git /opt/harbor
cd /opt/harbor

# Setup automático
./scripts/setup-production-server.sh production.example.com admin@example.com

# Verificar
harborctl status
```

### 2. Deploy de Microserviço Auth

```bash
# Criar microserviço (desenvolvedor)
./scripts/create-microservice.sh auth-service api

# Configurar GitHub Secrets no repositório
# DATABASE_PASSWORD, JWT_SECRET, etc.

# Deploy automático via push ou manual
harborctl deploy-service --service auth-service --repo https://github.com/company/auth-service.git
```

### 3. Escalabilidade

```bash
# Durante pico de tráfego
harborctl deploy-service --service auth-service --replicas 10
harborctl deploy-service --service payment-service --replicas 8

# Verificar recursos
harborctl status --details

# Voltar ao normal
harborctl deploy-service --service auth-service --replicas 2
```

## 🔍 Troubleshooting

### Problemas Comuns

#### Serviço não inicia
```bash
# Verificar logs
harborctl logs auth-service --tail 100

# Verificar configuração
harborctl validate -f deploy/stack.yml

# Verificar recursos do Docker
docker system df
docker stats
```

#### SSL não funciona
```bash
# Verificar Traefik
harborctl logs traefik --tail 50

# Verificar DNS
nslookup your-domain.com

# Forçar renovação certificado
docker exec traefik traefik acme --force
```

#### Deploy falha
```bash
# Verificar conectividade SSH
ssh harbor@production.example.com "harborctl status"

# Deploy com debug
harborctl deploy-service --service auth-service --debug

# Verificar GitHub Actions
# https://github.com/company/repo/actions
```

### Logs e Monitoramento

```bash
# Logs do sistema
journalctl -u docker -f

# Logs via web
# https://logs.yourdomain.com (Dozzle)

# Métricas via web  
# https://monitor.yourdomain.com (Beszel)

# Status detalhado
harborctl status --json
```

## 🔗 Links Úteis

- [Quick Start](QUICK_START.md) - Começar rapidamente
- [Scripts](../scripts/) - Scripts de automação  
- [Templates](../templates/) - Templates prontos
- [GitHub Issues](https://github.com/leandrodaf/harborctl/issues) - Support
