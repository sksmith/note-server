package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/note-server/api"
	"github.com/sksmith/note-server/config"
	"github.com/sksmith/note-server/core"
	"github.com/sksmith/note-server/core/note"
	"github.com/sksmith/note-server/core/user"
	"github.com/sksmith/note-server/repo/noterepo"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := loadConfigs()

	configLogging(cfg)
	printLogHeader(cfg)

	log.Info().Msg("creating note repository...")
	repo := createNoteRepo(cfg)

	log.Info().Msg("creating note service...")
	noteService := note.NewService(core.NewClock(), repo)

	log.Info().Msg("creating user service...")
	userService := user.NewService()

	log.Info().Msg("configuring router...")
	r := configureRouter(cfg, userService, noteService)

	log.Info().Str("port", cfg.Port).Msg("listening")
	log.Fatal().Err(http.ListenAndServe(":"+cfg.Port, r))
}

func createNoteRepo(cfg config.Config) note.Repository {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	}))
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	deleter := s3.New(sess)
	return noterepo.NewS3Repo(uploader, downloader, deleter, cfg.BucketName)
}

func loadConfigs() (cfg config.Config) {
	var err error

	log.Info().Msg("loading configurations...")
	cfg, err = config.LoadConfigs()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configurations")
	}

	return cfg
}

func printLogHeader(c config.Config) {
	if c.LogText {
		log.Info().Msg("=============================================")
		log.Info().Msg(fmt.Sprintf("    Application: %s", c.ApplicationName))
		log.Info().Msg(fmt.Sprintf("       Revision: %s", c.Revision))
		log.Info().Msg(fmt.Sprintf("        Profile: %s", c.Profile))
		log.Info().Msg(fmt.Sprintf("    Tag Version: %s", c.AppVersion))
		log.Info().Msg(fmt.Sprintf("   Sha1 Version: %s", c.Sha1Version))
		log.Info().Msg(fmt.Sprintf("     Build Time: %s", c.BuildTime))
		log.Info().Msg("=============================================")
	} else {
		log.Info().Str("application", c.ApplicationName).
			Str("revision", c.Revision).
			Str("profile", c.Profile).
			Str("version", c.AppVersion).
			Str("sha1ver", c.Sha1Version).
			Str("build-time", c.BuildTime).
			Send()
	}
}

func configureRouter(cfg config.Config, userService user.Service, service api.NoteService) chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*.seanksmith.me", "http://*.seanksmith.me", "http://localhost*", "https://localhost*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(api.Metrics)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(api.Logging)

	api.ConfigureMetrics()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("UP"))
	})

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/env", envApi(cfg))

	r.With(api.Authenticate(userService)).Route("/api/v1", func(r chi.Router) {
		r.Route("/note", noteApi(service))
	})

	return r
}

func envApi(cfg config.Config) func(r chi.Router) {
	envApi := api.NewEnvApi(cfg)
	return envApi.ConfigureRouter
}

func noteApi(s api.NoteService) func(r chi.Router) {
	noteApi := api.NewNoteApi(s)
	return noteApi.ConfigureRouter
}

func configLogging(cfg config.Config) {
	log.Info().Msg("configuring logging...")

	if cfg.LogText {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Warn().Str("loglevel", cfg.LogLevel).Err(err).Msg("defaulting to info")
		level = zerolog.InfoLevel
	}
	log.Info().Str("loglevel", level.String()).Msg("setting log level")
	zerolog.SetGlobalLevel(level)
}
