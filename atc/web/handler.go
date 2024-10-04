package web

import (
	"atc/atc"
	"atc/static"
	"log/slog"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandlerConfig struct {
	UseEmbeddedFS  bool
	UseBasicAuth   bool
	BasicAuthRealm string
	BasicAuthUser  string
	BasicAuthPass  string
}

type handler struct {
	cfg     *HandlerConfig
	logger  *slog.Logger
	plt     *atc.Platform
	mainjs  string
	maincss string
}

func Handler(logger *slog.Logger, plt *atc.Platform, cfg *HandlerConfig) http.Handler {
	h := &handler{
		cfg:    cfg,
		logger: logger,
		plt:    plt,
	}

	if cfg.UseEmbeddedFS {
		// Use hashfs to serve static files to allow for caching.
		h.mainjs = static.FS.HashName("dist/main.js")
		h.maincss = static.FS.HashName("dist/main.css")
	} else {
		h.mainjs = "dist/main.js"
		h.maincss = "dist/main.css"
	}

	r := chi.NewRouter()

	r.Use(
		gziphandler.GzipHandler,
		middleware.RequestSize(1<<20),
		middleware.StripSlashes,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)
	r.NotFound(h.handleNotFound)

	if cfg.UseBasicAuth {
		r.Use(middleware.BasicAuth(cfg.BasicAuthRealm, map[string]string{
			cfg.BasicAuthUser: cfg.BasicAuthPass,
		}))
	}

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/text/robots.txt")
	})

	if cfg.UseEmbeddedFS {
		r.Handle("/static/*", http.StripPrefix("/static/", hashfs.FileServer(static.FS)))
	} else {
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	}

	r.Get("/", h.handleIndexPage())
	return r
}
