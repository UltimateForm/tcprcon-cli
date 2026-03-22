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
		t.Fatal(errors.Join(errors.New("unexpected file stat error"), err))
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

func TestShouldReturnErrorWhenBasePathIsEmpty(t *testing.T) {
	_, err := Load("")
	if err == nil {
		t.Fatal("expected error for empty base path, got nil")
	}
	if !errors.Is(err, ErrUndefinedConfigBasePath) {
		t.Fatalf("expected ErrUndefinedConfigBasePath, got %v", err)
	}
}

func TestSaveCreatesConfigDirIfMissing(t *testing.T) {
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

	err := sourceConfig.Save(baseConfigPath)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected save error"), err))
	}

	configDirPath := filepath.Join(baseConfigPath, configDirName)
	dirInfo, err := os.Stat(configDirPath)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected dir stat error"), err))
	}
	if !dirInfo.IsDir() {
		t.Fatal("expected config dir to be a directory")
	}
	if dirInfo.Mode().Perm() != 0700 {
		t.Errorf("unexpected dir mode: %v, expected 0700", dirInfo.Mode())
	}

	filePath := filepath.Join(configDirPath, configFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("expected config file to exist after Save")
	}
}

func TestSaveOverwritesExistingProfile(t *testing.T) {
	baseConfigPath := t.TempDir()
	firstProfile := Profile{
		Address:  "localhost",
		Port:     7778,
		Password: "firstpassword",
	}
	cfg := Config{
		Profiles: map[string]Profile{
			"myserver": firstProfile,
		},
	}
	cfg.Save(baseConfigPath)

	updatedProfile := Profile{
		Address:  "192.168.1.50",
		Port:     9000,
		Password: "updatedpassword",
	}
	cfg.SetProfile("myserver", updatedProfile)
	err := cfg.Save(baseConfigPath)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected save error"), err))
	}

	loaded, err := Load(baseConfigPath)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected load error"), err))
	}
	profile, exists := loaded.GetProfile("myserver")
	if !exists {
		t.Fatal("expected profile 'myserver' to exist after overwrite")
	}
	if profile != updatedProfile {
		t.Fatalf("expected updated profile %+v, got %+v", updatedProfile, profile)
	}
}

func TestResolveShouldLoadFromConfig(t *testing.T) {
	baseConfigPath := t.TempDir()
	sourceProfile := Profile{
		Address:  "169.230.184.1",
		Port:     7482,
		Password: "mycloset",
	}
	sourceConfig := Config{
		Profiles: map[string]Profile{
			"pleyades": sourceProfile,
		},
	}
	sourceConfig.Save(baseConfigPath)
	resolvedProfile, err := Resolve(
		baseConfigPath,
		"pleyades",
		Profile{
			Address:  DefaultAddr,
			Port:     DefaultPort,
			Password: "",
		},
	)

	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected load error"), err))
	}

	if resolvedProfile.Address != sourceProfile.Address {
		t.Errorf("expected address %s, but got %s", sourceProfile.Address, resolvedProfile.Address)
	}
	if resolvedProfile.Port != sourceProfile.Port {
		t.Errorf("expected port %d, but got %d", sourceProfile.Port, resolvedProfile.Port)
	}
	if resolvedProfile.Password != sourceProfile.Password {
		t.Errorf("expected password %s, but got %s", sourceProfile.Password, resolvedProfile.Password)
	}
}

func TestResolveShouldNotOverrideExplicitFlags(t *testing.T) {
	baseConfigPath := t.TempDir()
	sourceProfile := Profile{
		Address:  "169.230.184.1",
		Port:     7482,
		Password: "profilepassword",
	}
	sourceConfig := Config{
		Profiles: map[string]Profile{
			"pleyades": sourceProfile,
		},
	}
	sourceConfig.Save(baseConfigPath)

	explicitAddr := "10.0.0.1"
	explicitPort := uint(9999)
	explicitPw := "explicitpassword"

	resolvedProfile, err := Resolve(
		baseConfigPath,
		"pleyades",
		Profile{
			Address:  explicitAddr,
			Port:     explicitPort,
			Password: explicitPw,
		},
	)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected resolve error"), err))
	}
	if resolvedProfile.Address != explicitAddr {
		t.Errorf("expected explicit address %s, but got %s", explicitAddr, resolvedProfile.Address)
	}
	if resolvedProfile.Port != explicitPort {
		t.Errorf("expected explicit port %d, but got %d", explicitPort, resolvedProfile.Port)
	}
	if resolvedProfile.Password != explicitPw {
		t.Errorf("expected explicit password %s, but got %s", explicitPw, resolvedProfile.Password)
	}
}

func TestResolveShouldReturnErrorForMissingProfile(t *testing.T) {
	baseConfigPath := t.TempDir()
	emptyConfig := Config{
		Profiles: map[string]Profile{},
	}
	emptyConfig.Save(baseConfigPath)

	_, err := Resolve(baseConfigPath, "nonexistent", Profile{
		Address:  DefaultAddr,
		Port:     DefaultPort,
		Password: "",
	},
	)
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}

func TestResolveShouldReturnDefaultsWhenNoProfile(t *testing.T) {
	baseConfigPath := t.TempDir()

	resolveProfile, err := Resolve(
		baseConfigPath,
		"",
		Profile{
			Address:  DefaultAddr,
			Port:     DefaultPort,
			Password: "",
		},
	)
	if err != nil {
		t.Fatal(errors.Join(errors.New("unexpected resolve error"), err))
	}
	if resolveProfile.Address != DefaultAddr {
		t.Errorf("expected default address %s, but got %s", DefaultAddr, resolveProfile.Address)
	}
	if resolveProfile.Port != DefaultPort {
		t.Errorf("expected default port %d, but got %d", DefaultPort, resolveProfile.Port)
	}
	if resolveProfile.Password != "" {
		t.Errorf("expected empty password, but got %s", resolveProfile.Password)
	}
}

func TestBuildConfigPathReturnsErrorOnEmptyBase(t *testing.T) {
	_, err := BuildConfigPath("")
	if err == nil {
		t.Fatal("expected error for empty base path, got nil")
	}
	if !errors.Is(err, ErrUndefinedConfigBasePath) {
		t.Fatalf("expected ErrUndefinedConfigBasePath, got %v", err)
	}
}
