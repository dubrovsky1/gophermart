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
	OrderID OrderID `json:"number"`
	UserID  UserID  `json:"userid,omitempty"`
	Upload  string  `json:"uploaded_at,omitempty"`
	Type    string  `json:"type,omitempty"`
	Status  string  `json:"status,omitempty"`
	Accrual int     `json:"accrual,omitempty"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetOrderListResult struct {
	OrderID string    `db:"orderid"`
	Status  string    `db:"status"`
	Accrual int       `db:"accrual"`
	Upload  time.Time `db:"uploaded_at"`
}

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

type Withdraw struct {
	OrderID OrderID `json:"order,string"`
	Sum     int     `json:"sum"`
	Upload  string  `json:"processed_at,omitempty"`
}

type WithdrawListResult struct {
	OrderID string    `db:"orderid"`
	Sum     int       `db:"accrual"`
	Upload  time.Time `db:"uploaded_at"`
}

type OrderAccrual struct {
	OrderID OrderID `json:"order"`
	Status  string  `json:"status"`
	Accrual int     `json:"accrual,omitempty"`
}
