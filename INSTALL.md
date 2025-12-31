# Installation Guide

## Prerequisites

- **Go 1.21 or higher** - [Download Go](https://golang.org/dl/)

Check your Go version:

```bash
go version
```

## Installation Methods

### Method 1: Install with `go install` (Easiest)

Once the repository is public on GitHub:

```bash
go install github.com/abdielwilsn/pidgin-lang@latest
```

This installs the `pidgin-lang` binary to your `$GOPATH/bin` directory.

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Method 2: Build from Source

#### Step 1: Clone the Repository

```bash
git clone https://github.com/abdielwilsn/pidgin-lang.git
cd pidgin-lang
```

#### Step 2: Build the Binary

```bash
go build -o pidgin
```

This creates a `pidgin` executable in the current directory.

#### Step 3: (Optional) Install System-Wide

**On macOS/Linux:**

```bash
sudo mv pidgin /usr/local/bin/
```

**On Windows:**
Move `pidgin.exe` to a directory in your PATH, or add the current directory to PATH.

### Method 3: Download Pre-built Binary

Download the latest release for your platform from the [Releases page](https://github.com/abdielwilsn/pidgin-lang/releases).

**Available platforms:**

- Linux (amd64, arm64)
- macOS (amd64, arm64/M1)
- Windows (amd64)

**Installation steps:**

1. Download the appropriate binary for your system
2. Extract the archive (if compressed)
3. Move the binary to a directory in your PATH
4. Make it executable (Linux/macOS only):
   ```bash
   chmod +x pidgin
   ```

## Verify Installation

Run the following command to verify Pidgin-Lang is installed correctly:

```bash
pidgin --version
```

Or start the REPL:

```bash
pidgin
```

You should see:

```
╔═══════════════════════════════════════════════════════════════╗
║                    PIDGIN-LANG v0.1                           ║
║         Na programming language wey dey use Pidgin            ║
╠═══════════════════════════════════════════════════════════════╣
...
```

## Development Installation

For contributing or development:

### 1. Clone and Setup

```bash
git clone https://github.com/abdielwilsn/pidgin-lang.git
cd pidgin-lang
go mod download
```

### 2. Run Tests

```bash
go test ./...
```

### 3. Run Benchmarks

```bash
go test -bench=. ./...
```

### 4. Build with Debug Info

```bash
go build -gcflags="all=-N -l" -o pidgin-debug
```

## Updating

### If installed with `go install`:

```bash
go install github.com/abdielwilsn/pidgin-lang@latest
```

### If built from source:

```bash
cd pidgin-lang
git pull
go build -o pidgin
```

## Uninstalling

### If installed with `go install`:

```bash
rm $(go env GOPATH)/bin/pidgin-lang
```

### If installed system-wide:

```bash
sudo rm /usr/local/bin/pidgin
```

### If built from source:

Simply delete the `pidgin-lang` directory.

## Troubleshooting

### "pidgin: command not found"

**Solution:** Add Go's bin directory to your PATH.

**On macOS/Linux**, add to your `~/.bashrc`, `~/.zshrc`, or `~/.profile`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Then reload:

```bash
source ~/.bashrc  # or ~/.zshrc
```

**On Windows**, add `%GOPATH%\bin` to your PATH environment variable.

### "go: cannot find main module"

**Solution:** Make sure you're in the `pidgin-lang` directory and `go.mod` exists.

```bash
cd pidgin-lang
ls go.mod  # Should exist
```

### Build fails with "missing dependencies"

**Solution:** Download dependencies first:

```bash
go mod download
go mod tidy
```

### Permission denied when running binary

**Solution (Linux/macOS):**

```bash
chmod +x pidgin
```

## Platform-Specific Notes

### macOS

- On Apple Silicon (M1/M2), use the `arm64` binary
- On Intel Macs, use the `amd64` binary
- You may need to allow the app in System Preferences → Security & Privacy

### Linux

- Ensure you have execute permissions: `chmod +x pidgin`
- For system-wide install, use `sudo`

### Windows

- Use PowerShell or Command Prompt
- Add the directory containing `pidgin.exe` to your PATH
- Alternatively, use WSL (Windows Subsystem for Linux)

## Next Steps

After installation, see:

- [README.md](README.md) - Quick start guide
- [examples/](examples/) - Example programs
- [PERFORMANCE.md](PERFORMANCE.md) - Performance details

---

**Wahala don finish! Start to code!** 🚀
