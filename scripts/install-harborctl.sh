#!/bin/bash
#
# HarborCtl Installation/Update Script
# 
# This script automatically detects your system architecture and installs 
# or updates HarborCtl to the latest version.
#
# Usage:
#   curl -sSLf https://raw.githubusercontent.com/leandrodaf/harborctl/main/scripts/install-harborctl.sh | bash
#   or
#   wget -qO- https://raw.githubusercontent.com/leandrodaf/harborctl/main/scripts/install-harborctl.sh | bash
#
# Options:
#   --force     Force reinstallation even if already installed
#   --version   Install specific version (e.g., --version v1.0.0)
#   --help      Show this help message
#

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="leandrodaf/harborctl"
BINARY_NAME="harborctl"
INSTALL_DIR="/usr/local/bin"
FORCE_INSTALL=false
SPECIFIC_VERSION=""

# Functions
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

show_help() {
    cat << EOF
HarborCtl Installation Script

This script automatically installs or updates HarborCtl to the latest version.

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --force          Force reinstallation even if already installed
    --version TAG    Install specific version (e.g., --version v1.0.0)
    --help           Show this help message

EXAMPLES:
    # Install latest version
    $0

    # Force reinstall latest version
    $0 --force

    # Install specific version
    $0 --version v1.2.0

    # Remote installation
    curl -sSLf https://raw.githubusercontent.com/${REPO}/main/scripts/install-harborctl.sh | bash

EOF
}

detect_architecture() {
    local arch
    arch=$(uname -m)
    
    case $arch in
        x86_64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            print_info "Supported architectures: x86_64 (amd64), aarch64/arm64"
            exit 1
            ;;
    esac
}

detect_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    case $os in
        linux)
            echo "linux"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            print_info "Currently only Linux is supported"
            exit 1
            ;;
    esac
}

check_dependencies() {
    local deps=("curl" "tar")
    local missing=()
    
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" >/dev/null 2>&1; then
            missing+=("$dep")
        fi
    done
    
    if [ ${#missing[@]} -ne 0 ]; then
        print_error "Missing required dependencies: ${missing[*]}"
        print_info "Please install the missing dependencies and try again"
        exit 1
    fi
}

get_latest_version() {
    print_info "Fetching latest version..."
    curl -sSLf "https://api.github.com/repos/${REPO}/releases/latest" | \
        grep -oE '"tag_name": "v[0-9]+\.[0-9]+\.[0-9]+"' | \
        cut -d'"' -f4
}

get_current_version() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        "$BINARY_NAME" --version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "unknown"
    else
        echo "not_installed"
    fi
}

check_sudo() {
    if [ "$EUID" -eq 0 ]; then
        return 0
    fi
    
    if ! command -v sudo >/dev/null 2>&1; then
        print_error "This script requires root privileges or sudo to install to $INSTALL_DIR"
        print_info "Please run as root or install sudo"
        exit 1
    fi
    
    if ! sudo -n true 2>/dev/null; then
        print_warning "This script will ask for sudo password to install to $INSTALL_DIR"
    fi
}

download_and_install() {
    local version="$1"
    local arch="$2"
    local os="$3"
    
    print_info "Downloading HarborCtl $version for ${os}/${arch}..."
    
    local download_url
    if [ "$version" = "latest" ]; then
        download_url="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}_latest_${os}_${arch}"
    else
        download_url="https://github.com/${REPO}/releases/download/${version}/${BINARY_NAME}_${version}_${os}_${arch}.tar.gz"
    fi
    
    local temp_dir
    temp_dir=$(mktemp -d)
    trap "rm -rf $temp_dir" EXIT
    
    if [ "$version" = "latest" ]; then
        # Direct binary download
        print_info "Downloading binary directly..."
        curl -sSLf "$download_url" -o "$temp_dir/$BINARY_NAME"
        chmod +x "$temp_dir/$BINARY_NAME"
    else
        # Compressed archive download
        print_info "Downloading and extracting archive..."
        curl -sSLf "$download_url" | tar -xzC "$temp_dir"
    fi
    
    if [ ! -f "$temp_dir/$BINARY_NAME" ]; then
        print_error "Downloaded file not found or extraction failed"
        exit 1
    fi
    
    print_info "Installing to $INSTALL_DIR..."
    if [ "$EUID" -eq 0 ]; then
        cp "$temp_dir/$BINARY_NAME" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        sudo cp "$temp_dir/$BINARY_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    print_success "Installation completed!"
}

verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version
        version=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        print_success "HarborCtl installed successfully!"
        print_info "Version: $version"
        print_info "Location: $(which $BINARY_NAME)"
        print_info ""
        print_info "Run '$BINARY_NAME --help' to get started!"
    else
        print_error "Installation verification failed"
        print_info "The binary was installed but is not in PATH"
        print_info "You may need to restart your shell or add $INSTALL_DIR to your PATH"
        exit 1
    fi
}

main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force)
                FORCE_INSTALL=true
                shift
                ;;
            --version)
                SPECIFIC_VERSION="$2"
                shift 2
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_info "HarborCtl Installation Script"
    print_info "=============================="
    
    # Check dependencies
    check_dependencies
    
    # Detect system
    local arch os
    arch=$(detect_architecture)
    os=$(detect_os)
    
    print_info "System detected: ${os}/${arch}"
    
    # Check current installation
    local current_version
    current_version=$(get_current_version)
    
    if [ "$current_version" != "not_installed" ]; then
        print_info "Current version: $current_version"
        
        if [ "$FORCE_INSTALL" = false ] && [ -z "$SPECIFIC_VERSION" ]; then
            local latest_version
            latest_version=$(get_latest_version)
            print_info "Latest version: $latest_version"
            
            if [ "$current_version" = "$latest_version" ]; then
                print_success "You already have the latest version installed!"
                print_info "Use --force to reinstall or --version to install a specific version"
                exit 0
            else
                print_info "An update is available: $current_version → $latest_version"
            fi
        fi
    else
        print_info "HarborCtl is not currently installed"
    fi
    
    # Check permissions
    check_sudo
    
    # Determine version to install
    local version_to_install
    if [ -n "$SPECIFIC_VERSION" ]; then
        version_to_install="$SPECIFIC_VERSION"
        print_info "Installing specific version: $version_to_install"
    else
        version_to_install="latest"
        print_info "Installing latest version"
    fi
    
    # Download and install
    download_and_install "$version_to_install" "$arch" "$os"
    
    # Verify installation
    verify_installation
}

# Run main function with all arguments
main "$@"
