package domain

import "time"

type OTP struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}
