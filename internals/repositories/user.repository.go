package repositories

import (
	"context"
	"fmt"
	"reflect"

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

func (u *UserRepository) GetUserinf(ctx context.Context, id int) (models.UserInf, error) {
	sql := `
		SELECT
			user_id, first_name, last_name, phone_number, point_count
		FROM
			personal_info
		WHERE
			user_id = $1
	`

	var userinf models.UserInf
	if err := u.dbpool.QueryRow(ctx, sql, id).Scan(&userinf.UID, &userinf.FirstName, &userinf.LastName, &userinf.PhoneNumber, &userinf.PointCount); err != nil {
		return models.UserInf{}, err
	}

	return userinf, nil
}

func (u *UserRepository) UpdateUserinf(newUserInf models.NewInf, ctx context.Context, id int) (pgconn.CommandTag, error) {
	rt := reflect.TypeOf(newUserInf)
	rv := reflect.ValueOf(newUserInf)

	var args []any
	var argIndex int = 1

	sql := "UPDATE personal_info SET "
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		if value.IsZero() {
			continue
		} else {
			args = append(args, value.Interface())
		}

		sql += fmt.Sprintf("%s = $%d", field.Tag.Get("db"), argIndex)
		sql += ", "

		argIndex++
	}

	sql += fmt.Sprintf(" updated_at = current_timestamp WHERE user_id = $%d", argIndex)
	args = append(args, id)

	fmt.Println(sql)

	return u.dbpool.Exec(ctx, sql, args...)
}

func (u *UserRepository) InitUpdateUserinf(newUserInf models.NewInf, ctx context.Context, id int) (pgconn.CommandTag, error) {
	rt := reflect.TypeOf(newUserInf)
	rv := reflect.ValueOf(newUserInf)

	var args []any
	var argIndex []int

	sql := "INSERT INTO personal_info (user_id, "
	args = append(args, id)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		if value.IsZero() {
			continue
		} else {
			args = append(args, value.Interface())
		}

		sql += field.Tag.Get("db")
		sql += ", "

		argIndex = append(argIndex, i+2)
	}

	sql += "updated_at) VALUES ($1, "
	for i, v := range argIndex {
		sql += fmt.Sprintf("$%d", v)
		if i < len(argIndex)-1 {
			sql += ", "
		} else {
			sql += ", current_timestamp"
		}
	}
	sql += ")"

	return u.dbpool.Exec(ctx, sql, args...)
}
