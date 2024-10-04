package httpmw

import (
	"net/http"
	"strings"
	"time"
)

// Unix epoch time
var epoch = time.Unix(0, 0).UTC().Format(http.TimeFormat)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
// a router (or subrouter) from being cached by an upstream proxy and/or client.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//
// Expires: Thu, 01 Jan 1970 00:00:00 UTC
// Cache-Control: no-cache, private, max-age=0
// X-Accel-Expires: 0
// Pragma: no-cache (for HTTP/1.0 proxies/clients)
func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignore requests to /static/* as these are assumed to be static assets that can and should be cached.
		if strings.HasPrefix(r.URL.Path, "/static") {
			h.ServeHTTP(w, r)
			return
		}
		// Delete any ETag headers that may have been set.
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}
		// Set the NoCache headers.
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}
		h.ServeHTTP(w, r)
	})
}
