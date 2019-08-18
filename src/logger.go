package main

import (
	"net/http"
	"time"
)

func Logger(c Controller, inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		c.LogInfo(r.Method, r.RequestURI, "from", r.RemoteAddr, "tool", time.Since(start))
	})
}
