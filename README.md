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
