package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
	"github.com/masfuulaji/store/internal/app/repositories"
	"github.com/masfuulaji/store/internal/utils"
)

type CategoryHandler interface {
	CreateCategory(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategory(w http.ResponseWriter, r *http.Request)
	GetCategories(w http.ResponseWriter, r *http.Request)
}

type CategoryHandlerImpl struct {
	categoryRepository *repositories.CategoryRepositoryImpl
}

func NewCategoryHandler(db *sqlx.DB) *CategoryHandlerImpl {
	return &CategoryHandlerImpl{categoryRepository: repositories.NewCategoryRepository(db)}
}

func (f *CategoryHandlerImpl) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	id, err := f.categoryRepository.CreateCategory(category)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": fmt.Sprintf("Category created successfully %d", id)})
}

func (f *CategoryHandlerImpl) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	category := models.Category{}
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	err = f.categoryRepository.UpdateCategory(category, id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Category updated successfully"})
}

func (f *CategoryHandlerImpl) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := f.categoryRepository.DeleteCategory(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Category deleted successfully"})
}

func (f *CategoryHandlerImpl) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	category, err := f.categoryRepository.GetCategory(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(category)
}

func (f *CategoryHandlerImpl) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := f.categoryRepository.GetCategories()
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(categories)
}
