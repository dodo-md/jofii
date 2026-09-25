# jofii

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/dodo-md/jofii)](https://go.dev/)

jofii is a j-rock radio for the terminal. Pick a station, and it plays in the background while you work.

Stations come from [Radio Browser](https://www.radio-browser.info/). Streams that do not return audio are skipped.

## Install

You need [mpv](https://mpv.io/) on your `PATH`. On macOS:

```bash
brew install mpv
```

Then install jofii:

```bash
go install github.com/dodo-md/jofii@latest
```

The binary is named `jofii` and lands in `$(go env GOPATH)/bin`. If the shell cannot find it, add that directory to your `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

To run this repo without installing:

```bash
go run .
```

## Usage

```bash
jofii
```

| Key | Action |
| --- | --- |
| `↑` / `↓` | change station |
| `p` | pause or resume |
| `q` | quit and stop playback |

The screen shows the station that is playing. At the end of a page, `↓` loads the next set of stations.

## Help

Questions and bugs go in [GitHub issues](https://github.com/dodo-md/jofii/issues).

## License

[MIT](LICENSE). Maintained by [doruk](https://github.com/dodo-md).
