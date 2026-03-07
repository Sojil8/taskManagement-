package usecase

import "context"

type EmailService interface {
	SendOTP(ctx context.Context, toEmail, code string) error
}
