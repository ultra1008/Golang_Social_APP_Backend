package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/niklod/highload-social-network/user/city"
)

type mysql struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return &mysql{db: db}
}

func (m *mysql) Create(user *User) (*User, error) {
	query := queryMap[createQuery]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	res, err := m.db.ExecContext(ctx, query.SQL,
		user.FirstName,
		user.Lastname,
		user.Age,
		user.Sex,
		user.City.ID,
		user.Login,
		user.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("creating new user: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("gettings last insert id: %v", err)
	}

	user.ID = int(id)

	return user, nil
}

func (m *mysql) List() ([]User, error) {
	query := queryMap[listQuery]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query.SQL)
	if err != nil {
		return nil, fmt.Errorf("list users: %v", err)
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		var cityName sql.NullString
		var cityID sql.NullInt64

		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.Lastname,
			&u.Age,
			&u.Sex,
			&u.Login,
			&cityID,
			&cityName,
		)
		if err != nil {
			log.Printf("scanning users list row: %v", err)
			continue
		}

		u.City = city.City{}

		if cityName.Valid && cityID.Valid {
			u.City.Name = cityName.String
			u.City.ID = int(cityID.Int64)
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating through rows: %v", err)
	}

	return users, nil
}

func (m *mysql) GetByID(id int) (*User, error) {
	query := queryMap[getByID]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	var user User
	var cityName sql.NullString
	var cityID sql.NullInt64

	row := m.db.QueryRowContext(ctx, query.SQL, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.Lastname,
		&user.Age,
		&user.Sex,
		&user.Login,
		&cityID,
		&cityName,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user id %d not found: %v", id, err)
		}
		return nil, fmt.Errorf("get user by id: scanning user sql row: %v", err)
	}

	user.City = city.City{}

	if cityName.Valid && cityID.Valid {
		user.City.Name = cityName.String
		user.City.ID = int(cityID.Int64)
	}

	return &user, nil
}

func (m *mysql) GetByFirstAndLastName(firstname, lastname string) ([]User, error) {
	query := queryMap[GetByFirstAndLastName]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	var cityName sql.NullString
	var cityID sql.NullInt64

	users := []User{}

	firstNameQuery := "%" + firstname + "%"
	lastNameQuery := "%" + lastname + "%"

	rows, err := m.db.QueryContext(ctx, query.SQL, firstNameQuery, lastNameQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.Lastname,
			&user.Age,
			&user.Sex,
			&user.Login,
			&cityID,
			&cityName,
		)
		if err != nil {
			return nil, fmt.Errorf("get user by id: scanning user sql row: %v", err)
		}

		user.City = city.City{}

		if cityName.Valid && cityID.Valid {
			user.City.Name = cityName.String
			user.City.ID = int(cityID.Int64)
		}

		users = append(users, user)
	}

	return users, nil
}

func (m *mysql) GetByLogin(login string) (*User, error) {
	query := queryMap[getByLogin]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	var user User
	var cityName sql.NullString
	var cityID sql.NullInt64

	row := m.db.QueryRowContext(ctx, query.SQL, login)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.Lastname,
		&user.Age,
		&user.Sex,
		&user.Login,
		&cityID,
		&cityName,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by login: scanning user sql row: %v", err)
	}

	user.City = city.City{}

	if cityName.Valid && cityID.Valid {
		user.City.Name = cityName.String
		user.City.ID = int(cityID.Int64)
	}

	return &user, nil
}

func (m *mysql) AddFriend(userId int, friendId int) error {
	query := queryMap[addFriend]
	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	if userId == friendId {
		return fmt.Errorf("user ID and friend ID are equal")
	}

	res, err := m.db.ExecContext(ctx, query.SQL, userId, friendId, friendId, userId)
	if err != nil {
		return fmt.Errorf("adding friend with ID %d to user ID %d: %v", friendId, userId, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting afected rows: %v", err)
	}

	if int(rowsAffected) == 0 {
		return fmt.Errorf("affected rows equal 0: %v", err)
	}

	return nil
}

func (m *mysql) DeleteFriend(userId int, friendId int) error {
	query := queryMap[deleteFriend]
	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	if userId == friendId {
		return fmt.Errorf("user ID and friend ID are equal")
	}

	res, err := m.db.ExecContext(ctx, query.SQL, userId, friendId, friendId, userId)
	if err != nil {
		return fmt.Errorf("deleting friend with ID %d from user ID %d: %v", friendId, userId, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting afected rows: %v", err)
	}

	if int(rowsAffected) == 0 {
		return fmt.Errorf("affected rows equal 0: %v", err)
	}

	return nil
}

func (m *mysql) Friends(userId int) ([]User, error) {
	query := queryMap[getFriends]
	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query.SQL, userId)
	if err != nil {
		return nil, fmt.Errorf("getting friends: %v", err)
	}
	defer rows.Close()

	friends := []User{}

	for rows.Next() {
		var friend User
		var cityID sql.NullInt64
		var cityName sql.NullString

		err := rows.Scan(
			&friend.ID,
			&friend.FirstName,
			&friend.Lastname,
			&friend.Age,
			&friend.Sex,
			&friend.Login,
			&cityID,
			&cityName,
		)
		if err != nil {
			log.Printf("scanning friends: %v", err)
			continue
		}

		friend.City = city.City{}

		if cityName.Valid && cityID.Valid {
			friend.City.Name = cityName.String
			friend.City.ID = int(cityID.Int64)
		}

		friends = append(friends, friend)
	}

	return friends, nil
}
