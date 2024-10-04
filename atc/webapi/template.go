package webapi

import (
	"atc/static"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Masterminds/sprig/v3"
)

// ParseTemplate parses the template at the specified path.
// It ensures that the base template and all components are included in the template.
// It also loads any layouts that are specified.
func ParseTemplate(path string, layouts ...string) *template.Template {
	pattern := ""
	if strings.HasSuffix(path, ".tmpl") {
		pattern = fmt.Sprintf("html/app/%s", path)
	} else {
		pattern = fmt.Sprintf("html/app/%s/*.tmpl", path)
	}

	args := []string{
		"html/app/base.tmpl",
		"html/components/*",
	}
	for _, layout := range layouts {
		layoutTmpl := fmt.Sprintf("html/app/%s/layout.tmpl", layout)
		if _, err := static.FS.Open(layoutTmpl); err == nil {
			args = append(args, layoutTmpl)
		}
	}
	args = append(args, pattern)

	return template.Must(
		template.New(path).Funcs(sprig.HtmlFuncMap()).ParseFS(static.FS,
			args...,
		))
}

func RenderTemplate(w http.ResponseWriter, tmp *template.Template, name string, data any) {
	w.Header().Set("Vary", "Hx-Request")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer
	if err := tmp.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
