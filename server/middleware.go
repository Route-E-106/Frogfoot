package server

import "net/http"

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
