package switcher

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
)

type fileOp struct {
	src, dst string
}

func performOperations(isCopy bool, operations []fileOp) error {
	if isCopy {
		for _, op := range operations {
			fmt.Println("Copying", op.src, "to", op.dst)
			if err := copyDir(op.src, op.dst); err != nil {
				return fmt.Errorf("failed to copy %s to %s: %w", op.src, op.dst, err)
			}
		}
		return nil
	}

	for _, op := range operations {
		fmt.Println("Renaming", op.src, "to", op.dst)
		if err := os.Rename(op.src, op.dst); err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", op.src, op.dst, err)
		}
	}
	return nil
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting source info: %w", err)
	}

	if !srcInfo.IsDir() {
		return copyFile(src, dst, srcInfo.Mode())
	}
	
	if runtime.GOOS == "windows" {
		// Use robocopy on Windows
		cmd := exec.Command("robocopy", src, dst, "/E", "/NFL", "/NDL")
		if err := cmd.Run(); err != nil {
			// robocopy returns non-zero exit codes even for successful copies
			// Exit codes 0-7 indicate successful copy with different warnings
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() <= 7 {
				return nil
			}
			return fmt.Errorf("robocopy failed: %w", err)
		}
		return nil
	}

	// Use cp on Unix-like systems
	cmd := exec.Command("cp", "-r", src, dst)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cp failed: %w", err)
	}
	return nil
}

func copyFile(src, dst string, mode fs.FileMode) error {
	// Copy file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error copying file contents: %w", err)
	}

	return nil
}
