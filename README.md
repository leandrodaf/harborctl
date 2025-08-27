# 🚢 Harbor CLI - Deployment Tool

> Uma ferramenta CLI moderna para orquestração e deploy de microserviços usando Docker Compose e Traefik.

## 🎯 O que é o Harbor CLI?

O Harbor CLI é uma ferramenta que automatiza o processo de deploy e gerenciamento de microserviços. Ele gera configurações Docker Compose otimizadas, configura roteamento automático com Traefik e oferece comandos simples para deploy remoto.

## 🏗️ Conceitos

### � Como Biblioteca (este repositório)
Este repositório contém o **código-fonte** do Harbor CLI:
- ✅ Build e release de binários
- ✅ Testes e validação
- ✅ Templates para microserviços
- ✅ Documentação da ferramenta

### 🚀 Como Ferramenta nos Microserviços
Os microserviços **usam** o Harbor CLI para deploy:
- ✅ GitHub Actions baixam binário do Harbor CLI
- ✅ Executam comandos de deploy remotamente
- ✅ Usam templates fornecidos por este repo

## 📥 Instalação

### Instalação Automática (Recomendado)
```bash
curl -sSL https://raw.githubusercontent.com/SEU-USUARIO/harbor-cli/main/scripts/install.sh | bash
```

### Download Manual
```bash
# Linux x64
curl -sSL https://github.com/SEU-USUARIO/harbor-cli/releases/latest/download/harborctl-linux-amd64 -o harborctl
chmod +x harborctl
sudo mv harborctl /usr/local/bin/

# macOS x64
curl -sSL https://github.com/SEU-USUARIO/harbor-cli/releases/latest/download/harborctl-darwin-amd64 -o harborctl
chmod +x harborctl
sudo mv harborctl /usr/local/bin/

# Windows x64
curl -sSL https://github.com/SEU-USUARIO/harbor-cli/releases/latest/download/harborctl-windows-amd64.exe -o harborctl.exe
```

### ✅ Verificar Instalação
```bash
harborctl --version
```

## � Quick Start

### 1️⃣ Servidor (Comando Local)
```bash
# Configurar servidor de produção
harborctl init-server --domain exemplo.com

# Subir infraestrutura
harborctl up

# Verificar status
harborctl status
```

### 2️⃣ Microserviço (Comando Remoto)
```bash
# Criar novo microserviço
harborctl init --name minha-api --type node

# Deploy de microserviço
harborctl deploy-service \
  --host servidor.exemplo.com \
  --service minha-api \
  --image ghcr.io/usuario/minha-api:latest
```

## 📚 Documentação

| Documento | Descrição |
|-----------|-----------|
| [📖 Quick Start](docs/QUICK_START.md) | Primeiros passos e exemplos práticos |
| [📘 Guia Completo](docs/GUIDE.md) | Documentação detalhada |
| [⚡ Guia de Comandos](docs/COMMAND_GUIDE.md) | Referência de todos os comandos |

## 🛠️ Comandos Principais

### 🖥️ Comandos do Servidor (Local)
```bash
# Inicializar servidor
harborctl init-server --domain exemplo.com

# Gerenciar infraestrutura
harborctl up          # Subir serviços
harborctl down        # Derrubar serviços
harborctl status      # Ver status
harborctl scale       # Escalar serviços
```

### 🚀 Comandos Remotos
```bash
# Deploy de microserviço
harborctl deploy-service \
  --host servidor.com \
  --service api-users \
  --image ghcr.io/company/api-users:v1.2.0

# Criar microserviço
harborctl init \
  --name nova-api \
  --type python \
  --template fastapi
```

## 🎨 Templates Disponíveis

### 📁 Microserviços
```bash
# Criar microserviço Node.js
harborctl init --name minha-api --type node

# Criar microserviço Python
harborctl init --name minha-api --type python --template fastapi

# Criar microserviço Go
harborctl init --name minha-api --type go
```

### ⚙️ GitHub Actions
Os templates de GitHub Actions estão em `templates/github-actions/`:

- **deploy.yml**: Pipeline completo de CI/CD
- **auto-scale.yml**: Monitoramento e auto-scaling

#### Como usar nos microserviços:
```bash
# Copiar template para seu microserviço
cp templates/github-actions/deploy.yml .github/workflows/

# Personalizar variáveis no arquivo
# Configurar secrets no GitHub:
# - PRODUCTION_HOST
# - PRODUCTION_USER  
# - PRODUCTION_SSH_KEY
```

## 🔧 Scripts de Automação

| Script | Descrição |
|--------|-----------|
| `scripts/install.sh` | Instalação automática do Harbor CLI |
| `scripts/setup-production-server.sh` | Configuração de servidor de produção |
| `scripts/create-microservice.sh` | Criação de microserviço completo |

## 🏗️ Desenvolvimento

### Requisitos
- Go 1.21+
- Docker
- Docker Compose

### Build Local
```bash
# Clone o repositório
git clone https://github.com/SEU-USUARIO/harbor-cli.git
cd harbor-cli

# Build
go build -o harborctl ./cmd/harborctl

# Testes
go test ./...
```

### Release
O release é automatizado via GitHub Actions:
1. Crie uma tag: `git tag v1.2.0`
2. Push da tag: `git push origin v1.2.0`
3. GitHub Actions gera binários para todas as plataformas

## 📄 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

---

## 🆘 Suporte

- 📖 [Documentação Completa](docs/)
- 🐛 [Reportar Bugs](https://github.com/SEU-USUARIO/harbor-cli/issues)
- 💡 [Solicitar Features](https://github.com/SEU-USUARIO/harbor-cli/issues/new)

---

<div align="center">
  <strong>🚢 Harbor CLI - Simplificando deploys de microserviços</strong>
</div>

- **Issues**: [GitHub Issues](https://github.com/company/harborctlr/issues)
- **Docs**: [Documentação](docs/)
- **Email**: devops@company.com
