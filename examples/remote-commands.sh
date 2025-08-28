#!/bin/bash

# HarborCtl Remote Commands Usage Examples
# Replace values with your actual server data

# Server configuration
SERVER_HOST="your-server.com"
SERVER_USER="deploy"
SSH_KEY_PATH="~/.ssh/id_rsa"

echo "=== REMOTE COMMANDS EXAMPLES ==="
echo

echo "1. 📋 View real-time logs from all services:"
echo "   ./harborctl remote-logs --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --follow"
echo

echo "2. 📋 View logs from specific service:"
echo "   ./harborctl remote-logs --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --service traefik --tail 50"
echo

echo "3. 📋 Follow logs from specific service:"
echo "   ./harborctl remote-logs --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --service dozzle --follow"
echo

echo "4. 🎛️  Check services status:"
echo "   ./harborctl remote-control --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --action status"
echo

echo "5. 🎛️  Restart specific service:"
echo "   ./harborctl remote-control --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --action restart --service traefik"
echo

echo "6. 🎛️  View complete service details:"
echo "   ./harborctl remote-control --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --action details --service my-api"
echo

echo "7. 🎛️  Check overall system health:"
echo "   ./harborctl remote-control --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --action health"
echo

echo "8. 🎛️  Stop a service:"
echo "   ./harborctl remote-control --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --action stop --service my-api"
echo

echo "9. 🎛️  Start a service:"
echo "   ./harborctl remote-control --host $SERVER_HOST --user $SERVER_USER --key $SSH_KEY_PATH --action start --service my-api"
echo

echo "=== QUICK SETUP ==="
echo
echo "For easier usage, create aliases in your ~/.bashrc or ~/.zshrc:"
echo
echo 'alias remote-logs="./harborctl remote-logs --host '$SERVER_HOST' --user '$SERVER_USER' --key '$SSH_KEY_PATH'"'
echo 'alias remote-control="./harborctl remote-control --host '$SERVER_HOST' --user '$SERVER_USER' --key '$SSH_KEY_PATH'"'
echo
echo "Then you can use:"
echo "  remote-logs --service traefik --follow"
echo "  remote-control --action restart --service my-api"
echo

echo "=== AVAILABLE PARAMETERS ==="
echo
echo "remote-logs:"
echo "  --host       : Server IP/hostname (required)"
echo "  --user       : SSH user (default: root)"
echo "  --key        : SSH private key file"
echo "  --port       : SSH port (default: 22)"
echo "  --service    : Specific service name"
echo "  --follow     : Follow logs in real time"
echo "  --tail       : Number of lines to show (default: 100)"
echo "  --compose    : Compose file path (default: .deploy/compose.generated.yml)"
echo
echo "remote-control:"
echo "  --host       : Server IP/hostname (required)"
echo "  --user       : SSH user (default: root)"
echo "  --key        : SSH private key file"
echo "  --port       : SSH port (default: 22)"
echo "  --action     : status, restart, stop, start, details, health"
echo "  --service    : Specific service name"
echo "  --verbose    : Show detailed information"
echo "  --compose    : Compose file path (default: .deploy/compose.generated.yml)"
echo
