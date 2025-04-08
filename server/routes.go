package server

import "net/http"

func (s *Server) Routes() http.Handler {

	r := http.NewServeMux()
	r.Handle("GET /resources/{userId}", s.handlerGetResources())

	return r

}
