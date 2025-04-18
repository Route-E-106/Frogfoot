package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Route-E-106/Frogfoot/server/helpers"
	"github.com/Route-E-106/Frogfoot/server/internal/database/models"
)

func (s *Server) handlerGetResources() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		s.Logger.Info("Listing income history")
		userId := s.sessionManager.GetInt64(r.Context(), "userAuthID")

		incomeHistory, err := s.Queries.ReturnIncomeHistory(s.ctx, userId)
		if err != nil {
			s.Logger.Error(err.Error())
		}
		data, err := json.Marshal(&incomeHistory)
		if err != nil {
			s.Logger.Error(err.Error())
		}
		_, err = w.Write(data)
		if err != nil {
			s.Logger.Error(err.Error())
		}
		s.Logger.Info("Income history", "history", incomeHistory)
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
		_, err = w.Write(data)

		if err != nil {
			s.Logger.Error(err.Error())
		}
	}
	return http.HandlerFunc(fn)
}

func (s *Server) handlerRegisterUser() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		var userParams models.CreateUserParams
		createdAt := time.Now().Unix()

		err := helpers.DecodeJSONBody(w, r, &userParams)
		if err != nil {
			s.Logger.Error(err.Error())
			helpers.ClientError(w, err, 400)
			return
		}
		userParams.CreatedAt = createdAt
		user, err := s.Queries.CreateUser(s.ctx, userParams)
		if err != nil {
			s.Logger.Error(err.Error())
			if strings.HasSuffix(err.Error(), "(2067)") {
				err = errors.New("User with provided name already exists")
				helpers.ClientError(w, err, 400)
			}
			return
		}

		var initialIncome int64 = 1000

		initalGasIncome := models.UpdateIncomeHistoryParams{
			ResourceName:    "Gas",
			Income:          initialIncome,
			UserID:          user.ID,
			ChangeTimestamp: createdAt,
		}
		initalMetalIncome := models.UpdateIncomeHistoryParams{
			ResourceName:    "Metal",
			Income:          initialIncome,
			UserID:          user.ID,
			ChangeTimestamp: createdAt,
		}
		err = s.Queries.UpdateIncomeHistory(s.ctx, initalGasIncome)
		if err != nil {
			s.Logger.Error(err.Error())
			helpers.ClientError(w, err, 400)
			return
		}
		err = s.Queries.UpdateIncomeHistory(s.ctx, initalMetalIncome)
		if err != nil {
			s.Logger.Error(err.Error())
			helpers.ClientError(w, err, 400)
			return
		}

		s.Logger.Info("Created user", "username", user.Username)

		msg := fmt.Sprintf("User %s created sucessfully", user.Username)

		w.Write([]byte(msg))
	}
	return http.HandlerFunc(fn)
}

func (s *Server) handlerLoginUser() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		var rcvUserParams models.CreateUserParams
		err := helpers.DecodeJSONBody(w, r, &rcvUserParams)
		if err != nil {
			s.Logger.Error(err.Error())
			return

		}
		user, err := s.Queries.GetUserByUserName(s.ctx, rcvUserParams.Username)
		if err != nil {
			s.Logger.Error(err.Error())
			if errors.Is(err, sql.ErrNoRows) {
				err = errors.New("User doesn't exist")
				s.Logger.Error(err.Error())
				helpers.ClientError(w, err, 400)
				return
			}
		}
		if rcvUserParams.Password == user.Password {
			s.Logger.Info("Logging user", "username", user.Username)
			s.sessionManager.Put(r.Context(), "userAuthID", user.ID)
		} else {
			s.Logger.Info("User provided wrong password", "provided password", rcvUserParams.Password)
			err = errors.New("Wrong password")
			helpers.ClientError(w, err, 401)
		}

	}
	return http.HandlerFunc(fn)
}
