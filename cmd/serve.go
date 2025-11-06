package cmd

import (
	"fmt"
	"goRoute/handlers"
	"goRoute/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Serve() {
	r := chi.NewRouter()

	//r.Use(middleware.Logger)
	
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.Useless)

	r.Post("/api/auth/signup", handlers.Signup)
	r.Post("/api/auth/login", handlers.Login)

	r.Get("/api/products", handlers.GetProducts)
	r.Get("/api/products/{id}", handlers.GetProductByID)
	r.Post("/api/CreateUser",handlers.CreateUser)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/api/auth/profile", handlers.GetProfile)

		r.Post("/api/products", handlers.CreateProduct)
		r.Put("/api/products/{id}", handlers.UpdateProduct)
		r.Delete("/api/products/{id}", handlers.DeleteProduct)
	})

	fmt.Println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", r)
}
