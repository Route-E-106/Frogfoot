package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Route-E-106/Frogfoot/server/buildings"
	"github.com/Route-E-106/Frogfoot/server/helpers"
	"github.com/Route-E-106/Frogfoot/server/internal/database/models"
)

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

		initalGasIncome := models.UpdateGasIncomeHistoryParams{
			Income:          initialIncome,
			UserID:          user.ID,
			ChangeTimestamp: createdAt,
		}
		initalMetalIncome := models.UpdateMetalIncomeHistoryParams{
			Income:          initialIncome,
			UserID:          user.ID,
			ChangeTimestamp: createdAt,
		}
		err = s.Queries.UpdateGasIncomeHistory(s.ctx, initalGasIncome)
		if err != nil {
			s.Logger.Error(err.Error())
			helpers.ClientError(w, err, 400)
			return
		}
		err = s.Queries.UpdateMetalIncomeHistory(s.ctx, initalMetalIncome)
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
			w.Write([]byte("User logged in"))
		} else {
			s.Logger.Info("User provided wrong password", "provided password", rcvUserParams.Password)
			err = errors.New("Wrong password")
			helpers.ClientError(w, err, 401)
		}

	}
	return http.HandlerFunc(fn)
}

func (s *Server) handlerGetResourcesHistory() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		userId := s.sessionManager.GetInt64(r.Context(), "userAuthID")
		err := s.resourcesHistory.RetriveResourcesHistory(userId, s.Queries, r.Context())
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
		data, err := json.Marshal(s.resourcesHistory)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}

		s.Logger.Info("Income history", "history", data)

		_, err = w.Write(data)
		if err != nil {
			s.Logger.Error(err.Error())
		}
	}
	return http.HandlerFunc(fn)
}

func (s *Server) handlerGetCurrentResources() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		userId := s.sessionManager.GetInt64(r.Context(), "userAuthID")
		err := s.resourcesHistory.RetriveResourcesHistory(userId, s.Queries, r.Context())
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
		res := s.resources.CalculateResources(s.resourcesHistory)
		data, err := json.Marshal(res)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}

		s.Logger.Info("Current resources", "resources", data)

		_, err = w.Write(data)
		if err != nil {
			s.Logger.Error(err.Error())
		}
	}
	return http.HandlerFunc(fn)
}

func (s *Server) handlerUpgradeMetalExtractor() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		userId := s.sessionManager.GetInt64(r.Context(), "userAuthID")
		err := s.resourcesHistory.RetriveResourcesHistory(userId, s.Queries, r.Context())
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
		res := s.resources.CalculateResources(s.resourcesHistory)
		lvl, err := s.Queries.GetUserMetalExtractorLevel(r.Context(), userId)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}

		cost := buildings.MetalBuildingCostPerLevel[lvl+1]

		if cost.GasCost <= res.Gas && cost.MetalCost <= res.Metal {
			currTime := time.Now().Unix()

			err = s.Queries.UpdateUserMetalExtractorLevel(r.Context(), userId)
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

			totalExpensesUpdate := models.UpdateTotalExpensesParams{
				TotalGasExpenses:   cost.GasCost,
				TotalMetalExpenses: cost.MetalCost,
				ID:                 userId,
			}

			err = s.Queries.UpdateTotalExpenses(r.Context(), totalExpensesUpdate)
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

			incomeUpdate := models.UpdateMetalIncomeHistoryParams{
				Income:          buildings.MetalIncomePerLevel[lvl+1],
				UserID:          userId,
				ChangeTimestamp: currTime,
			}
			err = s.Queries.UpdateMetalIncomeHistory(r.Context(), incomeUpdate)
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

		} else if cost.GasCost >= res.Gas || cost.MetalCost >= res.Metal {
			s.Logger.Info("User has too few resources to upgrade metal extractor")
			helpers.ClientError(w, errors.New("Too few resources"), 428)
		}

	}
	return http.HandlerFunc(fn)
}

func (s *Server) handlerUpgradeGasExtractor() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		userId := s.sessionManager.GetInt64(r.Context(), "userAuthID")
		err := s.resourcesHistory.RetriveResourcesHistory(userId, s.Queries, r.Context())
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
		res := s.resources.CalculateResources(s.resourcesHistory)
		lvl, err := s.Queries.GetUserGasExtractorLevel(r.Context(), userId)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}

		cost := buildings.GasBuildingCostPerLevel[lvl+1]

		if cost.GasCost <= res.Gas && cost.MetalCost <= res.Metal {
			currTime := time.Now().Unix()

			err = s.Queries.UpdateUserGasExtractorLevel(r.Context(), userId)
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

			totalExpensesUpdate := models.UpdateTotalExpensesParams{
				TotalGasExpenses:   cost.GasCost,
				TotalMetalExpenses: cost.MetalCost,
				ID:                 userId,
			}

			err = s.Queries.UpdateTotalExpenses(r.Context(), totalExpensesUpdate)
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

			incomeUpdate := models.UpdateGasIncomeHistoryParams{
				Income:          buildings.MetalIncomePerLevel[lvl+1],
				UserID:          userId,
				ChangeTimestamp: currTime,
			}
			err = s.Queries.UpdateGasIncomeHistory(r.Context(), incomeUpdate)
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

		} else if cost.GasCost >= res.Gas || cost.MetalCost >= res.Metal {
			s.Logger.Info("User has too few resources to upgrade metal extractor")
			helpers.ClientError(w, errors.New("Too few resources"), 428)
		}

	}
	return http.HandlerFunc(fn)
}
