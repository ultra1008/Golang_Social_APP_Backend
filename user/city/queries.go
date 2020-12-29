package city

import "time"

type SQLQuery struct {
	SQL     string
	Timeout time.Duration
}

const (
	createNewCity = iota
)

var queryMap map[int]SQLQuery

func init() {
	queryMap[createNewCity] = SQLQuery{
		SQL: `INSERT INTO citys (` + "`name`" + `) VALUES ('Екатеринбург')
				ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id);
				SELECT LAST_INSERT_ID();`,
		Timeout: 10 * time.Second,
	}
}
