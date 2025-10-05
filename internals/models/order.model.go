package models

import (
	"time"
)

type OrderHistory struct {
	OrderID    uint16     `db:"id" json:"id"`
	UserID     uint16     `db:"user_id" json:"user_id"`
	Title      string     `db:"title" json:"title" example:"Pulp Fiction"`
	Date       time.Time  `db:"date" json:"date"`
	Time       string     `db:"time" json:"time"`
	CinemaName string     `db:"cinema_name" json:"cinema_name" example:"ebv"`
	PaidAt     *time.Time `json:"paid_at"`
	Seats      []string   `json:"seats"`
}

// type CinemaOrder struct {
// 	UID           *uint16    `db:"user_id" json:"user_id,omitempty"`
// 	ScheduleID    uint16     `db:"schedule_id" json:"schedule_id" example:"12"`
// 	PaymentMethod string     `db:"payment_method" json:"payment_method" example:"PayPal"`
// 	Total         uint16     `db:"total" json:"total" example:"30"`
// 	Seats         []int      `json:"seats"`
// 	PaidAt        *time.Time `db:"paid_at" json:"paid_at"`
// }

type CinemaOrderBody struct {
	UID           *uint16 `db:"user_id" json:"user_id,omitempty"`
	ScheduleID    uint16  `db:"schedule_id" json:"schedule_id" example:"12"`
	PaymentMethod string  `db:"payment_method" json:"payment_method" example:"PayPal"`
	Total         uint16  `db:"total" json:"total" example:"30"`
	Seats         []int   `json:"seats"`
	PaidAt        *bool   `json:"paid_at"`
}

type OrderResponse struct {
	Result  string
	Success bool
	Error   string
}

type OrderHistoriesResponse struct {
	Result  []OrderHistory `json:"result"`
	Success bool           `json:"success"`
	Error   string         `json:"error"`
}

type SeatBody struct {
	ID int `json:"id"`
}
