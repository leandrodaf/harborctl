# Template: API Microservice

Template para criar uma API REST com Harbor CLI.

## 🚀 Uso Rápido

```bash
# Criar automaticamente
./scripts/create-microservice.sh my-api api yourdomain.com

# Ou copiar manualmente
cp -r templates/microservice/api/* my-api/
```

## 📂 Estrutura Criada

```
my-api/
├── deploy/
│   ├── stack.yml              # Configuração Harbor
│   ├── environments/          # Ambientes (dev/staging/prod)
│   └── secrets/               # Templates de secrets
├── src/
│   └── index.js              # Código da API
├── Dockerfile                # Build otimizado
├── package.json              # Dependencies Node.js
├── .github/workflows/
│   └── deploy.yml            # CI/CD automático
└── README.md                 # Documentação
```

## ⚙️ Configuração

### 1. GitHub Secrets

Configure no repositório:

```bash
DATABASE_PASSWORD=secure_password
JWT_SECRET=secure_jwt_key_32_chars
API_KEY=external_api_key
```

### 2. GitHub Variables

```bash
DATABASE_URL=postgresql://user:${DATABASE_PASSWORD}@postgres:5432/db
API_BASE_URL=https://api.yourdomain.com
LOG_LEVEL=info
```

### 3. Personalizar

Edite `deploy/stack.yml` conforme sua necessidade:

- **Port**: Altere `expose: 8080`
- **Replicas**: Altere `replicas: 2`
- **Resources**: Ajuste `memory` e `cpus`
- **Environment**: Adicione variáveis em `env:`

## 🚀 Deploy

```bash
# Deploy manual
harborctl deploy-service --service my-api --repo https://github.com/company/my-api.git

# Deploy automático via push para main
git push origin main
```

## 📊 Monitoramento

- **API**: https://my-api.yourdomain.com
- **Health**: https://my-api.yourdomain.com/health
- **Logs**: https://logs.yourdomain.com
- **Métricas**: https://monitor.yourdomain.com

## 🔧 Desenvolvimento

```bash
# Instalar dependencies
npm install

# Executar localmente
npm run dev

# Testar
npm test

# Build
npm run build
```
