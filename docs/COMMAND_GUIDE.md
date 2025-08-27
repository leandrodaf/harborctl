# Harbor CLI - Guia de Comandos

Este guia separa claramente os comandos por contexto de uso.

## 🏗️ Comandos do Servidor Base (Administrador)

Estes comandos devem ser executados **no servidor de produção** onde a infraestrutura base será instalada:

### Inicialização do Servidor

```bash
# Inicializar infraestrutura base do servidor
harborctl init-server --domain production.example.com --email admin@example.com

# Verificar configuração antes de aplicar
harborctl validate -f server-base.yml

# Renderizar configuração Docker Compose (debug)
harborctl render -f server-base.yml
```

### Gerenciamento da Infraestrutura Base

```bash
# Iniciar infraestrutura base (Traefik, Dozzle, Beszel)
harborctl up -f server-base.yml

# Parar infraestrutura base
harborctl down

# Ver status da infraestrutura
harborctl status

# Escalar serviços específicos
harborctl scale dozzle --replicas 2
harborctl scale beszel --replicas 1
```

### Utilitários do Servidor

```bash
# Gerar hash de senha para autenticação básica
harborctl hash-password --generate
harborctl hash-password --password "mypassword"

# Auditoria de segurança da configuração
harborctl security-audit -f server-base.yml

# Ver documentação dos comandos
harborctl docs
```

## 🚀 Comandos Remotos (Desenvolvedor)

Estes comandos podem ser executados **remotamente** (local ou CI/CD) para gerenciar microserviços:

### Setup Local de Microserviço

```bash
# Criar configuração inicial de microserviço (desenvolvimento)
harborctl init --project auth-service --domain localhost
```

### Deploy de Microserviços

```bash
# Deploy de microserviço para servidor remoto
harborctl deploy-service --service auth-service --repo https://github.com/company/auth-service.git

# Deploy com branch específica
harborctl deploy-service --service auth-service --repo https://github.com/company/auth-service.git --branch develop

# Deploy usando código local (se já clonado)
harborctl deploy-service --service auth-service

# Deploy com override de réplicas
harborctl deploy-service --service auth-service --replicas 5

# Deploy com arquivo de environment customizado
harborctl deploy-service --service auth-service --env-file production.env
```

### Gerenciamento de Microserviços

```bash
# Ver status de um microserviço específico
harborctl status --service auth-service

# Escalar microserviço específico
harborctl scale auth-service --replicas 10

# Parar microserviço específico
harborctl down auth-service

# Restartar microserviço específico
harborctl up auth-service
```

## 📋 Exemplos de Workflows

### 1. Setup Inicial do Servidor (Admin - Uma vez)

```bash
# No servidor de produção
cd /opt/harbor
git clone https://github.com/company/harborctlr.git .

# Inicializar infraestrutura base
harborctl init-server --domain production.company.com --email devops@company.com

# Aplicar infraestrutura
harborctl up -f server-base.yml

# Verificar se está funcionando
harborctl status
# Resultado esperado:
# ✅ traefik: 1/1 replicas running
# ✅ dozzle: 1/1 replicas running  
# ✅ beszel: 1/1 replicas running
```

### 2. Deploy de Novo Microserviço (Desenvolvedor)

```bash
# Em qualquer lugar (local, CI/CD, etc.)
# Deploy direto para servidor configurado
harborctl deploy-service \
  --service payment-service \
  --repo https://github.com/company/payment-service.git \
  --branch main \
  --replicas 3

# Verificar se deployou
harborctl status --service payment-service
```

### 3. Desenvolvimento Local

```bash
# No diretório do microserviço
harborctl init --project user-service --domain localhost

# Desenvolver e testar localmente
harborctl up -f deploy/stack.yml

# Deploy para servidor quando pronto
harborctl deploy-service --service user-service
```

### 4. Escalabilidade Sob Demanda

```bash
# Durante pico de tráfego
harborctl scale auth-service --replicas 20
harborctl scale payment-service --replicas 15

# Volta ao normal
harborctl scale auth-service --replicas 3
harborctl scale payment-service --replicas 2
```

### 5. Manutenção do Servidor

```bash
# Parar todos os microserviços (mantém infraestrutura base)
harborctl down --exclude-base

# Manutenção da infraestrutura base
harborctl down
# ... manutenção ...
harborctl up -f server-base.yml

# Restart de serviço específico
harborctl down dozzle
harborctl up dozzle
```

## 🔧 Configuração de Ambiente

### Variáveis de Ambiente para Deploy Remoto

```bash
# Para uso remoto, configure estas variáveis:
export HARBOR_SERVER_HOST="production.company.com"
export HARBOR_SERVER_USER="harbor"
export HARBOR_SSH_KEY_PATH="~/.ssh/harbor_key"

# Ou via arquivo de configuração ~/.harbor/config.yml
server:
  host: "production.company.com"
  user: "harbor"
  ssh_key: "~/.ssh/harbor_key"
  port: 22
```

### Estrutura de Arquivos

```
# No servidor (/opt/harbor/):
server-base.yml          # Configuração da infraestrutura base
.services/               # Códigos dos microserviços clonados
├── auth-service/
├── payment-service/
└── user-service/

# No desenvolvimento:
microservice-repo/
├── deploy/
│   ├── stack.yml        # Configuração do microserviço
│   ├── environments/    # Env por ambiente
│   └── secrets/         # Templates de secrets
├── src/                 # Código do microserviço
└── Dockerfile           # Build do microserviço
```

## 🛡️ Segurança e Boas Práticas

### Comandos de Servidor Base
- ✅ **Onde executar**: Apenas no servidor de produção
- ✅ **Quem executa**: Administrador de sistema
- ✅ **Quando executar**: Setup inicial e manutenção
- ✅ **SSH**: Acesso direto ao servidor

### Comandos Remotos
- ✅ **Onde executar**: Local, CI/CD, qualquer lugar
- ✅ **Quem executa**: Desenvolvedores, DevOps
- ✅ **Quando executar**: Deploy de features, escalabilidade
- ✅ **SSH**: Conexão remota configurada

### Isolamento
- 🔒 **Infraestrutura base**: Independente dos microserviços
- 🔒 **Microserviços**: Isolados uns dos outros
- 🔒 **Deploy**: Não afeta serviços rodando
- 🔒 **Configuração**: Secrets isolados por serviço

Este guia garante que você saiba exatamente onde e quando executar cada comando!
