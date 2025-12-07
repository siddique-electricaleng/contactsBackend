package users

import (
	"contacts/internal/env"
	"fmt"
	"net/smtp"
)

// generateEmailVerificationToken

var (
	smtpHost   = env.GetString("SMTP_HOST", "smtp.gmail.com")
	smtpPort   = env.GetInt("SMTP_PORT", 587)
	smtpUser   = env.GetString("SMTP_USERNAME", "")
	smtpPass   = env.GetString("SMTP_PASSWORD", "")
	appBaseURL = env.GetString("APP_BASE_URL", "http://localhost:8080")
)

func sendVerificationEmail(toEmail, token string) error {
	if smtpUser == "" || smtpPass == "" {
		return nil
	}

	link := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", appBaseURL, token)

	subject := "Verify your email"
	body := fmt.Sprintf(
		"Hi, \r\n\r\nPlease verify your email by clicking the link below:\r\n%s\r\n If you didn't sign up, ignore this.\r\n",
		link,
	)

	// RFC 5322 email format (Subject + blank line + body)
	msg := []byte("Subject: " + subject + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"\r\n" +
		body)

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	return smtp.SendMail(addr, auth, smtpUser, []string{toEmail}, msg)
}
