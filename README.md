# backporter

A CLI tool for backporting git commits and pull requests to target branches.

## Features

- Backport commits by SHA or pull requests by number
- Interactive mode with branch and PR selection
- Support for GitHub and Forgejo/Gitea forges
- Configurable target branches (supports regex patterns)
- Cache of backported commits/PRs for tracking
- Colored terminal output

## Installation

### Binary releases

Download the latest release from the [releases page](https://codefloe.com/pat-s/backporter/releases).

### Container image

```bash
docker pull codefloe.com/pat-s/backporter:latest
```

### From source

```bash
go install codefloe.com/pat-s/backporter/cmd/backporter@latest
```

### Via homebrew

```bash
brew tap pat-s/homebrew-tap https://codefloe.com/pat-s/homebrew-tap
brew install pat-s/tap/backporter
```

## Usage

### Interactive mode

Run without arguments to start the interactive wizard:

```bash
backporter
```

### Backport a commit

```bash
backporter backport commit <sha> <target-branch>

# Or directly with SHA:
backporter <sha> <target-branch>
```

### Backport a pull request

```bash
backporter backport pr <pr-number> <target-branch>

# Or directly with PR number:
backporter <pr-number> <target-branch>
```

### List backported items

```bash
backporter list
backporter list --clear  # Clear cache
```

## Configuration

Configuration can be set globally (`~/.config/backporter/config.yaml`) or per-repository (`.backporter.yaml`).

```yaml
# Forge type: "github" or "forgejo"
forge_type: github

# Forgejo/Gitea instance URL (only for forgejo)
# forgejo_url: https://codeberg.org

# Default target branches (supports regex)
target_branches:
  - release-1.x
  - release-2.x
  - stable

# Default branch to work from
default_branch: main

# Git remote name
remote: origin

# Number of recent PRs in interactive mode
recent_pr_count: 10

# Cache settings
cache:
  enabled: true
  path: '' # Defaults to ~/.cache/backporter/history.json
```

## Authentication

Set the appropriate environment variable for your forge:

```bash
# GitHub
export GITHUB_TOKEN=<your-token>

# Forgejo/Gitea
export FORGEJO_TOKEN=<your-token>
```

## Global options

| Option         | Description                       |
| -------------- | --------------------------------- |
| `--config, -c` | Path to config file               |
| `--remote`     | Git remote name (default: origin) |
| `--log-level`  | Logging level (default: info)     |
| `--pretty`     | Pretty-printed debug output       |
| `--nocolor`    | Disable colored output            |

## License

MIT
