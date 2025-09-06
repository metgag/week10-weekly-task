package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderHistory struct {
	OrderID    uint16      `db:"id" json:"id"`
	Title      string      `db:"title" json:"title"`
	Date       pgtype.Date `db:"date" json:"date"`
	Time       time.Time   `db:"time" json:"time"`
	CinemaName string      `db:"cinema_name" json:"cinema_name"`
	IsPaid     bool        `db:"is_paid" json:"is_paid"`
}
