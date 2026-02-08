package domain

// SummaryHariIni is the response for GET /api/report/hari-ini.
type SummaryHariIni struct {
	TotalRevenue   int            `json:"total_revenue"`
	TotalTransaksi int            `json:"total_transaksi"`
	ProdukTerlaris ProdukTerlaris `json:"produk_terlaris"`
}

// ProdukTerlaris holds the best-selling product for the day.
type ProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
