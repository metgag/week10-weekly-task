package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	dbpool *pgxpool.Pool
	rdb    *redis.Client
}

func NewAuthRepository(dbpool *pgxpool.Pool, rdb *redis.Client) *AuthRepository {
	return &AuthRepository{dbpool: dbpool, rdb: rdb}
}

func (a *AuthRepository) AddNewUser(ctx context.Context, email, password string) (uint16, error) {
	tx, err := a.dbpool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	sql := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id
	`
	var id uint16
	if err := tx.QueryRow(ctx, sql, email, password).Scan(&id); err != nil {
		return 0, err
	}

	ctag, err := a.initUserProfile(tx, ctx, id)
	if err != nil {
		return 0, err
	} else {
		if err := tx.Commit(ctx); err != nil {
			return 0, err
		}
	}
	log.Printf("%s: USER PROFILE", ctag.String())

	return id, nil
}

func (a *AuthRepository) GetUser(ctx context.Context, email string) (models.User, error) {
	sql := `
		SELECT id, email, password, role
		FROM users
		WHERE email = $1
	`

	var user models.User
	if err := a.dbpool.QueryRow(ctx, sql, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	); err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}

func (a *AuthRepository) initUserProfile(tx pgx.Tx, ctx context.Context, id uint16) (pgconn.CommandTag, error) {
	sql := `
		INSERT INTO 
			personal_info (user_id)
		VALUES
			($1)
	`
	return tx.Exec(ctx, sql, id)
}

func (a *AuthRepository) SetLogoutCache(ctx context.Context, token string, iAt time.Time) error {
	redisKey := fmt.Sprintf("archie:blacklist_%s", token)

	expiresAt := iAt.Add(40 * time.Minute)
	expiration := time.Until(expiresAt)
	if expiration <= 0 {
		return nil
	}

	// simpan boolean true sebagai value, durasi expiration
	err := a.rdb.Set(ctx, redisKey, true, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set redis key: %w", err)
	}

	log.Println("redis> TOKEN BLACKLISTED")
	return nil
}
