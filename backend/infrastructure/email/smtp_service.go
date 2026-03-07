package email

import (
	"context"
	"fmt"
	"net/smtp"
	"taskmanager/usecase"
)

type smtpService struct {
	host     string
	port     string
	username string
	password string
}

func NewSMTPService(host, port, username, password string) usecase.EmailService {
	return &smtpService{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *smtpService) SendOTP(ctx context.Context, toEmail, code string) error {
	from := s.username
	to := []string{toEmail}

	subject := "Subject: Your OTP for Task Manager\n"
	body := fmt.Sprintf("Your OTP is: %s. It is valid for 10 minutes.\n", code)
	msg := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, from, to, msg)
}
