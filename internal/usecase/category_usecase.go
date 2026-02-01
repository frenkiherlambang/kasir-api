package usecase

import (
	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
)

// CategoryUsecase holds business logic for categories.
type CategoryUsecase struct {
	repo repository.CategoryRepository
}

// NewCategoryUsecase creates a new category use case.
func NewCategoryUsecase(repo repository.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

// GetAll returns all categories.
func (u *CategoryUsecase) GetAll() ([]domain.Category, error) {
	return u.repo.GetAll()
}

// GetByID returns a category by ID. Returns repository.ErrNotFound if not found.
func (u *CategoryUsecase) GetByID(id int) (*domain.Category, error) {
	return u.repo.GetByID(id)
}

// Create creates a new category and assigns the next ID.
func (u *CategoryUsecase) Create(c domain.Category) (domain.Category, error) {
	all, err := u.repo.GetAll()
	if err != nil {
		return domain.Category{}, err
	}
	maxID := 0
	for _, cat := range all {
		if cat.ID > maxID {
			maxID = cat.ID
		}
	}
	c.ID = maxID + 1
	return u.repo.Create(c)
}

// Update updates an existing category by ID.
func (u *CategoryUsecase) Update(id int, c domain.Category) (domain.Category, error) {
	return u.repo.Update(id, c)
}

// Delete deletes a category by ID.
func (u *CategoryUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
