package productdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/core/product"
	"github.com/halilylm/micro/business/data/order"
	"github.com/halilylm/micro/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

const (
	createQuery = `
	INSERT INTO products
	    (product_id, user_id, name, cost, quantity, date_created, date_updated)
	VALUES 
	    (:product_id, :user_id, :name, :cost, :quantity, :date_created, :date_updated)`
	updateQuery = `
	UPDATE 
		products
	SET 
	    "name" = :name,
	    "cost" = :cost,
	    "quantity" = :quantity,
	    "date_updated": date_updated
	WHERE 
	    product_id = :product_id`
	deleteQuery = `
	DELETE FROM
		products 
	WHERE 
	    product_id = :product_id`
	filterQuery = `
	SELECT 
		p.*,
		COALESCE(SUM(s.quantity), 0) AS sold,
		COALESCE(SUM(s.paid), 0) AS revenue
	FROM 
	    products AS p 
	LEFT JOIN 
	        sales AS s ON p.product_id = s.product_id 
	`
	filterByID = `
	SELECT
		p.*,
		COALESCE(SUM(s.quantity), 0) AS sold, 
		COALESCE(SUM(s.paid), 0) AS revenue,
	FROM 
	    products as p 
	LEFT JOIN 
	        sales AS s on p.product_id = s.product_id 
	WHERE 
	    p.product_id = :product_id 
	GROUP BY 
	    p.product_id`
	filterByUserID = `
	SELECT 
		p.*,
		COALESCE(SUM(s.quantity), 0) AS sold,
		COALESCE(SUM(s.paid), 0) AS revenue 
	FROM
		products AS p 
	LEFT JOIN 
		    sales AS s ON p.product_id = s.product_id 
	WHERE 
	    p.user_id = :user_id 
	GROUP BY
	    p.product_id`
)

// Repository manages the set of APIs for product database access.
type Repository struct {
	log *zap.SugaredLogger
	db  sqlx.ExtContext
}

// NewRepository constructs the api for data access.
func NewRepository(log *zap.SugaredLogger, db *sqlx.DB) *Repository {
	if log == nil {
		log = zap.NewNop().Sugar()
	}
	return &Repository{
		log: log,
		db:  db,
	}
}

func (r *Repository) Create(ctx context.Context, prd product.Product) error {
	if err := database.NamedExecContext(ctx, r.log, r.db, createQuery, toDBProduct(prd)); err != nil {
		return fmt.Errorf("inserting product: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, prd product.Product) error {
	if err := database.NamedExecContext(ctx, r.log, r.db, updateQuery, toDBProduct(prd)); err != nil {
		return fmt.Errorf("updating product productID[%s]: %w", prd.ID, err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, prd product.Product) error {
	data := struct {
		ProductID string `db:"product_id"`
	}{
		ProductID: prd.ID.String(),
	}
	if err := database.NamedExecContext(ctx, r.log, r.db, deleteQuery, data); err != nil {
		return fmt.Errorf("deleting product productID[%s]: %w", prd.ID, err)
	}
	return nil
}

func (r *Repository) Query(ctx context.Context, filter product.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]product.Product, error) {
	data := struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Cost        int    `db:"cost"`
		Quantity    int    `db:"quantity"`
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	var wc []string
	if filter.ID != nil {
		data.ID = *filter.ID
		wc = append(wc, "id = :id")
	}
	if filter.Name != nil {
		data.Name = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}
	if filter.Cost != nil {
		data.Cost = *filter.Cost
		wc = append(wc, "cost = :cost")
	}
	if filter.Quantity != nil {
		data.Quantity = *filter.Quantity
		wc = append(wc, "quantity = :quantity")
	}
	buf := bytes.NewBufferString(filterQuery)
	if len(wc) > 0 {
		buf.WriteString("WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
	buf.WriteString(" GROUP BY p.product_id ")
	buf.WriteString(" ORDER BY ")
	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY ")

	var prds []dbProduct
	if err := database.NamedQuerySlice(ctx, r.log, r.db, buf.String(), data, &prds); err != nil {
		return nil, fmt.Errorf("selecting products: %w", err)
	}

	return toCoreProductSlice(prds), nil
}

func (r *Repository) QueryByID(ctx context.Context, productID uuid.UUID) (product.Product, error) {
	data := struct {
		ProductID string `db:"product_id"`
	}{
		ProductID: productID.String(),
	}
	var prd dbProduct
	if err := database.NamedQueryStruct(ctx, r.log, r.db, filterByID, data, &prd); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return product.Product{}, product.ErrNotFound
		}
		return product.Product{}, fmt.Errorf("selecting product productID[%q]: %w", productID, err)
	}
	return toCoreProduct(prd), nil
}

func (r *Repository) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]product.Product, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID.String(),
	}
	var prds []dbProduct
	if err := database.NamedQuerySlice(ctx, r.log, r.db, filterByUserID, data, &prds); err != nil {
		return nil, fmt.Errorf("select products userID[%s]: %w", userID, err)
	}
	return toCoreProductSlice(prds), nil
}
