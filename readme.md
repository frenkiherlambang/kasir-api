# Kasir API

API REST sederhana untuk sistem kasir (Point of Sale) yang dibangun menggunakan Go. API ini menyediakan endpoint untuk mengelola kategori dan produk.

## 📋 Fitur

- ✅ Manajemen Kategori (CRUD)
- ✅ Manajemen Produk (CRUD)
- ✅ Relasi antara Produk dan Kategori
- ✅ Health Check Endpoint
- ✅ RESTful API Design
- ✅ PostgreSQL (pgx) untuk persistence—kompatibel dengan Supabase
- ✅ Konfigurasi via Viper (.env + environment variables)

## 🚀 Memulai

### Prasyarat

- Go 1.22.2 atau lebih baru
- PostgreSQL ([Supabase](https://supabase.com))
- Terminal/Command Line

### Instalasi

1. Clone repository ini:
```bash
git clone <repository-url>
cd kasir-api
```

2. Install dependencies:
```bash
go mod download
```

3. Buat file `.env` di root proyek:
```
DB_CONN=postgres://user:password@host:port/database
PORT=8080
```
- `DB_CONN` wajib (connection string ke PostgreSQL/Supabase).
- `PORT` opsional; default `8080`.

4. Jalankan migrasi schema sekali (mis. di Supabase SQL Editor):
- Salin dan jalankan isi file `migrations/001_schema.sql`.

5. Jalankan aplikasi:
```bash
go run .
```

Server berjalan di `http://localhost:8080` (atau sesuai `PORT` di `.env`).

## 📚 Dokumentasi API

### Base URL
```
http://localhost:8080
```

### Endpoints

#### Health Check

**GET** `/health`

Mengecek status API.

**Response:**
```json
{
  "status": "ok",
  "message": "API is running"
}
```

---

### Kategori (Categories)

#### 1. Mendapatkan Semua Kategori

**GET** `/api/categories`

**Response:**
```json
[
  {
    "id": 1,
    "nama": "Sneakers"
  },
  {
    "id": 2,
    "nama": "Running"
  }
]
```

#### 2. Mendapatkan Kategori Berdasarkan ID

**GET** `/api/categories/{id}`

**Response:**
```json
{
  "id": 1,
  "nama": "Sneakers"
}
```

**Error Response (404):**
```json
{
  "status": "error",
  "message": "Category not found"
}
```

#### 3. Membuat Kategori Baru

**POST** `/api/categories`

**Request Body:**
```json
{
  "nama": "Casual"
}
```

**Response (201):**
```json
{
  "id": 4,
  "nama": "Casual"
}
```

#### 4. Mengupdate Kategori

**PUT** `/api/categories/{id}`

**Request Body:**
```json
{
  "nama": "Sneakers Premium"
}
```

**Response:**
```json
{
  "id": 1,
  "nama": "Sneakers Premium"
}
```

#### 5. Menghapus Kategori

**DELETE** `/api/categories/{id}`

**Response:**
```json
{
  "status": "success",
  "message": "Category deleted successfully"
}
```

---

### Produk (Products)

#### 1. Mendapatkan Semua Produk

**GET** `/api/products`

**Response:**
```json
[
  {
    "id": 1,
    "nama": "Nike Air Max",
    "harga": 35000000,
    "stok": 10,
    "category": {
      "id": 1,
      "nama": "Sneakers"
    }
  }
]
```

#### 2. Mendapatkan Produk Berdasarkan ID

**GET** `/api/products/{id}`

**Response:**
```json
{
  "id": 1,
  "nama": "Nike Air Max",
  "harga": 35000000,
  "stok": 10,
  "category": {
    "id": 1,
    "nama": "Sneakers"
  }
}
```

#### 3. Membuat Produk Baru

**POST** `/api/products`

**Request Body:**
```json
{
  "nama": "New Balance 990",
  "harga": 400000,
  "stok": 10,
  "category": {
    "id": 1
  }
}
```

**Response (201):**
```json
{
  "id": 4,
  "nama": "New Balance 990",
  "harga": 400000,
  "stok": 10,
  "category": {
    "id": 1,
    "nama": "Sneakers"
  }
}
```

**Error Response (400) - Kategori tidak ditemukan:**
```json
{
  "status": "error",
  "message": "Category not found"
}
```

#### 4. Mengupdate Produk

**PUT** `/api/products/{id}`

**Request Body:**
```json
{
  "nama": "Nike Air Max Pro",
  "harga": 40000000,
  "stok": 15
}
```

**Catatan:** ID tidak dapat diubah. Field `category` dapat diikutsertakan (mis. `"category": {"id": 2}`) untuk mengubah kategori produk.

**Response:**
```json
{
  "id": 1,
  "nama": "Nike Air Max Pro",
  "harga": 40000000,
  "stok": 15,
  "category": {
    "id": 1,
    "nama": "Sneakers"
  }
}
```

#### 5. Menghapus Produk

**DELETE** `/api/products/{id}`

**Response:**
```json
{
  "status": "success",
  "message": "Product deleted successfully"
}
```

---

## 📝 Model Data

### Category
```go
type Category struct {
    ID   int    `json:"id"`
    Nama string `json:"nama"`
}
```

### Product
```go
type Product struct {
    ID       int      `json:"id"`
    Nama     string   `json:"nama"`
    Harga    int      `json:"harga"`
    Stok     int      `json:"stok"`
    Category Category `json:"category"`
}
```

## 🧪 Testing

Anda dapat menggunakan file HTTP yang tersedia untuk testing:

- `category.http` - Contoh request untuk endpoint kategori
- `product.http` - Contoh request untuk endpoint produk

Atau gunakan tools seperti:
- Postman
- cURL
- HTTPie
- REST Client extension di VS Code

### Contoh cURL

**Mendapatkan semua kategori:**
```bash
curl http://localhost:8080/api/categories
```

**Membuat kategori baru:**
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{"nama": "Sport"}'
```

**Membuat produk baru:**
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Adidas Ultraboost",
    "harga": 5000000,
    "stok": 20,
    "category": {"id": 1}
  }'
```

## 📦 Struktur Proyek

```
kasir-api/
├── main.go              # Entry point, wiring HTTP + repos
├── go.mod
├── .env                 # DB_CONN, PORT (jangan di-commit)
├── internal/
│   ├── config/          # Viper, Load(), Config struct
│   ├── domain/          # Category, Product
│   ├── handler/         # HTTP handlers
│   ├── repository/      # Interface + memory + PostgreSQL (pgx)
│   └── usecase/         # Business logic
├── migrations/
│   └── 001_schema.sql   # Tabel categories & products
├── category.http
├── product.http
└── readme.md
```

