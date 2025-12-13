package main

import (
	repo "contacts/internal/adapters/postgresql/sqlc"
	"contacts/internal/auth"
	conf "contacts/internal/config"
	"contacts/internal/contacts"
	appMiddleware "contacts/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// mount - create all the endpoints
func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(chiMiddleware.RequestID) // important for rate limiting - avoids DDOS
	r.Use(chiMiddleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)                 // recover from crashes
	r.Use(chiMiddleware.Timeout(60 * time.Second)) // 60 second timeout

	// Routes

	userService := auth.NewService(repo.New(app.db))
	userHandler := auth.NewHandler(userService)

	// r.Route("/contacts/api/"+conf.APIVersion, func(r chi.Router) {
	r.Route("/api/"+conf.APIVersion, func(r chi.Router) {

		// ---------------------------- PUBLIC ROUTES ------------------------

		// ------------------ Auth Routes --------------------

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
			r.Post("/refresh", userHandler.Refresh)
			r.Post("/logout", userHandler.Logout)
			r.Get("/verify-email", userHandler.VerifyEmail)
		})

		// ------------------Swagger-------------------

		r.Get("/docs/*", httpSwagger.WrapHandler)

		// ------------------Protected Routes----------
		contactsService := contacts.NewService(repo.New(app.db))
		contactsHandler := contacts.NewHandler(contactsService)

		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Auth([]byte(conf.JWT_SECRET)))
			// --------------ALL PROTECTED ROUTES BELOW HERE ----------------

			// POST contacts & upsert phone numbers/emails
			r.Post("/contacts", contactsHandler.CreateContacts)

			// GET all contacts for that user
			r.Get("/contacts", contactsHandler.ListContactsForUserWithDetails)
		})

	})
	return r
}

// run - start server
// Graceful shutdown - not done yet
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
