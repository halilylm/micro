package restaurantdb

import (
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/core/restaurant"
	"time"
)

// dbRestaurant represents an individual restaurant.
type dbRestaurant struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

func toDBRestaurant(rest restaurant.Restaurant) dbRestaurant {
	restDB := dbRestaurant{
		ID:          rest.ID,
		Name:        rest.Name,
		Location:    rest.Location,
		DateCreated: rest.DateCreated.UTC(),
		DateUpdated: rest.DateUpdated.UTC(),
	}

	return restDB
}

func toCoreRestaurant(dbRest dbRestaurant) restaurant.Restaurant {
	rest := restaurant.Restaurant{
		ID:          dbRest.ID,
		Name:        dbRest.Name,
		Location:    dbRest.Location,
		DateCreated: dbRest.DateCreated.In(time.Local),
		DateUpdated: dbRest.DateUpdated.In(time.Local),
	}

	return rest
}

func toCoreRestaurantSlice(dbRests []dbRestaurant) []restaurant.Restaurant {
	rests := make([]restaurant.Restaurant, len(dbRests))
	for i, dbRest := range dbRests {
		rests[i] = toCoreRestaurant(dbRest)
	}
	return rests
}
