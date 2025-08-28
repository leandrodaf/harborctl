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

### 1️⃣ Setup (One Command)
```bash
# Interactive setup (recommended)
harborctl init --interactive

# Or direct setup
harborctl init --domain example.com --email admin@example.com --project my-app --env production

# Start infrastructure
harborctl up

# Check status
harborctl status
```

### 2️⃣ Deploy Services
```bash
# Deploy services from Git
harborctl deploy-service --service my-api --repo https://github.com/user/my-api.git

# Deploy with environment files
harborctl deploy-service --service my-api --env-file .env.prod

# Scale services
harborctl scale my-api --replicas 3
```

## 📚 Documentation

| Document | Description |
|-----------|-----------|
| [📖 Quick Start](docs/QUICK_START.md) | First steps and practical examples |
| [📘 Complete Guide](docs/GUIDE.md) | Detailed documentation |
| [⚡ Command Guide](docs/COMMAND_GUIDE.md) | Reference for all commands |

## 🛠️ Essential Commands
```bash
# Initialize project (interactive or direct)
harborctl init --interactive
harborctl init --domain example.com --email admin@example.com --project my-app

# Manage infrastructure
harborctl up          # Start services
harborctl down        # Stop services
harborctl status      # View status
harborctl scale       # Scale services

# Deploy services
harborctl deploy-service --service api-users --repo https://github.com/company/api-users.git

# Edit configuration
harborctl edit-server
```

## 🎨 Available Features

### 📁 Project Types
```bash
# Local development environment
harborctl init --env local --domain localhost

# Production environment
harborctl init --env production --domain example.com --email admin@example.com
```

### 🔧 Services Included
- **Traefik**: Automatic routing and SSL certificates
- **Dozzle**: Centralized log viewer 
- **Beszel**: Monitoring and metrics
- **Your Apps**: Custom microservices deployment

### ⚙️ Deployment Options
```bash
# Deploy from Git repository
harborctl deploy-service --service my-api --repo https://github.com/user/my-api.git

# Deploy with environment file
harborctl deploy-service --service my-api --env-file .env.prod

# Deploy with secrets
harborctl deploy-service --service my-api --secrets-file .secrets

# Scale service
harborctl scale my-api --replicas 3
```

## 🔧 Additional Commands

| Command | Description |
|---------|-----------|
| `harborctl remote-logs` | View logs from remote servers |
| `harborctl remote-control` | Control services on remote servers |
| `harborctl security-audit` | Run security audits |
| `harborctl hash-password` | Generate password hashes |
| `harborctl docs` | Show inline documentation |

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
