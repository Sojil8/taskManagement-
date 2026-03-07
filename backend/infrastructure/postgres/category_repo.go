package postgres

import (
	"context"
	"database/sql"
	"taskmanager/domain"
	"taskmanager/repository"
)

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) CreateCategory(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (user_id, name, created_at) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, category.UserID, category.Name, category.CreatedAt).Scan(&category.ID)
}

func (r *categoryRepo) GetCategoriesByUserID(ctx context.Context, userID int) ([]domain.Category, error) {
	query := `SELECT id, user_id, name, created_at FROM categories WHERE user_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
