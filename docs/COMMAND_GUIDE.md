# 🔧 Command Guide - HarborCtl

Quick reference for all commands.

## 🏗️ Server Commands

Execute on the server where infrastructure runs:

### Initial Setup
```bash
# Interactive server setup (recommended)
harborctl setup

# Direct server setup
harborctl init-server --domain yourdomain.com --email admin@yourdomain.com

# Start infrastructure
harborctl up

# Validate configuration
harborctl validate
```

### Lifecycle Management
```bash
# Service status
harborctl status

# Service control
harborctl stop      # Stop services (keep containers)
harborctl start     # Start stopped services
harborctl restart   # Restart all services
harborctl pause     # Pause execution
harborctl unpause   # Resume execution
harborctl down      # Stop and remove everything

# Scale services
harborctl scale SERVICE --replicas N
```

### Utilities
```bash
# Generate password for basic auth
harborctl hash-password --password "mypassword"

# Security audit
harborctl security-audit

# Render compose (debug)
harborctl render

# Documentation
harborctl docs
```

## 📱 Project Commands

### Project Creation
```bash
# Interactive project creation (recommended)
harborctl init --interactive

# Direct project creation
harborctl init --domain example.com --email admin@example.com --project my-app --env production

# Local development
harborctl init --env local --domain localhost --project my-app

# With observability services
harborctl init --domain example.com --email admin@example.com --project my-app
```

### Project Configuration
```bash
# Edit server configuration
harborctl edit-server

# Available options:
# - Change domain
# - Update email
# - Configure authentication
# - Manage observability services
```

## 🚀 Deployment Commands

### Service Deployment
```bash
# Deploy from Git repository
harborctl deploy-service --service my-api --repo https://github.com/user/my-api.git

# Deploy with specific branch
harborctl deploy-service --service my-api --repo https://github.com/user/my-api.git --branch develop

# Deploy with environment file
harborctl deploy-service --service my-api --env-file .env.prod

# Deploy with secrets
harborctl deploy-service --service my-api --secrets-file .secrets

# Deploy with replicas
harborctl deploy-service --service my-api --replicas 3

# Force deployment (ignore warnings)
harborctl deploy-service --service my-api --force

# Dry run (validate only)
harborctl deploy-service --service my-api --dry-run
```

## 🔍 Monitoring Commands

### Logs and Status
```bash
# View service logs
harborctl logs SERVICE_NAME

# Follow logs in real-time
harborctl logs SERVICE_NAME --follow

# Remote logs (from another server)
harborctl remote-logs --host server.com --service my-api

# Service status with details
harborctl status --verbose
```

### Remote Management
```bash
# Remote status check
harborctl remote-status --host server.com

# Remote command execution
harborctl remote-control --host server.com --command "status"
```

## 🔧 Configuration Commands

### Authentication
```bash
# Generate password hash
harborctl hash-password --password "mysecurepassword"

# Output: $2a$10$... (bcrypt hash)
```

### Security
```bash
# Full security audit
harborctl security-audit

# Quick security check
harborctl security-audit --quick

# Audit specific configuration
harborctl security-audit --config path/to/config.yml
```

### Rendering and Validation
```bash
# Render Docker Compose
harborctl render

# Render specific file
harborctl render -f custom-stack.yml

# Validate configuration
harborctl validate

# Validate specific file
harborctl validate -f custom-stack.yml
```

## 📋 Command Flags Reference

### Common Flags
```bash
# File specification
-f, --file STRING     # Configuration file (default: stack.yml)
-o, --output STRING   # Output file

# Environment control
--env STRING          # Environment (local/production)
--domain STRING       # Base domain
--email STRING        # Email for certificates

# Service control
--replicas INT        # Number of replicas
--force               # Force operation
--dry-run             # Validate only

# Observability
--no-dozzle           # Disable log viewer
--no-beszel           # Disable monitoring

# Verbosity
--verbose             # Detailed output
--quiet               # Minimal output
```

### Init Command Flags
```bash
harborctl init [flags]

--interactive         # Interactive mode
--domain STRING       # Base domain
--email STRING        # Email for SSL certificates
--project STRING      # Project name (default: app)
--env STRING          # Environment (local/production)
--no-dozzle           # Don't include Dozzle
--no-beszel           # Don't include Beszel
```

### Deploy Service Flags
```bash
harborctl deploy-service [flags]

--service STRING      # Microservice name (required)
--repo STRING         # Repository URL
--branch STRING       # Repository branch (default: main)
--env-file STRING     # Environment variables file
--secrets-file STRING # Secrets file
--replicas INT        # Number of replicas
--force               # Force deployment ignoring warnings
--dry-run             # Validate only without deploying
```

### Scale Command Flags
```bash
harborctl scale SERVICE [flags]

--replicas INT        # Number of replicas (required)
```

## 💡 Usage Examples

### Complete Workflow
```bash
# 1. Setup server
harborctl setup

# 2. Create project
harborctl init --interactive

# 3. Deploy services
harborctl deploy-service --service api --repo https://github.com/company/api.git
harborctl deploy-service --service frontend --repo https://github.com/company/frontend.git

# 4. Scale as needed
harborctl scale api --replicas 3
harborctl scale frontend --replicas 2

# 5. Monitor
harborctl status --verbose
harborctl security-audit
```

### Development Workflow
```bash
# Local development
harborctl init --env local --domain localhost --project my-dev-project

# Deploy local services
harborctl deploy-service --service api --path ./api
harborctl up

# Check logs
harborctl logs api --follow
```

### Production Workflow
```bash
# Production setup
harborctl init-server --domain production.com --email admin@production.com
harborctl up

# Deploy production services
harborctl deploy-service --service api --repo https://github.com/company/api.git --replicas 3
harborctl deploy-service --service frontend --repo https://github.com/company/frontend.git --replicas 2

# Monitor and maintain
harborctl status
harborctl security-audit
harborctl logs api
```

## 🆘 Help and Documentation

```bash
# General help
harborctl --help

# Command-specific help
harborctl COMMAND --help

# Examples:
harborctl init --help
harborctl deploy-service --help
harborctl scale --help
```