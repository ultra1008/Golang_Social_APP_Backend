package user

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/niklod/highload-social-network/config"
)

type mysql struct {
	db *sqlx.DB
}

func NewRepository(cfg *config.DBConfig) (repository, error) {
	db, err := sqlx.Connect("mysql", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping: %v", err)
	}

	return mysql{db: db}, nil
}
