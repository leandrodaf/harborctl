# Scripts do Harbor CLI

Scripts de automação para facilitar o uso do Harbor CLI.

## 📂 Scripts Disponíveis

### 🚀 install.sh
Instalador automático do Harbor CLI.

```bash
# Instalação automática
curl -sSL https://github.com/company/harborctlr/raw/main/scripts/install.sh | bash

# Ou download e execução manual
wget https://github.com/company/harborctlr/raw/main/scripts/install.sh
chmod +x install.sh
./install.sh
```

### 🏗️ setup-production-server.sh
Setup completo da infraestrutura do servidor.

```bash
# Setup do servidor
./scripts/setup-production-server.sh yourdomain.com admin@yourdomain.com

# O que faz:
# - Instala dependências (Docker, Docker Compose)
# - Configura usuário harbor
# - Gera server-base.yml
# - Inicia infraestrutura base
```

### 📦 create-microservice.sh
Cria estrutura completa de microserviço.

```bash
# Criar API
./scripts/create-microservice.sh my-api api yourdomain.com

# Criar Frontend
./scripts/create-microservice.sh my-app frontend yourdomain.com

# Criar Worker
./scripts/create-microservice.sh my-worker worker

# O que faz:
# - Cria estrutura de pastas
# - Copia templates apropriados
# - Configura GitHub Actions
# - Gera arquivos de configuração
```

## 🎯 Fluxo de Uso

### 1. Setup Inicial (Admin)
```bash
# Instalar Harbor CLI
curl -sSL https://github.com/company/harborctlr/raw/main/scripts/install.sh | bash

# Setup do servidor
./scripts/setup-production-server.sh production.example.com admin@example.com
```

### 2. Criar Microserviços (Dev)
```bash
# Criar estrutura
./scripts/create-microservice.sh auth-service api production.example.com

# Implementar código
cd auth-service/src/
# ... código da aplicação ...

# Deploy
git add . && git commit -m "Initial commit"
git push origin main  # Deploy automático
```

### 3. Deploy Manual (se necessário)
```bash
harborctl deploy-service --service auth-service --repo https://github.com/company/auth-service.git
```

## 🔧 Personalização

Todos os scripts podem ser editados conforme sua necessidade:

- **Modificar templates**: Edite `templates/microservice/`
- **Alterar configurações**: Edite os scripts diretamente
- **Adicionar tipos**: Crie novos templates em `templates/microservice/`

## 🆘 Troubleshooting

### Script não executa
```bash
# Verificar permissões
chmod +x scripts/*.sh

# Verificar dependências
# Docker, git, curl devem estar instalados
```

### Template não encontrado
```bash
# Verificar se templates existem
ls -la templates/microservice/

# Tipos suportados: api, frontend, worker
```

## 📚 Próximos Passos

- Consulte [README principal](../README.md) para visão geral
- Veja [documentação](../docs/) para detalhes
- Use [templates](../templates/) para personalizar
