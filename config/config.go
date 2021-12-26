package config

import "flag"

type Config struct {
	Port            string
	LogLevel        string
	LogText         bool
	Region          string
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
	profile := flag.String("P", "local", "profile for the application config")
	port := flag.String("p", "8080", "port for the application to listen to")
	region := flag.String("r", "us-east-1", "region the bucket resides in")
	bucket := flag.String("b", "sksmithnotes", "bucket name for the application to use")
	flag.Parse()

	cfg := Config{
		AppVersion:      AppVersion,
		ApplicationName: ApplicationName,
		BucketName:      *bucket,
		BuildTime:       BuildTime,
		Profile:         *profile,
		Port:            *port,
		Region:          *region,
		Revision:        Revision,
		Sha1Version:     Sha1Version,
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

	return nil
}

func loadEnvironmentConfigs(cfg *Config) error {
	cfg.LogLevel = "trace"
	cfg.LogText = false

	return nil
}
