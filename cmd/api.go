package main

import (
	repo "contacts/internal/adapters/postgresql/sqlc"
	"contacts/internal/users"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

// mount - create all the endpoints
func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting - avoids DDOS
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)                 // recover from crashes
	r.Use(middleware.Timeout(60 * time.Second)) // 60 second timeout

	// Routes

	userService := users.NewService(repo.New(app.db))
	userHandler := users.NewHandler(userService)

	r.Route("/api/v1", func(r chi.Router) {

		// ------------------ Auth Routes --------------------
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
			r.Post("/refresh", userHandler.Refresh)
			r.Post("/logout", userHandler.Logout)
			r.Get("/verify", userHandler.VerifyEmail)
		})
	})

	return r
}

// run - start server and do graceful shutdown
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at address: %s\n", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	host     string
	port     string
	user     string
	password string
	name     string
	sslmode  string
}

func (c dbConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.host, c.port, c.name, c.user, c.password, c.sslmode,
	)
}
