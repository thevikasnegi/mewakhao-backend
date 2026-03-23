package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringArray is a custom type for storing string slices as JSON in Postgres
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	data, err := json.Marshal(s)
	return string(data), err
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	bytes, ok := value.(string)
	if !ok {
		b, ok := value.([]byte)
		if !ok {
			return errors.New("failed to scan StringArray")
		}
		bytes = string(b)
	}
	return json.Unmarshal([]byte(bytes), s)
}

type Product struct {
	ID               string           `json:"id" gorm:"unique;not null;index;primary_key"`
	Name             string           `json:"name" gorm:"not null"`
	Slug             string           `json:"slug" gorm:"unique;not null;index"`
	Description      string           `json:"description"`
	ShortDescription string           `json:"short_description"`
	CategoryID       string           `json:"category_id" gorm:"index"`
	Category         *Category        `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Images           StringArray      `json:"images" gorm:"type:text"`
	BasePrice        float64          `json:"base_price"`
	Variants         []ProductVariant `json:"variants" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	NutritionalInfo  *NutritionalInfo `json:"nutritional_info,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Stock            int              `json:"stock"`
	Rating           float64          `json:"rating" gorm:"default:0"`
	ReviewCount      int              `json:"review_count" gorm:"default:0"`
	Featured         bool             `json:"featured" gorm:"default:false"`
	BestSeller       bool             `json:"best_seller" gorm:"default:false"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	DeletedAt        *time.Time       `json:"deleted_at" gorm:"index"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New().String()
	return nil
}

// Category is embedded here so the product module can reference it without circular imports
type Category struct {
	ID          string `json:"id" gorm:"unique;not null;index;primary_key"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

func (Category) TableName() string {
	return "categories"
}

type ProductVariant struct {
	ID        string  `json:"id" gorm:"unique;not null;index;primary_key"`
	ProductID string  `json:"product_id" gorm:"index"`
	Weight    string  `json:"weight"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

func (v *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	v.ID = uuid.New().String()
	return nil
}

type NutritionalInfo struct {
	ID        string `json:"id" gorm:"unique;not null;index;primary_key"`
	ProductID string `json:"product_id" gorm:"uniqueIndex"`
	Calories  string `json:"calories"`
	Protein   string `json:"protein"`
	Fat       string `json:"fat"`
	Carbs     string `json:"carbs"`
	Fiber     string `json:"fiber"`
}

func (n *NutritionalInfo) BeforeCreate(tx *gorm.DB) error {
	n.ID = uuid.New().String()
	return nil
}
