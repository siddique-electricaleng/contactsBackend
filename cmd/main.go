// @title Contacts Backend API
// @description API Endpoint Documentation for Contacts Module

package main

import (
	"contacts/docs"
	conf "contacts/internal/config"
	"contacts/internal/env"
	"context"
	"log/slog"
	"os"

	_ "contacts/docs" // <- module path + /docs

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	godotenv.Load(".env")

	baseURLPath := env.GetString("API_BASE_URL", "/contacts/api/") + conf.APIVersion
	// Override API version in Swagger documentation - dynamic API version management
	docs.SwaggerInfo.Version = conf.APIVersion // swagger API version
	docs.SwaggerInfo.BasePath = baseURLPath    // swagger API Base Path

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

	// Running the API
	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
