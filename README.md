# jofii

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/dodo-md/jofii)](https://go.dev/)

Blow your shell away. jofii brings j-rock radio to your terminal.

## Requirements

- [Go](https://go.dev/dl/)
- [mpv](https://mpv.io/) on your `PATH` (`brew install mpv`)

## Install

```bash
go install github.com/dodo-md/jofii@latest
```

`~/go/bin` needs to be on your `PATH`. After that, `jofii` runs from anywhere.

From this repo:

```bash
go run .
```

## Usage

```bash
jofii
```

- `↑` / `↓` change station
- `p` pause
- `q` quit

Stations come from [Radio Browser](https://www.radio-browser.info/). Streams that do not return audio are skipped.
