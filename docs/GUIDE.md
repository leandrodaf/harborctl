# 📖 Complete Guide - HarborCtl

Complete system for microservice deployment.

## 🎯 Architecture

**Base Server (one time):**
- Traefik: Reverse proxy + automatic SSL
- Dozzle: Centralized logs
- Beszel: Real-time monitoring
- Isolated networks and volumes

**Applications (multiple):**
- Easy deployment from Git repositories
- Automatic integration with infrastructure
- Independent scalability
- Environment variable management

## 🚀 Complete Setup

### 1. Install on Server
```bash
curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64 -o harborctl
chmod +x harborctl && sudo mv harborctl /usr/local/bin/
```

### 2. Configure Base Infrastructure

#### Interactive Setup (Recommended)
```bash
harborctl setup
```

This will guide you through:
- Domain configuration
- Email for SSL certificates
- Project name
- Environment selection
- Authentication setup
- Observability services

#### Direct Setup (Alternative)
```bash
harborctl init-server --domain yourdomain.com --email admin@yourdomain.com
harborctl up
```

### 3. Create Projects

#### Interactive Project Creation
```bash
harborctl init --interactive
```

Follow the prompts to configure:
- Project name (validated format)
- Environment (local or production)
- Domain configuration
- Email for certificates (production only)
- Dozzle log viewer inclusion
- Beszel monitoring inclusion

#### Direct Project Creation
```bash
# Production project
harborctl init --domain yourdomain.com --email admin@yourdomain.com --project my-app --env production

# Local development project
harborctl init --env local --domain localhost --project my-dev-app
```

## 🔧 Configuration Management

### Server Configuration
```bash
# Edit existing server configuration
harborctl edit-server
```

You can modify:
- Domain settings
- Email configuration
- Authentication credentials
- Observability service settings
- SSL certificate configuration

### Project Configuration
Projects are configured via `stack.yml` files that define:
- Services and their configurations
- Domain routing
- Environment variables
- Scaling parameters
- Resource limits

## 🚀 Service Deployment

### Basic Deployment
```bash
# Deploy from Git repository
harborctl deploy-service --service my-api --repo https://github.com/user/my-api.git

# Deploy specific branch
harborctl deploy-service --service my-api --repo https://github.com/user/my-api.git --branch develop
```

### Advanced Deployment
```bash
# Deploy with environment variables
harborctl deploy-service --service my-api --env-file .env.prod

# Deploy with secrets
harborctl deploy-service --service my-api --secrets-file .secrets

# Deploy with scaling
harborctl deploy-service --service my-api --replicas 3

# Force deployment (ignore warnings)
harborctl deploy-service --service my-api --force
```

### Environment Files
Create `.env` files for your services:
```env
# .env.prod
NODE_ENV=production
DATABASE_URL=postgres://user:pass@db:5432/myapp
API_KEY=your-api-key-here
LOG_LEVEL=info
```

### Secrets Files
Create `.secrets` files for sensitive data:
```env
# .secrets
DB_PASSWORD=super-secure-password
JWT_SECRET=jwt-secret-key
OAUTH_CLIENT_SECRET=oauth-secret
```

## 📊 Monitoring and Logs

### Built-in Services
- **Dozzle**: Web-based log viewer at `https://logs.yourdomain.com`
- **Beszel**: Monitoring dashboard at `https://monitor.yourdomain.com`

### Command Line Monitoring
```bash
# Check service status
harborctl status

# Detailed status
harborctl status --verbose

# View logs
harborctl logs my-api

# Follow logs in real-time
harborctl logs my-api --follow
```

### Remote Monitoring
```bash
# Check remote server status
harborctl remote-status --host production.com

# View remote logs
harborctl remote-logs --host production.com --service my-api

# Execute remote commands
harborctl remote-control --host production.com --command "status"
```

## 🔒 Security

### Security Audit
```bash
# Full security audit
harborctl security-audit

# Quick security check
harborctl security-audit --quick
```

The security audit checks:
- SSL/TLS configuration
- Authentication settings
- File permissions
- Network configuration
- Container security
- Repository security (if config files present)

