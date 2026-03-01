package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestShouldLoadWhenFileExists(t *testing.T) {
	sourceProfile := Profile{
		Address:  "localhost",
		Port:     7778,
		Password: "localpassword",
	}
	sourceConfig := Config{
		Profiles: map[string]Profile{
			"docker": sourceProfile,
		},
	}
	baseConfigPath := t.TempDir()
	configFolder := filepath.Join(baseConfigPath, configDirName)
	_ = os.Mkdir(configFolder, 0700)
	jsonStr, _ := json.MarshalIndent(sourceConfig, "", "  ")
	os.WriteFile(
		filepath.Join(configFolder, configFileName),
		[]byte(jsonStr),
		0600,
	)

	cfg, err := Load(baseConfigPath)
	if err != nil {
		t.Error(errors.Join(errors.New("unexpected error loading config"), err))
	}
	if cfg == nil {
		t.Error("unexpected nil config")
	}
	if len(cfg.Profiles) == 0 {
		t.Error("unexpected zero-length profiles")
	}
	profile, exists := cfg.GetProfile("docker")
	if !exists {
		t.Error("unexpected profile (docker) absence in loaded config")
	}
	if profile != sourceProfile {
		t.Errorf("unexpected mismatching loaded data, expected: %+v; received: %+v", sourceProfile, profile)
	}
}
