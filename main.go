package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// =======================
// CATEGORY MODEL
// =======================
type Category struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
}

// =======================
// PRODUCT MODEL
// =======================
type Product struct {
	ID       int      `json:"id"`
	Nama     string   `json:"nama"`
	Harga    int      `json:"harga"`
	Stok     int      `json:"stok"`
	Category Category `json:"category"`
}

// =======================
// CATEGORY DATA
// =======================
var categories = []Category{
	{ID: 1, Nama: "Sneakers"},
	{ID: 2, Nama: "Running"},
	{ID: 3, Nama: "Casual"},
}

// =======================
// PRODUCT DATA
// =======================
var products = []Product{
	{ID: 1, Nama: "Nike Air Max", Harga: 35000000, Stok: 10, Category: categories[0]},
	{ID: 2, Nama: "Adidas Superstar", Harga: 3000000, Stok: 20, Category: categories[1]},
	{ID: 3, Nama: "Converse Chuck Taylor", Harga: 4000000, Stok: 15, Category: categories[2]},
}

// =======================
// CATEGORY HANDLERS
// =======================

// get category by id
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid category ID",
		})
	}
	for _, category := range categories {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(category)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Category not found",
	})
}

func getAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}
	category.ID = getNextCategoryID()
	categories = append(categories, category)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid category ID",
		})
		return
	}
	for _, category := range categories {
		if category.ID == id {
			var updatedCategory Category
			err := json.NewDecoder(r.Body).Decode(&updatedCategory)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"status":  "error",
					"message": "Invalid request body",
				})
				return
			}
			for i, currentCategory := range categories {
				if currentCategory.ID == id {
					updatedCategory.ID = currentCategory.ID
					categories[i] = updatedCategory
					break
				}
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Category not found",
	})
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid category ID",
		})
	}
	for i, category := range categories {
		if category.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "Category deleted successfully",
			})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Category not found",
	})
}

func getNextCategoryID() int {
	maxID := 0
	for _, category := range categories {
		if category.ID > maxID {
			maxID = category.ID
		}
	}
	return maxID + 1
}

// =======================
// PRODUCT HANDLERS
// =======================

// get product by id
func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}
	for _, product := range products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Product not found",
	})
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}
	for _, product := range products {
		if product.ID == id {
			var updatedProduct Product
			err := json.NewDecoder(r.Body).Decode(&updatedProduct)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"status":  "error",
					"message": "Invalid request body",
				})
			}
			for i, currentProduct := range products {
				if currentProduct.ID == id {
					//prevent update ID
					updatedProduct.ID = currentProduct.ID
					updatedProduct.Category = currentProduct.Category
					products[i] = updatedProduct
					break
				}
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Product not found",
	})
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}
	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "Product deleted successfully",
			})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Product not found",
	})
}

func getNextProductID() int {
	maxID := 0
	for _, product := range products {
		if product.ID > maxID {
			maxID = product.ID
		}
	}
	return maxID + 1
}

func main() {
	// =======================
	// CATEGORY ROUTES
	// =======================
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategoryByID(w, r)
		case "PUT":
			updateCategory(w, r)
		case "DELETE":
			deleteCategory(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Method not allowed",
			})
		}
	})

	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllCategories(w, r)
		case "POST":
			createCategory(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Method not allowed",
			})
		}
	})

	// =======================
	// PRODUCT ROUTES
	// =======================
	// Get product by ID
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getProductByID(w, r)
		case "PUT":
			updateProduct(w, r)
		case "DELETE":
			deleteProduct(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Method not allowed",
			})
		}
	})

	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(products)
		case "POST":
			var product Product
			err := json.NewDecoder(r.Body).Decode(&product)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"status":  "error",
					"message": "Invalid request body",
				})
			}

			// find category by id
			for _, category := range categories {
				if category.ID == product.Category.ID {
					product.Category = category
					break
				}
			}
			if product.Category.Nama == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"status":  "error",
					"message": "Category not found",
				})
				return
			}

			product.ID = getNextProductID()
			products = append(products, product)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(product)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Method not allowed",
			})
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "API is running",
		})
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
