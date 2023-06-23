package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ps *ProductSource) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	ps.ID = uuid.NewString()
	return
}

type ProductSource struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	IsActive      bool      `json:"is_active" gorm:"default:true;not null"`
	Slug          string    `json:"slug" gorm:"not null"`
	ProductSource string    `json:"product_source" gorm:"not null"`
}
