package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"taskmanager/domain"
	"taskmanager/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, name, email, password string) error
	SendOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, email, code string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authUseCase struct {
	userRepo  repository.UserRepository
	otpRepo   repository.OTPRepository
	emailSvc  EmailService
	jwtSecret string
}

func NewAuthUseCase(userRepo repository.UserRepository, otpRepo repository.OTPRepository, emailSvc EmailService, jwtSecret string) AuthUseCase {
	return &authUseCase{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		emailSvc:  emailSvc,
		jwtSecret: jwtSecret,
	}
}

func (u *authUseCase) Register(ctx context.Context, name, email, password string) error {
	existing, _ := u.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
		CreatedAt:    time.Now(),
	}

	return u.userRepo.Create(ctx, user)
}

func (u *authUseCase) SendOTP(ctx context.Context, email string) error {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	otp := &domain.OTP{
		UserID:    user.ID,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := u.otpRepo.Create(ctx, otp); err != nil {
		return err
	}

	return u.emailSvc.SendOTP(ctx, email, code)
}

func (u *authUseCase) VerifyOTP(ctx context.Context, email, code string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	otp, err := u.otpRepo.GetLatestValidOTP(ctx, user.ID)
	if err != nil || otp == nil {
		return "", errors.New("invalid or expired OTP")
	}

	if otp.Code != code {
		return "", errors.New("incorrect OTP")
	}

	if err := u.otpRepo.MarkAsUsed(ctx, otp.ID); err != nil {
		return "", err
	}

	return u.generateJWT(user.ID)
}

func (u *authUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return u.generateJWT(user.ID)
}

func (u *authUseCase) generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(u.jwtSecret))
}
