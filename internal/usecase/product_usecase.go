package usecase

import (
	"errors"
	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
)

// ErrCategoryNotFound is returned when a product references a non-existent category.
var ErrCategoryNotFound = errors.New("category not found")

// ProductUsecase holds business logic for products.
type ProductUsecase struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductUsecase creates a new product use case.
func NewProductUsecase(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// GetAll returns all products.
func (u *ProductUsecase) GetAll() ([]domain.Product, error) {
	return u.productRepo.GetAll()
}

// GetByID returns a product by ID. Returns repository.ErrNotFound if not found.
func (u *ProductUsecase) GetByID(id int) (*domain.Product, error) {
	return u.productRepo.GetByID(id)
}

// Create creates a new product. Resolves category by ID; returns error if category not found.
// The repository assigns and returns the new product ID.
func (u *ProductUsecase) Create(p domain.Product) (domain.Product, error) {
	cat, err := u.categoryRepo.GetByID(p.Category.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Product{}, ErrCategoryNotFound
		}
		return domain.Product{}, err
	}
	p.Category = *cat
	return u.productRepo.Create(p)
}

// Update updates an existing product by ID. ID and Category are preserved.
func (u *ProductUsecase) Update(id int, p domain.Product) (domain.Product, error) {
	return u.productRepo.Update(id, p)
}

// Delete deletes a product by ID.
func (u *ProductUsecase) Delete(id int) error {
	return u.productRepo.Delete(id)
}
