package user

import "time"

const (
	createQuery = iota
	listQuery
	getByID
	getByLogin
	GetByFirstAndLastName
	addFriend
	getFriends
	deleteFriend
)

type Query struct {
	SQL     string
	Timeout time.Duration
}

var queryMap map[int]Query

func init() {
	queryMap = make(map[int]Query)

	queryMap[createQuery] = Query{
		SQL: `INSERT INTO users (first_name, last_name, age, sex, city_id, login, password)
				VALUES (?, ?, ?, ?, ?, ?, ?);`,
		Timeout: 10 * time.Second,
	}

	queryMap[listQuery] = Query{
		SQL: `SELECT u.id
			, u.first_name
			, u.last_name
			, u.age
			, u.sex
			, u.login
			, u.city_id
			, c.city_name
				FROM users as u
						LEFT JOIN citys as c ON u.city_id = c.id
				ORDER BY u.first_name`,
		Timeout: 10 * time.Second,
	}

	queryMap[getByID] = Query{
		SQL: `SELECT u.id
			, u.first_name
			, u.last_name
			, u.age
			, u.sex
			, u.login
			, u.city_id
			, c.city_name
			, u.password
				FROM users as u
						LEFT JOIN citys as c ON u.city_id = c.id
				ORDER BY u.first_name
				WHERE u.id = ?`,
		Timeout: 10 * time.Second,
	}

	queryMap[getByLogin] = Query{
		SQL: `SELECT u.id
			, u.first_name
			, u.last_name
			, u.age
			, u.sex
			, u.login
			, u.city_id
			, c.city_name
			, u.password
				FROM users as u
						LEFT JOIN citys as c ON u.city_id = c.id
				WHERE u.login = ?
				ORDER BY u.first_name`,
		Timeout: 10 * time.Second,
	}

	queryMap[GetByFirstAndLastName] = Query{
		SQL: `SELECT u.id
				, u.first_name
				, u.last_name
				, u.age
				, u.sex
				, u.login
				, u.city_id
				, c.city_name
			FROM users as u
					LEFT JOIN citys as c ON u.city_id = c.id
			WHERE u.first_name like ?
			AND u.last_name like ?
			ORDER BY u.id`,
		Timeout: 2 * time.Minute,
	}

	queryMap[addFriend] = Query{
		SQL:     `INSERT INTO friends (user_id, friend_id) VALUES (?, ?), (?, ?);`,
		Timeout: 10 * time.Second,
	}
	queryMap[deleteFriend] = Query{
		SQL:     `DELETE FROM friends WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)`,
		Timeout: 10 * time.Second,
	}
	queryMap[getFriends] = Query{
		SQL: `SELECT u.id
					, u.first_name
					, u.last_name
					, u.age
					, u.sex
					, u.login
					, u.city_id
					, c.city_name
				FROM friends f
					LEFT JOIN users u on f.friend_id = u.id
					LEFT JOIN citys c on u.city_id = c.id
				WHERE f.user_id = ?`,
		Timeout: 10 * time.Second,
	}
}
