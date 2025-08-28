# 🚀 Quick Start - HarborCtl

Funcione em 3 minutos!

## 📦 1. Instalar no Servidor

```bash
# Download e instalar
curl -sSLf https://github.com/leandrodaf/harborctl/releases/latest/download/harborctl_linux_amd64 -o harborctl
chmod +x harborctl
sudo mv harborctl /usr/local/bin/
```

## 🏗️ 2. Configurar Servidor (Uma vez)

```bash
# Criar infraestrutura base
harborctl init-server --domain seudominio.com --email admin@seudominio.com
harborctl up -f server-base.yml

# ✅ Pronto! Servidor configurado com:
# • Traefik (proxy + SSL automático)
# • Logs: https://logs.seudominio.com
# • Monitor: https://monitor.seudominio.com
```

## 🚀 3. Deploy de Apps (GitHub Actions)

### Configurar Repositório da App

**1. Copiar template:**
```bash
# No seu repositório de microserviço
cp templates/microservice/api/deploy/stack.yml deploy/stack.yml
cp templates/github-actions/deploy.yml .github/workflows/deploy.yml
```

**2. Configurar GitHub Secrets:**
```
PRODUCTION_HOST=seuservidor.com
PRODUCTION_USER=deploy
PRODUCTION_SSH_KEY=sua-chave-ssh-privada
```

**3. Push = Deploy Automático!**
```bash
git add .
git commit -m "Setup deploy"
git push origin main
# ✅ Deploy automático ativado!
```

## 📱 4. Deploy Manual (Opcional)

```bash
# Deploy direto do repositório
harborctl deploy-service --service minha-api --repo https://github.com/usuario/minha-api.git

# Deploy local (para testes)
harborctl deploy-service --service minha-api --path deploy
```

## 🎯 Resultado Final

- **✅ Servidor:** Infraestrutura rodando
- **✅ Apps:** Deploy automático via Git
- **✅ SSL:** Certificados automáticos  
- **✅ Logs:** Centralizados e acessíveis
- **✅ Monitor:** Métricas em tempo real

**🔗 URLs de acesso:**
- Sua app: `https://app.seudominio.com`
- Logs: `https://logs.seudominio.com`
- Monitor: `https://monitor.seudominio.com`
