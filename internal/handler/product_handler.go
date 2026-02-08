package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
	"kasir-api/internal/usecase"
)

// ProductHandler handles HTTP for products.
type ProductHandler struct {
	uc *usecase.ProductUsecase
}

// NewProductHandler creates a new product HTTP handler.
func NewProductHandler(uc *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

// GetByID handles GET /api/products/:id
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath(r.URL.Path, "/api/products/")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	prod, err := h.uc.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, prod)
}

// GetAll handles GET /api/products. Optional query: name=Nike to filter by product name.
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	prods, err := h.uc.GetAll(name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, prods)
}

// Create handles POST /api/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	created, err := h.uc.Create(p)
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			writeError(w, http.StatusBadRequest, "Category not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// Update handles PUT /api/products/:id
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath(r.URL.Path, "/api/products/")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updated, err := h.uc.Update(id, p)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// Delete handles DELETE /api/products/:id
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath(r.URL.Path, "/api/products/")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	err := h.uc.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}
