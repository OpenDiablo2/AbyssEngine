package configuration

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	od2ConfigDirName  = "OpenDiablo2"
	od2ConfigFileName = "config.json"
)

// Load loads the configuration file
func Load() (*Configuration, error) {
	configFileFullPath := ""

	if configDir, err := os.UserConfigDir(); err == nil {
		configFileFullPath = path.Join(configDir, od2ConfigDirName, od2ConfigFileName)
	} else {
		configFileFullPath = path.Join(path.Dir(os.Args[0]), od2ConfigFileName)
	}

	result := &Configuration{}

	configAsset, err := os.Open(configFileFullPath)
	if err != nil {
		return nil, err
	}

	// create the default if not found
	if configAsset != nil {
		if err := json.NewDecoder(configAsset).Decode(result); err != nil {
			return nil, err
		}

		if err := configAsset.Close(); err != nil {
			return nil, err
		}

		return result, nil
	}

	result = DefaultConfig()

	fullPath := filepath.Join(result.Dir(), result.Base())
	result.SetPath(fullPath)

	log.Warnf("creating default configuration file at %s...", fullPath)

	saveErr := result.Save()

	return result, saveErr
}
