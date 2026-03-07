package usecase

import (
	"context"
	"errors"
	"time"

	"taskmanager/domain"
	"taskmanager/repository"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, userID int, name string) (*domain.Category, error)
	GetCategories(ctx context.Context, userID int) ([]domain.Category, error)
}

type categoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUseCase) CreateCategory(ctx context.Context, userID int, name string) (*domain.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}

	category := &domain.Category{
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := u.categoryRepo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (u *categoryUseCase) GetCategories(ctx context.Context, userID int) ([]domain.Category, error) {
	return u.categoryRepo.GetCategoriesByUserID(ctx, userID)
}
