package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"product-catalogue-Telkom-LKPP/internal/models"
)

type ReviewRepository interface {
	CreateProduct(product *models.Product) error
}

type reviewRepository struct {
	DB *sql.DB
}

func NewReviewRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (repo *productRepository) CreateReview(product *models.Product) error {
	// Convert images slice to JSONB data
	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return fmt.Errorf("failed to marshal images to JSON: %v", err)
	}

	// Insert new product record into the database
	_, err = repo.DB.Exec(`
		INSERT INTO products (id, sku, title, description, category, etalase, images, weight, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, product.ID, product.SKU, product.Title, product.Description, product.Category, product.Etalase, imagesJSON, product.Weight, product.Price)
	if err != nil {
		return fmt.Errorf("failed to insert product: %v", err)
	}

	return nil
}
