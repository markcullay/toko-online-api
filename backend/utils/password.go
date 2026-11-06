package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mengenkripsi kata sandi menggunakan algoritma Bcrypt.
func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword membandingkan hash Bcrypt dengan kata sandi mentah.
func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}
