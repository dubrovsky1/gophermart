package models

import "time"

type (
	UserID  string
	OrderID int
)

type User struct {
	UserID   UserID
	Login    string
	Password string
}

type Order struct {
	OrderID     OrderID `json:"number"`
	UserID      UserID  `json:"userid,omitempty"`
	Upload      string  `json:"uploaded_at,omitempty"`
	Type        string  `json:"type,omitempty"`
	Status      string  `json:"status,omitempty"`
	Accrual     float64 `json:"accrual,omitempty"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetOrderListResult struct {
	OrderID string    `db:"orderid"`
	Status  string    `db:"status"`
	Accrual float64   `db:"accrual"`
	Upload  time.Time `db:"uploaded_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	OrderID     OrderID `json:"order,string"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

type WithdrawListResult struct {
	OrderID     string    `db:"orderid"`
	Sum         float64   `db:"accrual"`
	ProcessedAt time.Time `db:"processed_at"`
}

type OrderAccrual struct {
	OrderID OrderID `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
