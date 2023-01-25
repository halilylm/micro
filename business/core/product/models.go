package product

import (
	"github.com/google/uuid"
	"time"
)

// Product represents an in individual product.
type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Cost        int       `json:"cost"`
	Quantity    int       `json:"quantity"`
	Sold        int       `json:"sold"`
	Revenue     int       `json:"revenue"`
	UserID      uuid.UUID `json:"user_id"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewProduct is required fields from the clients adding a product
type NewProduct struct {
	Name     string    `json:"name" validate:"required"`
	Cost     int       `json:"cost" validate:"required,gte=0"`
	Quantity int       `json:"quantity" validate:"gte=1"`
	UserID   uuid.UUID `json:"user_id" validate:"required,uuid4"`
}

// UpdateProduct defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields, so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types, but
// we make exceptions around marshalling/unmarshalling.
type UpdateProduct struct {
	Name     *string `json:"name"`
	Cost     *int    `json:"cost" validate:"omitempty,gte=0"`
	Quantity *int    `json:"quantity" validate:"omitempty,gte=1"`
}
