package api

import (
	"atc/atc"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	logger *slog.Logger
	plt    *atc.Platform
}

func Handler(logger *slog.Logger, plt *atc.Platform) http.Handler {
	h := &handler{
		logger: logger,
		plt:    plt,
	}

	r := chi.NewRouter()
	r.Post("/reset", h.handleResetEncounters())
	return r
}
