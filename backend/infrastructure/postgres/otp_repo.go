package postgres

import (
	"context"
	"database/sql"
	"errors"

	"taskmanager/domain"
	"taskmanager/repository"
)

type otpRepo struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) repository.OTPRepository {
	return &otpRepo{db: db}
}

func (r *otpRepo) Create(ctx context.Context, otp *domain.OTP) error {
	query := `INSERT INTO otps (user_id, code, expires_at, used) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, otp.UserID, otp.Code, otp.ExpiresAt, otp.Used).Scan(&otp.ID)
}

func (r *otpRepo) GetLatestValidOTP(ctx context.Context, userID int) (*domain.OTP, error) {
	query := `
		SELECT id, user_id, code, expires_at, used 
		FROM otps 
		WHERE user_id = $1 AND used = false AND expires_at > NOW() 
		ORDER BY created_at DESC LIMIT 1`

	otp := &domain.OTP{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.ExpiresAt, &otp.Used)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return otp, nil
}

func (r *otpRepo) MarkAsUsed(ctx context.Context, otpID int) error {
	query := `UPDATE otps SET used = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, otpID)
	return err
}
