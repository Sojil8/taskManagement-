package repository

import (
	"context"
	"taskmanager/domain"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoriesByUserID(ctx context.Context, userID int) ([]domain.Category, error)
}
