# 📦 Inventory Management REST API (Enterprise Edition)

Aplikasi **Backend RESTful API** untuk sistem manajemen inventaris barang, rak, gudang, dan transaksi penjualan dengan standar **Senior Backend Engineering**.

Project ini dibangun menggunakan **Go (Golang)** dengan pendekatan **Clean Architecture** (Handler -> Service -> Repository), dilengkapi **JWT Authentication**, **Redis Caching**, **ACID Transactions**, **Soft Deletes**, **Panic Recovery**, dan **Graceful Shutdown**.

---

## ✨ Fitur & Peningkatan Utama (Senior Backend Standards)

* 🔐 **Autentikasi & Otorisasi Modern**
  * **JWT (JSON Web Token)**: Bearer token format (`Authorization: Bearer <token>`) & Cookie fallback untuk web browser.
  * Stateless token verification dengan context injection (`r.Context()`), bebas dari session DB lock / nil pointer panics.
  * **Role-Based Access Control (RBAC)** dengan caching izin via Redis.
  * Password hashing aman menggunakan **Bcrypt** (tanpa logging plaintext password).

* ⚡ **Redis Caching Layer (High Performance)**
  * Caching otomatis untuk dashboard reporting & permissions check.
  * Cache invalidation otomatis saat terjadi transaksi penjualan atau perubahan inventaris.
  * Graceful fallback: aplikasi tetap berjalan normal jika Redis tidak aktif.

* 💰 **Transaksi Penjualan ACID & Data Integrity**
  * Pessimistic Row Locking (`SELECT ... FOR UPDATE`) untuk mencegah *race conditions* saat banyak user memesan barang bersamaan.
  * **Stock Restoration**: Membatalkan/menghapus transaksi otomatis mengembalikan stok barang ke inventaris.
  * Validasi kuantitas ketat (`quantity > 0`).

* 🛡️ **Reliability & Observability**
  * **Graceful Shutdown**: Menunggu request yang sedang berjalan selesai (hingga 10 detik) saat aplikasi menerima `SIGINT`/`SIGTERM`.
  * **Panic Recovery**: `chiMiddleware.Recoverer` melindungi proses server agar tidak *crash* jika terjadi panic tak terduga.
  * **CORS Support**: Siap dihubungkan ke frontend modern (React, Next.js, Vue, mobile).
  * **Structured Logging**: Logging zap lengkap mencatat status code HTTP, latency, dan client IP.

* 📦 **Manajemen Inventaris & Soft Deletes**
  * Konsisten menggunakan **Soft Delete** (`deleted_at = NOW()`).
  * Perhitungan pagination akurat (`total_items`, `total_pages`, `current_page`, `limit`).
  * Low Stock Alert (`/api/v1/stock`) untuk stok <= 5.

---

## 🛠️ Tech Stack

