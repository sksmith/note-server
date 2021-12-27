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
	// Runtime flags
	bucket  *string
	port    *string
	profile *string
	region  *string

	// Build time arguments
	AppVersion  string
	Sha1Version string
	BuildTime   string
)

const (
	// Build time arguments
	ApplicationName = "note-server"
	Revision        = "1"

	// Default runtime arguments
	DefaultBucket  = "sksmithnotes"
	DefaultPort    = "8080"
	DefaultProfile = "local"
	DefaultRegion  = "us-east-1"

	// Default runtime arguments when running locally
	DefaultLocalLogLevel = "trace"
	DefaultLocalLogText  = true

	// Default runtime arguments when running remotely
	DefaultEnvironmentLogLevel = "trace"
	DefaultEnvironmentLogText  = false
)

func LoadConfigs() (Config, error) {
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
	cfg.LogLevel = DefaultLocalLogLevel
	cfg.LogText = DefaultLocalLogText

	return nil
}

func loadEnvironmentConfigs(cfg *Config) error {
	cfg.LogLevel = DefaultEnvironmentLogLevel
	cfg.LogText = DefaultEnvironmentLogText

	return nil
}

func init() {
	// Set runtime arguments
	profile = flag.String("P", DefaultProfile, "profile for the application config")
	port = flag.String("p", DefaultPort, "port for the application to listen to")
	region = flag.String("r", DefaultRegion, "region the bucket resides in")
	bucket = flag.String("b", DefaultBucket, "bucket name for the application to use")
}
