package switcher

import (
	"fmt"
	"github.com/test/warno-utils/pkg/utils"
	"os/exec"
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

	isVip, err := utils.ReadManifest(cfg.GetCurManifest())
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

	if !updatedCfg.CheckVersionExists(current) {
		return fmt.Errorf("warno not installed\n\n %s", cfg)
	}

	if err := switchVersion(isVip, updatedCfg); err != nil {
		return fmt.Errorf("failed to switch profiles: %w", err)
	}

	if err := exec.Command(updatedCfg.SteamExecutablePath).Start(); err != nil {
		return fmt.Errorf("An error occurred while starting the process: %v\n", err)
	}
	return nil
}

func switchVersion(isVip bool, cfg Config) error {
	var operations []fileOp
	// we haven't copied either version out and need to start with that
	isCopy := cfg.NeitherVipNorLive()
	if isCopy {
		if isVip {
			operations = []fileOp{
				{cfg.GetCurManifest(), cfg.GetVipManifest()},
				{cfg.GetCurProfile(), cfg.GetVipManifest()},
				{cfg.GetCurDir(), cfg.GetVipDir()},
			}
		} else {
			operations = []fileOp{
				{cfg.GetCurManifest(), cfg.GetLiveManifest()},
				{cfg.GetCurProfile(), cfg.GetLiveProfile()},
				{cfg.GetCurDir(), cfg.GetLiveDir()},
			}
		}
		return performOperations(true, operations)
	}

	if isVip {
		operations = []fileOp{
			{cfg.GetCurManifest(), cfg.GetVipManifest()},
			{cfg.GetCurProfile(), cfg.GetVipProfile()},
			{cfg.GetCurDir(), cfg.GetVipDir()},
			{cfg.GetLiveManifest(), cfg.GetCurManifest()},
			{cfg.GetLiveProfile(), cfg.GetCurProfile()},
			{cfg.GetLiveDir(), cfg.GetCurDir()},
		}
	} else {
		operations = []fileOp{
			{cfg.GetCurManifest(), cfg.GetLiveManifest()},
			{cfg.GetCurProfile(), cfg.GetLiveProfile()},
			{cfg.GetCurDir(), cfg.GetLiveDir()},
			{cfg.GetVipManifest(), cfg.GetCurManifest()},
			{cfg.GetVipProfile(), cfg.GetCurProfile()},
			{cfg.GetVipDir(), cfg.GetCurDir()},
		}
	}
	return performOperations(false, operations)
}
