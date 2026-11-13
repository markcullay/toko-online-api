API ini menggunakan standar HTTP Response Code:

200 OK : Request berhasil.
201 Created : Data berhasil dibuat/ditambahkan.
400 Bad Request : Validasi gagal, input tidak sesuai, atau logika bisnis ditolak (misal: stok habis).
401 Unauthorized : Belum login atau token kadaluarsa.
403 Forbidden : Sudah login, tapi tidak punya hak akses (misal: customer mencoba akses endpoint admin).
404 Not Found : Data tidak ditemukan.
409 Conflict : Terjadi duplikasi data (misal: nama kategori/produk sudah ada).
500 Internal Server Error : Terjadi kesalahan pada server/database.