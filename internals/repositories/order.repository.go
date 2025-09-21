package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/models"
)

type OrderRepository struct {
	dbpool *pgxpool.Pool
}

func NewOrderRepository(dbpool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{dbpool: dbpool}
}

func (o *OrderRepository) GetOrderHistories(ctx context.Context) ([]models.OrderHistory, error) {
	sql := `
		SELECT 
			b.id "order_id", b.user_id, m.title, s.date, t.time, ct.name, b.is_paid
		FROM
			book_ticket AS b
		JOIN
			cinema_schedule AS c ON b.schedule_id = c.id
		JOIN
			movies AS m ON c.movie_id = m.id
		JOIN
			cinema_schedule AS s ON b.schedule_id = s.id
		JOIN
			jam_tayang AS t ON s.time_id = t.id
		JOIN
			cinema_tayang AS ct ON s.cinema_id = ct.id
	`
	rows, err := o.dbpool.Query(ctx, sql)
	if err != nil {
		return []models.OrderHistory{}, err
	}
	defer rows.Close()

	var histories []models.OrderHistory
	for rows.Next() {
		var history models.OrderHistory
		if err := rows.Scan(&history.OrderID, &history.UserID, &history.Title, &history.Date, &history.Time, &history.CinemaName, &history.IsPaid); err != nil {
			return []models.OrderHistory{}, err
		}

		seats, err := o.getOrderSeats(ctx, int(history.OrderID))
		if err != nil {
			return nil, err
		}
		history.Seats = append(history.Seats, seats...)

		histories = append(histories, history)
	}

	return histories, nil
}

func (o *OrderRepository) getOrderSeats(ctx context.Context, bookId int) ([]string, error) {
	sql := `
		SELECT s.pos FROM orders_seats bs
		JOIN seats s ON s.id = bs.seat_id
		WHERE bs.order_id = $1;
	`
	rows, err := o.dbpool.Query(ctx, sql, bookId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []string
	for rows.Next() {
		var seat string
		if err := rows.Scan(&seat); err != nil {
			return nil, err
		}

		seats = append(seats, seat)
	}

	return seats, nil
}

func (o *OrderRepository) CreateOrder(ctx context.Context, body models.CinemaOrder, uid uint16, seats ...int) (string, error) {
	tx, err := o.dbpool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	sql := `
		INSERT INTO
			orders (user_id, schedule_id, payment_method, total, is_paid)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id
	`

	var orderId int
	if err := tx.QueryRow(ctx, sql, uid, body.ScheduleID, body.PaymentMethod, body.Total, body.IsPaid).Scan(&orderId); err != nil {
		return sql, nil
	}

	ctag, err := o.createBookSeats(tx, ctx, orderId, body.Seats)
	if ctag.Insert() {
		tx.Commit(ctx)
		return fmt.Sprintf("%s: CREATE ORDER", ctag.String()), nil
	}

	return "", err
}

func (o *OrderRepository) createBookSeats(tx pgx.Tx, ctx context.Context, orderId int, seats []int) (pgconn.CommandTag, error) {
	sql := `
		INSERT INTO
			orders_seats (order_id, seat_id)
		VALUES
	`
	args := []any{}
	for i, v := range seats {
		sql += fmt.Sprintf("(%d, $%d)", orderId, i+1)
		if i < len(seats)-1 {
			sql += ", "
		}
		args = append(args, v)
	}

	log.Println(sql)

	return tx.Exec(ctx, sql, args...)
	// return sql
}

// func (o *OrderRepository) getPaymentInfo(ctx context.Context)

// func (o *OrderRepository) GetMovieSchedule(ctx context.Context, movieId int) ([]models.MovieSchedule, error) {
// 	sql := `
// 		SELECT
// 			cs.id, m.title, cs.date, t.time, l.location, c.name
// 		FROM
// 			cinema_schedule cs
// 		JOIN
// 			movies m ON cs.movie_id = m.id
// 		JOIN
// 			jam_tayang t ON cs.time_id = t.id
// 		JOIN
// 			lokasi_tayang l ON cs.location_id = l.id
// 		JOIN
// 			cinema_tayang c ON cs.cinema_id = cs.cinema_id
// 		WHERE m.id = $1

// 	`
// 	rows, err := o.dbpool.Query(ctx, sql, movieId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var schedules []models.MovieSchedule
// 	for rows.Next() {
// 		var schedule models.MovieSchedule
// 		if err := rows.Scan(
// 			&schedule.ScheduleID,
// 			&schedule.Title,
// 			&schedule.Date,
// 			&schedule.Time,
// 			&schedule.Location,
// 			&schedule.Cinema,
// 		); err != nil {
// 			return nil, err
// 		}

// 		schedules = append(schedules, schedule)
// 	}

// 	return schedules, nil
// }
