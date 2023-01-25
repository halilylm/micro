package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/data/order"
	"github.com/halilylm/micro/business/sys/validate"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
)

// Set of error variables for CRUD operations
var (
	ErrNotFound              = errors.New("user not found")
	ErrInvalidEmail          = errors.New("email is not valid")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrInvalidOrder          = errors.New("validating order by")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Repository interface declares the behaviour this package needs to
// persists and retrieve data.
type Repository interface {
	WithinTran(ctx context.Context, fn func(r Repository) error) error
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// Core manages the set of APIs for user access.
type Core struct {
	repo Repository
}

// NewCore constructs a core for user api access.
func NewCore(repo Repository) *Core {
	return &Core{repo: repo}
}

// Create inserts a new user into the database.
func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	now := time.Now()

	user := User{
		ID:           uuid.New(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}
	if err := c.repo.Create(ctx, user); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}
	return user, nil
}

func (c *Core) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
	if err := validate.Check(uu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generating password hash: %w", err)
		}
		usr.PasswordHash = pw
	}
	if uu.Enabled != nil {
		usr.Enabled = *uu.Enabled
	}
	usr.DateUpdated = time.Now()

	if err := c.repo.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	return usr, nil
}

// Delete removes a user from the database
func (c *Core) Delete(ctx context.Context, usr User) error {
	if err := c.repo.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

// Query retrieves a list of existing users from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error) {
	if err := validate.Check(filter); err != nil {
		return nil, fmt.Errorf("validating filter: %w", err)
	}

	if err := ordering.Check(orderBy); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidOrder, err.Error())
	}

	users, err := c.repo.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

// QueryByID gets the specified user from the database.
func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := c.repo.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: %w", err)
	}

	return user, nil
}

// QueryByEmail gets the specified user from the database by email
func (c *Core) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	user, err := c.repo.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: %w", err)
	}

	return user, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success, it returns a Claims User representing this user. The claims can be
// used to generate a token for authentication.
func (c *Core) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {
	usr, err := c.repo.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return User{}, ErrAuthenticationFailure
	}

	return usr, nil
}
