# Installing Pidgin-Lang VSCode Extension

## Quick Install (Copy to Extensions Folder)

This is the easiest way to install the extension locally.

### macOS/Linux

```bash
# From the pidgin-lang repository root
cp -r .vscode-extension ~/.vscode/extensions/pidgin-lang-0.1.0

# Reload VSCode
# Press Cmd+Shift+P and type "Reload Window"
```

### Windows (PowerShell)

```powershell
# From the pidgin-lang repository root
Copy-Item -Recurse -Force .vscode-extension $env:USERPROFILE\.vscode\extensions\pidgin-lang-0.1.0

# Reload VSCode
# Press Ctrl+Shift+P and type "Reload Window"
```

### Windows (Command Prompt)

```cmd
xcopy /E /I .vscode-extension %USERPROFILE%\.vscode\extensions\pidgin-lang-0.1.0

REM Reload VSCode
REM Press Ctrl+Shift+P and type "Reload Window"
```

## Building and Installing VSIX Package

For a more permanent installation, create a VSIX package:

### Prerequisites

Install `vsce` (Visual Studio Code Extension Manager):

```bash
npm install -g @vscode/vsce
```

### Build VSIX

```bash
cd .vscode-extension

# Package the extension
vsce package

# This creates pidgin-lang-0.1.0.vsix
```

### Install VSIX

1. Open VSCode
2. Press `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)
3. Type "Extensions: Install from VSIX..."
4. Select `pidgin-lang-0.1.0.vsix`
5. Reload VSCode

## Verify Installation

1. Open or create a file with `.pdg` extension
2. Check the bottom-right corner of VSCode - it should show "Pidgin" as the language
3. Type some Pidgin code and verify syntax highlighting:

```pidgin
make x be 10
suppose x big pass 5 {
    yarn("E plenty!")
}
```

Keywords like `make`, `suppose`, `yarn` should be highlighted.

## Uninstalling

### If installed via extensions folder:

```bash
# macOS/Linux
rm -rf ~/.vscode/extensions/pidgin-lang-0.1.0

# Windows (PowerShell)
Remove-Item -Recurse -Force $env:USERPROFILE\.vscode\extensions\pidgin-lang-0.1.0
```

### If installed via VSIX:

1. Open VSCode
2. Go to Extensions (Cmd+Shift+X / Ctrl+Shift+X)
3. Search for "Pidgin-Lang"
4. Click Uninstall

## Publishing to VSCode Marketplace (Optional)

To publish the extension for all users:

1. Create a [Visual Studio Marketplace publisher account](https://marketplace.visualstudio.com/manage)

2. Get a Personal Access Token from Azure DevOps

3. Login with vsce:
   ```bash
   vsce login abdielwilsn
   ```

4. Publish:
   ```bash
   cd .vscode-extension
   vsce publish
   ```

After publishing, users can install via:
- VSCode Marketplace UI
- Command: `code --install-extension abdielwilsn.pidgin-lang`

## Troubleshooting

### Extension not showing up

- Make sure you reloaded VSCode after copying files
- Check the extension folder name matches exactly: `pidgin-lang-0.1.0`
- Verify the folder is in the correct location: `~/.vscode/extensions/`

### Syntax highlighting not working

- Verify the file extension is `.pdg` or `.pidgin`
- Check the language indicator in VSCode's bottom-right corner
- Try manually selecting the language:
  1. Click the language indicator
  2. Type "Pidgin" and select it

### VSIX build fails

Make sure you have Node.js and npm installed:
```bash
node --version
npm --version
```

Install vsce globally:
```bash
npm install -g @vscode/vsce
```

---

**Installation complete! Start coding in Pidgin! 🚀**
