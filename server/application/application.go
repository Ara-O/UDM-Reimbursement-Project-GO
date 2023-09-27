package application

import (
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/router"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type App struct {
	Addr   string
	Routes *chi.Mux
}

func New(port string) *App {
	app := &App{
		Addr:   port,
		Routes: loadRoutes(),
	}

	return app
}

func loadRoutes() *chi.Mux {
	r := chi.NewRouter()
	//Defining middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//Defining routes
	router.DefineRoutes(r)

	return r
}
