package context

import (
	"fmt"

	"github.com/1outres/cldenv/internal/config"
	"github.com/1outres/cldenv/pkg/symlink"
)

// IsFirstRun checks if this is the first run of cldenv
func IsFirstRun() bool {
	claudeFilePath, err := config.GetClaudeFilePath()
	if err != nil {
		return false
	}

	settingsFilePath, err := config.GetSettingsFilePath()
	if err != nil {
		return false
	}

	// Check if either file exists and is NOT a symlink
	claudeExists := config.FileExists(claudeFilePath)
	settingsExists := config.FileExists(settingsFilePath)

	if claudeExists && !config.IsSymlink(claudeFilePath) {
		return true
	}

	if settingsExists && !config.IsSymlink(settingsFilePath) {
		return true
	}

	return false
}

// MigrateToDefault migrates existing files to the default context
func MigrateToDefault() error {
	claudeFilePath, err := config.GetClaudeFilePath()
	if err != nil {
		return fmt.Errorf("failed to get claude file path: %w", err)
	}

	settingsFilePath, err := config.GetSettingsFilePath()
	if err != nil {
		return fmt.Errorf("failed to get settings file path: %w", err)
	}

	// Create default context directory
	defaultContextPath, err := config.GetContextDir(config.DefaultContext)
	if err != nil {
		return fmt.Errorf("failed to get default context path: %w", err)
	}

	if err := config.CreateDir(defaultContextPath); err != nil {
		return fmt.Errorf("failed to create default context directory: %w", err)
	}

	// Migrate CLAUDE.md if it exists and is not a symlink
	if config.FileExists(claudeFilePath) && !config.IsSymlink(claudeFilePath) {
		defaultClaudeFile, err := config.GetContextFilePath(config.DefaultContext, config.ClaudeFile)
		if err != nil {
			return fmt.Errorf("failed to get default claude file path: %w", err)
		}

		if err := config.MoveFile(claudeFilePath, defaultClaudeFile); err != nil {
			return fmt.Errorf("failed to move CLAUDE.md to default context: %w", err)
		}

		// Create symlink
		if err := symlink.CreateSymlink(defaultClaudeFile, claudeFilePath); err != nil {
			return fmt.Errorf("failed to create symlink for CLAUDE.md: %w", err)
		}
	}

	// Migrate settings.json if it exists and is not a symlink
	if config.FileExists(settingsFilePath) && !config.IsSymlink(settingsFilePath) {
		defaultSettingsFile, err := config.GetContextFilePath(config.DefaultContext, config.SettingsFile)
		if err != nil {
			return fmt.Errorf("failed to get default settings file path: %w", err)
		}

		if err := config.MoveFile(settingsFilePath, defaultSettingsFile); err != nil {
			return fmt.Errorf("failed to move settings.json to default context: %w", err)
		}

		// Create symlink
		if err := symlink.CreateSymlink(defaultSettingsFile, settingsFilePath); err != nil {
			return fmt.Errorf("failed to create symlink for settings.json: %w", err)
		}
	}

	return nil
}

// EnsureDefaultContext ensures the default context exists with proper symlinks
func EnsureDefaultContext() error {
	// Check if default context directory exists
	defaultContextPath, err := config.GetContextDir(config.DefaultContext)
	if err != nil {
		return fmt.Errorf("failed to get default context path: %w", err)
	}

	if err := config.CreateDir(defaultContextPath); err != nil {
		return fmt.Errorf("failed to create default context directory: %w", err)
	}

	claudeFilePath, err := config.GetClaudeFilePath()
	if err != nil {
		return fmt.Errorf("failed to get claude file path: %w", err)
	}

	settingsFilePath, err := config.GetSettingsFilePath()
	if err != nil {
		return fmt.Errorf("failed to get settings file path: %w", err)
	}

	defaultClaudeFile, err := config.GetContextFilePath(config.DefaultContext, config.ClaudeFile)
	if err != nil {
		return fmt.Errorf("failed to get default claude file path: %w", err)
	}

	defaultSettingsFile, err := config.GetContextFilePath(config.DefaultContext, config.SettingsFile)
	if err != nil {
		return fmt.Errorf("failed to get default settings file path: %w", err)
	}

	// Create empty files in default context if they don't exist
	if !config.FileExists(defaultClaudeFile) {
		if err := config.EnsureDir(defaultClaudeFile); err != nil {
			return fmt.Errorf("failed to ensure directory for CLAUDE.md: %w", err)
		}
		// Create empty file
		if err := config.CopyFile("/dev/null", defaultClaudeFile); err != nil {
			// If /dev/null doesn't work, create empty file manually
			if err := config.CreateDir(defaultContextPath); err != nil {
				return fmt.Errorf("failed to create default context directory: %w", err)
			}
		}
	}

	if !config.FileExists(defaultSettingsFile) {
		if err := config.EnsureDir(defaultSettingsFile); err != nil {
			return fmt.Errorf("failed to ensure directory for settings.json: %w", err)
		}
		// Create empty JSON file
		if err := config.CopyFile("/dev/null", defaultSettingsFile); err != nil {
			// If /dev/null doesn't work, create empty file manually
			if err := config.CreateDir(defaultContextPath); err != nil {
				return fmt.Errorf("failed to create default context directory: %w", err)
			}
		}
	}

	// Create symlinks if they don't exist or are broken
	if !symlink.IsValidSymlink(claudeFilePath) {
		if err := symlink.CreateSymlink(defaultClaudeFile, claudeFilePath); err != nil {
			return fmt.Errorf("failed to create symlink for CLAUDE.md: %w", err)
		}
	}

	if !symlink.IsValidSymlink(settingsFilePath) {
		if err := symlink.CreateSymlink(defaultSettingsFile, settingsFilePath); err != nil {
			return fmt.Errorf("failed to create symlink for settings.json: %w", err)
		}
	}

	return nil
}