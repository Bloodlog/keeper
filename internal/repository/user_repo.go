package repository

import (
	"context"
	"fmt"
	"keeper/internal/entity"

	pgx "github.com/jackc/pgx/v5"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	GetByLogin(ctx context.Context, tx pgx.Tx, login string) (entity.User, error)
	Register(ctx context.Context, tx pgx.Tx, user entity.User) (entity.User, error)
}

type userRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
	return &userRepository{
		Pool: db,
	}
}

func (r *userRepository) GetByLogin(ctx context.Context, tx pgx.Tx, login string) (entity.User, error) {
	var user entity.User
	query := `
		SELECT id, login, password
		FROM users
		WHERE login = $1
	`
	err := tx.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return user, fmt.Errorf("failed to get login: %w", err)
	}
	return user, nil
}

func (r *userRepository) Register(ctx context.Context, tx pgx.Tx, user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id
	`
	err := tx.QueryRow(ctx, query, user.Login, user.Password).Scan(&user.ID)
	if err != nil {
		return user, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}
