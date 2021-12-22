package config

import "os"

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
)

const (
	ApplicationName = "note-server"
	Revision        = "1"
)

func LoadConfigs() (Config, error) {
	cfg := Config{
		ApplicationName: ApplicationName,
		Revision:        Revision,
		Profile:         os.Getenv("PROFILE"),
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
	// Log Configs
	cfg.LogLevel = "trace"
	cfg.LogText = true

	// DB Configs
	cfg.BucketName = "sksmithnotes"
	cfg.InMemoryDb = false

	return nil
}

func loadEnvironmentConfigs(cfg *Config) error {
	// Log Configs
	cfg.LogLevel = "trace"
	cfg.LogText = true

	// DB Configs
	cfg.BucketName = "sksmithnotes"
	cfg.InMemoryDb = true

	return nil
}
