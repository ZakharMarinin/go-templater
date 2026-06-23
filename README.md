# go-templater

`go-templater` is a lightweight, blazing-fast CLI tool written in Go to manage, scaffold, and share architecture project templates and dependencies. Built with a beautiful, interactive TUI (Terminal User Interface) driven by Cobra and Bubble Tea.

## Features

* **Structure Scaffolding:** Create architectural layout templates directly from your existing directories and insert them into new projects.
* **Dependency Automation:** Extract, save, and install Go dependencies (`go.mod`) seamlessly with interactive feedback.
* **Interactive TUI:** Smooth list selectors, clean spinners, and self-cleaning automated status notices using Charmbracelet components.
* **Zero Clutter:** Internal errors and technical logs are isolated completely in `~/.go-templater/logs/app.log`, keeping your standard output beautiful and professional.

---

## Installation

### For Developers (Via Go)

If you have Go installed on your machine, you can install the binary directly into your `$GOPATH/bin`:

```bash
go install [github.com/yourusername/go-templater/cmd/](https://github.com/yourusername/go-templater/cmd/)...

```

Make sure your shell profile (`.zshrc` or `.bashrc`) includes your Go bin directory:

```bash
export PATH=$PATH:$(go env GOPATH)/bin

```

### For Production Users (Standalone Binary)

Download the pre-compiled binary for your specific OS (Linux, macOS, Windows) from the **Releases** page, unpack it, and move it to your local executable path:

```bash
sudo mv go-templater /usr/local/bin/

```

---

## Usage & Commands

The tool uses a highly intuitive structural routing system: `go-templater <command> <subcommand> [flags]`.

### General Help

```bash
go-templater --help

```

### 1. Generating Templates (`make`)

Extract patterns from your current workspace or specific directories.

* **Create architecture structure template:**
```bash
go-templater make struct --dir ./my-boilerplate

```


* **Create a blueprint for dependencies:**
```bash
go-templater make deps -d ./existing-core-app

```


*Note: If the `--dir` (or `-d`) flag is omitted, the tool automatically fallbacks to your current working directory.*

### 2. Injecting Templates (`insert`)

Inject pre-saved structural baselines or load up requirements.

* **Insert boilerplate folder architectures:**
```bash
go-templater insert struct --dir ./new-project

```


* **Download and run background routines for dependency installation:**
```bash
go-templater insert deps

```



### 3. Managing Assets (`remove`)

Launches an interactive, beautiful TUI list selector allowing you to navigate through your configurations and clean out old templates smoothly.

```bash
go-templater remove struct
go-templater remove deps

```

---

## Configuration & Architecture

`go-templater` follows a standalone, self-sufficient configuration model. It doesn't rely on flaky `.env` local file lookups.

Upon first launch, it automatically initializes an isolated environment directory inside your user home scope:

```text
~/.go-templater/
├── templates/
│   ├── structs/      # Saved architecture templates (.json)
│   └── deps/         # Saved dependency files (.json)
└── logs/
    └── app.log       # App silent error output tracking (JSON format)

```

---

## Quality of Life (Shell Aliases)

Speed up your daily scaffolding workflow by adding the following aliases to your `~/.zshrc` or `~/.bashrc`:

```bash
alias gt="go-templater"
alias gti="go-templater insert"
alias gtm="go-templater make"
alias gtr="go-templater remove"

```

Now you can run operations with lightning speed:

```bash
gtr struct
gti deps -d ./microservice

```