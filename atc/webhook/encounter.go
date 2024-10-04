package webhook

import (
	"encoding/json"
	"net/http"
)

func (h *handler) handleUpdateEncounter() http.HandlerFunc {
	type request struct{}

	type response struct{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response{}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
