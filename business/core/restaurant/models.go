package restaurant

import (
	"github.com/google/uuid"
	"time"
)

// Restaurant represents an individual restaurant.
type Restaurant struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewRestaurant is that what we require from clients when adding a Product.
type NewRestaurant struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

// UpdateRestaurant defines what information may be provided to modify an
// existing Restaurant. All fields are optional so clients can send just the fields
// they want changed. It uses pointer fields, so we can differentiate between a field
// that was not provided and a field that was provided as explicitly blank. Normally we
// do not want to user pointers to basic types, but we make exceptions around
// marshalling and unmarshalling
type UpdateRestaurant struct {
	Name     *string `json:"name" validate:"omitempty,gte=0"`
	Location *string `json:"location" validate:"omitempty,gte=1"`
}
