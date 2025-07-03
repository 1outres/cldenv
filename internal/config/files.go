package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsSymlink checks if a file is a symbolic link
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// CreateDir creates a directory if it doesn't exist
func CreateDir(path string) error {
	if !FileExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// MoveFile moves a file from src to dst
func MoveFile(src, dst string) error {
	// First try to rename (fastest if on same filesystem)
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// If rename fails, copy and delete
	if err := CopyFile(src, dst); err != nil {
		return fmt.Errorf("failed to copy file during move: %w", err)
	}

	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to remove source file during move: %w", err)
	}

	return nil
}

// RemoveFile removes a file if it exists
func RemoveFile(path string) error {
	if FileExists(path) {
		return os.Remove(path)
	}
	return nil
}

// EnsureDir ensures that the directory exists for the given file path
func EnsureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return CreateDir(dir)
}