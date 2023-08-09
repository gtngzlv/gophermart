package model

import "time"

type GetOrdersResponse struct {
	UserID     int
	Number     uint      `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
