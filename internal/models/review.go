package models

import (
	"github.com/google/uuid"
)

type Review struct {
	ID      uuid.UUID `json:"id"`
	Rating  int       `json:"rating"`
	Comment string    `json:"comment"`
}
