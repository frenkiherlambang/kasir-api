package repository

import "kasir-api/internal/domain"

// CategoryRepository defines the interface for category data access.
type CategoryRepository interface {
	GetAll() ([]domain.Category, error)
	GetByID(id int) (*domain.Category, error)
	Create(c domain.Category) (domain.Category, error)
	Update(id int, c domain.Category) (domain.Category, error)
	Delete(id int) error
}
