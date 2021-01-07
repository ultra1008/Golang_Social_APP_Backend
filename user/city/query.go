package city

import "time"

type SQLQuery struct {
	SQL     string
	Timeout time.Duration
}

const (
	createNewCity = iota
	listCitys
	getByID
)

var queryMap map[int]SQLQuery

func init() {
	queryMap = make(map[int]SQLQuery)

	queryMap[createNewCity] = SQLQuery{
		SQL: `INSERT INTO citys (` + "`city_name`, " + "`created_by_user`" + `) VALUES (?, 1)
				ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id);`,
		Timeout: 10 * time.Second,
	}

	queryMap[listCitys] = SQLQuery{
		SQL:     `SELECT id, city_name, created_by_user FROM citys WHERE created_by_user = 0 ORDER BY city_name`,
		Timeout: 10 * time.Second,
	}

	queryMap[getByID] = SQLQuery{
		SQL:     `SELECT id, city_name, created_by_user FROM citys WHERE id = ?`,
		Timeout: 10 * time.Second,
	}
}
