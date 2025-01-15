package switcher

import (
	"fmt"
	"github.com/test/warno-utils/pkg/utils"
	"path/filepath"
)

func performMigration(cfg Config, isVip bool) error {
	stagedDir := filepath.Join(cfg.SteamAppsPath, "common", stagedDir)
	stagedManifest := filepath.Join(cfg.SteamAppsPath, stagedManifest)
	stagedProfile := filepath.Join(cfg.remotePath, stagedProfile)

	if !utils.PathExists(stagedDir) || !utils.PathExists(stagedManifest) || !utils.PathExists(stagedProfile) {
		return nil // No migration needed
	}

	manifestIsVip, err := utils.ReadManifest(stagedManifest)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	if manifestIsVip != isVip {
		return nil // No migration needed for this type
	}

	targetDir := vipDir
	targetManifest := vipManifest
	targetProfile := vipProfile
	if !isVip {
		targetDir = liveDir
		targetManifest = liveManifest
		targetProfile = liveProfile
	}

	operations := []fileOp{
		{stagedDir, filepath.Join(cfg.SteamAppsPath, "common", targetDir)},
		{stagedManifest, filepath.Join(cfg.SteamAppsPath, targetManifest)},
		{stagedProfile, filepath.Join(cfg.remotePath, targetProfile)},
	}

	return performOperations(operations)
}
