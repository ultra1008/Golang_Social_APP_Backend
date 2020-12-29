package city

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type mysql struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository {
	return &mysql{db: db}
}

func (m *mysql) Create(city string) (*City, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryMap[createNewCity].Timeout)
	defer cancel()

	var c City
	err := m.db.GetContext(ctx, &c, queryMap[createNewCity].SQL)
	if err != nil {
		return nil, fmt.Errorf("gettings city: %v", err)
	}

	return &c, nil
}
