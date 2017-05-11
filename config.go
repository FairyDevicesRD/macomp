package macomp

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

//GetConfigPath returns config file path
func GetConfigPath() string {
	if path := os.Getenv("MACOMP"); len(path) != 0 {
		return path
	} else if path, err := homedir.Dir(); err == nil && len(path) != 0 {
		return filepath.Join(path, ".macomp.json")
	}
	return ""
}

//GetStaticRootPath returns the static root path
func GetStaticRootPath() string {
	if path := os.Getenv("MACOMP_STATIC"); len(path) != 0 {
		return path
	}
	if path := os.Getenv("GOPATH"); len(path) != 0 {
		spath := filepath.Join(path, "src", "github.com", "FairyDevicesRD", "macomp", "static")
		if _, err := os.Stat(spath); err == nil {
			return spath
		}
	}
	return ""
}

//PrintDefaultPath prints default paths
func PrintDefaultPath() {
	fmt.Printf("Default Path:\n")
	fmt.Printf("  MA Config path  : %s\n", GetConfigPath())
	fmt.Printf("  Static root path: %s\n", GetStaticRootPath())
}

//MaSetting is a setting structure
type MaSetting struct {
	MaType  string   `json:"type"`
	Path    string   `json:"path"`
	Disable bool     `json:"disable"`
	Aliases []string `json:"aliases"`
	Options map[string]string
}
