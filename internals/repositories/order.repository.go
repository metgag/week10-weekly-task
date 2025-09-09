package repositories

import (
	"context"

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

		histories = append(histories, history)
	}

	return histories, nil
}

func (o *OrderRepository) CreateOrder(ctx context.Context, body models.CinemaOrder, uid uint16) (string, error) {
	sql := `
		INSERT INTO
			book_ticket (user_id, schedule_id, payment_method, total, is_paid)
		VALUES
			($1, $2, $3, $4, $5)
	`

	ctag, err := o.dbpool.Exec(ctx, sql, uid, body.ScheduleID, body.PaymentMethod, body.Total, body.IsPaid)
	if err != nil {
		return "", err
	}

	return ctag.String(), nil
}
