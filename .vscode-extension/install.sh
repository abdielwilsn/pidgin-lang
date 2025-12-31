#!/bin/bash

# Pidgin-Lang VSCode Extension Installer

set -e

EXTENSION_DIR="$HOME/.vscode/extensions/pidgin-lang-0.1.0"
SOURCE_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "🚀 Installing Pidgin-Lang VSCode Extension..."

# Remove old version if exists
if [ -d "$EXTENSION_DIR" ]; then
    echo "📦 Removing old version..."
    rm -rf "$EXTENSION_DIR"
fi

# Create extension directory
echo "📁 Creating extension directory..."
mkdir -p "$EXTENSION_DIR"

# Copy extension files
echo "📋 Copying extension files..."
cp -r "$SOURCE_DIR"/* "$EXTENSION_DIR/"

echo "✅ Installation complete!"
echo ""
echo "Next steps:"
echo "1. Reload VSCode: Press Cmd+Shift+P (Mac) or Ctrl+Shift+P (Windows/Linux)"
echo "2. Type 'Reload Window' and press Enter"
echo "3. Open a .pdg file to see syntax highlighting!"
echo ""
echo "Example test file: examples/hello.pdg"
