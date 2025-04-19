package resources

import (
	"context"

	"github.com/Route-E-106/Frogfoot/server/internal/database/models"
)

type Resources struct {
	Gas   []models.ReturnGasIncomeHistoryRow   `json:"gas"`
	Metal []models.ReturnMetalIncomeHistoryRow `json:"metal"`
}

func (r *Resources) RetriveResourcesHistory(userId int64, queries *models.Queries, ctx context.Context) error {

	var err error
	r.Gas, err = queries.ReturnGasIncomeHistory(ctx, userId)
	if err != nil {
		return err
	}
	r.Metal, err = queries.ReturnMetalIncomeHistory(ctx, userId)
	if err != nil {
		return err
	}

	// buffer, err := json.Marshal(r)
	// log.Println(buffer)
	// if err != nil {
	// 	return err
	// }

	return nil
}
