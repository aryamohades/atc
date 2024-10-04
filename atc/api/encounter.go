package api

import (
	"log/slog"
	"net/http"
)

func (h *handler) handleResetEncounters() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := h.plt.ResetEncounters(ctx); err != nil {
			h.logger.Error("reset encounters failed", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
