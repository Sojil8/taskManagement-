package repository

import (
	"context"

	"taskmanager/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

type OTPRepository interface {
	Create(ctx context.Context, otp *domain.OTP) error
	GetLatestValidOTP(ctx context.Context, userID int) (*domain.OTP, error)
	MarkAsUsed(ctx context.Context, otpID int) error
}
