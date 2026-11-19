# Battertui

A lightweight TUI for viewing and switching Linux power profiles via `powerprofilectl` built with Gl and ["Bubble Tea v2"](https://charm.land/bubbletea/v2).

![`$ battertui`](images/battertui.png)

## Features

- Displays the currently active power profile (`[*]`) and available profiles (`[ ]`)
- Switch power frofiles directly from the terminal
- Manual refresh support and clear status/error feedback
- Runs in an alternate terminal screen (`AltScreen`)

## Prerequisites

- **Linux** with `power-profiles-daemon` installed and running (providing the `powerprofilesctl` CLI)
- **Go** (if building directly with `go build` / `go run`) or **Nix** with flakes enabled

## Installation & Running

### Using Nix Flakes

Run directly without installing:

```bash
nix run .
```

Or build the package:

```bash
nix build .
./result/bin/battertui
```

Enter the development shell (includes `go`, `gopls`, `gotools`, `staticcheck`, and `power-profiles-daemon`):

```bash
nix develop
```

### Using Go

```bash
go run ./src
```

Or build a binary:

```bash
go build -o bin/battertui ./src
./bin/battertui
```

## Keybindings

| Key | Action |
| --- | --- |
| `up` / `down` | Move the cursor through available profiles |
| `enter` / `space` | Apply the selected power profile |
| `r` | Refresh active and available profiles |
| `q` / `ctrl+c` | Quit Battertui |

### How It Works

- Queries the active profile with `powerprofilesctl get` and parses available profiles from `powerprofilesctl list`.
- Applies changes via `powerprofilesctl set <profile>` and re-fetches state upon completion.
