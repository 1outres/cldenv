# cldenv

A CLI tool for managing multiple Claude Code environments by switching between different `~/.claude/CLAUDE.md` and `~/.claude/settings.json` configurations.

## Installation

```bash
go install github.com/1outres/cldenv/cmd/cldenv@latest
```

## Usage

### List contexts
```bash
cldenv
```

### Switch context
```bash
cldenv use <context-name>
```

### Create new context
```bash
cldenv create <context-name>
```

### Remove context
```bash
cldenv remove <context-name>
```

### Show help
```bash
cldenv help
```

## How it works

cldenv manages multiple Claude Code configurations by storing them in separate directories under `~/.cldenv/` and creating symbolic links in `~/.claude/` to the active context.

Each context contains:
- `CLAUDE.md` - Claude Code instructions
- `settings.json` - Claude Code settings

## Example

```bash
# Create work context
cldenv create work

# Switch to work context
cldenv use work

# List all contexts
cldenv

# Switch back to default
cldenv use default
```

## License

MIT
