# 📦 Inventory Management REST API

Aplikasi **Backend RESTful API** untuk sistem manajemen inventaris barang, gudang, dan transaksi penjualan.

Project ini dibangun menggunakan **Go (Golang)** dengan pendekatan **Clean Architecture** untuk memastikan kode yang rapi, mudah dites, dan *maintainable*.

---

## ✨ Fitur Utama

* 🔐 **Autentikasi & Otorisasi**

  * Login berbasis Session/Cookie
  * Role-Based Access Control (RBAC) untuk setiap endpoint

* 👤 **Manajemen User**

  * CRUD (Create, Read, Update, Delete) pengguna
  * Pengelolaan role user

* 📦 **Manajemen Inventaris**

  * Pengelolaan data barang dengan stok real-time
  * Manajemen kategori, rak penyimpanan, dan gudang
  * **Low Stock Alert** untuk mendeteksi stok menipis

* 💰 **Transaksi Penjualan**

  * Pencatatan barang keluar (penjualan)
  * Update stok otomatis saat transaksi dibuat, diubah, atau dihapus
  * Validasi stok untuk mencegah penjualan melebihi persediaan

* 📊 **Reporting**

  * Ringkasan total transaksi
  * Total barang terjual
  * Total pendapatan

* 📝 **Logging Terpusat**

  * Pencatatan request dan error menggunakan **Zap Logger**

---

## 🛠️ Tech Stack

* **Language**: Go (Golang)
* **Router**: [go-chi/chi v5](https://github.com/go-chi/chi)
* **Database Driver**: [jackc/pgx](https://github.com/jackc/pgx)
* **Configuration**: [spf13/viper](https://github.com/spf13/viper)
* **Logging**: [uber-go/zap](https://github.com/uber-go/zap)
* **Database**: PostgreSQL

---

## 📂 Struktur Project

Struktur folder mengikuti prinsip **Clean Architecture**:

```text
.
├── database/                 # Konfigurasi koneksi database (Postgres Pool)
├── dto/                      # Data Transfer Objects (Request & Response)
├── handler/                  # HTTP Handlers & validasi input
├── middleware/               # Auth, Logging, Permission
├── model/                    # Representasi tabel database
├── repository/               # Query & akses database
├── router/                   # Routing & grouping endpoint
├── service/                  # Business logic aplikasi
├── utils/                    # Helper (config, response, logger)
├── logs/                     # File log aplikasi
├── inventory_management.sql  # Dump SQL inisialisasi database
├── main.go                   # Entry point aplikasi
└── README.md                 # Dokumentasi project
```

---

## 🚀 Instalasi & Menjalankan Aplikasi

### 1️⃣ Prasyarat

Pastikan sudah terinstall:

* Go **v1.20+**
* PostgreSQL

---

### 2️⃣ Clone Repository

```bash
git clone https://github.com/username/project-app-inventory-restapi.git
cd project-app-inventory-restapi
```

---

### 3️⃣ Setup Database

Buat database PostgreSQL kosong, lalu import file SQL:

```bash
psql -U postgres -d nama_database_anda -f inventory_management.sql
```

> File `inventory_management.sql` berisi struktur tabel dan data awal (jika ada).

---

### 4️⃣ Konfigurasi Environment (.env)

Buat file `.env` di root project:

```env
APP_NAME=AppInventory
PORT=8080
DEBUG=true
LIMIT=10
PATH_LOGGING=./logs/

# Database Configuration
DATABASE_NAME=nama_database_anda
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=password_db_anda
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_MAX_CONN=10
```

---

### 5️⃣ Menjalankan Aplikasi

```bash
go mod tidy
go run main.go
```

Jika berhasil, server akan berjalan di:

```text
http://localhost:8080
```

---

## 📚 Dokumentasi API

**Base URL**

```text
http://localhost:8080/api/v1
```

---

### 🔐 Auth

| Method | Endpoint  | Deskripsi                            |
| ------ | --------- | ------------------------------------ |
| POST   | `/login`  | Login dan mendapatkan session cookie |
| POST   | `/logout` | Logout dari sistem                   |

---

### 👤 User

| Method | Endpoint     | Deskripsi        |
| ------ | ------------ | ---------------- |
| GET    | `/user`      | List semua user  |
| POST   | `/user`      | Tambah user baru |
| PUT    | `/user/{id}` | Update user      |
| DELETE | `/user/{id}` | Hapus user       |

---

### 📦 Inventory (Barang)

| Method | Endpoint          | Deskripsi            |
| ------ | ----------------- | -------------------- |
| GET    | `/inventory`      | List barang + lokasi |
| POST   | `/inventory`      | Tambah barang        |
| PUT    | `/inventory/{id}` | Update barang        |
| DELETE | `/inventory/{id}` | Hapus barang         |
| GET    | `/stock`          | Low stock alert      |

---

### 🏷️ Category & Rack

| Method | Endpoint         | Deskripsi       |
| ------ | ---------------- | --------------- |
| GET    | `/category`      | List kategori   |
| POST   | `/category`      | Tambah kategori |
| PUT    | `/category/{id}` | Update kategori |
| DELETE | `/category/{id}` | Hapus kategori  |
| GET    | `/rack`          | List rak        |
| POST   | `/rack`          | Tambah rak      |
| PUT    | `/rack/{id}`     | Update rak      |
| DELETE | `/rack/{id}`     | Hapus rak       |

---

### 🏭 Warehouse (Gudang)

| Method | Endpoint          | Deskripsi     |
| ------ | ----------------- | ------------- |
| GET    | `/warehouse`      | List gudang   |
| POST   | `/warehouse`      | Tambah gudang |
| PUT    | `/warehouse/{id}` | Update gudang |
| DELETE | `/warehouse/{id}` | Hapus gudang  |

---

### 💰 Transaction

| Method | Endpoint            | Deskripsi                         |
| ------ | ------------------- | --------------------------------- |
| GET    | `/transaction`      | Riwayat transaksi                 |
| POST   | `/transaction`      | Tambah transaksi (stok berkurang) |
| PUT    | `/transaction/{id}` | Update transaksi                  |
| DELETE | `/transaction/{id}` | Batalkan transaksi                |

---

### 📊 Report

| Method | Endpoint  | Deskripsi           |
| ------ | --------- | ------------------- |
| GET    | `/report` | Dashboard ringkasan |

---

## 🧪 Pengujian API (Postman)

1. Import file `inventory-items.postman_collection.json`
2. Login terlebih dahulu
3. Gunakan semua endpoint sesuai role

---

## 📌 Catatan

Project ini cocok sebagai:

* Portfolio backend Golang
* Referensi Clean Architecture
* Dasar pengembangan sistem inventory skala menengah
