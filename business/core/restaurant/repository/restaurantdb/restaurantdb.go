// Package restaurantdb contains product related CRUD functionality.
package restaurantdb

import (
	"context"
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/core/restaurant"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

const index = "restaurant"

// Repository manages the set of APIs for restaurant database access.
type Repository struct {
	log    *zap.SugaredLogger
	client *elastic.Client
}

func (r *Repository) Create(ctx context.Context, rest restaurant.Restaurant) error {
	_, err := r.client.Index().
		Index(index).
		Id(rest.ID.String()).
		BodyJson(rest).
		Do(ctx)
	return err
}

func (r *Repository) Update(ctx context.Context, rest restaurant.Restaurant) error {
	_, err := r.client.
		Update().
		Index("restaurant").
		Id(rest.ID.String()).
		Upsert(rest).
		Do(ctx)
	return err
}

func (r *Repository) Delete(ctx context.Context, rest restaurant.Restaurant) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Query(ctx context.Context, filter restaurant.QueryFilter, pageNumber int, rowsPerPage int) ([]restaurant.Restaurant, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) QueryByID(ctx context.Context, restID uuid.UUID) (restaurant.Restaurant, error) {
	//TODO implement me
	panic("implement me")
}
