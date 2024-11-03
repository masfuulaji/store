package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
	"github.com/masfuulaji/store/internal/app/repositories"
	"github.com/masfuulaji/store/internal/utils"
)

type ProductHandler interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	GetProducts(w http.ResponseWriter, r *http.Request)
}

type ProductHandlerImpl struct {
	productRepository *repositories.ProductRepositoryImpl
}

func NewProductHandler(db *sqlx.DB) *ProductHandlerImpl {
	return &ProductHandlerImpl{productRepository: repositories.NewProductRepository(db)}
}

func (f *ProductHandlerImpl) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := models.Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	err = f.productRepository.CreateProduct(product)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Product created successfully"})
}

func (f *ProductHandlerImpl) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product := models.Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	err = f.productRepository.UpdateProduct(product, id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Product updated successfully"})
}

func (f *ProductHandlerImpl) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := f.productRepository.DeleteProduct(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Product deleted successfully"})
}

func (f *ProductHandlerImpl) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := f.productRepository.GetProduct(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (f *ProductHandlerImpl) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := f.productRepository.GetProducts()
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(products)
}

func (f *ProductHandlerImpl) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	products, err := f.productRepository.GetProductsByCategory(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(products)
}
