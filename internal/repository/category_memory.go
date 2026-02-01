package repository

import (
	"kasir-api/internal/domain"
	"sync"
)

// CategoryMemoryRepo is an in-memory implementation of CategoryRepository.
type CategoryMemoryRepo struct {
	mu   sync.RWMutex
	data []domain.Category
}

// NewCategoryMemoryRepo creates a new in-memory category repository with optional initial data.
func NewCategoryMemoryRepo(initial []domain.Category) *CategoryMemoryRepo {
	data := make([]domain.Category, len(initial))
	copy(data, initial)
	return &CategoryMemoryRepo{data: data}
}

// GetAll returns all categories.
func (r *CategoryMemoryRepo) GetAll() ([]domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Category, len(r.data))
	copy(out, r.data)
	return out, nil
}

// GetByID returns a category by ID or ErrNotFound.
func (r *CategoryMemoryRepo) GetByID(id int) (*domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.data {
		if r.data[i].ID == id {
			c := r.data[i]
			return &c, nil
		}
	}
	return nil, ErrNotFound
}

// Create adds a new category. If c.ID is 0, assigns the next ID.
func (r *CategoryMemoryRepo) Create(c domain.Category) (domain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		maxID := 0
		for _, cat := range r.data {
			if cat.ID > maxID {
				maxID = cat.ID
			}
		}
		c.ID = maxID + 1
	}
	r.data = append(r.data, c)
	return c, nil
}

// Update updates an existing category by ID.
func (r *CategoryMemoryRepo) Update(id int, c domain.Category) (domain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data {
		if r.data[i].ID == id {
			c.ID = id
			r.data[i] = c
			return c, nil
		}
	}
	return domain.Category{}, ErrNotFound
}

// Delete removes a category by ID.
func (r *CategoryMemoryRepo) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data {
		if r.data[i].ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
