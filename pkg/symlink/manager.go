package symlink

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrSymlinkFailed = errors.New("symlink operation failed")
)

// CreateSymlink creates a symbolic link from src to dst
func CreateSymlink(src, dst string) error {
	// Remove existing file/symlink if it exists
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%w: failed to remove existing file: %v", ErrSymlinkFailed, err)
	}

	// Create the symlink
	if err := os.Symlink(src, dst); err != nil {
		return fmt.Errorf("%w: failed to create symlink: %v", ErrSymlinkFailed, err)
	}

	return nil
}

// RemoveSymlink removes a symbolic link if it exists
func RemoveSymlink(path string) error {
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, nothing to remove
		}
		return fmt.Errorf("%w: failed to stat file: %v", ErrSymlinkFailed, err)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("%w: failed to remove symlink: %v", ErrSymlinkFailed, err)
	}

	return nil
}

// ReadSymlink reads the target of a symbolic link
func ReadSymlink(path string) (string, error) {
	target, err := os.Readlink(path)
	if err != nil {
		return "", fmt.Errorf("%w: failed to read symlink: %v", ErrSymlinkFailed, err)
	}

	// Convert to absolute path if it's relative
	if !filepath.IsAbs(target) {
		dir := filepath.Dir(path)
		target = filepath.Join(dir, target)
	}

	return target, nil
}

// IsValidSymlink checks if a symlink exists and points to a valid target
func IsValidSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	// Check if it's a symlink
	if info.Mode()&os.ModeSymlink == 0 {
		return false
	}

	// Check if target exists
	target, err := ReadSymlink(path)
	if err != nil {
		return false
	}

	_, err = os.Stat(target)
	return err == nil
}

// BackupFile creates a backup of a file before replacing it with a symlink
func BackupFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to backup
	}

	backupPath := path + ".backup"
	if err := os.Rename(path, backupPath); err != nil {
		return fmt.Errorf("%w: failed to backup file: %v", ErrSymlinkFailed, err)
	}

	return nil
}

// RestoreBackup restores a backup file
func RestoreBackup(path string) error {
	backupPath := path + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return nil // Backup doesn't exist, nothing to restore
	}

	// Remove the current file if it exists
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%w: failed to remove current file: %v", ErrSymlinkFailed, err)
	}

	// Restore the backup
	if err := os.Rename(backupPath, path); err != nil {
		return fmt.Errorf("%w: failed to restore backup: %v", ErrSymlinkFailed, err)
	}

	return nil
}