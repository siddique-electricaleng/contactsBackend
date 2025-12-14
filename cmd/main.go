// @title Contacts Backend API
// @description API Endpoint Documentation for Contacts Module

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

package main

import (
	"contacts/docs"
	conf "contacts/internal/config"
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

	baseURLPath := env.GetString("API_BASE_URL", "/api/") + conf.APIVersion
	// Override API version in Swagger documentation - dynamic API version management
	docs.SwaggerInfo.Version = conf.APIVersion // swagger API version
	docs.SwaggerInfo.BasePath = baseURLPath    // swagger API Base Path

	cfg := config{
		addr: env.GetStringNoFallback("PORT"),
		db: dbConfig{
			host:     env.GetStringNoFallback("DB_HOST"),
			port:     env.GetStringNoFallback("DB_PORT"),
			user:     env.GetStringNoFallback("DB_USER"),
			password: env.GetStringNoFallback("DB_PASSWORD"),
			name:     env.GetStringNoFallback("DB_NAME"),
			sslmode:  env.GetStringNoFallback("DB_SSLMODE"),
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
