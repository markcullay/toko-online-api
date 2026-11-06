package models

import "time"

type Sale struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	InvoiceNo string     `gorm:"size:50;uniqueIndex;not null" json:"invoice_no"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	BuyerName string     `gorm:"size:100;not null" json:"buyer"`
	Method    string     `gorm:"size:50;not null" json:"method"`
	Total     int64      `gorm:"not null" json:"total"`
	Items     []SaleItem `gorm:"foreignKey:SaleID" json:"items"`
	CreatedAt time.Time  `json:"date"`
}

func (Sale) TableName() string { return "sales" }

type SaleItem struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	SaleID      uint   `gorm:"not null;index" json:"-"`
	ProductID   uint   `gorm:"not null;index" json:"product_id"`
	ProductName string `gorm:"size:150;not null" json:"name"` // snapshot nama saat beli
	Price       int64  `gorm:"not null" json:"price"`         // snapshot harga saat beli
	Qty         int    `gorm:"not null" json:"qty"`
	Subtotal    int64  `gorm:"not null" json:"subtotal"`
}

func (SaleItem) TableName() string { return "sale_items" }
