package model

type GetBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

type WithdrawBalanceRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}
