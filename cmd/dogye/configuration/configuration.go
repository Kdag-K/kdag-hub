package configuration

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// DogyeConfigDir is the absolute path of the dogye configuration directory
var DogyeConfigDir = defaultDogyeDir()

const (
	// DogyeNetworkDir is the networks subfolder of the Dogye config folder
	DogyeNetworkDir = "networks"
)

// defaultDogyeDir returns the full path for Dogye's data directory.
func defaultDogyeDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Dogye")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Dogye")
		} else {
			return filepath.Join(home, ".Dogye")
		}
	}
	return ""
}

// Guess a sensible default location from OS and environment variables.
func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}