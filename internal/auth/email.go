package auth

import (
	"contacts/internal/env"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
)

// generateEmailVerificationToken

var (
	appBaseURL            = os.Getenv("APP_BASE_URL")
	emaiLVerificationPath = env.GetString("EMAIL_API_ENDPOINT", "/api/v1/auth/verify-email")
)

func SendEmailVerification(to, subject, body string) error {
	host := os.Getenv("SMTP_HOST") // smtp.gmail.com
	port := os.Getenv("SMTP_PORT") // 587
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// addr := fmt.Sprintf("%s:%s", host, port)
	addr := net.JoinHostPort(host, port)

	// 1. Connect to SMTP server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	// 2. Upgrade to TLS — STARTTLS
	tlsConfig := &tls.Config{
		ServerName: host,
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(tlsConfig); err != nil {
			return err
		}
	}

	// 3. Authentication
	auth := smtp.PlainAuth("", username, password, host)
	if err = client.Auth(auth); err != nil {
		return err
	}

	// 4. Set the sender and receiver
	if err = client.Mail(username); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	// 5. Write email body
	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	_, err = wc.Write([]byte(message))
	if err != nil {
		return err
	}

	return client.Quit()
}
