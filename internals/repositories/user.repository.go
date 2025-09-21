package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/models"
)

type UserRepository struct {
	dbpool *pgxpool.Pool
}

func NewUserRepository(dbpool *pgxpool.Pool) *UserRepository {
	return &UserRepository{dbpool: dbpool}
}

func (u *UserRepository) GetUserinf(ctx context.Context, id uint16) (models.UserInf, error) {
	sql := `
		SELECT
			user_id, first_name, last_name, phone_number, point_count, avatar
		FROM
			personal_info
		WHERE
			user_id = $1
	`

	var userinf models.UserInf
	if err := u.dbpool.QueryRow(ctx, sql, uint16(id)).Scan(
		&userinf.UID,
		&userinf.FirstName,
		&userinf.LastName,
		&userinf.PhoneNumber,
		&userinf.PointCount,
		&userinf.Avatar,
	); err != nil {
		return models.UserInf{}, err
	}

	return userinf, nil
}

func (u *UserRepository) UpdateUserinf(newUserInf models.NewInf, ctx context.Context, id uint16, avatarPath string) (pgconn.CommandTag, error) {
	tx, err := u.dbpool.Begin(ctx)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	defer tx.Rollback(ctx)

	rt := reflect.TypeOf(newUserInf)
	rv := reflect.ValueOf(newUserInf)

	var setClauses []string
	var args []any
	argIndex := 1

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		dbTag := field.Tag.Get("db")
		if dbTag == "" || value.IsZero() || dbTag == "avatar" {
			continue
		}

		setClause := fmt.Sprintf("%s = $%d", dbTag, argIndex)
		setClauses = append(setClauses, setClause)
		args = append(args, value.Interface())
		argIndex++
	}

	if avatarPath != "" {
		setClauses = append(setClauses, fmt.Sprintf("avatar = $%d", argIndex))
		args = append(args, avatarPath)
		argIndex++
	}

	setClauses = append(setClauses, "updated_at = current_timestamp")

	// Final query
	sql := fmt.Sprintf("UPDATE personal_info SET %s WHERE user_id = $%d", strings.Join(setClauses, ", "), argIndex)
	args = append(args, id)

	log.Println(sql)

	ctag, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	if ctag.RowsAffected() > 0 {
		tx.Commit(ctx)
		return ctag, nil
	}

	return pgconn.CommandTag{}, nil
	// return
}

func (u *UserRepository) GetUserOrderHistory(ctx context.Context, id uint16) (models.UserOrder, error) {
	sql := `
		SELECT
			b.id "order_id", u.id "user_id", m.title, s.date, t.time, ct.name, b.is_paid
		FROM
			orders AS b
		JOIN
			users AS u ON b.user_id = u.id
		JOIN
			schedule AS s ON b.schedule_id = s.id
		JOIN
			movies AS m ON s.movie_id = m.id
		JOIN
			jam_tayang AS t ON s.time_id = t.id
		JOIN
			cinema_tayang AS ct ON s.cinema_id = ct.id
		WHERE
			u.id = $1
	`
	rows, err := u.dbpool.Query(ctx, sql, id)
	if err != nil {
		return models.UserOrder{}, err
	}
	defer rows.Close()

	var histories []models.OrderHistory
	for rows.Next() {
		var history models.OrderHistory
		if err := rows.Scan(&history.OrderID, &history.UserID, &history.Title, &history.Date, &history.Time, &history.CinemaName, &history.IsPaid); err != nil {
			return models.UserOrder{}, err
		}

		seats, err := u.getOrderSeats(ctx, int(history.OrderID))
		if err != nil {
			return models.UserOrder{}, err
		}
		history.Seats = append(history.Seats, seats...)

		histories = append(histories, history)
	}

	return models.UserOrder{UID: uint16(id), OrderHistory: histories}, nil
}

func (u *UserRepository) getOrderSeats(ctx context.Context, bookId int) ([]string, error) {
	sql := `
		SELECT s.pos FROM orders_seats bs
		JOIN seats s ON s.id = bs.seat_id
		WHERE order_id = $1;
	`
	rows, err := u.dbpool.Query(ctx, sql, bookId)
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

func (u *UserRepository) UpdateUserPassword(ctx context.Context, newPassword string, id uint16, currTime time.Time) (pgconn.CommandTag, error) {
	tx, err := u.dbpool.Begin(ctx)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	defer tx.Rollback(ctx)

	sql := `
		UPDATE users 
		SET password = $1
		WHERE id = $2
	`
	ctagPwd, err := tx.Exec(ctx, sql, newPassword, id)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	if !ctagPwd.Update() {
		return pgconn.CommandTag{}, errors.New("unable to edit password")
	}
	if ctagPwd.RowsAffected() < 1 {
		return pgconn.CommandTag{}, errors.New("user id not found")
	}

	ctagInf, err := u.fixUpdateAt(ctx, tx, id)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	if !ctagInf.Update() {
		return pgconn.CommandTag{}, errors.New("unable to update user profile")
	}
	if ctagInf.RowsAffected() != 0 {
		tx.Commit(ctx)
	}

	return ctagInf, nil
}

func (u *UserRepository) GetLastUpdated(ctx context.Context, id uint16) (time.Time, error) {
	sql := `
		SELECT 
			updated_at
		FROM
			personal_info
		WHERE
			user_id = $1
	`
	var lastUpdate time.Time
	if err := u.dbpool.QueryRow(ctx, sql, id).Scan(&lastUpdate); err != nil {
		return time.Time{}, err
	}

	return lastUpdate, nil
}

func (u *UserRepository) fixUpdateAt(ctx context.Context, tx pgx.Tx, id uint16) (pgconn.CommandTag, error) {
	sql := `
		UPDATE
			personal_info
		SET
			updated_at = current_timestamp
		WHERE
			user_id = $1
	`
	return tx.Exec(ctx, sql, id)
}
