# @eeshm/typr

npm wrapper for the `typr` Go CLI.

## Install

```bash
npm i -g @eeshm/typr
```

or run directly:

```bash
npx @eeshm/typr
```

## Usage

```bash
typr
```

## How it works

During install, this package downloads a prebuilt binary from:

- `https://github.com/eeshm/typer-cli/releases`

The package version should match a Git tag/release version (e.g. `1.0.0` package -> `v1.0.0` release).
