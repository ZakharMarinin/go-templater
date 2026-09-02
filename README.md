# go-templater

`go-templater` is a fast CLI tool for scaffolding and sharing Go project layouts. Turn any directory into a reusable template — file structure, file contents, and `go.mod` dependencies — and drop it into new projects in seconds, through a clean interactive TUI.

## Features

- **Structure templates** — capture a directory tree from any existing project and replay it into a new one.
- **File content capture** — templates aren't just empty folders. When you create a structure template, file contents (`.golangci.yaml`, configs, boilerplate code, etc.) are captured automatically. An interactive tree lets you exclude entire subtrees from content capture with a single toggle — no need to click through every file one by one. Dotfiles (`.env`, `.env.local`, ...) are excluded from content capture by default to avoid leaking secrets, but can be toggled back on.
- **Binary/large-file safe** — files over 256 KB or detected as binary are never embedded; they're kept as empty placeholders.
- **Dependency templates** — snapshot the dependencies of an existing module (via `go list -m all`) or type them in by hand, then install them into any project with one command.
- **Interactive TUI** — list pickers, text inputs, spinners, and confirmation prompts, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- **Zero clutter** — internal errors and debug logs go to `~/.go-templater/logs/app.log` in JSON, never to your terminal.

## Installation

### Download a release binary (recommended)

Grab the archive for your OS/arch from the [Releases](https://github.com/ZakharMarinin/go-templater/releases) page, unpack it, and put the binary on your `PATH`:

```bash
tar -xzf go-templater_Linux_x86_64.tar.gz
sudo mv go-templater /usr/local/bin/
```

### Build from source

Requires Go (see `go.mod` for the minimum version).

```bash
git clone https://github.com/ZakharMarinin/go-templater.git
cd go-templater
make install   # go install ./cmd/go-templater -> $GOPATH/bin/go-templater
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

```bash
go-templater <command> <subcommand> [flags]
go-templater --help
```

### Create templates (`make`)

```bash
# Save the structure (and file contents) of a directory as a template
go-templater make struct --dir ./my-boilerplate

# Save the go.mod dependencies of a directory as a template
go-templater make deps -d ./existing-core-app
```

`--dir` / `-d` defaults to the current working directory when omitted.

`make struct` walks the target directory, then opens an interactive screen where you can exclude entire folders from content capture with a single toggle (`space`) before saving — the folder structure is always kept, only the file contents inside it are dropped. Press `enter` to confirm.

`make deps` with no `--dir` reads dependency lines from stdin (`github.com/user/pkg v1.2.3`) instead of scanning a module.

### Insert templates (`insert`)

```bash
# Replay a saved structure into a new project (also runs `go mod init`)
go-templater insert struct --dir ./new-project

# Install a saved set of dependencies into an existing module (needs go.mod)
go-templater insert deps -d ./microservice
```

### Manage templates (`remove`, `list`)

```bash
go-templater list           # table of every saved structure/deps template
go-templater remove struct  # pick a structure template to delete
go-templater remove deps    # pick a dependencies template to delete
```

## Where things live

`go-templater` is self-contained — no config files to manage. On first run it creates:

```text
~/.go-templater/
├── templates/
│   ├── structs/      # saved structure templates (.json)
│   └── deps/         # saved dependency templates (.json)
└── logs/
    └── app.log        # JSON error/debug log
```

## Shell aliases

```bash
alias gt="go-templater"
alias gti="go-templater insert"
alias gtm="go-templater make"
alias gtr="go-templater remove"
```

```bash
gtm struct -d ./my-app
gti deps -d ./microservice
gtr struct
```

## Development

```bash
make          # run unit tests (default target)
make lint     # golangci-lint
make vet      # go vet
make cover    # tests + per-function coverage report
make build    # build ./go-templater in the repo root
make run ARGS="make struct -d ./example"
make help     # list all targets
```
