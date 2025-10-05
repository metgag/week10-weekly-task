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
	Result  []CinemaSchedule `json:"schedules,omitempty"`
	Success bool             `json:"success" example:"true"`
	Error   string           `json:"message,omitempty"`
}

type Seat struct {
	ID  uint8  `db:"id" json:"id" example:"32"`
	Pos string `db:"pos" json:"pos" example:"C4"`
}

type AvailSeatsResponse struct {
	Result  []Seat `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type CinemaAndTime struct {
	CinemaName string `json:"cinema_name"`
	Time       string `json:"time"`
	CinemaImg  string `json:"cinema_img"`
}

type CinemaAndTimeResponse struct {
	Result  CinemaAndTime `json:"result"`
	Success bool          `json:"success"`
	Error   string        `json:"error"`
}
