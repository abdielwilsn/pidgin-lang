# Changelog

All notable changes to the Pidgin-Lang VSCode extension will be documented in this file.

## [0.1.0] - 2025-12-25

### Added
- Initial release of Pidgin-Lang syntax highlighting
- Support for `.pdg` and `.pidgin` file extensions
- Syntax highlighting for all Pidgin keywords:
  - Control flow: `suppose`, `abi`, `dey do while`, `bring`, `comot`
  - Declarations: `make`
  - Operators: `be`, `no be`, `big pass`, `no reach`, `and`
- Boolean literals: `tru`, `lie`
- Numeric literals (integers and decimals)
- String literals with escape sequences
- Built-in function names: `yarn`, `len`, `type`
- Line comment support (`#`)
- Auto-closing pairs for brackets and quotes
- Bracket matching for `{}` and `()`

### Language Features
- Works with all VSCode color themes
- TextMate grammar-based highlighting
- Configurable brackets and auto-closing pairs
- File association for Pidgin files

## [Unreleased]

### Planned Features
- Code snippets for common patterns
- IntelliSense/autocomplete
- Hover documentation
- Go to definition
- Find all references
- Code formatting
- Linting integration
- Debugger support
- Semantic highlighting

---

## Version History

- **0.1.0** - Initial release with syntax highlighting
