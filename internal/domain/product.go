package domain

// Product is the domain entity for a product.
type Product struct {
	ID       int      `json:"id"`
	Nama     string   `json:"nama"`
	Harga    int      `json:"harga"`
	Stok     int      `json:"stok"`
	Category Category `json:"category"`
}
