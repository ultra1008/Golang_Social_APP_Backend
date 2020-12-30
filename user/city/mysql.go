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

	rows, err := m.db.NamedQueryContext(ctx, queryMap[createNewCity].SQL, &City{Name: city})
	if err != nil {
		return nil, fmt.Errorf("gettings city: %v", err)
	}

	var c City

	for rows.Next() {
		err := rows.StructScan(&c)
		if err != nil {
			return nil, fmt.Errorf("scanning sql result: %v", err)
		}
	}

	return &c, nil
}

func (m *mysql) List() ([]City, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryMap[listCitys].Timeout)
	defer cancel()

	var citys []City

	err := m.db.SelectContext(ctx, &citys, queryMap[listCitys].SQL)
	if err != nil {
		return nil, fmt.Errorf("getting citys list: %v", err)
	}

	return citys, nil
}
