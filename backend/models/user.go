package models

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      Role      `gorm:"size:20;not null;default:customer" json:"role"`
	CreatedAt time.Time `json:"joined"`
	UpdatedAt time.Time `json:"-"`
}

func (User) TableName() string { return "users" }
