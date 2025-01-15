package update

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Release struct {
	TagName string `json:"tag_name"`
}

const (
	githubOwner = "mrpaintedtank"
	githubRepo  = "warno-util"
)

type Updater struct {
	Version string
}

func (u Updater) RunUpdate() error {
	// Check latest version
	latest, err := u.getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	// Parse versions for comparison
	currentVer, err := semver.NewVersion(u.Version)
	if err != nil {
		return fmt.Errorf("invalid current version format: %w", err)
	}

	latestVer, err := semver.NewVersion(latest)
	if err != nil {
		return fmt.Errorf("invalid latest version format: %w", err)
	}

	// Compare versions
	if currentVer.GreaterThanEqual(latestVer) {
		fmt.Println("Already running the latest version!")
		return nil
	}

	fmt.Printf("Updating from %s to %s...\n", currentVer, latestVer)

	// Get current executable path
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Download new version
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s_%s_windows_%s.zip",
		githubOwner, githubRepo, latest, githubRepo, strings.TrimPrefix(latest, "v"), runtime.GOARCH)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	// Read zip file into memory
	zipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read zip file: %w", err)
	}

	// Open zip archive
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}

	// Find the executable in the zip
	var exeFile *zip.File
	for _, f := range zipReader.File {
		if strings.HasSuffix(f.Name, ".exe") {
			exeFile = f
			break
		}
	}

	if exeFile == nil {
		return fmt.Errorf("no executable found in update package")
	}

	// Create a temporary file for the new executable
	tmpDir := filepath.Dir(executable)
	tmpFile, err := os.CreateTemp(tmpDir, "*.exe")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Extract the new executable
	rc, err := exeFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open new executable: %w", err)
	}
	defer rc.Close()

	if _, err := io.Copy(tmpFile, rc); err != nil {
		return fmt.Errorf("failed to write new executable: %w", err)
	}
	tmpFile.Close()

	// Create batch file for replacing executable
	batchFile, err := os.CreateTemp(tmpDir, "*.bat")
	if err != nil {
		return fmt.Errorf("failed to create batch file: %w", err)
	}

	batchCommands := fmt.Sprintf(`@echo off
setlocal
set maxAttempts=10
set attempt=0

:retry
set /a attempt+=1
if %%attempt%% gtr %%maxAttempts%% (
    echo Failed to update after %%maxAttempts%% attempts. Please ensure you have the necessary permissions and try again.
    exit /b 1
)
timeout /t 1 /nobreak > nul
del "%s"
if exist "%s" goto retry
move "%s" "%s"
del "%s"
echo Update completed successfully!`,
		executable, executable, tmpFile.Name(), executable, batchFile.Name())

	if err := os.WriteFile(batchFile.Name(), []byte(batchCommands), 0755); err != nil {
		return fmt.Errorf("failed to write batch file: %w", err)
	}

	if err := exec.Command("cmd", "/C", batchFile.Name()).Start(); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}

func (u Updater) getLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubOwner, githubRepo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}