### Authentication Management
```bash
# Generate password hash for basic auth
harborctl hash-password --password "mysecurepassword"

# Edit authentication in server
harborctl edit-server  # Follow prompts for auth setup
```

### Best Practices
1. **Use strong passwords** for all authentication
2. **Keep certificates updated** (automatic with Let's Encrypt)
3. **Regular security audits** with `harborctl security-audit`
4. **Environment separation** (use different domains for dev/prod)
5. **Secrets management** (use secrets files, not environment variables)

## ⚡ Scaling and Performance

### Service Scaling
```bash
# Scale up
harborctl scale my-api --replicas 5

# Scale down
harborctl scale my-api --replicas 1

# Check current scaling
harborctl status
```

### Load Balancing
Traefik automatically load balances across replicas:
- Sticky sessions available
- Health checks included
- Automatic failover

### Resource Management
Configure in your `stack.yml`:
```yaml
services:
  - name: my-api
    resources:
      memory: "512m"
      cpus: "0.5"
      reserve_mem: "256m"
      reserve_cpu: "0.25"
```

## 🔄 Workflow Examples

### Development Workflow
```bash
# 1. Setup local environment
harborctl init --env local --domain localhost --project my-dev-project

# 2. Deploy development services
harborctl deploy-service --service api --path ./api-dev
harborctl deploy-service --service frontend --path ./frontend-dev

# 3. Start services
harborctl up

# 4. Monitor during development
harborctl logs api --follow
harborctl status
```

### Production Workflow
```bash
# 1. Setup production server (one time)
harborctl setup  # Interactive production setup

# 2. Deploy production services
harborctl deploy-service --service api --repo https://github.com/company/api.git --replicas 3
harborctl deploy-service --service frontend --repo https://github.com/company/frontend.git --replicas 2

# 3. Configure monitoring
harborctl security-audit
harborctl status --verbose

# 4. Ongoing maintenance
harborctl logs api  # Check logs
harborctl scale api --replicas 5  # Scale as needed
```

### CI/CD Integration
```bash
# In your CI/CD pipeline
harborctl deploy-service \
  --service $SERVICE_NAME \
  --repo $GITHUB_REPOSITORY \
  --branch $GITHUB_REF_NAME \
  --replicas $REPLICAS \
  --env-file .env.${ENVIRONMENT}
```

## 🛠️ Troubleshooting

### Common Issues

**Service won't start:**
```bash
# Check logs
harborctl logs service-name

# Check configuration
harborctl validate

# Render and inspect compose
harborctl render
```

**SSL issues:**
```bash
# Check domain configuration
harborctl edit-server

# Restart Traefik
harborctl restart traefik
```

**Deployment failures:**
```bash
# Try with force flag
harborctl deploy-service --service my-api --force

# Check security audit
harborctl security-audit
```

### Debug Commands
```bash
# Validate configuration
harborctl validate

# Render Docker Compose (for inspection)
harborctl render

# Security audit
harborctl security-audit

# Detailed status
harborctl status --verbose
```

## 📚 Advanced Topics

### Custom Domains
Configure multiple domains in your stack configuration or use the edit server command to manage domain routing.

### SSL Certificates
- Automatic Let's Encrypt certificates for production
- Custom certificates supported
- Local development uses HTTP

### Network Configuration
- Isolated networks per project
- Automatic service discovery
- Traefik routing integration

### Volume Management
- Persistent volumes for data
- Automatic backup strategies
- Volume mounting for development

## 🔗 Integration

### Git Integration
- Direct deployment from Git repositories
- Branch-specific deployments
- Automatic updates via webhooks

### Docker Integration
- Docker Compose generation
- Multi-stage builds supported
- Image optimization

### Monitoring Integration
- Prometheus metrics (via Beszel)
- Log aggregation (via Dozzle)
- Custom monitoring endpoints

## 📞 Support

For additional help:
- Check the [Command Guide](COMMAND_GUIDE.md) for specific commands
- Run `harborctl COMMAND --help` for detailed command help
- Use `harborctl docs` for inline documentation
- File issues on GitHub for bugs or feature requests