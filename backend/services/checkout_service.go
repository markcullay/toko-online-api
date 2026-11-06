package services

import (
	"fmt"
	"strings"
	"time"

	"backend/dto"
	"backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CheckoutService struct {
	DB *gorm.DB
}

func NewCheckoutService(db *gorm.DB) *CheckoutService {
	return &CheckoutService{DB: db}
}

// Checkout menjalankan seluruh logika bisnis pembelian:
// validasi, row lock, cek stok, hitung total, buat sale, kurangi stok.
// Mengembalikan *models.Sale yang sudah tersimpan, atau error.
func (s *CheckoutService) Checkout(userID uint, req dto.CheckoutRequest) (*models.Sale, error) {
	req.Buyer = strings.TrimSpace(req.Buyer)
	if req.Buyer == "" {
		return nil, fmt.Errorf("nama pembeli wajib diisi")
	}

	// Gabungkan qty jika ada product_id duplikat di request
	qtyMap := make(map[uint]int)
	for _, it := range req.Items {
		qtyMap[it.ProductID] += it.Qty
	}

	var sale models.Sale

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Ambil produk dengan row lock (FOR UPDATE).
		//    Mencegah race condition ketika 2 pembeli membeli stok terakhir bersamaan.
		ids := make([]uint, 0, len(qtyMap))
		for id := range qtyMap {
			ids = append(ids, id)
		}

		var products []models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ?", ids).
			Find(&products).Error; err != nil {
			return err
		}

		if len(products) != len(ids) {
			return fmt.Errorf("ada produk yang tidak ditemukan")
		}

		// 2. Validasi stok + siapkan item
		var items []models.SaleItem
		var total int64
		for _, p := range products {
			qty := qtyMap[p.ID]
			if p.Stock < qty {
				return fmt.Errorf("stok %q tidak cukup (tersisa %d, diminta %d)", p.Name, p.Stock, qty)
			}

			sub := p.Price * int64(qty)
			items = append(items, models.SaleItem{
				ProductID:   p.ID,
				ProductName: p.Name,
				Price:       p.Price,
				Qty:         qty,
				Subtotal:    sub,
			})
			total += sub
		}

		// 3. Buat sale header
		now := time.Now()
		sale = models.Sale{
			InvoiceNo: generateInvoiceNo(now),
			UserID:    userID,
			BuyerName: req.Buyer,
			Method:    req.Method,
			Total:     total,
			Items:     items,
			CreatedAt: now,
		}

		// GORM otomatis insert sale lalu sale_items karena relasi has-many
		if err := tx.Create(&sale).Error; err != nil {
			return err
		}

		// 4. Kurangi stok tiap produk
		for _, it := range items {
			// Update bersyarat: hanya kurangi jika stok masih cukup.
			// Ini lapisan kedua pengaman selain row lock.
			res := tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", it.ProductID, it.Qty).
				UpdateColumn("stock", gorm.Expr("stock - ?", it.Qty))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("stok %q berubah saat transaksi, silakan coba lagi", it.ProductName)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func generateInvoiceNo(t time.Time) string {
	return fmt.Sprintf("INV-%d", t.UnixNano()/1e6) // milidetik, unik
}
