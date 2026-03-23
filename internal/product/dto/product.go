package dto

import "time"

type VariantReq struct {
	Weight string  `json:"weight" validate:"required"`
	Price  float64 `json:"price" validate:"required"`
	Stock  int     `json:"stock"`
}

type NutritionalInfoReq struct {
	Calories string `json:"calories"`
	Protein  string `json:"protein"`
	Fat      string `json:"fat"`
	Carbs    string `json:"carbs"`
	Fiber    string `json:"fiber"`
}

type CreateProductReq struct {
	Name             string              `json:"name" validate:"required"`
	Description      string              `json:"description"`
	ShortDescription string              `json:"short_description"`
	CategoryID       string              `json:"category_id" validate:"required"`
	Images           []string            `json:"images"`
	BasePrice        float64             `json:"base_price" validate:"required"`
	Variants         []VariantReq        `json:"variants"`
	NutritionalInfo  *NutritionalInfoReq `json:"nutritional_info"`
	Featured         bool                `json:"featured"`
	BestSeller       bool                `json:"best_seller"`
}

type UpdateProductReq struct {
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	ShortDescription string              `json:"short_description"`
	CategoryID       string              `json:"category_id"`
	Images           []string            `json:"images"`
	BasePrice        float64             `json:"base_price"`
	Variants         []VariantReq        `json:"variants"`
	NutritionalInfo  *NutritionalInfoReq `json:"nutritional_info"`
	Featured         *bool               `json:"featured"`
	BestSeller       *bool               `json:"best_seller"`
}

type UpdateInventoryReq struct {
	VariantID string `json:"variant_id" validate:"required"`
	Stock     int    `json:"stock" validate:"required"`
}

type VariantRes struct {
	ID     string  `json:"id"`
	Weight string  `json:"weight"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

type NutritionalInfoRes struct {
	Calories string `json:"calories"`
	Protein  string `json:"protein"`
	Fat      string `json:"fat"`
	Carbs    string `json:"carbs"`
	Fiber    string `json:"fiber"`
}

type CategoryRef struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type ProductRes struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Slug             string              `json:"slug"`
	Description      string              `json:"description"`
	ShortDescription string              `json:"shortDescription"`
	Category         CategoryRef         `json:"category"`
	Images           []string            `json:"images"`
	BasePrice        float64             `json:"basePrice"`
	Variants         []VariantRes        `json:"variants"`
	NutritionalInfo  *NutritionalInfoRes `json:"nutritionalInfo"`
	Stock            int                 `json:"stock"`
	Rating           float64             `json:"rating"`
	ReviewCount      int                 `json:"reviewCount"`
	Featured         bool                `json:"featured"`
	BestSeller       bool                `json:"bestSeller"`
	CreatedAt        time.Time           `json:"createdAt"`
}

type ProductListRes struct {
	Products []ProductRes `json:"products"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	Limit    int          `json:"limit"`
}
