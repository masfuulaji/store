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
}
