package handlers

import (
	"encoding/json"
	"net/http"

	"product-catalogue-Telkom-LKPP/internal/models"
	"product-catalogue-Telkom-LKPP/internal/repositories"

	"github.com/google/uuid"

	"fmt"
)

type ReviewHandler struct {
	ReviewRepo repositories.ReviewRepository
}

func NewReviewHandler(reviewRepo repositories.ReviewRepository) *ReviewHandler {
	return &ReviewHandler{
		ReviewRepo: reviewRepo,
	}
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	// Parse JSON data from the request body
	var requestBody models.Review

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Generate UUID for the review
	reviewID := uuid.New()

	// Create a Review struct with the extracted data
	review := &models.Review{
		ID:        reviewID,
		ProductID: requestBody.ProductID,
		Rating:    requestBody.Rating,
		Comment:   requestBody.Comment,
	}

	// Call the CreateReview method of the repository to insert the review into the database
	err = h.ReviewRepo.CreateReview(review)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Review created successfully %s", reviewID)))
}
