# Quick Start - Harbor CLI

Comece a usar o Harbor CLI em 5 minutos.

## ⚡ Setup Rápido

### 1. Instalar Harbor CLI

```bash
curl -sSL https://github.com/company/harborctlr/raw/main/scripts/install.sh | bash
```

### 2. Setup do Servidor (Admin)

```bash
# No servidor de produção
harborctl init-server --domain yourdomain.com --email admin@yourdomain.com
harborctl up -f server-base.yml
```

### 3. Deploy de Microserviço (Dev)

```bash
# Criar microserviço
./scripts/create-microservice.sh my-service api

# Deploy
harborctl deploy-service --service my-service --repo https://github.com/company/my-service.git
```

## 🎯 Exemplo Prático

### Microserviço de API

```bash
# 1. Criar estrutura
./scripts/create-microservice.sh auth-api api yourdomain.com

# 2. Configurar secrets no GitHub
# DATABASE_PASSWORD, JWT_SECRET, API_KEY

# 3. Implementar código em src/
# 4. Commit e push = deploy automático

# Ou deploy manual
harborctl deploy-service --service auth-api
```

### Resultado

- **API**: https://auth-api.yourdomain.com
- **Logs**: https://logs.yourdomain.com  
- **Métricas**: https://monitor.yourdomain.com

## 🔧 Comandos Essenciais

```bash
# Status geral
harborctl status

# Deploy microserviço  
harborctl deploy-service --service NAME

# Escalar
harborctl scale NAME --replicas 5

# Logs
harborctl logs NAME --tail 50

# Help
harborctl docs
```

## 📚 Próximos Passos

- [📖 Guia Completo](GUIDE.md) - Documentação detalhada
- [🔧 Scripts](../scripts/) - Automação
- [📋 Templates](../templates/) - Templates prontos
