# gum - Go Utility Manager

![gum version](https://img.shields.io/github/v/release/baj-/gum?label=version&color=77DD77)
![License: MIT](https://img.shields.io/badge/License-MIT-8BC6FC.svg)


![Build Status](https://github.com/baj-/gum/actions/workflows/release.yml/badge.svg)
![Tests](https://github.com/baj-/gum/actions/workflows/test.yml/badge.svg)

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

You can install Go versions in several ways:

```bash
# Install a specific version (major.minor.patch)
gum install 1.24.2

# Install the latest patch version for a major.minor release
gum install 1.24    # Will install the latest 1.24.x version

# Install with or without the "go" prefix
gum install go1.24
gum install 1.24
```

When you specify only a major.minor version (like `1.24`), `gum` will automatically find and install the latest patch version available for that release.

### Use a specific Go version

```bash
gum use 1.24.2
```

or

```bash
gum use
```
The version specified in your `go.mod` will be set as active.

### Uninstall a Go version

```bash
gum uninstall 1.24.2
```

### List installed versions

```bash
gum list
```

## License

[MIT License](LICENSE)
