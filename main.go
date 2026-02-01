package main

import (
	"encoding/json"
	"log"
	"net/http"

	"kasir-api/internal/domain"
	"kasir-api/internal/handler"
	"kasir-api/internal/repository"
	"kasir-api/internal/usecase"
)

func main() {
	// Initial data
	initialCategories := []domain.Category{
		{ID: 1, Nama: "Sneakers"},
		{ID: 2, Nama: "Running"},
		{ID: 3, Nama: "Casual"},
	}
	initialProducts := []domain.Product{
		{ID: 1, Nama: "Nike Air Max", Harga: 35000000, Stok: 10, Category: initialCategories[0]},
		{ID: 2, Nama: "Adidas Superstar", Harga: 3000000, Stok: 20, Category: initialCategories[1]},
		{ID: 3, Nama: "Converse Chuck Taylor", Harga: 4000000, Stok: 15, Category: initialCategories[2]},
	}

	// Repositories
	categoryRepo := repository.NewCategoryMemoryRepo(initialCategories)
	productRepo := repository.NewProductMemoryRepo(initialProducts)

	// Use cases
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)
	productUC := usecase.NewProductUsecase(productRepo, categoryRepo)

	// Handlers
	categoryHandler := handler.NewCategoryHandler(categoryUC)
	productHandler := handler.NewProductHandler(productUC)

	// Method not allowed response
	methodNotAllowed := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Method not allowed",
		})
	}

	// Category routes
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetByID(w, r)
		case http.MethodPut:
			categoryHandler.Update(w, r)
		case http.MethodDelete:
			categoryHandler.Delete(w, r)
		default:
			methodNotAllowed(w)
		}
	})
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetAll(w, r)
		case http.MethodPost:
			categoryHandler.Create(w, r)
		default:
			methodNotAllowed(w)
		}
	})

	// Product routes
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetByID(w, r)
		case http.MethodPut:
			productHandler.Update(w, r)
		case http.MethodDelete:
			productHandler.Delete(w, r)
		default:
			methodNotAllowed(w)
		}
	})
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetAll(w, r)
		case http.MethodPost:
			productHandler.Create(w, r)
		default:
			methodNotAllowed(w)
		}
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "API is running",
		})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
