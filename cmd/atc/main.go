package main

import (
	"atc/atc"
	"atc/atc/api"
	"atc/atc/httpapi"
	"atc/atc/postgres"
	"atc/atc/web"
	"atc/atc/webhook"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	defer func() {
		if err := recover(); err != nil {
			logger.Error("recovered from panic",
				slog.Any("error", err),
				slog.String("stack", string(debug.Stack())),
			)
		}
	}()

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	db, err := postgres.Connect(ctx, &postgres.ConnectConfig{
		ConnectionURL:     cfg.Database.ConnectionURL,
		MinOpenConns:      cfg.Database.MinOpenConns,
		MaxOpenConns:      cfg.Database.MaxOpenConns,
		MaxConnLifeTime:   cfg.Database.MaxConnLifeTime,
		MaxConnIdleTime:   cfg.Database.MaxConnIdleTime,
		HealthCheckPeriod: cfg.Database.HealthCheckPeriod,
		Logger:            logger,
	})
	if err != nil {
		return fmt.Errorf("postgres connect: %w", err)
	}
	defer db.Close()

	plt := atc.NewPlatform(logger, db, &atc.PlatformConfig{})

	if err := plt.SimulateEncounters(ctx); err != nil {
		logger.Error("failed to simulate encounters", slog.String("error", err.Error()))
		return err
	}

	webHandler := web.Handler(logger, plt, &web.HandlerConfig{
		UseEmbeddedFS: cfg.Env != "development",
		UseBasicAuth:  cfg.Env == "staging",
		BasicAuthUser: cfg.Web.BasicAuthUser,
		BasicAuthPass: cfg.Web.BasicAuthPassword,
	})
	apiHandler := api.Handler(logger, plt)
	webhookHandler := webhook.Handler(logger, plt)

	handler := chi.NewRouter()
	handler.Mount("/", webHandler)
	handler.Mount("/api", apiHandler)
	handler.Mount("/webhooks", webhookHandler)

	logger.Info("starting server", slog.String("port", fmt.Sprintf("%d", cfg.Server.Port)))

	err = httpapi.StartServer(ctx, handler, &httpapi.ServerConfig{
		Port: cfg.Server.Port,
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	plt.Done()
	return nil
}

type config struct {
	Env    string `env:"ENV,unset" envDefault:"development"`
	Server struct {
		Port int `env:"SERVER_PORT,unset" envDefault:"8080"`
	}
	Database struct {
		ConnectionURL     string        `env:"DATABASE_CONNECTION_URL,unset" envDefault:"postgres://postgres:postgres@localhost:5432/atc?sslmode=disable"`
		MinOpenConns      int           `env:"DATABASE_MIN_OPEN_CONNS,unset" envDefault:"1"`
		MaxOpenConns      int           `env:"DATABASE_MAX_OPEN_CONNS,unset" envDefault:"10"`
		MaxConnLifeTime   time.Duration `env:"DATABASE_MAX_CONN_LIFE_TIME,unset" envDefault:"30m"`
		MaxConnIdleTime   time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME,unset" envDefault:"10m"`
		HealthCheckPeriod time.Duration `env:"DATABASE_HEALTH_CHECK_PERIOD,unset" envDefault:"3m"`
	}
	Web struct {
		BasicAuthRealm    string `env:"WEB_BASIC_AUTH_REALM,unset"`
		BasicAuthUser     string `env:"WEB_BASIC_AUTH_USER,unset"`
		BasicAuthPassword string `env:"WEB_BASIC_AUTH_PASSWORD,unset"`
	}
}
