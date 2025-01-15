package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func FindRemotePath(steamUserDataPath string) (string, error) {
	entries, err := os.ReadDir(steamUserDataPath)
	if err != nil {
		return "", fmt.Errorf("error reading Steam user data directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			remotePath := filepath.Join(steamUserDataPath, entry.Name(), "1611600", "remote")
			if _, err := os.Stat(remotePath); err == nil {
				return remotePath, nil
			}
		}
	}

	return "", fmt.Errorf("remote path not found in %s", steamUserDataPath)
}

func GetUserDocsDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	// Construct the path to the Documents folder
	return filepath.Join(usr.HomeDir, "Documents"), nil
}

func GetBinaryDir() (string, error) {
	// Get the absolute path of the executable
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)

	// Handle symbolic links
	realPath, err := filepath.EvalSymlinks(execDir)
	if err != nil {
		return "", fmt.Errorf("failed to eval symlinks: %w", err)
	}

	return realPath, nil
}
