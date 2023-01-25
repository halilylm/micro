// Package usercache contains user related CRUD functionality with caching
package usercache

import (
	"context"
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/core/user"
	"github.com/halilylm/micro/business/data/order"
	"go.uber.org/zap"
	"net/mail"
	"sync"
)

type Repository struct {
	log  *zap.SugaredLogger
	repo user.Repository

	mu    sync.RWMutex
	cache map[string]*user.User
}

// NewRepository constructs the api for data and caching access.
func NewRepository(log *zap.SugaredLogger, repo user.Repository) *Repository {
	if log == nil {
		log = zap.NewNop().Sugar()
	}
	return &Repository{
		log:   log,
		repo:  repo,
		cache: map[string]*user.User{},
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (r *Repository) WithinTran(ctx context.Context, fn func(r user.Repository) error) error {
	return r.repo.WithinTran(ctx, fn)
}

// Create insert a new user into the database.
func (r *Repository) Create(ctx context.Context, usr user.User) error {
	if err := r.repo.Create(ctx, usr); err != nil {
		return err
	}

	r.writeCache(usr)

	return nil
}

// Update replaces a user document in the database.
func (r *Repository) Update(ctx context.Context, usr user.User) error {
	if err := r.repo.Update(ctx, usr); err != nil {
		return err
	}

	r.writeCache(usr)

	return nil
}

func (r *Repository) Delete(ctx context.Context, usr user.User) error {
	if err := r.repo.Delete(ctx, usr); err != nil {
		return err
	}

	r.deleteCache(usr)

	return nil
}

func (r *Repository) Query(ctx context.Context, filter user.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]user.User, error) {
	return r.repo.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
}

func (r *Repository) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	cachedUsr, ok := r.readCache(userID.String())
	if ok {
		return cachedUsr, nil
	}

	usr, err := r.repo.QueryByID(ctx, userID)
	if err != nil {
		return user.User{}, err
	}

	r.writeCache(usr)

	return usr, nil
}

func (r *Repository) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	cachedUsr, ok := r.readCache(email.Address)
	if ok {
		return cachedUsr, nil
	}

	usr, err := r.repo.QueryByEmail(ctx, email)
	if err != nil {
		return user.User{}, err
	}

	r.writeCache(usr)

	return usr, nil
}

// readCache performs a safe search in the cache for specified key.
func (r *Repository) readCache(key string) (user.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	usr, exists := r.cache[key]
	if !exists {
		return user.User{}, false
	}

	return *usr, true
}

// writeCache performs a safe write to the cache for specified user.
func (r *Repository) writeCache(usr user.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[usr.ID.String()] = &usr
	r.cache[usr.Email.Address] = &usr
}

// deleteCache performs a safe removal from the cache for the specified user.
func (r *Repository) deleteCache(usr user.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cache, usr.ID.String())
	delete(r.cache, usr.Email.Address)
}
