package postgres

import (
	"context"
	"database/sql"
	"errors"

	"taskmanager/domain"
	"taskmanager/repository"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.PasswordHash, user.CreatedAt).Scan(&user.ID)
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
