package usergrp

import (
	"github.com/google/uuid"
	"github.com/halilylm/micro/business/core/user"
	"net/http"
	"net/mail"
)

func getFilter(r *http.Request) (user.QueryFilter, error) {
	values := r.URL.Query()

	var filter user.QueryFilter
	if id, err := uuid.Parse(values.Get("id")); err == nil {
		filter.ByID(id)
	}

	filter.ByName(values.Get("name"))

	if email, err := mail.ParseAddress(values.Get("email")); err == nil {
		filter.ByEmail(*email)
	}

	return filter, nil
}
