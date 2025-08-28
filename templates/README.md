# HarborCtl Templates

Ready-to-use templates for HarborCtl.

## 📂 Estrutura

```
templates/
├── microservice/           # Templates de microserviços
│   ├── api/               # API REST
│   ├── frontend/          # Frontend SPA
│   ├── worker/            # Background worker
│   └── database/          # Banco de dados
├── github-actions/        # GitHub Actions workflows
│   ├── deploy.yml         # Deploy automático
│   └── build-test.yml     # Build e test
└── config/               # Configurações
    ├── server-base.yml    # Infraestrutura base
    └── stack-examples/    # Exemplos de stack.yml
```

## 🚀 Como Usar

### Criar Microserviço

```bash
# API REST
./scripts/create-microservice.sh my-api api yourdomain.com

# Frontend SPA  
./scripts/create-microservice.sh my-app frontend yourdomain.com

# Background Worker
./scripts/create-microservice.sh my-worker worker
```

### Copiar Templates Manualmente

```bash
# Copiar template de API
cp -r templates/microservice/api/* my-new-service/

# Copiar GitHub Action
cp templates/github-actions/deploy.yml my-new-service/.github/workflows/

# Editar configurações conforme necessário
```

## 📚 Documentação

- [API Template](microservice/api/README.md)
- [Frontend Template](microservice/frontend/README.md)  
- [Worker Template](microservice/worker/README.md)
- [GitHub Actions](github-actions/README.md)
