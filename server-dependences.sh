# 1) Instalação no servidor Linux (passo a passo + script)
# Requisitos

# Ubuntu 22.04+ (ou Debian 12+)

# Acesso sudo

# Domínio apontando para o IP do servidor (para emitir TLS via Let’s Encrypt)

# Instalação rápida (tudo de uma vez)

# Cole isso no seu servidor (edite a variável HARBORCTL_URL para o seu release do GitHub quando publicar o binário):


#!/usr/bin/env bash
set -euo pipefail

# ======= ajuste isto para o seu release =======
HARBORCTL_URL="${HARBORCTL_URL:-https://github.com/sua-org/harborctl/releases/download/v0.1.0/harborctl_linux_amd64}"
# ==============================================

sudo apt-get update -y
sudo apt-get install -y ca-certificates curl gnupg git

# Docker Engine + Compose plugin (repo oficial)
if ! command -v docker >/dev/null 2>&1; then
  sudo install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
    | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
    https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" \
    | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
  sudo apt-get update -y
  sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
  sudo usermod -aG docker "$USER" || true
fi

# Baixa e instala o CLI (harborctl)
TMP="$(mktemp -d)"
curl -fL "$HARBORCTL_URL" -o "$TMP/harborctl"
chmod +x "$TMP/harborctl"
sudo mv "$TMP/harborctl" /usr/local/bin/harborctl
echo "harborctl instalado em /usr/local/bin/harborctl"

# Diretório padrão do projeto (você pode trocar)
sudo mkdir -p /opt/harbor && sudo chown -R "$USER":"$USER" /opt/harbor

# (Opcional) firewall básico liberando 80/443
if command -v ufw >/dev/null 2>&1; then
  sudo ufw allow 80/tcp || true
  sudo ufw allow 443/tcp || true
fi

echo "OK. Saia e entre novamente no shell para o grupo 'docker' aplicar (ou 'newgrp docker')."
echo "Exemplo de uso:"
echo "  cd /opt/harbor"
echo "  harborctl init --domain seu-dominio.com --email voce@seu-dominio.com"
echo "  harborctl up"


# O Traefik sobe em container (nada para “instalar”). Certificados são salvos no volume nomeado do Traefik.
# Dozzle e Beszel também sobem em containers por padrão:

# Dozzle = painel de logs em tempo real (precisa montar docker.sock). 
# GitHub
# dozzle.dev

# Beszel = monitoramento (hub + agent). O agent precisa do TOKEN e KEY fornecidos pelo hub (ou token universal). 
# Beszel

# Traefik = reverse-proxy com provider Docker e roteamento por labels/Host.