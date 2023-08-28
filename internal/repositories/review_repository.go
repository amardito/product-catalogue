package repositories

import (
	"database/sql"
	"fmt"
	"product-catalogue-Telkom-LKPP/internal/models"
)

type ReviewRepository interface {
	CreateReview(review *models.Review) error
}

type reviewRepository struct {
	DB *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &reviewRepository{
		DB: db,
	}
}

func (repo *reviewRepository) CreateReview(review *models.Review) error {
	// Insert new review record into the database
	_, err := repo.DB.Exec(`
		INSERT INTO product_reviews (id, product_id, rating, review_comment)
		VALUES ($1, $2, $3, $4)
	`, review.ID, review.ProductID, review.Rating, review.Comment)
	if err != nil {
		return fmt.Errorf("failed to insert review: %v", err)
	}

	return nil
}
