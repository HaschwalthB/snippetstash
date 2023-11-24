package main

import (
	"net/http"
)

func secureHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self'; fonts.googleapis.com; font-src font.gstatic.com")
		w.Header().Set("referrer-policy", "origin-when-cross-origin")
		w.Header().Set("x-content-type-options", "nosniff")
		w.Header().Set("x-frame-options", "deny")
		w.Header().Set("X-XSS-protection", "0")

    next.ServeHTTP(w, r)
	})
}

