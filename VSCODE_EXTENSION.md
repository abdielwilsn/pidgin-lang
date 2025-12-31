# VSCode Syntax Highlighting for Pidgin-Lang

Your Pidgin-Lang repository now includes a VSCode extension for syntax highlighting!

## Quick Install (Already Done!)

The extension has been installed to your VSCode. Just reload VSCode to activate it:

1. **Press `Cmd+Shift+P`** (Mac) or **`Ctrl+Shift+P`** (Windows/Linux)
2. **Type "Reload Window"** and press Enter
3. **Open any `.pdg` file** - syntax highlighting will work automatically!

## What Gets Highlighted

### Keywords (Purple/Blue)
- `make` - variable declaration
- `suppose` - if statement
- `abi` - else
- `dey do while` - while loop
- `bring` - return statement
- `comot` - exit

### Operators (Cyan)
- `be` - equals (==)
- `no be` - not equals (!=)
- `big pass` - greater than (>)
- `no reach` - less than (<)
- `and` - logical AND
- `+`, `-`, `*`, `/` - arithmetic

### Literals
- **Numbers**: `42`, `3.14`, `-10` (Green)
- **Booleans**: `tru`, `lie` (Orange)
- **Strings**: `"Hello World"` (Red/Orange)

### Built-in Functions (Yellow)
- `yarn()` - print
- `len()` - length
- `type()` - type inspection

### Comments (Gray/Green)
- `# This is a comment`

## Example with Highlighting

Create a test file to see it in action:

**test.pdg**
```pidgin
# Counter example
make count be 0

dey do while count no reach 5 {
    yarn("Count: ")
    yarn(count)
    make count be count + 1
}

suppose count big pass 3 {
    yarn("We don plenty!")
}
```

## Manual Installation (For Other Users)

Other users can install the extension by:

### Method 1: Copy Extension Folder
```bash
# Clone your repo
git clone https://github.com/abdielwilsn/pidgin-lang.git
cd pidgin-lang

# Run install script
.vscode-extension/install.sh

# Reload VSCode
```

### Method 2: Build VSIX Package (Recommended for Distribution)

1. Install vsce:
   ```bash
   npm install -g @vscode/vsce
   ```

2. Build the package:
   ```bash
   cd .vscode-extension
   vsce package
   ```
   This creates `pidgin-lang-0.1.0.vsix`

3. Install in VSCode:
   - Open VSCode
   - Press `Cmd+Shift+P` / `Ctrl+Shift+P`
   - Type "Extensions: Install from VSIX..."
   - Select the `.vsix` file

### Method 3: Publish to VSCode Marketplace (For Wide Distribution)

To make it available to all VSCode users:

1. Create a publisher account at https://marketplace.visualstudio.com/manage

2. Get a Personal Access Token from Azure DevOps

3. Login and publish:
   ```bash
   vsce login abdielwilsn
   cd .vscode-extension
   vsce publish
   ```

After publishing, anyone can install via:
- VSCode Extensions Marketplace (search "Pidgin-Lang")
- Command line: `code --install-extension abdielwilsn.pidgin-lang`

## Extension Files

The extension is located in `.vscode-extension/`:

```
.vscode-extension/
├── package.json                      # Extension metadata
├── language-configuration.json       # Brackets, auto-closing pairs
├── syntaxes/
│   └── pidgin.tmLanguage.json       # Syntax highlighting rules
├── README.md                         # Extension documentation
├── INSTALL.md                        # Installation guide
└── install.sh                        # Quick install script
```

## Customizing Syntax Highlighting

To modify syntax highlighting:

1. Edit `.vscode-extension/syntaxes/pidgin.tmLanguage.json`
2. Run the install script again:
   ```bash
   .vscode-extension/install.sh
   ```
3. Reload VSCode

### Adding New Keywords

Add to the `keywords` section in `pidgin.tmLanguage.json`:

```json
{
  "name": "keyword.control.pidgin",
  "match": "\\b(suppose|abi|dey do while|make|bring|comot|YOUR_NEW_KEYWORD)\\b"
}
```

## Color Themes

The extension uses TextMate scopes that work with any VSCode theme:

- **Dark themes**: Keywords are typically blue/purple, strings orange, comments gray
- **Light themes**: Similar but with darker shades
- **Custom themes**: Respect the theme's color scheme

## File Associations

The extension automatically activates for:
- `.pdg` files
- `.pidgin` files

To associate other file extensions, add to VSCode settings:

```json
{
  "files.associations": {
    "*.pidgin-lang": "pidgin"
  }
}
```

## Troubleshooting

### Extension not loading
- Verify installation: Check `~/.vscode/extensions/pidgin-lang-0.1.0` exists
- Reload VSCode: `Cmd+Shift+P` → "Reload Window"
- Check VSCode version: Requires VSCode 1.60.0 or higher

### No syntax highlighting
- Verify file extension is `.pdg` or `.pidgin`
- Check language mode (bottom-right corner of VSCode)
- Manually select language: Click language indicator → "Pidgin"

### Changes not appearing
- Reinstall: `rm -rf ~/.vscode/extensions/pidgin-lang-0.1.0` then run install script
- Hard reload: `Cmd+Shift+P` → "Reload Window"

## Features Roadmap

Future enhancements:
- [ ] Code snippets (type `make` → auto-complete variable declaration)
- [ ] IntelliSense (autocomplete for keywords and built-ins)
- [ ] Bracket matching for `{}` blocks
- [ ] Error checking (linting)
- [ ] Go to definition
- [ ] Format on save
- [ ] Debugger integration

## Distribution Checklist

When sharing the extension:

- [ ] Test on clean VSCode install
- [ ] Create `.vsix` package: `cd .vscode-extension && vsce package`
- [ ] Add `.vsix` to GitHub Releases
- [ ] Update README with installation instructions
- [ ] Consider publishing to VSCode Marketplace
- [ ] Add screenshots to extension README

## Links

- **Extension Source**: `.vscode-extension/`
- **Installation Guide**: `.vscode-extension/INSTALL.md`
- **TextMate Language Grammar**: https://macromates.com/manual/en/language_grammars
- **VSCode Extension API**: https://code.visualstudio.com/api

---

**Your .pdg files now have beautiful syntax highlighting! 🎨**
