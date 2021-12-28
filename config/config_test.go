package config_test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sksmith/note-server/config"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}

func TestLoadLocalConfigs(t *testing.T) {
	config.AppVersion = "appversion"
	config.Sha1Version = "sha1version"
	config.BuildTime = "buildtime"

	cfg, err := config.LoadConfigs()
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	expect(cfg.AppVersion, "appversion", t)
	expect(cfg.ApplicationName, config.ApplicationName, t)
	expect(cfg.BucketName, config.DefaultBucket, t)
	expect(cfg.BuildTime, "buildtime", t)
	expect(cfg.Profile, config.DefaultProfile, t)
	expect(cfg.Port, config.DefaultPort, t)
	expect(cfg.Region, config.DefaultRegion, t)
	expect(cfg.Revision, config.Revision, t)
	expect(cfg.Sha1Version, "sha1version", t)
	expect(cfg.LogLevel, config.DefaultLocalLogLevel, t)
	expect(cfg.LogText, config.DefaultLocalLogText, t)
}

func TestLoadOverriddenConfigs(t *testing.T) {
	const (
		expProfile = "dev"
		expPort    = "9999"
		expRegion  = "some-region"
		expBucket  = "some-bucket"
	)
	addArg("-P", expProfile)
	addArg("-p", expPort)
	addArg("-r", expRegion)
	addArg("-b", expBucket)

	config.AppVersion = "appversion"
	config.Sha1Version = "sha1version"
	config.BuildTime = "buildtime"

	cfg, err := config.LoadConfigs()
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	expect(cfg.AppVersion, "appversion", t)
	expect(cfg.ApplicationName, config.ApplicationName, t)
	expect(cfg.BucketName, expBucket, t)
	expect(cfg.BuildTime, "buildtime", t)
	expect(cfg.Profile, expProfile, t)
	expect(cfg.Port, expPort, t)
	expect(cfg.Region, expRegion, t)
	expect(cfg.Revision, config.Revision, t)
	expect(cfg.Sha1Version, "sha1version", t)
	expect(cfg.LogLevel, config.DefaultEnvironmentLogLevel, t)
	expect(cfg.LogText, config.DefaultEnvironmentLogText, t)
}

func addArg(flag, value string) {
	os.Args = append(os.Args, flag)
	os.Args = append(os.Args, value)
}

func expect(got, want interface{}, t *testing.T) {
	if got != want {
		t.Errorf("got=[%v] want=[%v]", got, want)
	}
}
