package webhook

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
	r.Post("/encounter", h.handleUpdateEncounter())
	return r
}
