package config

import (
	"os"
	"path/filepath"
)

const (
	ClaudeDir     = ".claude"
	CldenvDir     = ".cldenv"
	ClaudeFile    = "CLAUDE.md"
	SettingsFile  = "settings.json"
	DefaultContext = "default"
)

// GetClaudeDir returns the Claude configuration directory path
func GetClaudeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ClaudeDir), nil
}

// GetCldenvDir returns the cldenv configuration directory path
func GetCldenvDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, CldenvDir), nil
}

// GetContextDir returns the path for a specific context
func GetContextDir(contextName string) (string, error) {
	cldenvDir, err := GetCldenvDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cldenvDir, contextName), nil
}

// GetClaudeFilePath returns the path to CLAUDE.md in Claude directory
func GetClaudeFilePath() (string, error) {
	claudeDir, err := GetClaudeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(claudeDir, ClaudeFile), nil
}

// GetSettingsFilePath returns the path to settings.json in Claude directory
func GetSettingsFilePath() (string, error) {
	claudeDir, err := GetClaudeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(claudeDir, SettingsFile), nil
}

// GetContextFilePath returns the path to a specific file in a context directory
func GetContextFilePath(contextName, filename string) (string, error) {
	contextDir, err := GetContextDir(contextName)
	if err != nil {
		return "", err
	}
	return filepath.Join(contextDir, filename), nil
}