package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/models"
)

type CinemaRepository struct {
	dbpool *pgxpool.Pool
}

func NewCinemaRepository(dbpool *pgxpool.Pool) *CinemaRepository {
	return &CinemaRepository{dbpool: dbpool}
}

func (c *CinemaRepository) GetSchedule(ctx context.Context) ([]models.CinemaSchedule, error) {
	sql := `
		SELECT 
			s.id "schedule_id", m.title, s.date, t.time, l.location, c.name "cinema_name"
		FROM 
			schedule AS s
		JOIN
			movies AS m ON s.movie_id = m.id
		JOIN
			jam_tayang AS t ON s.time_id = t.id
		JOIN
			lokasi_tayang AS l ON s.location_id = l.id
		JOIN
			cinema_tayang AS c ON s.cinema_id = c.id;
	`
	rows, err := c.dbpool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.CinemaSchedule
	for rows.Next() {
		var schedule models.CinemaSchedule

		if err := rows.Scan(&schedule.ID, &schedule.Title, &schedule.Date, &schedule.Time, &schedule.Location, &schedule.CinemaName); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (c *CinemaRepository) GetCinemaNameAndTime(ctx context.Context, scheduleId int) (models.CinemaAndTime, error) {
	sql := `
		SELECT ct.cinema_name, jt.show_time, ct.cinema_img
		FROM schedule s
		JOIN cinema_tayang ct ON ct.id = s.cinema_id 
		JOIN jam_tayang jt ON jt.id = s.time_id
		WHERE s.id = $1
	`

	var result models.CinemaAndTime
	if err := c.dbpool.QueryRow(ctx, sql, scheduleId).Scan(
		&result.CinemaName,
		&result.Time,
		&result.CinemaImg,
	); err != nil {
		return models.CinemaAndTime{}, err
	}

	return result, nil
}

func (c *CinemaRepository) GetAvailSeats(ctx context.Context, scheduleId int) ([]models.Seat, error) {
	var scheduleExists bool
	checkScheduleSql := `
		SELECT EXISTS (
			SELECT 1
			FROM "orders"
			WHERE schedule_id = $1
		)
	`
	if err := c.dbpool.QueryRow(ctx, checkScheduleSql, scheduleId).Scan(&scheduleExists); err != nil {
		return nil, err
	}

	// if schedule_id not found
	if !scheduleExists {
		return []models.Seat{}, nil
	}

	sql := `
		SELECT 
			s.id, s.pos
		FROM 
			seats AS s
		WHERE 
			s.id NOT IN (
				SELECT bs.seat_id
				FROM orders_seats bs
				JOIN orders bt ON bs.order_id = bt.id
				WHERE bt.schedule_id = $1
			)
		ORDER BY 
			s.id ASC

	`
	rows, err := c.dbpool.Query(ctx, sql, scheduleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(
			&seat.ID,
			&seat.Pos,
		); err != nil {
			return nil, err
		}

		seats = append(seats, seat)
	}
	return seats, nil
}
