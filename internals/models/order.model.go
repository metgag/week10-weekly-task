package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderHistory struct {
	OrderID    uint16      `db:"id" json:"id"`
	UserID     uint16      `db:"user_id" json:"user_id"`
	Title      string      `db:"title" json:"title"`
	Date       pgtype.Date `db:"date" json:"date"`
	Time       time.Time   `db:"time" json:"time"`
	CinemaName string      `db:"cinema_name" json:"cinema_name"`
	IsPaid     bool        `db:"is_paid" json:"is_paid"`
}

type CinemaOrder struct {
	UID           uint16 `db:"user_id" json:"user_id"`
	ScheduleID    uint16 `db:"schedule_id" json:"schedule_id"`
	PaymentMethod string `db:"payment_method" json:"payment_method"`
	Total         uint16 `db:"total" json:"total"`
	IsPaid        bool   `db:"is_paid" json:"is_paid"`
}

type OrderResponse struct {
	Result  string
	Success bool
	Error   string
}
