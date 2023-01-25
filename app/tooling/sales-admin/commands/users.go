package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/halilylm/micro/business/core/user"
	"github.com/halilylm/micro/business/core/user/repository/userdb"
	"github.com/halilylm/micro/business/sys/database"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

func Users(log *zap.SugaredLogger, cfg database.Config, pageNumber string, rowsPerPage string) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := strconv.Atoi(pageNumber)
	if err != nil {
		return fmt.Errorf("converting page number: %w", err)
	}

	rows, err := strconv.Atoi(rowsPerPage)
	if err != nil {
		return fmt.Errorf("converting rows per page: %w", err)
	}

	core := user.NewCore(userdb.NewRepository(log, db))

	users, err := core.Query(ctx, user.QueryFilter{}, user.DefaultOrderBy, page, rows)
	if err != nil {
		return fmt.Errorf("retrieve users: %w", err)
	}

	return json.NewEncoder(os.Stdout).Encode(users)
}
