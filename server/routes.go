package server

import "net/http"

func (s *Server) Routes() http.Handler {

	r := http.NewServeMux()
	r.Handle("GET /resources/{userId}", s.handlerGetResources())
	r.Handle("GET /users", s.handlerGetUsers())
	r.Handle("POST /users/createUser", s.handlerCreateUser())

	return r

}
