package productdb

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/core/product"
	"github.com/halilylm/micro/business/data/order"
	"time"
)

// dbProduct represents an individual product.
type dbProduct struct {
	ID          uuid.UUID `db:"product_id"`
	Name        string    `db:"name"`
	Cost        int       `db:"cost"`
	Quantity    int       `db:"quantity"`
	Sold        int       `db:"sold"`
	Revenue     int       `db:"revenue"`
	UserID      uuid.UUID `db:"user_id"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBProduct(prd product.Product) dbProduct {
	return dbProduct{
		ID:          prd.ID,
		Name:        prd.Name,
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		Sold:        prd.Sold,
		Revenue:     prd.Revenue,
		UserID:      prd.UserID,
		DateCreated: prd.DateCreated,
		DateUpdated: prd.DateUpdated,
	}
}

func toCoreProduct(dbPrd dbProduct) product.Product {
	return product.Product{
		ID:          dbPrd.ID,
		Name:        dbPrd.Name,
		Cost:        dbPrd.Cost,
		Quantity:    dbPrd.Quantity,
		Sold:        dbPrd.Sold,
		Revenue:     dbPrd.Revenue,
		UserID:      dbPrd.UserID,
		DateCreated: dbPrd.DateCreated,
		DateUpdated: dbPrd.DateUpdated,
	}
}

func toCoreProductSlice(dbProducts []dbProduct) []product.Product {
	prds := make([]product.Product, len(dbProducts))
	for i, dbPrd := range dbProducts {
		prds[i] = toCoreProduct(dbPrd)
	}
	return prds
}

// orderByFields is the map of fields that is used to translate between the
// application layer names and the database
var orderByFields = map[string]string{
	product.OrderByID:       "product_id",
	product.OrderByName:     "name",
	product.OrderByCost:     "cost",
	product.OrderByQuantity: "quantity",
	product.OrderBySold:     "sold",
	product.OrderByRevenue:  "revenue",
	product.OrderByUserID:   "user_id",
}

// orderByClause validates the order by for correct fields and sql injection.
func orderByClause(orderBy order.By) (string, error) {
	if err := order.Validate(orderBy.Field, orderBy.Direction); err != nil {
		return "", err
	}

	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exists", orderBy.Field)
	}

	return by + " " + orderBy.Direction, nil
}
