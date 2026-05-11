---
title: Installation
order: 2
template: docs
---

## Using Nix

If you have Nix installed, the easiest way is to use `nix run`:

```bash
nix run github:oojdir/oojsite -- --help
```

## Manual Installation

If you have Go 1.24 or later:

```bash
go install github.com/oojdir/oojsite/cmd/oojsite@latest
```

Then run:

```bash
oojsite --help
```

## From Source

Clone the repository and build:

```bash
git clone https://github.com/oojdir/oojsite.git
cd oojsite
go build -o oojsite .
```

The binary will be created in the current directory.

## Using `allDir` for Convenience

The `--allDir` flag is a shortcut for projects where all content lives under one directory:

```bash
oojsite --allDir docs --dev
```

This automatically sets:
- `--postDir docs/posts`
- `--pageDir docs/site`
- `--templateDir docs/templates`
- `--componentDir docs/components`
- `--staticDir docs/static`

All relative to the `docs` directory.
