package switcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/test/warno-utils/pkg/utils"
)

func Switcher(cfg Config) error {
	if utils.ProcessRunning("WARNO") {
		if err := utils.StopProcess("WARNO"); err != nil {
			return fmt.Errorf("failed to stop 'WARNO' process: %w", err)
		}
	}

	if utils.ProcessRunning("steam") {
		if err := utils.StopProcess("steam"); err != nil {
			return fmt.Errorf("failed to stop 'steam' process: %w", err)
		}
	}

	if utils.ProcessRunning("WARNO") {
		return fmt.Errorf("the 'WARNO' process is still running. Please close the game manually and try again")
	}

	if utils.ProcessRunning("steam") {
		return fmt.Errorf("the 'steam' process is still running. Please close Steam manually and try again")
	}

	isVip, err := utils.ReadManifest(filepath.Join(cfg.SteamAppsPath, curManifest))
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	remotePath, err := utils.FindRemotePath(cfg.SteamUserDataPath)
	if err != nil {
		return err
	}

	updatedCfg := Config{
		SteamAppsPath:       cfg.SteamAppsPath,
		SteamUserDataPath:   cfg.SteamUserDataPath,
		SteamExecutablePath: cfg.SteamExecutablePath,
		remotePath:          remotePath,
	}

	if err := performMigration(updatedCfg, isVip); err != nil {
		return err
	}

	if err := checkFileExistence(updatedCfg, isVip); err != nil {
		return err
	}

	if err := switchProfiles(isVip, updatedCfg); err != nil {
		return fmt.Errorf("failed to switch profiles: %w", err)
	}

	if err := exec.Command(updatedCfg.SteamExecutablePath).Start(); err != nil {
		return fmt.Errorf("An error occurred while starting the process: %v\n", err)
	}
	return nil
}

func switchProfiles(isVip bool, cfg Config) error {
	var operations []fileOp
	if isVip {
		operations = []fileOp{
			{filepath.Join(cfg.SteamAppsPath, "common", curDir), filepath.Join(cfg.SteamAppsPath, "common", vipDir)},
			{filepath.Join(cfg.SteamAppsPath, curManifest), filepath.Join(cfg.SteamAppsPath, vipManifest)},
			{filepath.Join(cfg.remotePath, curProfile), filepath.Join(cfg.remotePath, vipProfile)},
			{filepath.Join(cfg.SteamAppsPath, "common", liveDir), filepath.Join(cfg.SteamAppsPath, "common", curDir)},
			{filepath.Join(cfg.SteamAppsPath, liveManifest), filepath.Join(cfg.SteamAppsPath, curManifest)},
			{filepath.Join(cfg.remotePath, liveProfile), filepath.Join(cfg.remotePath, curProfile)},
		}
	} else {
		operations = []fileOp{
			{filepath.Join(cfg.SteamAppsPath, "common", curDir), filepath.Join(cfg.SteamAppsPath, "common", liveDir)},
			{filepath.Join(cfg.SteamAppsPath, curManifest), filepath.Join(cfg.SteamAppsPath, liveManifest)},
			{filepath.Join(cfg.remotePath, curProfile), filepath.Join(cfg.remotePath, liveProfile)},
			{filepath.Join(cfg.SteamAppsPath, "common", vipDir), filepath.Join(cfg.SteamAppsPath, "common", curDir)},
			{filepath.Join(cfg.SteamAppsPath, vipManifest), filepath.Join(cfg.SteamAppsPath, curManifest)},
			{filepath.Join(cfg.remotePath, vipProfile), filepath.Join(cfg.remotePath, curProfile)},
		}
	}

	return performOperations(operations)
}

func checkFileExistence(cfg Config, isVip bool) error {
	// Determine which files should exist based on isVip
	var dirsToCheck, manifestsToCheck, profilesToCheck []string

	// Common files that should always exist
	dirsToCheck = append(dirsToCheck, filepath.Join("common", curDir))
	manifestsToCheck = append(manifestsToCheck, curManifest)
	profilesToCheck = append(profilesToCheck, curProfile)

	if isVip {
		dirsToCheck = append(dirsToCheck, filepath.Join("common", liveDir))
		manifestsToCheck = append(manifestsToCheck, liveManifest)
		profilesToCheck = append(profilesToCheck, liveProfile)
	} else {
		dirsToCheck = append(dirsToCheck, filepath.Join("common", vipDir))
		manifestsToCheck = append(manifestsToCheck, vipManifest)
		profilesToCheck = append(profilesToCheck, vipProfile)
	}

	// Check directories in SteamAppsPath
	for _, dir := range dirsToCheck {
		path := filepath.Join(cfg.SteamAppsPath, dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
	}

	// Check manifests in SteamAppsPath
	for _, manifest := range manifestsToCheck {
		path := filepath.Join(cfg.SteamAppsPath, manifest)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("manifest does not exist: %s", path)
		}
	}

	// Check profiles in remotePath
	for _, profile := range profilesToCheck {
		path := filepath.Join(cfg.remotePath, profile)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("profile does not exist: %s", path)
		}
	}

	return nil
}
