// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package models

type GasIncomeHistory struct {
	ID              int64 `json:"id"`
	Income          int64 `json:"income"`
	UserID          int64 `json:"user_id"`
	ChangeTimestamp int64 `json:"change_timestamp"`
}

type MetalIncomeHistory struct {
	ID              int64 `json:"id"`
	Income          int64 `json:"income"`
	UserID          int64 `json:"user_id"`
	ChangeTimestamp int64 `json:"change_timestamp"`
}

type Session struct {
	Token  interface{} `json:"token"`
	Data   []byte      `json:"data"`
	Expiry interface{} `json:"expiry"`
}

type User struct {
	ID                 int64  `json:"id"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	CreatedAt          int64  `json:"created_at"`
	MetalExtractorLvl  int64  `json:"metal_extractor_lvl"`
	GasExtractorLvl    int64  `json:"gas_extractor_lvl"`
	TotalGasExpenses   int64  `json:"total_gas_expenses"`
	TotalMetalExpenses int64  `json:"total_metal_expenses"`
}
