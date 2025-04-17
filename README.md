# gum - Go Utility Manager

`gum` a utility manager for Go.

## Installation

You can install `gum` with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/baj-/gum/main/scripts/install.sh | bash
```

After installation, add `gum` to your PATH by adding this line to your shell profile (`.bashrc`, `.zshrc`, etc.):

```bash
export PATH="$HOME/.gum/bin:$PATH"
```

## Usage

### Install a Go version

```bash
gum install 1.24.2
```

### Use a specific Go version

```bash
gum use 1.24.2
```

### Uninstall a Go version

```bash
gum uninstall 1.24.2
```

## License

[MIT License](LICENSE)
