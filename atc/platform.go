package atc

import (
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlatformConfig struct {
}

type Platform struct {
	logger *slog.Logger
	cfg    *PlatformConfig
	db     *pgxpool.Pool
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

func NewPlatform(logger *slog.Logger, db *pgxpool.Pool, cfg *PlatformConfig) *Platform {
	return &Platform{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}

func (p *Platform) Done() {
	p.logger.Info("stopping all tasks...")
	p.wg.Wait()
	p.logger.Info("all tasks stopped.")
}
