1. Auth (Autentikasi)

Method	        Endpoint	            Akses	        Deskripsi
POST	        /auth/register	        Publik	        Mendaftarkan akun customer baru.
POST	        /auth/login	            Publik	        Login untuk mendapatkan JWT Token.
GET	            /auth/me	            Login	        Mengambil profil user yang sedang login.

2. Categories (Kategori Produk)

Method	        Endpoint	            Akses	        Deskripsi
GET	            /categories	            Publik	        Melihat daftar semua kategori.
POST	        /categories	            Admin	        Menambahkan kategori baru.
PUT	            /categories/:id	        Admin	        Mengubah nama kategori.
DELETE	        /categories/:id	        Admin	        Menghapus kategori (gagal jika masih dipakai produk).

3. Products (Produk)

Method	        Endpoint	            Akses	        Deskripsi
GET	            /products	            Publik	        Mendapatkan daftar produk (mendukung Pagination, Search, Filter).
GET	            /products/:id	        Publik	        Mendapatkan detail satu produk.
POST	        /products	            Admin	        Menambahkan produk baru.
PUT	            /products/:id	        Admin	        Mengubah data produk.
DELETE	        /products/:id	        Admin	        Menghapus produk.

4. Checkout (Transaksi)
Method	        Endpoint	            Akses	        Deskripsi
POST	        /checkout	            Login	        Memproses pembelian. Memotong stok produk secara real-time dan membuat invoice.

5. Sales & Stats (Penjualan & Statistik)
Method	        Endpoint	            Akses	        Deskripsi
GET	            /sales	                Admin	        Melihat riwayat transaksi terbaru (maks 500 data).
GET	            /sales/summary	        Admin	        Ringkasan total produk terjual dan pendapatan per produk.
GET	            stats	                Admin	        Statistik umum (total produk, total user, total terjual, total pendapatan) untuk dashboard.

6. Users (Manajemen Pengguna)
Method	        Endpoint	            Akses	        Deskripsi
GET	            /users	                Admin	        Mendapatkan semua daftar user (mendukung filter ?role= dan ?q=).
GET	            /users/:id	            Admin / Self	Melihat detail user. Admin bisa melihat semua, user hanya bisa melihat dirinya sendiri.
PUT	            /users/:id	            Admin / Self	Mengubah profil user. Hanya admin yang bisa mengubah role.
PUT	            /users/:id/password	    Self	        Mengganti kata sandi (wajib memasukkan password lama).
DELETE	        /users/:id	            Admin	        Menghapus user (admin tidak bisa menghapus dirinya sendiri / admin terakhir).

7. Query Parameters untuk GET /products

page (int): Halaman yang ingin ditampilkan (default: 1).
limit (int): Jumlah data per halaman (default: 10, max: 100).
q (string): Pencarian produk berdasarkan nama.
category (string): Filter berdasarkan nama kategori.
in_stock (boolean): Jika true, hanya menampilkan produk yang stoknya > 0.
Contoh Request:
GET /api/products?page=1&limit=5&q=bajuA&category=KategoriA&in_stock=true