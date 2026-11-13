# ==========================================
# Tahap 1: Build Aplikasi
# ==========================================
# Gunakan image Go versi terbaru yang ringan (alpine)
FROM golang:1.26-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# Install dependensi sistem yang dibutuhkan (git untuk go mod)
# (kuota saya habis ) RUN apk add --no-cache git

# Copy file go.mod dan go.sum terlebih dahulu (untuk caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh kode sumber
COPY . .

# Build aplikasi menjadi binary bernama 'toko-api'
# CGO_ENABLED=0 membuat binary statis agar bisa berjalan di image alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o toko-api .

# ==========================================
# Tahap 2: Jalankan Aplikasi (Image Final)
# ==========================================
# Gunakan image alpine yang sangat kecil untuk produksi
FROM alpine:latest

# Install CA-Certificates dan Timezone data (PENTING untuk format tanggal)
RUN apk add --no-cache ca-certificates tzdata

# Set timezone ke Asia/Jakarta
ENV TZ=Asia/Jakarta

WORKDIR /app

# Copy binary dari tahap build sebelumnya
COPY --from=builder /app/toko-api .

# Ekspos port 8080 (sesuai default aplikasi)
EXPOSE 8080

# Perintah untuk menjalankan aplikasi
CMD ["./toko-api"]