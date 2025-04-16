package server

import (
	"net/http"

	"github.com/Route-E-106/Frogfoot/server/helpers"
)

func (s *Server) logRequests(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		s.Logger.Info("Received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
func (s *Server) requireAuthentication(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(s.sessionManager, r) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
