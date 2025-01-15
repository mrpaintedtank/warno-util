package switcher

// Switcher struct to hold configuration variables
type Config struct {
	SteamAppsPath       string `yaml:"steamAppsPath"`
	SteamUserDataPath   string `yaml:"steamUserDataPath"`
	SteamExecutablePath string `yaml:"steamExecutablePath"`
	remotePath          string
}
