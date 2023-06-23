package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	product.ID = uuid.NewString()
	return
}

type Product struct {
	ID          string `json:"id" gorm:"primaryKey"`
	IdToped     string `json:"id_toped" gorm:"null"`
	CreatedAt   time.Time
	CreatedBy   string `json:"created_by"`
	UpdatedAt   time.Time
	UpdatedBy   string  `json:"updated_by"`
	IsActive    bool    `json:"is_active" gorm:"default:true;not null"`
	Name        string  `json:"name" gorm:"not null"`
	Weight      int32   `json:"weight" gorm:"not null"`
	Detail      string  `json:"detail" gorm:"not null"`
	UrlVideo    string  `json:"url_video"`
	Brand       string  `json:"brand" gorm:"not null"`
	BoxItem     string  `json:"box_item" gorm:"not null"`
	SKU         string  `json:"sku" gorm:"not null"`
	Slug        string  `json:"slug" gorm:"not null"`
	Image1      string  `json:"image_1" gorm:"not null"`
	Image2      string  `json:"image_2"`
	Image3      string  `json:"image_3"`
	Image4      string  `json:"image_4"`
	Image5      string  `json:"image_5"`
	Prices      float64 `json:"prices" gorm:"not null"`
	Stock       int64   `json:"stock" gorm:"not null"`
	IsAvailable bool    `json:"is_available" gorm:"not null"`
	Location    string  `json:"location" gorm:"not null"`
	UrlProduct  string  `json:"url_product" gorm:"null"`
}
