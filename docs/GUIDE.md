# 📖 Guia Completo - HarborCtl

Sistema completo para deploy de microserviços.

## 🎯 Arquitetura

**Servidor Base (uma vez):**
- Traefik: Proxy reverso + SSL automático
- Dozzle: Logs centralizados
- Beszel: Monitoramento em tempo real
- Redes e volumes isolados

**Apps (múltiplas):**
- Deploy via GitHub Actions
- Integração automática com infraestrutura
- Escalabilidade independente

## 🚀 Setup Completo

### 1. Instalar no Servidor
```bash
curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64 -o harborctl
chmod +x harborctl && sudo mv harborctl /usr/local/bin/
```

### 2. Configurar Infraestrutura Base
```bash
harborctl init-server --domain seudominio.com --email admin@seudominio.com
harborctl up -f server-base.yml
```

### 3. Configurar App para Deploy Automático

**No repositório da sua app:**
```bash
# Copiar templates
mkdir -p deploy .github/workflows
cp templates/microservice/api/deploy/stack.yml deploy/stack.yml
cp templates/github-actions/deploy.yml .github/workflows/deploy.yml

# Editar deploy/stack.yml conforme sua app
# Configurar secrets no GitHub
```

### 4. GitHub Secrets Necessários
```
PRODUCTION_HOST=seuservidor.com
PRODUCTION_USER=deploy  
PRODUCTION_SSH_KEY=sua-chave-ssh-privada
```

### 5. Deploy Automático Ativado!
```bash
git push origin main  # ← Deploy automático!
```

## 🔧 Comandos Principais

### Gerenciar Servidor Base
```bash
# Status da infraestrutura
harborctl status

# Parar/iniciar infraestrutura  
harborctl stop
harborctl start
harborctl restart

# Desligar tudo
harborctl down
```

### Deploy Manual de Apps
```bash
# Deploy via repositório
harborctl deploy-service --service minha-api --repo https://github.com/usuario/minha-api.git

# Deploy local (para testes)
harborctl deploy-service --service minha-api --path deploy

# Escalar app específica
harborctl scale minha-api --replicas 3
```

## ⚙️ Configuração da App

### Stack.yml Básico
```yaml
version: 1
project: minha-api

services:
  - name: minha-api
    subdomain: api
    image: node:18-alpine
    expose: 3000
    replicas: 2
    
    env:
      NODE_ENV: production
      API_PORT: 3000
      
    resources:
      memory: 512m
      cpus: "0.5"
      
    health_check:
      enabled: true
      path: /health
      
    traefik: true

volumes:
  - name: minha_api_data
```

### GitHub Action Básico
```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Deploy to production
      run: |
        curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64 -o harborctl
        chmod +x harborctl
        
        echo "${{ secrets.PRODUCTION_SSH_KEY }}" > key
        chmod 600 key
        
        ./harborctl deploy-service \
          --host "${{ secrets.PRODUCTION_HOST }}" \
          --user "${{ secrets.PRODUCTION_USER }}" \
          --key key \
          --service "${{ github.event.repository.name }}" \
          --repo "${{ github.server_url }}/${{ github.repository }}"
```

## 🎯 Exemplos Práticos

### API Node.js
```bash
# 1. Criar repositório com:
# - src/app.js (sua API)
# - Dockerfile 
# - deploy/stack.yml
# - .github/workflows/deploy.yml

# 2. Configurar GitHub Secrets

# 3. Push = Deploy automático!
git push origin main
```

### Frontend React
```yaml
# deploy/stack.yml
services:
  - name: meu-frontend
    subdomain: app
    image: nginx:alpine
    expose: 80
    build:
      context: .
      dockerfile: Dockerfile.prod
```

### Worker/Background Job
```yaml
# deploy/stack.yml  
services:
  - name: meu-worker
    image: node:18-alpine
    replicas: 1
    # Sem traefik (não precisa de acesso web)
    traefik: false
```

## 🔍 Troubleshooting

### App não responde
```bash
# Ver logs
harborctl logs minha-api --tail 50

# Verificar status
harborctl status

# Reiniciar app específica
harborctl restart minha-api
```

### SSL não funciona
```bash
# Verificar configuração
harborctl validate -f server-base.yml

# Ver logs do Traefik
harborctl logs traefik --tail 100
```

### Deploy falha
```bash
# Deploy com debug
harborctl deploy-service --service minha-api --dry-run --verbose

# Verificar secrets no GitHub
# Verificar conectividade SSH
```

## 📚 Links Úteis

- **Quick Start**: [QUICK_START.md](QUICK_START.md)
- **Comandos**: [COMMAND_GUIDE.md](COMMAND_GUIDE.md)  
- **Templates**: [../templates/](../templates/)
- **Exemplos**: [../examples/](../examples/)