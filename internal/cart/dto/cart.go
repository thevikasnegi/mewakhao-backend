package dto

type AddItemReq struct {
	ProductID string `json:"product_id" validate:"required"`
	VariantID string `json:"variant_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type UpdateItemReq struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type CartItemRes struct {
	ID       string     `json:"id"`
	Product  ProductRef `json:"product"`
	Variant  VariantRef `json:"variant"`
	Quantity int        `json:"quantity"`
}

type ProductRef struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Slug             string   `json:"slug"`
	ShortDescription string   `json:"shortDescription"`
	Images           []string `json:"images"`
	BasePrice        float64  `json:"basePrice"`
}

type VariantRef struct {
	ID     string  `json:"id"`
	Weight string  `json:"weight"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

type CartRes struct {
	ID       string        `json:"id"`
	Items    []CartItemRes `json:"items"`
	Subtotal float64       `json:"subtotal"`
	Shipping float64       `json:"shipping"`
	Total    float64       `json:"total"`
}
