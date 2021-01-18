package configuration

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

// Configuration represents the engine's configuration file.
type Configuration struct {
	MpqLoadOrder    []string
	MpqPath         string
	TicksPerSecond  int
	FpsCap          int
	SfxVolume       float64
	BgmVolume       float64
	FullScreen      bool
	RunInBackground bool
	VsyncEnabled    bool
	Backend         string
	filePath        string
}

// Save saves the configuration object to disk
func (c *Configuration) Save() error {
	configDir := path.Dir(c.filePath)
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return err
	}

	configFile, err := os.Create(c.filePath)
	if err != nil {
		return err
	}

	buf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if _, err := configFile.Write(buf); err != nil {
		return err
	}

	return configFile.Close()
}

// Dir returns the directory component of the path
func (c *Configuration) Dir() string {
	return filepath.Dir(c.filePath)
}

// Base returns the base component of the path
func (c *Configuration) Base() string {
	return filepath.Base(c.filePath)
}

// Path returns the configuration file path
func (c *Configuration) Path() string {
	return c.filePath
}

// SetPath sets where the configuration file is saved to (a full path)
func (c *Configuration) SetPath(p string) {
	c.filePath = p
}

func New() *Configuration {
	result, err := Load()
	if err != nil {
		panic(err)
	}

	return result
}
