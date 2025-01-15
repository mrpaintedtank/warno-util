package switcher

import (
	"fmt"
	"github.com/test/warno-utils/pkg/utils"
	"path/filepath"
	"strings"
)

const liveDir = "WARNO_LIVE"
const vipDir = "WARNO_VIP"
const curDir = "WARNO"
const curManifest = "appmanifest_1611600.acf"
const liveManifest = "live_appmanifest_1611600.acf"
const vipManifest = "vip_appmanifest_1611600.acf"
const curProfile = "PROFILE.profile2"
const liveProfile = "live_PROFILE.profile2"
const vipProfile = "vip_PROFILE.profile2"
const stagedDir = "WARNOSTAGED"
const stagedManifest = "stagedappmanifest_1611600.acf"
const stagedProfile = "stagedPROFILE.profile2"

// Switcher struct to hold configuration variables
type Config struct {
	SteamAppsPath       string `yaml:"steamAppsPath"`
	SteamUserDataPath   string `yaml:"steamUserDataPath"`
	SteamExecutablePath string `yaml:"steamExecutablePath"`
	remotePath          string
}

func (c Config) GetLiveDir() string {
	return filepath.Join(c.SteamAppsPath, "common", liveDir)
}

func (c Config) GetVipDir() string {
	return filepath.Join(c.SteamAppsPath, "common", vipDir)
}

func (c Config) GetCurDir() string {
	return filepath.Join(c.SteamAppsPath, "common", curDir)
}

func (c Config) GetLiveManifest() string {
	return filepath.Join(c.SteamAppsPath, liveManifest)
}

func (c Config) GetVipManifest() string {
	return filepath.Join(c.SteamAppsPath, vipManifest)
}

func (c Config) GetCurManifest() string {
	return filepath.Join(c.SteamAppsPath, curManifest)
}

func (c Config) GetLiveProfile() string {
	return filepath.Join(c.remotePath, liveProfile)
}

func (c Config) GetVipProfile() string {
	return filepath.Join(c.remotePath, vipProfile)
}

func (c Config) GetCurProfile() string {
	return filepath.Join(c.remotePath, curProfile)
}

func (c Config) GetStagedDir() string {
	return filepath.Join(c.SteamAppsPath, "common", stagedDir)
}

func (c Config) GetStagedManifest() string {
	return filepath.Join(c.SteamAppsPath, stagedManifest)
}

func (c Config) GetStagedProfile() string {
	return filepath.Join(c.remotePath, stagedProfile)
}

type warnoVersion int

const (
	current warnoVersion = iota
	live
	vip
)

// CheckVersionExists checks if the specified version exists and returns true if it does
func (c Config) CheckVersionExists(version warnoVersion) bool {
	switch version {
	case live:
		return checkExistence(c.GetLiveDir(), c.GetLiveManifest(), c.GetLiveProfile())
	case vip:
		return checkExistence(c.GetVipDir(), c.GetVipManifest(), c.GetVipProfile())
	case current:
		return checkExistence(c.GetCurDir(), c.GetCurManifest(), c.GetCurProfile())
	default:
		return false
	}
}

// NeitherVipNorLive checks if neither VIP nor LIVE versions exist, returns true if they don't
func (c Config) NeitherVipNorLive() bool {
	return !c.CheckVersionExists(live) && !c.CheckVersionExists(vip)
}

func checkExistence(path ...string) bool {
	for _, p := range path {
		if !utils.PathExists(p) {
			return false
		}
	}
	return true
}

func (c Config) String() string {
	var result strings.Builder

	// Define consistent column widths
	const (
		typeCol = 12 // LIVE, VIP, CURRENT
		nameCol = 15 // Directory, Manifest, Profile
		pathCol = 50 // Path value
	)

	result.WriteString("Configuration Values:\n")
	result.WriteString(strings.Repeat("-", typeCol+nameCol+pathCol+2) + "\n")
	result.WriteString(fmt.Sprintf("%-*s %-*s %-*s\n", typeCol, "TYPE", nameCol, "NAME", pathCol, "PATH"))
	result.WriteString(strings.Repeat("-", typeCol+nameCol+pathCol+2) + "\n")

	// Configuration values
	formats := []struct {
		typeName string
		name     string
		path     string
	}{
		{"LIVE", "Directory", c.GetLiveDir()},
		{"VIP", "Directory", c.GetVipDir()},
		{"CURRENT", "Directory", c.GetCurDir()},
		{"LIVE", "Manifest", c.GetLiveManifest()},
		{"VIP", "Manifest", c.GetVipManifest()},
		{"CURRENT", "Manifest", c.GetCurManifest()},
		{"LIVE", "Profile", c.GetLiveProfile()},
		{"VIP", "Profile", c.GetVipProfile()},
		{"CURRENT", "Profile", c.GetCurProfile()},
	}

	for _, f := range formats {
		result.WriteString(fmt.Sprintf("%-*s %-*s %-*s\n",
			typeCol, f.typeName,
			nameCol, f.name,
			pathCol, f.path))
	}

	// Base paths
	result.WriteString("\nBase Paths:\n")
	result.WriteString(strings.Repeat("-", typeCol+nameCol+pathCol+2) + "\n")

	basePaths := []struct {
		name string
		path string
	}{
		{"Steam Apps", c.SteamAppsPath},
		{"Steam User Data", c.SteamUserDataPath},
		{"Steam Executable", c.SteamExecutablePath},
		{"Remote", c.remotePath},
	}

	for _, bp := range basePaths {
		result.WriteString(fmt.Sprintf("%-*s %-*s\n",
			typeCol+nameCol, bp.name,
			pathCol, bp.path))
	}

	return result.String()
}
