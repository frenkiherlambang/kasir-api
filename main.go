package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"kasir-api/internal/config"
	"kasir-api/internal/handler"
	"kasir-api/internal/repository"
	"kasir-api/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DBConn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	categoryRepo := repository.NewCategoryPG(pool)
	productRepo := repository.NewProductPG(pool)

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

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
