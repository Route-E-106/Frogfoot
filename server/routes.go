package server

import "net/http"

func (s *Server) Routes() http.Handler {

	r := http.NewServeMux()
	r.Handle("GET /resources/history", s.requireAuthentication(s.handlerGetResourcesHistory()))
	r.Handle("GET /resources", s.requireAuthentication(s.handlerGetCurrentResources()))
	r.Handle("GET /users/list/{userId}", s.requireAuthentication(s.handlerGetUsers()))
	r.Handle("GET /buildings/metalExtractor", s.requireAuthentication(s.handlerGetMetalExtractorNextLevelCost()))
	r.Handle("GET /buildings/gasExtractor", s.requireAuthentication(s.handlerGetGasExtractorNextLevelCost()))
	r.Handle("POST /buildings/metalExtractor/upgrade", s.requireAuthentication(s.handlerUpgradeMetalExtractor()))
	r.Handle("POST /buildings/gasExtractor/upgrade", s.requireAuthentication(s.handlerUpgradeGasExtractor()))
	r.Handle("GET /users/list", s.requireAuthentication(s.handlerGetUsers()))
	r.Handle("POST /users/register", s.handlerRegisterUser())
	r.Handle("POST /users/login", s.handlerLoginUser())

	return s.logRequests(s.sessionManager.LoadAndSave(r))

}
