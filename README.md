# Domic

## Installation

```bash
go install github.com/cqroot/domic@latest

# Install using a proxy
env GOPROXY=https://goproxy.cn,direct go install github.com/cqroot/domic@latest

# PowerShell
$env:GOPROXY="https://goproxy.cn,direct"; go install github.com/cqroot/domic@latest
```

## Shell integration

Add the following to your shell rc to define `dmcd`, which switches the
current directory to the domic config directory:

### bash / zsh

```sh
function dmcd() {
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
