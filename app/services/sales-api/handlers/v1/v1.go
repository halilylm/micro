// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"github.com/halilylm/micro/app/services/sales-api/handlers/v1/productgrp"
	"github.com/halilylm/micro/app/services/sales-api/handlers/v1/usergrp"
	"github.com/halilylm/micro/business/core/product"
	"github.com/halilylm/micro/business/core/product/repository/productdb"
	"github.com/halilylm/micro/business/core/user"
	"github.com/halilylm/micro/business/core/user/repository/usercache"
	"github.com/halilylm/micro/business/core/user/repository/userdb"
	"github.com/halilylm/micro/business/web/auth"
	"github.com/halilylm/micro/business/web/v1/mid"
	"github.com/halilylm/micro/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *zap.SugaredLogger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)
	admin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)

	ugh := usergrp.Handlers{
		User: user.NewCore(usercache.NewRepository(cfg.Log, userdb.NewRepository(cfg.Log, cfg.DB))),
		Auth: cfg.Auth,
	}
	app.Handle(http.MethodGet, version, "/users/token/:kid", ugh.Token)
	app.Handle(http.MethodGet, version, "/users/:page/:rows", ugh.Query, authen, admin)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/users", ugh.Create, authen, admin)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update, authen, admin)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete, authen, admin)

	pgh := productgrp.Handlers{
		Product: product.NewCore(productdb.NewRepository(cfg.Log, cfg.DB)),
		Auth:    cfg.Auth,
	}
	app.Handle(http.MethodGet, version, "/products/:page/:rows", pgh.Query, authen)
	app.Handle(http.MethodGet, version, "/products/:id", pgh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/products", pgh.Create, authen)
	app.Handle(http.MethodPut, version, "/products/:id", pgh.Update, authen)
	app.Handle(http.MethodDelete, version, "/products/:id", pgh.Delete, authen)
}
