package main

import (
	"log/slog"
	"os"
)

func main() {
	// ctx := context.Background()

	cfg := config{
		addr: ":8080",
	}

	// structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// DB Pooling
	// conn, err := pgx.Connect(ctx, cfg.db.dsn)
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close(ctx)

	// logger.Info("Connected to database", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		// db:     conn,
	}

	// Running the API
	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
