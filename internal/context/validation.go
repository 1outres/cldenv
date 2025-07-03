package context

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Reserved context names that cannot be used
	reservedNames = map[string]bool{
		"help":    true,
		"version": true,
		"use":     true,
		"create":  true,
		"remove":  true,
		"list":    true,
		"switch":  true,
	}

	// Valid context name pattern: alphanumeric, dash, underscore
	validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// ValidateContextName validates a context name
func ValidateContextName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: context name cannot be empty", ErrInvalidContextName)
	}

	if len(name) > 50 {
		return fmt.Errorf("%w: context name too long (max 50 characters)", ErrInvalidContextName)
	}

	if !validNamePattern.MatchString(name) {
		return fmt.Errorf("%w: context name can only contain letters, numbers, dashes, and underscores", ErrInvalidContextName)
	}

	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("%w: context name cannot start or end with dash", ErrInvalidContextName)
	}

	if strings.HasPrefix(name, "_") || strings.HasSuffix(name, "_") {
		return fmt.Errorf("%w: context name cannot start or end with underscore", ErrInvalidContextName)
	}

	if reservedNames[strings.ToLower(name)] {
		return fmt.Errorf("%w: '%s' is a reserved name", ErrInvalidContextName, name)
	}

	return nil
}