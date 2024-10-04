package web

import (
	"atc/atc/webapi"
	"net/http"
)

var notFoundTemplate = webapi.ParseTemplate("not_found.tmpl")

var notFoundMeta = &webapi.Meta{
	Title:       "ATC | Not Found",
	Description: "The page you are looking for does not exist.",
	Keywords:    "404, Not Found",
}

func (h *handler) handleNotFound(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Meta":    notFoundMeta,
		"Mainjs":  h.mainjs,
		"Maincss": h.maincss,
	}
	webapi.RenderTemplate(w, notFoundTemplate, "not_found.tmpl", data)
}
