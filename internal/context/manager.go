package context

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/1outres/cldenv/internal/config"
	"github.com/1outres/cldenv/pkg/symlink"
)

var (
	ErrContextNotFound      = errors.New("context not found")
	ErrContextAlreadyExists = errors.New("context already exists")
	ErrInvalidContextName   = errors.New("invalid context name")
)

// Context represents a cldenv context
type Context struct {
	Name     string
	Path     string
	IsActive bool
	Files    []string
}

// Manager manages cldenv contexts
type Manager struct {
	claudeDir  string
	cldenvDir  string
	contexts   []Context
}

// NewManager creates a new context manager
func NewManager() (*Manager, error) {
	claudeDir, err := config.GetClaudeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get claude directory: %w", err)
	}

	cldenvDir, err := config.GetCldenvDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cldenv directory: %w", err)
	}

	return &Manager{
		claudeDir: claudeDir,
		cldenvDir: cldenvDir,
	}, nil
}

// LoadContexts loads all available contexts
func (m *Manager) LoadContexts() error {
	if !config.FileExists(m.cldenvDir) {
		return nil // No contexts directory exists yet
	}

	entries, err := os.ReadDir(m.cldenvDir)
	if err != nil {
		return fmt.Errorf("failed to read contexts directory: %w", err)
	}

	m.contexts = nil
	activeContext := m.getActiveContext()

	for _, entry := range entries {
		if entry.IsDir() {
			// Skip .git directory
			if entry.Name() == ".git" {
				continue
			}
			
			contextPath := filepath.Join(m.cldenvDir, entry.Name())
			files := m.getContextFiles(contextPath)
			
			context := Context{
				Name:     entry.Name(),
				Path:     contextPath,
				IsActive: entry.Name() == activeContext,
				Files:    files,
			}
			
			m.contexts = append(m.contexts, context)
		}
	}

	return nil
}

// GetContexts returns all loaded contexts
func (m *Manager) GetContexts() []Context {
	return m.contexts
}

// GetActiveContext returns the name of the currently active context
func (m *Manager) GetActiveContext() string {
	return m.getActiveContext()
}

// getActiveContext determines which context is currently active
func (m *Manager) getActiveContext() string {
	claudeFilePath, err := config.GetClaudeFilePath()
	if err != nil {
		return ""
	}

	settingsFilePath, err := config.GetSettingsFilePath()
	if err != nil {
		return ""
	}

	// Check if files are symlinks and point to a context
	if symlink.IsValidSymlink(claudeFilePath) {
		target, err := symlink.ReadSymlink(claudeFilePath)
		if err == nil {
			return m.extractContextName(target)
		}
	}

	if symlink.IsValidSymlink(settingsFilePath) {
		target, err := symlink.ReadSymlink(settingsFilePath)
		if err == nil {
			return m.extractContextName(target)
		}
	}

	return ""
}

// extractContextName extracts context name from a file path
func (m *Manager) extractContextName(path string) string {
	// Path should be like: ~/.cldenv/context-name/filename
	rel, err := filepath.Rel(m.cldenvDir, path)
	if err != nil {
		return ""
	}
	
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) >= 1 {
		return parts[0]
	}
	
	return ""
}

// getContextFiles returns the list of files in a context directory
func (m *Manager) getContextFiles(contextPath string) []string {
	var files []string
	
	claudeFile := filepath.Join(contextPath, config.ClaudeFile)
	if config.FileExists(claudeFile) {
		files = append(files, config.ClaudeFile)
	}
	
	settingsFile := filepath.Join(contextPath, config.SettingsFile)
	if config.FileExists(settingsFile) {
		files = append(files, config.SettingsFile)
	}
	
	return files
}

// CreateContext creates a new context
func (m *Manager) CreateContext(name string) error {
	if err := ValidateContextName(name); err != nil {
		return err
	}

	contextPath := filepath.Join(m.cldenvDir, name)
	if config.FileExists(contextPath) {
		return ErrContextAlreadyExists
	}

	if err := config.CreateDir(contextPath); err != nil {
		return fmt.Errorf("failed to create context directory: %w", err)
	}

	return nil
}

// RemoveContext removes a context
func (m *Manager) RemoveContext(name string) error {
	if name == config.DefaultContext {
		return fmt.Errorf("cannot remove default context")
	}

	contextPath := filepath.Join(m.cldenvDir, name)
	if !config.FileExists(contextPath) {
		return ErrContextNotFound
	}

	// Don't remove if it's the active context
	if m.getActiveContext() == name {
		return fmt.Errorf("cannot remove active context")
	}

	if err := os.RemoveAll(contextPath); err != nil {
		return fmt.Errorf("failed to remove context directory: %w", err)
	}

	return nil
}

// SwitchContext switches to a different context
func (m *Manager) SwitchContext(name string) error {
	contextPath := filepath.Join(m.cldenvDir, name)
	if !config.FileExists(contextPath) {
		return ErrContextNotFound
	}

	claudeFilePath, err := config.GetClaudeFilePath()
	if err != nil {
		return fmt.Errorf("failed to get claude file path: %w", err)
	}

	settingsFilePath, err := config.GetSettingsFilePath()
	if err != nil {
		return fmt.Errorf("failed to get settings file path: %w", err)
	}

	// Create symlinks to the new context
	contextClaudeFile := filepath.Join(contextPath, config.ClaudeFile)
	contextSettingsFile := filepath.Join(contextPath, config.SettingsFile)

	if err := symlink.CreateSymlink(contextClaudeFile, claudeFilePath); err != nil {
		return fmt.Errorf("failed to create symlink for CLAUDE.md: %w", err)
	}

	if err := symlink.CreateSymlink(contextSettingsFile, settingsFilePath); err != nil {
		// Try to rollback the first symlink
		symlink.RemoveSymlink(claudeFilePath)
		return fmt.Errorf("failed to create symlink for settings.json: %w", err)
	}

	return nil
}

// ContextExists checks if a context exists
func (m *Manager) ContextExists(name string) bool {
	contextPath := filepath.Join(m.cldenvDir, name)
	return config.FileExists(contextPath)
}