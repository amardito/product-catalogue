package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID       `json:"id"`
	SKU         string          `json:"sku"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Etalase     string          `json:"etalase"`
	Images      []*ProductImage `json:"images"`
	Weight      float64         `json:"weight"`
	Price       float64         `json:"price"`
	Rating      float64         `json:"rating"`
}

type ProductImage struct {
	ID          uuid.UUID `json:"id"`
	FilePath    string    `json:"file_path"` // Store the file path here
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Type        string    `json:"type"`
}

type ProductQuery struct {
	Title    string `json:"title"`
	Etalase  string `json:"etalase"`
	Category string `json:"category"`
	SKU      string `json:"sku"`
	SortBy   string `json:"sortBy"`
}

type ProductRequest struct {
	SKU         string   `json:"sku"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Etalase     string   `json:"etalase"`
	Weight      float64  `json:"weight"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"` // Base64-encoded image strings
}
