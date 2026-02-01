package domain

// Category is the domain entity for a product category.
type Category struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
}