* **Language**: Go (Golang 1.20+)
* **Web Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
* **Auth**: [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)
* **Cache**: [redis/go-redis v9](https://github.com/redis/go-redis)
* **Database Driver**: [jackc/pgx v5](https://github.com/jackc/pgx)
* **Validation**: [go-playground/validator v10](https://github.com/go-playground/validator)
* **Configuration**: [spf13/viper](https://github.com/spf13/viper)
* **Logging**: [uber-go/zap](https://github.com/uber-go/zap)
* **Database**: PostgreSQL

---

## 📂 Struktur Project

```text
.
├── database/                 # PostgreSQL pool (pgxpool) & Redis client
├── dto/                      # Data Transfer Objects (Request, Response, Pagination)
├── handler/                  # HTTP Handlers (Gin context, JSON binding & validation)
├── middleware/               # Gin Middlewares (JWT Auth, Logging, RBAC Permission, CORS)
├── model/                    # Representasi entitas domain & database
├── repository/               # SQL queries, ACID transactions, & Redis caching
├── router/                   # Gin Router Engine, grouping, & endpoint mapping
├── service/                  # Business logic aplikasi & cache invalidation
├── utils/                    # JWT helper, bcrypt, response helpers, validator
├── logs/                     # File log harian (Lumberjack rotation)
├── backup_inventory_management.sql # Schema & dump data inisialisasi database
├── .env.example              # Template variabel lingkungan
├── main.go                   # Entry point aplikasi & graceful shutdown
└── README.md                 # Dokumentasi project
```

---

## 🚀 Panduan Menjalankan Aplikasi

### 1️⃣ Prasyarat
* Go 1.20+ terinstall
* PostgreSQL 14+ terinstall & aktif
* Redis (Opsional, untuk caching)

### 2️⃣ Setup Database
Import database schema dan data awal:
```bash
psql -U postgres -d your_database_name -f backup_inventory_management.sql
```

### 3️⃣ Konfigurasi Environment (`.env`)
Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
Sesuaikan kredensial database PostgreSQL dan JWT Secret:
```env
APP_NAME=AppInventory
PORT=8080
DEBUG=true
LIMIT=10
PATH_LOGGING=./logs/

# Database PostgreSQL
DATABASE_NAME=inventory_management
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=password_db_anda
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_MAX_CONN=10

# JWT Auth
JWT_SECRET=super_rahasia_ganti_di_production
JWT_EXPIRY_MINUTES=1440

# Redis Cache (Set false jika belum ingin memakai Redis)
REDIS_ENABLED=false
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 4️⃣ Jalankan Server
```bash
go mod tidy
go run main.go
```
Server akan aktif di `http://localhost:8080`.

---

## 📚 Endpoint API & Autentikasi

Header untuk endpoint yang diproteksi:
```http
Authorization: Bearer <access_token>
```
*(Atau otomatis menggunakan cookie `session` jika mengakses dari browser)*.

### 🔐 Autentikasi
| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` | `/api/v1/login` | Login (mengembalikan token JWT & user info) |
| `POST` | `/api/v1/logout` | Logout |

### 👤 Pengguna (User)
| Method | Endpoint | Permission | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/user` | `user:view` | List pengguna dengan pagination |
| `GET` | `/api/v1/user/{id}` | `user:view` | Detail pengguna berdasarkan ID |
| `POST` | `/api/v1/user` | `user:manage` | Tambah pengguna baru |
| `PATCH/PUT`| `/api/v1/user/{id}` | `user:manage` | Perbarui data pengguna |
| `DELETE` | `/api/v1/user/{id}` | `user:manage` | Soft delete pengguna |

### 📦 Inventaris (Inventory)
| Method | Endpoint | Permission | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/inventory` | `inventory:view` | List inventaris barang lengkap |
| `GET` | `/api/v1/inventory/{id}` | `inventory:view` | Detail barang |
| `POST` | `/api/v1/inventory` | `inventory:create` | Tambah inventaris baru |
| `PATCH/PUT`| `/api/v1/inventory/{id}` | `inventory:edit` | Perbarui harga/stok/kategori |
| `DELETE` | `/api/v1/inventory/{id}` | `inventory:delete` | Soft delete inventaris |
| `GET` | `/api/v1/stock` | `stock:view` | Low stock alert (stok <= 5) |

### 💰 Transaksi Penjualan (Transactions)
| Method | Endpoint | Permission | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/transaction` | `transaction:manage` | Riwayat transaksi penjualan |
| `GET` | `/api/v1/transaction/{id}`| `transaction:manage`| Detail transaksi |
| `POST` | `/api/v1/transaction` | `transaction:manage` | Buat transaksi (stok berkurang otomatis) |
| `PATCH/PUT`| `/api/v1/transaction/{id}`| `transaction:manage`| Edit transaksi (penyesuaian delta stok) |
| `DELETE` | `/api/v1/transaction/{id}`| `transaction:manage`| Batalkan transaksi (stok kembali ke gudang) |

### 📊 Laporan & Monitoring
| Method | Endpoint | Permission | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/report` | `report:view` | Ringkasan omset & barang terjual (Redis cached) |
| `GET` | `/health` | Public | Status kesehatan server |
