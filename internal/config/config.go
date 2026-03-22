package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func BuildConfigPath(basePath string) (string, error) {
	if basePath == "" {
		return "", ErrUndefinedConfigBasePath
	}

	fullPath := filepath.Join(basePath, configDirName)
	return filepath.Join(fullPath, configFileName), nil
}

func Load(baseConfigPath string) (*Config, error) {

	path, err := BuildConfigPath(baseConfigPath)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{Profiles: make(map[string]Profile)}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Profiles == nil {
		// should never happen but jic
		cfg.Profiles = make(map[string]Profile)
	}

	return &cfg, nil
}

func (source *Config) Save(configBasePath string) error {
	path, err := BuildConfigPath(configBasePath)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	// 0700 sets permissions so that, (U)ser / owner can read, can write and can execute. (G)roup can't read, can't write and can't execute. (O)thers can't read, can't write and can't execute.
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.Join(errors.New("failed to create config directory"), err)
	}

	data, err := json.MarshalIndent(source, "", "  ")
	if err != nil {
		return err
	}

	// 0600 sets permissions so that, (U)ser / owner can read, can write and can't execute. (G)roup can't read, can't write and can't execute. (O)thers can't read, can't write and can't execute.
	return os.WriteFile(path, data, 0600)
}

func (source *Config) GetProfile(name string) (Profile, bool) {
	p, ok := source.Profiles[name]
	return p, ok
}

func (source *Config) SetProfile(name string, p Profile) {
	if source.Profiles == nil {
		// again, just being defensive
		source.Profiles = make(map[string]Profile)
	}
	source.Profiles[name] = p
}

func Resolve(configBasePath string, profileName string, prof Profile) (Profile, error) {
	cfg, err := Load(configBasePath)
	if err != nil {
		return Profile{}, err
	}

	if profileName != "" {
		p, ok := cfg.GetProfile(profileName)
		if !ok {
			return Profile{}, fmt.Errorf("profile '%s' not found", profileName)
		}

		// only override if the flags are still at their default values
		// NOTE: this logic assumes defaults are "localhost", 7778, and ""
		// TODO: this "default" handling can introduce bugs, rethink this at some point
		if prof.Address == DefaultAddr && p.Address != "" {
			prof.Address = p.Address
		}
		if prof.Port == DefaultPort && p.Port != 0 {
			prof.Port = p.Port
		}
		if prof.Password == "" && p.Password != "" {
			prof.Password = p.Password
		}
		if prof.Pulse == "" && p.Pulse != "" {
			prof.Pulse = p.Pulse
		}
		if prof.PulseInterval == DefaultPulseInterval && p.PulseInterval != 0 {
			prof.PulseInterval = p.PulseInterval
		}
	}

	return prof, nil
}
