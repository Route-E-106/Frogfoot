package server

import "net/http"

func (s *Server) Routes() http.Handler {

	r := http.NewServeMux()
	r.Handle("GET /resources", s.requireAuthentication(s.handlerGetResources()))
	r.Handle("GET /users/list/{userId}", s.requireAuthentication(s.handlerGetUsers()))
	r.Handle("GET /users/list", s.requireAuthentication(s.handlerGetUsers()))
	r.Handle("POST /users/register", s.handlerRegisterUser())
	r.Handle("POST /users/login", s.handlerLoginUser())

	return s.logRequests(s.sessionManager.LoadAndSave(r))

}
