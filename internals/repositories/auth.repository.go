package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/models"
)

type AuthRepository struct {
	dbpool *pgxpool.Pool
}

func NewAuthRepository(dbpool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{dbpool: dbpool}
}

func (a *AuthRepository) AddNewUser(ctx context.Context, email, password string) (uint16, error) {
	sql := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id
	`
	var id uint16
	if err := a.dbpool.QueryRow(ctx, sql, email, password).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (a *AuthRepository) GetUser(ctx context.Context, email string) (models.User, error) {
	sql := `
		SELECT id, email, password, role
		FROM users
		WHERE email = $1
	`

	var user models.User
	if err := a.dbpool.QueryRow(ctx, sql, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role); err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}
