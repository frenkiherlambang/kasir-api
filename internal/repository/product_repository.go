package repository

import "kasir-api/internal/domain"

// ProductRepository defines the interface for product data access.
type ProductRepository interface {
	GetAll(name string) ([]domain.Product, error)
	GetByID(id int) (*domain.Product, error)
	Create(p domain.Product) (domain.Product, error)
	Update(id int, p domain.Product) (domain.Product, error)
	Delete(id int) error
}
