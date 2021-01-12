package interest

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type mysql struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return &mysql{
		db: db,
	}
}

func (m *mysql) CreateIfNotExists(i *Interest) error {
	query := queryMap[createIfNotExists]
	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	res, err := m.db.ExecContext(ctx, query.SQL, i.Name)
	if err != nil {
		return fmt.Errorf("creating interest: %v", err)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %v", err)
	}

	i.ID = int(lastInsertId)

	return nil
}

func (m *mysql) List() ([]Interest, error) {
	query := queryMap[listInterests]
	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query.SQL)
	if err != nil {
		return nil, fmt.Errorf("get interests list: %v", err)
	}
	defer rows.Close()

	interests := []Interest{}

	for rows.Next() {
		var i Interest

		err := rows.Scan(&i.ID, &i.Name)
		if err != nil {
			log.Printf("scanning interest rows: %v", err)
			continue
		}

		interests = append(interests, i)
	}

	return interests, nil
}

func (m *mysql) InterestsByUserId(id int) ([]Interest, error) {
	query := queryMap[getUserInterests]
	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query.SQL, id)
	if err != nil {
		return nil, fmt.Errorf("get interests list by user id: %v", err)
	}
	defer rows.Close()

	interests := []Interest{}

	for rows.Next() {
		var i Interest

		err := rows.Scan(&i.ID, &i.Name)
		if err != nil {
			log.Printf("scanning interest by user id rows: %v", err)
			continue
		}

		interests = append(interests, i)
	}

	return interests, nil
}

func (m *mysql) AddInterestToUser(userId, interestId int) error {
	query := queryMap[addInterestToUser]
	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query.SQL, userId, interestId)
	if err != nil {
		return fmt.Errorf("adding interest to user: %v", err)
	}

	return nil
}
