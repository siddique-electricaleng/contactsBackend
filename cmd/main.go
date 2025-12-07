package main

import (
	"contacts/internal/env"
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	godotenv.Load(".env")

	cfg := config{
		addr: env.GetString("PORT", ":8080"),
		db: dbConfig{
			host:     env.GetString("DB_HOST", "localhost"),
			port:     env.GetString("DB_PORT", "5432"),
			user:     env.GetString("DB_USER", "postgres"),
			password: env.GetString("DB_PASSWORD", ""),
			name:     env.GetString("DB_NAME", "products"),
			sslmode:  env.GetString("DB_SSLMODE", "disable"),
		},
	}

	dsn := cfg.db.DSN()

	// structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("Connecting to database")

	// DB Pooling
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Connected to database")

	api := application{
		config: cfg,
		db:     conn,
	}
	/*
		from := "ecsesiddique.297@gmail.com"
		pass := "wvyk izjv qfcv ecic"
		to := "ecsesiddique.297@gmail.com"

		msg := []byte("Subject: Test Mail\r\n\r\nThis is a test email.")

		// Gmail SMTP server config
		// server := "smtp.gmail.com:587"
		host := "smtp.gmail.com"

		auth := smtp.PlainAuth("", from, pass, host)

		// TLS config is REQUIRED
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		connEmail, err := tls.Dial("tcp", "smtp.gmail.com:465", tlsconfig)
		if err != nil {
			log.Fatal("TLS Dial:", err)
		}

		c, err := smtp.NewClient(connEmail, host)
		if err != nil {
			log.Fatal("NewClient:", err)
		}

		if err = c.Auth(auth); err != nil {
			log.Fatal("Auth:", err)
		}

		if err = c.Mail(from); err != nil {
			log.Fatal("Mail:", err)
		}

		if err = c.Rcpt(to); err != nil {
			log.Fatal("Rcpt:", err)
		}

		w, err := c.Data()
		if err != nil {
			log.Fatal("Data:", err)
		}

		_, err = w.Write(msg)
		if err != nil {
			log.Fatal("Write:", err)
		}

		err = w.Close()
		if err != nil {
			log.Fatal("Close:", err)
		}

		c.Quit()
		log.Println("Email sent successfully!")
	*/
	// Running the API
	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
