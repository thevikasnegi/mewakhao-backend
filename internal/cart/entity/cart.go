package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
	UserID    string     `json:"user_id" gorm:"uniqueIndex;not null"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New().String()
	return nil
}

type CartItem struct {
	ID        string          `json:"id" gorm:"unique;not null;index;primary_key"`
	CartID    string          `json:"cart_id" gorm:"index"`
	ProductID string          `json:"product_id"`
	VariantID string          `json:"variant_id"`
	Quantity  int             `json:"quantity" gorm:"default:1"`
	Product   *Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Variant   *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	ci.ID = uuid.New().String()
	return nil
}

// Product is a read-only reference to avoid circular imports
type Product struct {
	ID               string      `json:"id" gorm:"primary_key"`
	Name             string      `json:"name"`
	Slug             string      `json:"slug"`
	ShortDescription string      `json:"short_description"`
	Images           StringArray `json:"images" gorm:"type:text"`
	BasePrice        float64     `json:"base_price"`
	Stock            int         `json:"stock"`
}

func (Product) TableName() string {
	return "products"
}

type ProductVariant struct {
	ID        string  `json:"id" gorm:"primary_key"`
	ProductID string  `json:"product_id"`
	Weight    string  `json:"weight"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

// Reuse the StringArray type from the product entity
type StringArray = []string
