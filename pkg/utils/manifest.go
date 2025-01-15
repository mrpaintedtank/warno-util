package utils

import (
	"bytes"
	"fmt"
	"os"
)

// ReadManifest reads the manifest file and returns true if VIP and false if LIVE. False is also valid if the file is not found so check error first.
func ReadManifest(profile string) (bool, error) {
	file, err := os.ReadFile(profile)
	if err != nil {
		return false, err
	}
	if bytes.Contains(file, []byte("BetaKey")) && bytes.Contains(file, []byte("WARNO")) && bytes.Contains(file, []byte("full_vip")) {
		return true, nil
	}
	if bytes.Contains(file, []byte("BetaKey")) && bytes.Contains(file, []byte("WARNO")) && bytes.Contains(file, []byte("public")) {
		return false, nil
	}
	return false, fmt.Errorf("manifest file does not contain expected values")
}
