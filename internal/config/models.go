package config

import "time"

type Profile struct {
	Address       string        `json:"address"`
	Port          uint          `json:"port"`
	Password      string        `json:"password,omitempty"` // omitempty so we don't save empty strings
	Pulse         string        `json:"pulse,omitempty"`
	PulseInterval time.Duration `json:"pulseInterval"` // i do like using time.Duration, but beware that this means it will be serialized in nanoseconds in the json file
}

type Config struct {
	Profiles map[string]Profile `json:"profiles"`
	// profiles could be root but preparing for potential expansion, jic
}

const (
	configDirName        = "tcprcon"
	configFileName       = "config.json"
	DefaultAddr          = "localhost"
	DefaultPort          = 7778
	DefaultPulseInterval = time.Second * 60
)
