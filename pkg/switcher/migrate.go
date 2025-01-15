package switcher

import (
	"fmt"
	"github.com/test/warno-utils/pkg/utils"
	"path/filepath"
)

// migrates from the old staged directory to live and vip directories
func performMigration(cfg Config, isVip bool) error {
	if !utils.PathExists(cfg.GetStagedDir()) || !utils.PathExists(cfg.GetStagedManifest()) || !utils.PathExists(cfg.GetStagedProfile()) {
		return nil // No migration needed
	}

	manifestIsVip, err := utils.ReadManifest(cfg.GetStagedManifest())
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	if manifestIsVip != isVip {
		return nil // No migration needed for this type
	}
	var operations []fileOp
	if isVip {
		operations = []fileOp{
			{stagedDir, filepath.Join(cfg.SteamAppsPath, "common", vipDir)},
			{stagedManifest, filepath.Join(cfg.SteamAppsPath, vipManifest)},
			{stagedProfile, filepath.Join(cfg.remotePath, vipProfile)},
		}
	} else {
		operations = []fileOp{
			{stagedDir, filepath.Join(cfg.SteamAppsPath, "common", liveDir)},
			{stagedManifest, filepath.Join(cfg.SteamAppsPath, liveManifest)},
			{stagedProfile, filepath.Join(cfg.remotePath, liveProfile)},
		}
	}

	return performOperations(false, operations)
}
