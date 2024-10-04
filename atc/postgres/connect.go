package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultConnectionURL     = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	defaultMinOpenConns      = 5
	defaultMaxOpenConns      = 10
	defaultMaxConnLifetime   = 10 * time.Minute
	defaultMaxConnIdleTime   = 5 * time.Minute
	defaultHealthCheckPeriod = 5 * time.Minute
)

type ConnectConfig struct {
	ConnectionURL     string
	MinOpenConns      int
	MaxOpenConns      int
	MaxConnLifeTime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	Logger            *slog.Logger
}

func Connect(ctx context.Context, cfg *ConnectConfig) (*pgxpool.Pool, error) {
	if cfg.ConnectionURL == "" {
		cfg.ConnectionURL = defaultConnectionURL
	}
	if cfg.MinOpenConns == 0 {
		cfg.MinOpenConns = defaultMinOpenConns
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = defaultMaxOpenConns
	}
	if cfg.MaxConnLifeTime == 0 {
		cfg.MaxConnLifeTime = defaultMaxConnLifetime
	}
	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = defaultMaxConnIdleTime
	}
	if cfg.HealthCheckPeriod == 0 {
		cfg.HealthCheckPeriod = defaultHealthCheckPeriod
	}

	poolConf, err := pgxpool.ParseConfig(cfg.ConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	poolConf.MinConns = int32(cfg.MinOpenConns)
	poolConf.MaxConns = int32(cfg.MaxOpenConns)
	poolConf.MaxConnLifetime = cfg.MaxConnLifeTime
	poolConf.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConf.HealthCheckPeriod = cfg.HealthCheckPeriod
	if cfg.Logger != nil {
		poolConf.ConnConfig.Tracer = &pgxTracer{logger: cfg.Logger}
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}

type ctxKeyQuerySQL struct{}
type ctxKeyQueryArgs struct{}

type pgxTracer struct {
	logger *slog.Logger
}

func (t *pgxTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, ctxKeyQuerySQL{}, data.SQL)
	ctx = context.WithValue(ctx, ctxKeyQueryArgs{}, data.Args)
	return ctx
}

func (t *pgxTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil && !errors.Is(data.Err, context.Canceled) {
		sql := ctx.Value(ctxKeyQuerySQL{}).(string)
		args := ctx.Value(ctxKeyQueryArgs{})
		t.logger.Error("postgres query failed", slog.String("sql", sql), slog.Any("args", args), slog.String("error", data.Err.Error()))
	}
}
