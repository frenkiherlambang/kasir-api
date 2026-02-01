package repository

import (
	"context"
	"errors"

	"kasir-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductPG is a PostgreSQL implementation of ProductRepository.
type ProductPG struct {
	pool *pgxpool.Pool
}

// NewProductPG creates a new PostgreSQL product repository.
func NewProductPG(pool *pgxpool.Pool) *ProductPG {
	return &ProductPG{pool: pool}
}

func scanProduct(scan func(...any) error) (domain.Product, error) {
	var p domain.Product
	var catID int
	var catNama string
	err := scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &catID, &catNama)
	if err != nil {
		return domain.Product{}, err
	}
	p.Category = domain.Category{ID: catID, Nama: catNama}
	return p, nil
}

// GetAll returns all products with their category.
func (r *ProductPG) GetAll() ([]domain.Product, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT p.id, p.nama, p.harga, p.stok, c.id, c.nama
		 FROM products p
		 JOIN categories c ON p.category_id = c.id
		 ORDER BY p.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Product
	for rows.Next() {
		p, err := scanProduct(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetByID returns a product by ID with its category, or ErrNotFound.
func (r *ProductPG) GetByID(id int) (*domain.Product, error) {
	row := r.pool.QueryRow(context.Background(),
		`SELECT p.id, p.nama, p.harga, p.stok, c.id, c.nama
		 FROM products p
		 JOIN categories c ON p.category_id = c.id
		 WHERE p.id = $1`, id)
	p, err := scanProduct(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

// Create inserts a product and returns it with the generated ID and category.
func (r *ProductPG) Create(p domain.Product) (domain.Product, error) {
	var id int
	err := r.pool.QueryRow(context.Background(),
		`INSERT INTO products (nama, harga, stok, category_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		p.Nama, p.Harga, p.Stok, p.Category.ID).Scan(&id)
	if err != nil {
		return domain.Product{}, err
	}
	created, err := r.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	return *created, nil
}

// Update updates a product by ID and returns the full product, or ErrNotFound.
func (r *ProductPG) Update(id int, p domain.Product) (domain.Product, error) {
	cmd, err := r.pool.Exec(context.Background(),
		`UPDATE products SET nama = $2, harga = $3, stok = $4, category_id = $5 WHERE id = $1`,
		id, p.Nama, p.Harga, p.Stok, p.Category.ID)
	if err != nil {
		return domain.Product{}, err
	}
	if cmd.RowsAffected() == 0 {
		return domain.Product{}, ErrNotFound
	}
	updated, err := r.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	return *updated, nil
}

// Delete removes a product by ID. Returns ErrNotFound if no rows affected.
func (r *ProductPG) Delete(id int) error {
	cmd, err := r.pool.Exec(context.Background(), "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
