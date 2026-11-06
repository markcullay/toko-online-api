package models

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `gorm:"index" json:"category_id"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // Relasi
	Name        string    `gorm:"size:150;not null;uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"desc"` // "desc" agar cocok dengan frontend
	Price       int64     `gorm:"not null;check:price > 0" json:"price"`
	Stock       int       `gorm:"not null;default:0;check:stock >= 0" json:"stock"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (Product) TableName() string { return "products" }
