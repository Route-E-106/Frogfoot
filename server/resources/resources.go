package resources

import (
	"context"
	"sort"
	"time"

	"github.com/Route-E-106/Frogfoot/server/internal/database/models"
)

type ResourcesHistory struct {
	Gas           []models.ReturnGasIncomeHistoryRow   `json:"gas"`
	Metal         []models.ReturnMetalIncomeHistoryRow `json:"metal"`
	TotalExpenses models.GetUserExpensesRow            `json:"expenses"`
}

func (r *ResourcesHistory) RetriveResourcesHistory(userId int64, queries *models.Queries, ctx context.Context) error {

	var err error
	r.Gas, err = queries.ReturnGasIncomeHistory(ctx, userId)
	if err != nil {
		return err
	}
	r.Metal, err = queries.ReturnMetalIncomeHistory(ctx, userId)
	if err != nil {
		return err
	}
	r.TotalExpenses, err = queries.GetUserExpenses(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

type Resources struct {
	Gas   int64 `json:"gas"`
	Metal int64 `json:"metal"`
}

func (r *Resources) CalculateResources(resourcesHistory ResourcesHistory) Resources {

	currentTime := time.Now().Unix()

	gas := resourcesHistory.Gas
	metal := resourcesHistory.Metal
	sort.Slice(gas, func(i, j int) bool {
		return gas[i].ChangeTimestamp < (gas[j].ChangeTimestamp)
	})
	sort.Slice(metal, func(i, j int) bool {
		return metal[i].ChangeTimestamp < (metal[j].ChangeTimestamp)
	})

	var (
		gasIncomeTotal   int64
		metalIncomeTotal int64
	)

	for i := range gas {
		start := gas[i].ChangeTimestamp

		if currentTime < start {
			break
		}

		var end int64
		if i+1 < len(gas) {
			end = gas[i+1].ChangeTimestamp
		} else {
			end = currentTime
		}

		if end > currentTime {
			end = currentTime
		}

		duration := float64(end-start) / 3600
		if duration > 0 {
			gasIncomeTotal += int64(duration * float64(gas[i].Income))
		}
	}

	for i := range metal {
		start := metal[i].ChangeTimestamp

		if currentTime < start {
			break
		}

		var end int64
		if i+1 < len(metal) {
			end = metal[i+1].ChangeTimestamp
		} else {
			end = currentTime
		}

		if end > currentTime {
			end = currentTime
		}

		duration := float64(end-start) / 3600
		if duration > 0 {
			metalIncomeTotal += int64(duration * float64(metal[i].Income))
		}
	}

	r.Gas = gasIncomeTotal - resourcesHistory.TotalExpenses.TotalGasExpenses
	r.Metal = metalIncomeTotal - resourcesHistory.TotalExpenses.TotalMetalExpenses

	return *r
}
