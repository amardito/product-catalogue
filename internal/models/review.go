package models

import (
	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
}
