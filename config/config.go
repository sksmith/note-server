package config

type Config struct {
	Port            string
	GenerateRoutes  bool
	LogLevel        string
	LogText         bool
	InMemoryDb      bool
	BucketName      string
	Revision        string
	ApplicationName string
	AppVersion      string
	Sha1Version     string
	BuildTime       string
	Profile         string
}

var (
	AppVersion  string
	Sha1Version string
	BuildTime   string
	Profile     string
)

const (
	ApplicationName = "note-server"
	Revision        = "1"
)

func LoadConfigs() (Config, error) {
	cfg := Config{
		ApplicationName: ApplicationName,
		Revision:        Revision,
		Profile:         Profile,
		Port:            "80",
		GenerateRoutes:  false,
		AppVersion:      AppVersion,
		Sha1Version:     Sha1Version,
		BuildTime:       BuildTime,
	}

	if cfg.Profile == "local" {
		if err := loadLocalConfigs(&cfg); err != nil {
			return Config{}, err
		}
	} else {
		if err := loadEnvironmentConfigs(&cfg); err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}

func loadLocalConfigs(cfg *Config) error {
	cfg.LogLevel = "trace"
	cfg.LogText = true
	cfg.BucketName = "sksmithnotes"
	cfg.InMemoryDb = false

	return nil
}

func loadEnvironmentConfigs(cfg *Config) error {
	cfg.LogLevel = "trace"
	cfg.LogText = false
	cfg.BucketName = "sksmithnotes"
	cfg.InMemoryDb = true

	return nil
}
