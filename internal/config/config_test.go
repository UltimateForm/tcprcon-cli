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
		t.Fatal(errors.Join(errors.New("unexpected error loading config"), err))
	}
	if cfg == nil {
		t.Fatal("unexpected nil config")
	}
	profile, exists := cfg.GetProfile("docker")
	if !exists {
		t.Fatal("unexpected profile (docker) absence in loaded config")
	}
	if profile != sourceProfile {
		t.Fatalf("unexpected mismatching loaded data, expected: %+v; received: %+v", sourceProfile, profile)
	}
}

func TestShouldSaveToConfig(t *testing.T) {
	baseConfigPath := t.TempDir()
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
	sourceConfig.Save(baseConfigPath)
	filePath := filepath.Join(baseConfigPath, configDirName, configFileName)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpect file stat error"), err))
	}
	if fileMode := fileInfo.Mode(); fileMode.Perm() != 0600 {
		t.Errorf("unxpected file mode: %v, expected 0600", fileMode)
	}
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected file read error"), err))
	}
	var cfg Config
	err = json.Unmarshal(fileBytes, &cfg)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected json unmarshal error"), err))
	}
	profile, exists := cfg.GetProfile("docker")
	if !exists {
		t.Fatal("unexpected profile (docker) absence in loaded config")
	}
	if profile != sourceProfile {
		t.Fatalf("unexpected mismatching loaded data, expected: %+v; received: %+v", sourceProfile, profile)
	}
}
