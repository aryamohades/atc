package webapi

import "net/http"

func IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get("Hx-Request") == "true"
}

func IsHtmxBoosted(r *http.Request) bool {
	return r.Header.Get("Hx-Boosted") == "true"
}
