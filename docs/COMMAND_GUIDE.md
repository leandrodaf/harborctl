# 🔧 Guia de Comandos - HarborCtl

Referência rápida de todos os comandos.

## 🏗️ Comandos do Servidor Base

Execute no servidor onde a infraestrutura roda:

### Setup Inicial
```bash
# Configurar infraestrutura base
harborctl init-server --domain seudominio.com --email admin@seudominio.com

# Subir infraestrutura
harborctl up -f server-base.yml

# Validar configuração
harborctl validate -f server-base.yml
```

### Gerenciamento
```bash
# Status dos serviços
harborctl status

# Controle de ciclo de vida
harborctl stop      # Para serviços (mantém containers)
harborctl start     # Inicia serviços parados
harborctl restart   # Reinicia todos os serviços
harborctl pause     # Pausa execução
harborctl unpause   # Resume execução
harborctl down      # Para e remove tudo

# Escalar serviços
harborctl scale SERVICO --replicas N
```

### Utilitários
```bash
# Gerar senha para auth básica
harborctl hash-password --password "minhasenha"

# Auditoria de segurança
harborctl security-audit

# Renderizar compose (debug)
harborctl render -f server-base.yml

# Documentação
harborctl docs
```

## 🚀 Comandos de Deploy

Execute de qualquer lugar (local, CI/CD, etc):

### Deploy de Apps
```bash
# Deploy via repositório
harborctl deploy-service \
  --service minha-api \
  --repo https://github.com/usuario/minha-api.git

# Deploy com branch específica
harborctl deploy-service \
  --service minha-api \
  --repo https://github.com/usuario/minha-api.git \
  --branch develop

# Deploy local (código já baixado)
harborctl deploy-service \
  --service minha-api \
  --path deploy

# Deploy com overrides
harborctl deploy-service \
  --service minha-api \
  --replicas 5 \
  --env-file .env.prod \
  --force

# Dry run (apenas validar)
harborctl deploy-service \
  --service minha-api \
  --dry-run
```

### Desenvolvimento Local
```bash
# Inicializar novo projeto
harborctl init --project meu-projeto --domain localhost

# Testar localmente
harborctl up -f deploy/stack.yml

# Validar configuração
harborctl validate -f deploy/stack.yml
```

## 📊 Comandos de Monitoramento

```bash
# Status geral
harborctl status

# Logs de serviço específico  
harborctl logs SERVICO --tail 50 --follow

# Métricas de uso
harborctl stats SERVICO

# Health check
harborctl health SERVICO
```

## 🔐 Comandos de Segurança

```bash
# Auditoria completa
harborctl security-audit

# Validar configuração
harborctl validate -f stack.yml

# Gerar senhas seguras
harborctl hash-password --generate
harborctl hash-password --password "senha123"

# Verificar permissões
harborctl check-permissions
```

## 🎛️ Parâmetros Globais

Todos os comandos aceitam:

```bash
# Arquivo de configuração específico
--config stack.yml
-f stack.yml

# Modo verboso
--verbose
-v

# Modo silencioso
--quiet
-q

# Ajuda do comando
--help
-h

# Versão
--version
```

## 🔄 Fluxos Comuns

### Setup Inicial Completo
```bash
# 1. No servidor
harborctl init-server --domain seudominio.com --email admin@seudominio.com
harborctl up -f server-base.yml

# 2. Configurar app (local)
mkdir minha-api && cd minha-api
cp templates/microservice/api/deploy/stack.yml deploy/stack.yml
cp templates/github-actions/deploy.yml .github/workflows/deploy.yml

# 3. Deploy automático
git add . && git commit -m "Setup" && git push origin main
```

### Deploy Manual
```bash
# Desenvolvimento
harborctl deploy-service --service minha-api --path deploy

# Produção
harborctl deploy-service --service minha-api --repo https://github.com/usuario/minha-api.git --replicas 3
```

### Troubleshooting
```bash
# Verificar status
harborctl status

# Ver logs
harborctl logs minha-api --tail 100

# Reiniciar serviço
harborctl restart minha-api

# Auditoria
harborctl security-audit

# Validar config
harborctl validate -f deploy/stack.yml
```

### Manutenção
```bash
# Parar tudo
harborctl down

# Backup de configs
cp server-base.yml backup/

# Atualizar infraestrutura
harborctl up -f server-base.yml

# Verificar saúde
harborctl security-audit
harborctl status
```

## 📚 Exemplos por Cenário

### GitHub Actions
```yaml
- name: Deploy
  run: |
    harborctl deploy-service \
      --service ${{ github.event.repository.name }} \
      --repo ${{ github.server_url }}/${{ github.repository }} \
      --replicas 2
```

### Jenkins Pipeline
```groovy
stage('Deploy') {
    steps {
        sh '''
            harborctl deploy-service \
              --service ${JOB_NAME} \
              --repo ${GIT_URL} \
              --branch ${BRANCH_NAME}
        '''
    }
}
```

### Docker Compose Local
```bash
# Testar antes do deploy
harborctl validate -f deploy/stack.yml
harborctl up -f deploy/stack.yml
curl http://localhost:3000/health
harborctl down
```

## 🆘 Comandos de Emergência

```bash
# Parar tudo imediatamente
harborctl down --force

# Logs de emergência
harborctl logs --all --tail 1000 > emergency.log

# Status detalhado
harborctl status --verbose > status.log

# Backup rápido
tar -czf backup-$(date +%Y%m%d).tar.gz server-base.yml .deploy/

# Restaurar do backup
harborctl down
tar -xzf backup-20240828.tar.gz
harborctl up -f server-base.yml
```