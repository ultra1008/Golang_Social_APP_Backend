package city

import (
	"context"
	"database/sql"
	"fmt"
)

type mysql struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return &mysql{db: db}
}

func (m *mysql) Create(cityName string) (*City, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryMap[createNewCity].Timeout)
	defer cancel()

	res, err := m.db.ExecContext(ctx, queryMap[createNewCity].SQL, cityName)
	if err != nil {
		return nil, fmt.Errorf("creating city sql query: %v", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last insert id; create city request: %v", err)
	}

	c := City{
		ID:            int(lastID),
		Name:          cityName,
		CreatedByUser: true,
	}

	return &c, nil
}

func (m *mysql) List() ([]City, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryMap[listCitys].Timeout)
	defer cancel()

	var citys []City

	rows, err := m.db.QueryContext(ctx, queryMap[listCitys].SQL)
	if err != nil {
		return nil, fmt.Errorf("getting citys list: %v", err)
	}

	for rows.Next() {
		var city City
		err := rows.Scan(&city.ID, &city.Name, &city.CreatedByUser)
		if err != nil {
			return nil, fmt.Errorf("scanning city list: %v", err)
		}
		citys = append(citys, city)
	}

	return citys, nil
}

func (m *mysql) GetByID(id int) (*City, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryMap[listCitys].Timeout)
	defer cancel()

	var city City

	row := m.db.QueryRowContext(ctx, queryMap[getByID].SQL, id)

	err := row.Scan(&city.ID, &city.Name, &city.CreatedByUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scanning city list: %v", err)
	}

	return &city, nil
}
