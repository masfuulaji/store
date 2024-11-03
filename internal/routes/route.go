package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/masfuulaji/store/internal/app/handlers"
	"github.com/masfuulaji/store/internal/database"
)

func SetupRoutes(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println(err)
	}

	categoryHandler := handlers.NewCategoryHandler(db.DB)
	r.Route("/category", func(r chi.Router) {
		r.Get("/", categoryHandler.GetCategories)
		r.Get("/{id}", categoryHandler.GetCategory)
		r.Post("/", categoryHandler.CreateCategory)
		r.Put("/{id}", categoryHandler.UpdateCategory)
		r.Delete("/{id}", categoryHandler.DeleteCategory)
	})

	productHandler := handlers.NewProductHandler(db.DB)
	r.Route("/product", func(r chi.Router) {
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/category/{id}", productHandler.GetProductsByCategory)
		r.Post("/", productHandler.CreateProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	userHandler := handlers.NewUserHandler(db.DB)
	r.Route("/user", func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
		r.Get("/{id}", userHandler.GetUser)
		r.Post("/", userHandler.CreateUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	loginHandler := handlers.NewLoginHandler(db.DB)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/", loginHandler.Login)
		r.Get("/logout", loginHandler.Logout)
	})

	cartHandler := handlers.NewCartHandler(db.DB)
	r.Route("/cart", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Auth!"))
		})
		r.Post("/add", cartHandler.AddCartItem)
		r.Delete("/delete/{id}", cartHandler.DeleteCart)
		r.Get("/{id}", cartHandler.ReadCart)
	})

	orderHandler := handlers.NewOrderHandler(db.DB)
	r.Route("/order", func(r chi.Router) {
		r.Post("/add", orderHandler.CreateOrder)
	})
}
