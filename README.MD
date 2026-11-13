# 🛒 Toko Online API

Proyek ini adalah sistem *E-Commerce* (Toko Online) *Fullstack* yang dibangun menggunakan **Golang** untuk Backend API dan **Vanilla JavaScript** untuk Frontend. Proyek ini didesain dengan arsitektur berlapis (*Layered Architecture*) untuk memastikan kode yang bersih, modular, dan mudah dipelihara.

## 🚀 Fitur Utama

- **Autentikasi & Otorisasi:** Menggunakan JWT (JSON Web Token) dengan Role-Based Access Control (Admin & Customer).
- **Manajemen Produk:** CRUD Produk lengkap dengan relasi Kategori.
- **Pagination, Search, & Filter:** Pencarian produk yang efisien dengan pagination standar.
- **Sistem Checkout (Transaksi):** Menggunakan *Database Transaction* dan *Row Locking* (`FOR UPDATE`) untuk mencegah *race condition* saat pembeli membeli stok terakhir secara bersamaan.
- **Dashboard Admin:** Laporan statistik penjualan, pendapatan, dan daftar transaksi.
- **Dockerized:** Siap dijalankan di mana saja dengan satu perintah berkat Docker & Docker Compose.

## 🛠️ Teknologi yang Digunakan

**Backend:**
- [Go (Golang)](https://golang.org/) (v1.26)
- [Gin Gonic](https://gin-gonic.com/) (Web Framework)
- [GORM](https://gorm.io/) (ORM untuk PostgreSQL)
- PostgreSQL (Database)
- JWT (Autentikasi)

**Frontend:**
- HTML5, CSS3, Vanilla JavaScript (ES6)

**DevOps & Tools:**
- Docker & Docker Compose
- Git

---

## 📂 Struktur Direktori (Backend)

Proyek ini memisahkan antara logika bisnis dan lalu lintas HTTP:
```text
backend/
├── config/      # Koneksi database & environment
├── dto/         # Data Transfer Objects (Validasi request/response)
├── handlers/    # Menerima request HTTP & mengembalikan response JSON
├── middleware/  # Penjaga gerbang API (Auth, Admin, CORS)
├── models/      # Skema tabel database (GORM)
├── routes/      # Daftar endpoint API
├── services/    # JANTUNG APLIKASI (Logika bisnis & transaksi)
└── utils/       # Fungsi bantuan (JWT, Hash Password)