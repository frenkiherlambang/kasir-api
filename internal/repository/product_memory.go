package repository

import (
	"kasir-api/internal/domain"
	"strings"
	"sync"
)

// ProductMemoryRepo is an in-memory implementation of ProductRepository.
type ProductMemoryRepo struct {
	mu   sync.RWMutex
	data []domain.Product
}

// NewProductMemoryRepo creates a new in-memory product repository with optional initial data.
func NewProductMemoryRepo(initial []domain.Product) *ProductMemoryRepo {
	data := make([]domain.Product, len(initial))
	copy(data, initial)
	return &ProductMemoryRepo{data: data}
}

// GetAll returns all products. If name is non-empty, filters by product name (substring, case-insensitive).
func (r *ProductMemoryRepo) GetAll(name string) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if name == "" {
		out := make([]domain.Product, len(r.data))
		copy(out, r.data)
		return out, nil
	}
	lower := strings.ToLower(name)
	var out []domain.Product
	for _, p := range r.data {
		if strings.Contains(strings.ToLower(p.Nama), lower) {
			out = append(out, p)
		}
	}
	return out, nil
}

// GetByID returns a product by ID or ErrNotFound.
func (r *ProductMemoryRepo) GetByID(id int) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.data {
		if r.data[i].ID == id {
			p := r.data[i]
			return &p, nil
		}
	}
	return nil, ErrNotFound
}

// Create adds a new product. Category must be resolved by caller. If p.ID is 0, assigns the next ID.
func (r *ProductMemoryRepo) Create(p domain.Product) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.ID == 0 {
		maxID := 0
		for _, prod := range r.data {
			if prod.ID > maxID {
				maxID = prod.ID
			}
		}
		p.ID = maxID + 1
	}
	r.data = append(r.data, p)
	return p, nil
}

// Update updates an existing product by ID. ID and Category are preserved.
func (r *ProductMemoryRepo) Update(id int, p domain.Product) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data {
		if r.data[i].ID == id {
			p.ID = id
			p.Category = r.data[i].Category
			r.data[i] = p
			return p, nil
		}
	}
	return domain.Product{}, ErrNotFound
}

// Delete removes a product by ID.
func (r *ProductMemoryRepo) Delete(id int) error {
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
