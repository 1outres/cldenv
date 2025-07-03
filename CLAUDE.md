# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`cldenv` is a CLI tool for managing multiple Claude Code environments by switching between different `~/.claude/CLAUDE.md` and `~/.claude/settings.json` configurations using symbolic links.

## Core Architecture

- **Context Management**: Each environment (default, hobby, work) is stored in `~/.cldenv/[context-name]/`
- **Symlink Strategy**: Active configuration files are symlinked from `~/.claude/` to the appropriate context directory
- **CLI Interface**: Simple command structure for listing, creating, removing, and switching contexts

## Key Commands (Design Phase)

```bash
# Switch to a context
cldenv use [context-name]

# List all contexts and show active one
cldenv

# Create new context
cldenv create [context-name]

# Remove context
cldenv remove [context-name]

# Show help
cldenv help
```

## Development Approach

This project will be implemented in Go following the patterns established in the user's global CLAUDE.md:

- Use Go modules for dependency management
- Implement comprehensive error handling with `errors.Is/As`
- Include context.Context for operations that might be cancelled
- Write tests for all core functionality
- Use `golangci-lint` for code quality
- Follow Go naming conventions and package structure

## File Structure (Planned)

- `cmd/` - CLI command implementations
- `internal/` - Core business logic for context management
- `pkg/` - Public API if needed
- `docs/` - Design documents and specifications

## Testing Strategy

- Unit tests for context management logic
- Integration tests for file system operations
- CLI command testing with test fixtures
- Error handling validation for edge cases

## Key Implementation Considerations

1. **Initial Migration**: Detect non-symlinked files and migrate them to `default` context
2. **File Safety**: Validate symlink targets and handle broken links gracefully  
3. **Cross-Platform**: Ensure symlink operations work across different operating systems
4. **Atomic Operations**: Prevent corruption during context switches
5. **Backup Strategy**: Consider backup mechanisms for configuration files