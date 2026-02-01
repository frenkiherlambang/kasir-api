package repository

import (
	"context"
	"errors"

	"kasir-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CategoryPG is a PostgreSQL implementation of CategoryRepository.
type CategoryPG struct {
	pool *pgxpool.Pool
}

// NewCategoryPG creates a new PostgreSQL category repository.
func NewCategoryPG(pool *pgxpool.Pool) *CategoryPG {
	return &CategoryPG{pool: pool}
}

// GetAll returns all categories.
func (r *CategoryPG) GetAll() ([]domain.Category, error) {
	rows, err := r.pool.Query(context.Background(),
		"SELECT id, nama FROM categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Nama); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// GetByID returns a category by ID or ErrNotFound.
func (r *CategoryPG) GetByID(id int) (*domain.Category, error) {
	var c domain.Category
	err := r.pool.QueryRow(context.Background(),
		"SELECT id, nama FROM categories WHERE id = $1", id).Scan(&c.ID, &c.Nama)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

// Create inserts a category and returns it with the generated ID.
func (r *CategoryPG) Create(c domain.Category) (domain.Category, error) {
	var out domain.Category
	err := r.pool.QueryRow(context.Background(),
		"INSERT INTO categories (nama) VALUES ($1) RETURNING id, nama", c.Nama).Scan(&out.ID, &out.Nama)
	if err != nil {
		return domain.Category{}, err
	}
	return out, nil
}

// Update updates a category by ID and returns it, or ErrNotFound.
func (r *CategoryPG) Update(id int, c domain.Category) (domain.Category, error) {
	var out domain.Category
	err := r.pool.QueryRow(context.Background(),
		"UPDATE categories SET nama = $2 WHERE id = $1 RETURNING id, nama", id, c.Nama).Scan(&out.ID, &out.Nama)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Category{}, ErrNotFound
		}
		return domain.Category{}, err
	}
	return out, nil
}

// Delete removes a category by ID. Returns ErrNotFound if no rows affected.
func (r *CategoryPG) Delete(id int) error {
	cmd, err := r.pool.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
