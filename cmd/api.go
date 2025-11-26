package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/wind1102/ecom/internal/adapters/postgresql/sqlc"
	"github.com/wind1102/ecom/internal/orders"
	"github.com/wind1102/ecom/internal/products"
)

type application struct {
	config config
	db     *pgx.Conn
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // important for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crash

	// Set timeout value on the request context ctx, that will signal through ctx.Done()
	// that the request has timed out and further processing should be stopped
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	productsService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productsService)

	ordersService := orders.NewService(repo.New(app.db))
	ordersHandler := orders.NewHandler(orders.NewService(ordersService))

	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", productHandler.FindProductById)
	r.Post("/orders", ordersHandler.PlaceOrder)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Server has started at addr %s", app.config.addr)
	return srv.ListenAndServe()
}

type config struct {
	addr string //
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
