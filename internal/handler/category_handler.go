package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
	"kasir-api/internal/usecase"
)

// CategoryHandler handles HTTP for categories.
type CategoryHandler struct {
	uc *usecase.CategoryUsecase
}

// NewCategoryHandler creates a new category HTTP handler.
func NewCategoryHandler(uc *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

// GetByID handles GET /api/categories/:id
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath(r.URL.Path, "/api/categories/")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	cat, err := h.uc.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Category not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

// GetAll handles GET /api/categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	cats, err := h.uc.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cats)
}

// Create handles POST /api/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c domain.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	created, err := h.uc.Create(c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// Update handles PUT /api/categories/:id
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath(r.URL.Path, "/api/categories/")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	var c domain.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updated, err := h.uc.Update(id, c)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Category not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// Delete handles DELETE /api/categories/:id
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath(r.URL.Path, "/api/categories/")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	err := h.uc.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Category not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Category deleted successfully",
	})
}
