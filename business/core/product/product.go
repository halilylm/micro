package product

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/data/order"
	"github.com/halilylm/micro/business/sys/validate"
	"time"
)

var (
	ErrNotFound     = errors.New("product not found")
	ErrInvalidID    = errors.New("ID is not in its proper form")
	ErrInvalidOrder = errors.New("validating order by")
)

// Repository interface declares the behaviour this package needs to persist
// and retrieve data
type Repository interface {
	Create(ctx context.Context, prd Product) error
	Update(ctx context.Context, prd Product) error
	Delete(ctx context.Context, prd Product) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Product, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error)
}

// Core manages the set of APIs for product access.
type Core struct {
	repo Repository
}

// NewCore constructs a core for product api access.
func NewCore(repo Repository) *Core {
	return &Core{repo: repo}
}

// Create adds a Product to the database. It returns the created Product with
// fields Like ID and DateCreated populated.
func (c *Core) Create(ctx context.Context, np NewProduct) (Product, error) {
	if err := validate.Check(np); err != nil {
		return Product{}, fmt.Errorf("validating data: %w", err)
	}

	now := time.Now()

	prd := Product{
		ID:          uuid.New(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      np.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.repo.Create(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func (c *Core) Update(ctx context.Context, prd Product, up UpdateProduct) (Product, error) {
	if err := validate.Check(up); err != nil {
		return Product{}, fmt.Errorf("validating data: %w", err)
	}

	if up.Name != nil {
		prd.Name = *up.Name
	}

	if up.Cost != nil {
		prd.Cost = *up.Cost
	}

	if up.Quantity != nil {
		prd.Quantity = *up.Quantity
	}

	prd.DateUpdated = time.Now()

	if err := c.repo.Update(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}

// Delete removes the product identified by a given ID.
func (c *Core) Delete(ctx context.Context, prd Product) error {
	if err := c.repo.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query gets all Products from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Product, error) {
	if err := validate.Check(filter); err != nil {
		return nil, fmt.Errorf("validating filter: %w", err)
	}

	if err := ordering.Check(orderBy); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidOrder, err.Error())
	}

	prds, err := c.repo.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// QueryByID finds the product identified by a given ID.
func (c *Core) QueryByID(ctx context.Context, productID uuid.UUID) (Product, error) {
	prd, err := c.repo.QueryByID(ctx, productID)
	if err != nil {
		return Product{}, fmt.Errorf("query: %w", err)
	}

	return prd, nil
}

// QueryByUserID finds the products by a given User ID.
func (c *Core) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error) {
	prds, err := c.repo.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}
