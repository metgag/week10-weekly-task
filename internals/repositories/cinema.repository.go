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
			s.id, m.title, s.date, t.time, l.location, c.name "cinema_name"
		FROM 
			cinema_schedule AS s
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

func (c *CinemaRepository) GetAvailSeats(ctx context.Context) ([]models.AvailSeat, error) {
	sql := `
		SELECT 
			s.id, s.pos
		FROM 
			seats AS s
		LEFT JOIN 
			books_seats AS b ON b.seat_id = s.id
		WHERE 
			b.seat_id IS NULL
		ORDER BY 
			s.id ASC
	`
	rows, err := c.dbpool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.AvailSeat
	for rows.Next() {
		var seat models.AvailSeat
		if err := rows.Scan(&seat.ID, &seat.Pos); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, nil
}
