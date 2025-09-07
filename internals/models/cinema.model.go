package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CinemaSchedule struct {
	ID         uint16      `db:"schedule_id" json:"schedule_id"`
	Title      string      `db:"title" json:"title"`
	Date       pgtype.Date `db:"date" json:"date"`
	Time       time.Time   `db:"time" json:"time"`
	Location   string      `db:"location" json:"location"`
	CinemaName string      `db:"cinema_name" json:"cinema_name"`
}

type ScheduleResponse struct {
	Result  []CinemaSchedule
	Success bool
	Error   string
}

type AvailSeat struct {
	ID  uint8  `db:"id" json:"id"`
	Pos string `db:"pos" json:"pos"`
}

type AvailSeatsResponse struct {
	Result  []AvailSeat
	Success bool
	Error   string
}
