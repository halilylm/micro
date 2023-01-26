// Package restaurant provides an example of a core business API.
package restaurant

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/sys/validate"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("invalid ID")
)

// Repository interface declares the behaviour this package needs to persists and
// retrieve data.
type Repository interface {
	Create(ctx context.Context, rest Restaurant) error
	Update(ctx context.Context, rest Restaurant) error
	Delete(ctx context.Context, rest Restaurant) error
	Query(ctx context.Context, filter QueryFilter, pageNumber int, rowsPerPage int) ([]Restaurant, error)
	QueryByID(ctx context.Context, restID uuid.UUID) (Restaurant, error)
}

// Core manages the set of APIs for restaurant access.
type Core struct {
	repo Repository
}

// NewCore constructs a core for restaurant api access.
func NewCore(repo Repository) *Core {
	return &Core{repo}
}

// Create adds a Restaurant to the database. It returns the created Restaurant with
// fields like ID and DateCreated populated.
func (c *Core) Create(ctx context.Context, nr NewRestaurant) (Restaurant, error) {
	if err := validate.Check(nr); err != nil {
		return Restaurant{}, fmt.Errorf("validating data: %w", err)
	}

	now := time.Now()

	rest := Restaurant{
		ID:          uuid.New(),
		Name:        nr.Name,
		Location:    nr.Location,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.repo.Create(ctx, rest); err != nil {
		return Restaurant{}, fmt.Errorf("create: %w", err)
	}

	return rest, nil
}

// Update modifies data about a Restaurant. It will error if the specified ID is
// invalid or does not reference an existing Restaurant.
func (c *Core) Update(ctx context.Context, rest Restaurant, ur UpdateRestaurant) (Restaurant, error) {
	if err := validate.Check(ur); err != nil {
		return Restaurant{}, fmt.Errorf("validating data: %w", err)
	}

	if ur.Name != nil {
		rest.Name = *ur.Name
	}
	if ur.Location != nil {
		rest.Location = *ur.Location
	}
	rest.DateUpdated = time.Now()

	if err := c.repo.Update(ctx, rest); err != nil {
		return Restaurant{}, fmt.Errorf("update: %w", err)
	}

	return rest, nil
}

// Delete removes the Restaurant identified by a given ID.
func (c *Core) Delete(ctx context.Context, rest Restaurant) error {
	if err := c.repo.Delete(ctx, rest); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query gets all Restaurant from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter, pageNumber int, rowsPerPage int) ([]Restaurant, error) {
	if err := validate.Check(filter); err != nil {
		return nil, fmt.Errorf("validating filter: %w", err)
	}

	rests, err := c.repo.Query(ctx, filter, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return rests, err
}

// QueryByID finds the Restaurant identified by a given ID.
func (c *Core) QueryByID(ctx context.Context, restID uuid.UUID) (Restaurant, error) {
	rest, err := c.repo.QueryByID(ctx, restID)
	if err != nil {
		return Restaurant{}, fmt.Errorf("query: %w", err)
	}

	return rest, nil
}
