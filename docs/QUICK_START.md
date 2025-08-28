# 🚀 Quick Start - HarborCtl

Get up and running in 3 minutes!

## 📦 1. Install on Server

```bash
# Download and install
curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64 -o harborctl
chmod +x harborctl
sudo mv harborctl /usr/local/bin/
```

## 🏗️ 2. Server Setup (One Time)

### Interactive Setup (Recommended)
```bash
# Guided setup wizard
harborctl setup

# ✅ Done! Server configured with:
# • Traefik (proxy + automatic SSL)
# • Logs: https://logs.yourdomain.com
# • Monitor: https://monitor.yourdomain.com
```

### Direct Setup (Alternative)
```bash
# Create base infrastructure
harborctl init-server --domain yourdomain.com --email admin@yourdomain.com
harborctl up

# ✅ Done! Server configured
```

## 🚀 3. Create and Deploy Projects

### Interactive Project Creation
```bash
# Step-by-step project setup
harborctl init --interactive

# Follow the prompts:
# 1. Project name
# 2. Environment (Local/Production)
# 3. Domain configuration
# 4. Email for SSL certificates
# 5. Include Dozzle (log viewer)
# 6. Include Beszel (monitoring)

# Start your project
harborctl up
```

### Deploy Microservices
```bash
# Deploy from Git repository
harborctl deploy-service --service my-api --repo https://github.com/user/my-api.git

# Deploy with environment variables
harborctl deploy-service --service my-api --env-file .env.prod

# Deploy with scaling
harborctl deploy-service --service my-api --replicas 3
```

## 📱 4. Management Commands

```bash
# Check status
harborctl status

# View logs
harborctl logs my-api

# Scale services
harborctl scale my-api --replicas 5

# Security audit
harborctl security-audit

# Edit server configuration
harborctl edit-server
```

## 🎯 Final Result

- **✅ Server:** Infrastructure running
- **✅ Apps:** Easy deployment from Git
- **✅ SSL:** Automatic certificates  
- **✅ Logs:** Centralized and accessible
- **✅ Monitor:** Real-time metrics
- **✅ Scaling:** Simple service scaling

**🔗 Access URLs:**
- Your app: `https://app.yourdomain.com`
- Logs: `https://logs.yourdomain.com`
- Monitor: `https://monitor.yourdomain.com`

## 🆘 Next Steps

- Read the [Complete Guide](GUIDE.md) for advanced features
- Check the [Command Reference](COMMAND_GUIDE.md) for all commands
- Set up automated deployments with GitHub Actions