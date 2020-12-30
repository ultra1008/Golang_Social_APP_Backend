package city

import "time"

type SQLQuery struct {
	SQL     string
	Timeout time.Duration
}

const (
	createNewCity = iota
	listCitys
)

var queryMap map[int]SQLQuery

func init() {
	queryMap = make(map[int]SQLQuery)

	queryMap[createNewCity] = SQLQuery{
		SQL: `INSERT INTO citys (` + "`name`" + `) VALUES (':name')
				ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id);
				SELECT LAST_INSERT_ID();`,
		Timeout: 10 * time.Second,
	}

	queryMap[listCitys] = SQLQuery{
		SQL:     `SELECT id, name, created_by_user FROM citys`,
		Timeout: 10 * time.Second,
	}
}
