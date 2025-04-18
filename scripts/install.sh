#!/bin/bash
set -e

# Detect OS and arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Find architecture
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Only support macOS and Linux
if [ "$OS" != "darwin" ] && [ "$OS" != "linux" ]; then
  echo "Unsupported operating system: $OS"
  exit 1
fi

# Get latest release version
VERSION=$(curl -s https://api.github.com/repos/baj-/gum/releases/latest | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4 | cut -c2-)

# Create installation directory
INSTALL_DIR="$HOME/.gum/bin"
mkdir -p "$INSTALL_DIR"

# Download and extract
echo "Downloading gum $VERSION for $OS/$ARCH..."
DOWNLOAD_URL="https://github.com/baj-/gum/releases/download/v$VERSION/gum-$VERSION-$OS-$ARCH.tar.gz"

# Create a temporary directory for extraction
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download and extract to temp directory
curl -L "$DOWNLOAD_URL" | tar xz -C "$TMP_DIR"

# Move binary and make it executable
mv "$TMP_DIR/gum" "$INSTALL_DIR/gum"
chmod +x "$INSTALL_DIR/gum"

echo "gum has been installed to $INSTALL_DIR/gum"
echo ""
echo "To use gum, add the following to your shell profile and source it:"
echo "  export PATH=\"\$HOME/.gum/bin:\$PATH\""
echo ""
echo "Available commands:"
echo "  gum install <version>   # Install Go version"
echo "  gum uninstall <version> # Uninstall Go version"
echo "  gum use <version>       # Use Go version (uses go.mod if no version provided)"
echo "  gum list                # List installed Go versions"
