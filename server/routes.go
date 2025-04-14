package server

import "net/http"

func (s *Server) Routes() http.Handler {

	r := http.NewServeMux()
	r.Handle("GET /resources/{userId}", s.handlerGetResources())
	r.Handle("GET /users/list/{userId}", s.handlerGetUsers())
	r.Handle("GET /users/list", s.handlerGetUsers())
	r.Handle("POST /users/register", s.handlerRegisterUser())
	r.Handle("POST /users/login", s.handlerLoginUser())

	return s.logRequests(r)

}
