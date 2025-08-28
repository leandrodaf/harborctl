# Quick Start - HarborCtl

Get started with HarborCtl in 5 minutes.

## ⚡ Quick Setup

### 1. Install HarborCtl

**Super Quick Installation (Direct Binary):**
```bash
sudo curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64 -o /usr/local/bin/harborctl && sudo chmod +x /usr/local/bin/harborctl
```

**Auto-detect Architecture:**
```bash
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

curl -sSLf "https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_${ARCH}.tar.gz" | sudo tar -xzC /usr/local/bin harborctl
```

### 2. Server Setup (Admin)

```bash
# On production server
harborctl init-server --domain yourdomain.com --email admin@yourdomain.com
harborctl up -f server-base.yml
```

### 3. Microservice Deploy (Dev)

```bash
# Create microservice
./scripts/create-microservice.sh my-service api

# Deploy
harborctl deploy-service --service my-service --repo https://github.com/company/my-service.git
```

## 🎯 Practical Example

### API Microservice

```bash
# 1. Create structure
./scripts/create-microservice.sh auth-api api yourdomain.com

# 2. Configure secrets in GitHub
# DATABASE_PASSWORD, JWT_SECRET, API_KEY

# 3. Implement code in src/
# 4. Commit and push = automatic deploy

# Or manual deploy
harborctl deploy-service --service auth-api
```

### Result

- **API**: https://auth-api.yourdomain.com
- **Logs**: https://logs.yourdomain.com  
- **Metrics**: https://monitor.yourdomain.com

## 🔧 Essential Commands

```bash
# General status
harborctl status

# Deploy microservice  
harborctl deploy-service --service NAME

# Scale
harborctl scale NAME --replicas 5

# Logs
harborctl logs NAME --tail 50

# Help
harborctl docs
```

## 📚 Next Steps

- [📖 Complete Guide](GUIDE.md) - Detailed documentation
- [🔧 Scripts](../scripts/) - Automation
- [📋 Templates](../templates/) - Ready-to-use templates
