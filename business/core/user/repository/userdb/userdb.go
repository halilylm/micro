// Package userdb contains user related CRUD functionality
package userdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/core/user"
	"github.com/halilylm/micro/business/data/order"
	"github.com/halilylm/micro/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/mail"
	"strings"
)

// Repository manages the set of APIs for user database accesr.
type Repository struct {
	log    *zap.SugaredLogger
	db     sqlx.ExtContext
	inTran bool
}

// NewRepository constructs the api for data accesr.
func NewRepository(log *zap.SugaredLogger, db *sqlx.DB) *Repository {
	return &Repository{
		log: log,
		db:  db,
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (r *Repository) WithinTran(ctx context.Context, fn func(s user.Repository) error) error {
	if r.inTran {
		return fn(r)
	}

	f := func(tx *sqlx.Tx) error {
		s := &Repository{
			log:    r.log,
			db:     tx,
			inTran: true,
		}
		return fn(s)
	}

	return database.WithinTran(ctx, r.log, r.db.(*sqlx.DB), f)
}

// Create inserts a new user into the database.
func (r *Repository) Create(ctx context.Context, usr user.User) error {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, enabled, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :enabled, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, r.log, r.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return fmt.Errorf("create: %w", user.ErrUniqueEmail)
		}
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

// Update replaces a user document in the database.
func (r *Repository) Update(ctx context.Context, usr user.User) error {
	const q = `
	UPDATE
		users
	SET 
		"name" = :name,
		"email" = :email,
		"roles" = :roles,
		"password_hash" = :password_hash,
		"date_updated" = :date_updated
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, r.log, r.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return user.ErrUniqueEmail
		}
		return fmt.Errorf("updating userID[%s]: %w", usr.ID, err)
	}

	return nil
}

// Delete removes a user from the database.
func (r *Repository) Delete(ctx context.Context, usr user.User) error {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: usr.ID.String(),
	}

	const q = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, r.log, r.db, q, data); err != nil {
		return fmt.Errorf("deleting userID[%s]: %w", usr.ID, err)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (r *Repository) Query(ctx context.Context, filter user.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]user.User, error) {
	data := struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Email       string `db:"email"`
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
		data.ID = (*filter.ID).String()
		wc = append(wc, "id = :id")
	}
	if filter.Name != nil {
		data.Name = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}
	if filter.Email != nil {
		data.Email = (*filter.Email).String()
		wc = append(wc, "email = :email")
	}

	const q = `
	SELECT
		*
	FROM
		users
	`
	buf := bytes.NewBufferString(q)

	if len(wc) > 0 {
		buf.WriteString("WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
	buf.WriteString(" ORDER BY ")
	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var usrs []dbUser
	if err := database.NamedQuerySlice(ctx, r.log, r.db, buf.String(), data, &usrs); err != nil {
		return nil, fmt.Errorf("selecting users: %w", err)
	}

	return toCoreUserSlice(usrs), nil
}

// QueryByID gets the specified user from the database.
func (r *Repository) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID.String(),
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = :user_id`

	var usr dbUser
	if err := database.NamedQueryStruct(ctx, r.log, r.db, q, data, &usr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, fmt.Errorf("selecting userID[%q]: %w", userID, err)
	}

	return toCoreUser(usr), nil
}

// QueryByEmail gets the specified user from the database by email.
func (r *Repository) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email.Address,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		email = :email`

	var usr dbUser
	if err := database.NamedQueryStruct(ctx, r.log, r.db, q, data, &usr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}

	return toCoreUser(usr), nil
}