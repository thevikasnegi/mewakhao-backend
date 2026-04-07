package dto

import "time"

type CreateCategoryReq struct {
	Name        string `form:"name" validate:"required"`
	Description string `form:"description"`
	// Image is set by the controller after uploading to Cloudinary.
	Image string
}

type UpdateCategoryReq struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	// Image is set by the controller after uploading to Cloudinary (optional on update).
	Image string
}

type CategoryRes struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
