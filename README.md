# 🚢 HarborCtl - Docker Compose Deployment Tool

> A modern CLI tool for orchestrating and deploying microservices using Docker Compose and Traefik.

## 🎯 What is HarborCtl?

HarborCtl is a tool that automates the process of deploying and managing microservices. It generates optimized Docker Compose configurations, sets up automatic routing with Traefik, and provides simple commands for remote deployment.

## 🏗️ Concepts

### 📚 As a Library (this repository)
This repository contains the **source code** of HarborCtl:
- ✅ Build and release of binaries
- ✅ Testing and validation
- ✅ Templates for microservices
- ✅ Tool documentation

### 🚀 As a Tool in Microservices
Microservices **use** HarborCtl for deployment:
- ✅ GitHub Actions download HarborCtl binary
- ✅ Execute deployment commands remotely
- ✅ Use templates provided by this repo

## 📥 Installation

### Super Quick Installation (Direct Binary)

**For amd64 (Intel/AMD):**
```bash
sudo curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64 -o /usr/local/bin/harborctl && sudo chmod +x /usr/local/bin/harborctl
```

**For arm64 (ARM64):**
```bash
sudo curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_arm64 -o /usr/local/bin/harborctl && sudo chmod +x /usr/local/bin/harborctl
```

### Automatic Installation (Compressed Archive)

**For amd64 (Intel/AMD):**
```bash
curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64.tar.gz | sudo tar -xzC /usr/local/bin harborctl
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

### ✅ Verify Installation
```bash
harborctl --version
harborctl --help
```

## 🚀 Quick Start

### 1️⃣ Server (Local Command)
```bash
# Configure production server
harborctl init-server --domain example.com

# Start infrastructure
harborctl up

# Check status
harborctl status
```

### 2️⃣ Microservice (Remote Command)
```bash
# Create new microservice
harborctl init --name my-api --type node

# Deploy microservice
harborctl deploy-service \
  --host server.example.com \
  --service my-api \
  --image ghcr.io/user/my-api:latest
```

## 📚 Documentation

| Document | Description |
|-----------|-----------|
| [📖 Quick Start](docs/QUICK_START.md) | First steps and practical examples |
| [📘 Complete Guide](docs/GUIDE.md) | Detailed documentation |
| [⚡ Command Guide](docs/COMMAND_GUIDE.md) | Reference for all commands |

## 🛠️ Main Commands

### 🖥️ Server Commands (Local)
```bash
# Initialize server
harborctl init-server --domain example.com

# Manage infrastructure
harborctl up          # Start services
harborctl down        # Stop services
harborctl status      # View status
harborctl scale       # Scale services
```

### 🚀 Remote Commands
```bash
### 🚀 Remote Commands
```bash
# Deploy microservice
harborctl deploy-service 
  --host server.com 
  --service api-users 
  --image ghcr.io/company/api-users:v1.2.0

# Create microservice
harborctl init 
  --name new-api 
  --type python 
  --template fastapi
```

## 🎨 Available Templates

### 📁 Microservices
```bash
# Create Node.js microservice
harborctl init --name my-api --type node

# Create Python microservice
harborctl init --name my-api --type python --template fastapi

# Create Go microservice
harborctl init --name my-api --type go
```

### ⚙️ GitHub Actions
GitHub Actions templates are in `templates/github-actions/`:

- **deploy.yml**: Complete CI/CD pipeline
- **auto-scale.yml**: Monitoring and auto-scaling

#### How to use in microservices:
```bash
# Copy template to your microservice
cp templates/github-actions/deploy.yml .github/workflows/

# Customize variables in the file
# Configure secrets in GitHub:
# - PRODUCTION_HOST
# - PRODUCTION_USER  
# - PRODUCTION_SSH_KEY
```

## 🔧 Automation Scripts

| Script | Description |
|--------|-----------|
| `scripts/install.sh` | Automatic HarborCtl installation |
| `scripts/setup-production-server.sh` | Production server configuration |
| `scripts/create-microservice.sh` | Complete microservice creation |

## 🏗️ Development

### Requirements
- Go 1.24+
- Docker
- Docker Compose

### Local Build
```bash
# Clone the repository
git clone https://github.com/leandrodaf/harborctl.git
cd harborctl

# Build
go build -o harborctl ./cmd/harborctl

# Tests
go test ./...
```

### Release
Release is automated via GitHub Actions:
1. Create a tag: `git tag v1.2.0`
2. Push the tag: `git push origin v1.2.0`
3. GitHub Actions generates binaries for all platforms

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

## 🆘 Support

- 📖 [Complete Documentation](docs/)
- 🐛 [Report Bugs](https://github.com/leandrodaf/harborctl/issues)
- 💡 [Request Features](https://github.com/leandrodaf/harborctl/issues/new)

---

<div align="center">
  <strong>🚢 HarborCtl - Simplifying microservice deployments</strong>
</div>
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

## 🔧 Automation Scripts

| Script | Description |
|--------|-----------|
| `scripts/install.sh` | Automatic HarborCtl installation |
| `scripts/setup-production-server.sh` | Production server configuration |
| `scripts/create-microservice.sh` | Complete microservice creation |

## 🏗️ Development

### Requirements
- Go 1.24+
- Docker
- Docker Compose

### Local Build
```bash
# Clone the repository
git clone https://github.com/leandrodaf/harborctl.git
cd harborctl

# Build
go build -o harborctl ./cmd/harborctl

# Tests
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

## 🆘 Support

- 📖 [Complete Documentation](docs/)
- 🐛 [Report Bugs](https://github.com/leandrodaf/harborctl/issues)
- 💡 [Request Features](https://github.com/leandrodaf/harborctl/issues/new)

---

<div align="center">
  <strong>🚢 HarborCtl - Simplifying microservice deployments</strong>
</div>
