package server

import "net/http"

func (s *Server) handlerGetResources() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

	}
	return http.HandlerFunc(fn)
}
