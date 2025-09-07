package models

import (
	"time"
)

type CinemaSchedule struct {
	ID         uint16    `db:"schedule_id" json:"schedule_id"`
	Title      string    `db:"title" json:"title" example:"Pulp Fiction"`
	Date       time.Time `db:"date" json:"date"`
	Time       time.Time `db:"time" json:"time"`
	Location   string    `db:"location" json:"location" example:"Bogor"`
	CinemaName string    `db:"cinema_name" json:"cinema_name" example:"ebv"`
}

type ScheduleResponse struct {
	Result  []CinemaSchedule
	Success bool
	Error   string
}

type AvailSeat struct {
	ID  uint8  `db:"id" json:"id" example:"32"`
	Pos string `db:"pos" json:"pos" example:"C4"`
}

type AvailSeatsResponse struct {
	Result  []AvailSeat
	Success bool
	Error   string
}
