package server

import (
	"encoding/json"
	"net/http"

	"github.com/Route-E-106/Frogfoot/server/internal/database/models"
)

func (s *Server) handlerGetResources() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Listing users")
	}
	return http.HandlerFunc(fn)
}

func (s *Server) handlerGetUsers() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		users, err := s.Queries.ListUsers(s.ctx)
		s.Logger.Info("Users", "user", users)
		if err != nil {
			s.Logger.Error(err.Error())
		}
		data, err := json.Marshal(&users)
		if err != nil {
			s.Logger.Error(err.Error())
		}
		s.Logger.Info("Listing users")
		w.Write(data)

	}
	return http.HandlerFunc(fn)
}

func (s *Server) handlerCreateUser() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		userParams := models.CreateUserParams{}
		err := json.NewDecoder(r.Body).Decode(&userParams)
		if err != nil {
			s.Logger.Error(err.Error())

		}
		user, err := s.Queries.CreateUser(s.ctx, userParams)
		if err != nil {
			s.Logger.Error(err.Error())

		}
		s.Logger.Info("Creating user", "username", user.Username)
	}
	return http.HandlerFunc(fn)
}
