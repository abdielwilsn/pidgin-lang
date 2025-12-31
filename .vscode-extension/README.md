# Pidgin-Lang VSCode Extension

Syntax highlighting for Pidgin-Lang - the Nigerian Pidgin programming language.

## Features

- **Syntax Highlighting** for `.pdg` and `.pidgin` files
- **Keyword Recognition**: `make`, `suppose`, `abi`, `dey do while`, `bring`
- **Operator Highlighting**: `be`, `no be`, `big pass`, `no reach`, `and`
- **Boolean Values**: `tru`, `lie`
- **Built-in Functions**: `yarn`, `len`, `type`
- **String Literals** with escape sequences
- **Numeric Literals**
- **Comment Support** (`#` line comments)

## Installation

### Option 1: Install from VSIX (Recommended)

1. Download the latest `.vsix` file from the [releases page](https://github.com/abdielwilsn/pidgin-lang/releases)
2. Open VSCode
3. Press `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)
4. Type "Extensions: Install from VSIX..."
5. Select the downloaded `.vsix` file

### Option 2: Install from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/abdielwilsn/pidgin-lang.git
   cd pidgin-lang/.vscode-extension
   ```

2. Copy to VSCode extensions folder:

   **macOS/Linux:**
   ```bash
   cp -r ../pidgin-lang/.vscode-extension ~/.vscode/extensions/pidgin-lang
   ```

   **Windows:**
   ```powershell
   xcopy /E /I .vscode-extension %USERPROFILE%\.vscode\extensions\pidgin-lang
   ```

3. Reload VSCode

### Option 3: Install from Marketplace (After Publishing)

Once published to the VSCode Marketplace:

1. Open VSCode
2. Go to Extensions (Cmd+Shift+X / Ctrl+Shift+X)
3. Search for "Pidgin-Lang"
4. Click Install

## Usage

Open any `.pdg` or `.pidgin` file in VSCode, and syntax highlighting will be applied automatically.

### Example

```pidgin
make name be "Chidi"
yarn("How far, " + name + "!")

suppose name be "Chidi" {
    yarn("Na correct guy!")
} abi {
    yarn("Who be this?")
}

make counter be 0
dey do while counter no reach 5 {
    yarn(counter)
    make counter be counter + 1
}
```

## Language Features

The extension provides syntax highlighting for:

- **Keywords**: Control flow and declarations
- **Operators**: Comparison and arithmetic
- **Literals**: Strings, numbers, booleans
- **Comments**: Line comments starting with `#`
- **Functions**: Built-in function names

## Color Themes

The extension uses semantic token types that work with any VSCode color theme:

- Keywords: Purple/Blue (theme dependent)
- Strings: Orange/Red
- Numbers: Green
- Comments: Gray/Green
- Functions: Yellow
- Operators: Cyan/Blue

## Development

To modify the extension:

1. Edit `.vscode-extension/syntaxes/pidgin.tmLanguage.json` for syntax rules
2. Edit `.vscode-extension/language-configuration.json` for brackets, auto-closing pairs
3. Reload VSCode to see changes

## Contributing

Contributions are welcome! Please submit issues and pull requests to the [main repository](https://github.com/abdielwilsn/pidgin-lang).

## License

MIT License - see [LICENSE](../LICENSE) for details.

## Links

- [Pidgin-Lang Repository](https://github.com/abdielwilsn/pidgin-lang)
- [Language Documentation](https://github.com/abdielwilsn/pidgin-lang#readme)
- [Report Issues](https://github.com/abdielwilsn/pidgin-lang/issues)

---

**Happy coding in Pidgin! 🚀**
