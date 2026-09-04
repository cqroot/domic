# Domic

A config-based, cross-platform dotfiles manager.

## Installation

```bash
go install github.com/cqroot/domic@latest

# Install using a proxy
env GOPROXY=https://goproxy.cn,direct go install github.com/cqroot/domic@latest

# PowerShell
$env:GOPROXY="https://goproxy.cn,direct"; go install github.com/cqroot/domic@latest
```

## Quick start

```bash
# 1. Clone your dotfiles repo into the config directory
domic init git@github.com:you/dotfiles.git

# 2. Edit domic.toml to declare your apps (see Configuration below)

# 3. See what would change
domic status -v
domic diff

# 4. Apply
domic apply
```

The config directory is `$XDG_CONFIG_HOME/domic` on Linux/macOS (typically
`~/.config/domic`) and `%AppData%\domic` on Windows. Print it with
`domic configdir`.

## Commands

| Command            | Description                                              |
| ------------------ | -------------------------------------------------------- |
| `domic init <src>` | Clone a dotfiles repo into the config directory          |
| `domic status`     | List configured apps and their sync status               |
| `domic apply`      | Apply configured apps to their target paths              |
| `domic diff [app]` | Show diff between source and target for each app         |
| `domic configdir`  | Print the OS-specific config directory                   |

Global flags:

| Flag             | Description                                              |
| ---------------- | -------------------------------------------------------- |
| `--config <path>` | Use a specific config file (default `domic.toml`)        |
| `-v, --verbose`  | Show per-file status rows (overrides `verbose` in config) |

`init` clones a repository into the config directory. `apply` copies each
file from source to target, using `.domic-cache.json` to detect local
edits and refuse to overwrite them. `status` shows whether each app and
file is `ok`, `missing`, `modified`, `partial`, or `skipped`. `diff`
renders any template files and invokes the configured external diff
command (default `diff -u`) for every changed file; missing targets are
diffed against `/dev/null`. `configdir` prints the absolute path of the
OS-specific config directory.

Skipped apps are those whose `bin` is not on `$PATH` or which have no
`target` for the current OS.

## Configuration

### `domic.toml` — main config

Located at `$XDG_CONFIG_HOME/domic/domic.toml`. Maps from app name to a
table of options:

```toml
[nvim]
path = "nvim"
target.linux = "~/.config/nvim"

[cargo]
path = "cargo"
target.linux = "~/.cargo/config.toml"

[niri]
path = "niri"
target.linux = "~/.config/niri"
target.darwin = "~/Library/Application Support/niri"
template = true

[pip]
path = "pip"
target.linux = "~/.config/pip/pip.conf"
bin = "pip"   # skip this app if `pip` is not on $PATH
```

Per-app fields:

| Field            | Type   | Description                                                                  |
| ---------------- | ------ | ---------------------------------------------------------------------------- |
| `path`           | string | Path to the app's source directory or file, relative to the config directory |
| `target.linux`   | string | Destination path on Linux (`~` is expanded)                                  |
| `target.darwin`  | string | Destination path on macOS                                                    |
| `target.windows` | string | Destination path on Windows                                                  |
| `template`       | bool   | Render `*.tmpl` files as Go templates (default `false`)                      |
| `bin`            | string | If set, skip this app when the binary is not on `$PATH`                      |

An app is skipped when:

- It has no `target` for the current OS, or
- `bin` is set and that binary is not found.

### `domic.config.toml` — user preferences

Optional, located next to `domic.toml`. Used for personal preferences that
should not live in the shared dotfiles repo:

```toml
# Show per-file status rows by default (overridden by -v on the CLI)
verbose = true

# Use delta with side-by-side layout for `domic diff`
diff = ["delta", "--side-by-side"]
```

| Field     | Type     | Default          | Description                                                             |
| --------- | -------- | ---------------- | ----------------------------------------------------------------------- |
| `verbose` | bool     | `false`          | Show per-file status rows by default; overridden by `-v` on the CLI     |
| `diff`    | string[] | `["diff", "-u"]` | External diff command. First element is the program, rest are arguments |

## Templates

When `template = true` is set on an app, source files ending in `.tmpl` are
rendered as Go [`text/template`](https://pkg.go.dev/text/template) before
being written to the target, and the `.tmpl` suffix is stripped from the
target filename. Other files are copied verbatim.

One variable is exposed:

| Expression        | Value                                     |
| ----------------- | ----------------------------------------- |
| `{{ .domic.os }}` | Current OS (`linux`, `darwin`, `windows`) |

Example `~/.config/domic/alacritty/alacritty.toml.tmpl`:

```toml
# Generated for {{ .domic.os }}
[window]
padding = 8
```

This produces `alacritty.toml` (no `.tmpl`) at the target with the OS line
expanded.

## Caching

To safely decide whether to overwrite a target that has been edited locally,
`domic apply` records the MD5 of every file it distributes in
`.domic-cache.json` at the root of the config directory. You should add
this file to your dotfiles repo's `.gitignore`.

## Shell integration

Add the following to your shell rc to define `dmcd`, which switches the
current directory to the domic config directory:

### bash / zsh

```sh
dmcd() {
  cd "$(domic configdir)" || return
}
```

### fish

```fish
function dmcd
    cd (domic configdir)
end
```

After reloading the shell (or sourcing the rc file), running `dmcd` will
`cd` into `~/.config/domic` (or the platform-equivalent path).
