package interest

import "time"

type Query struct {
	SQL     string
	Timeout time.Duration
}

const (
	createIfNotExists = iota
	listInterests
	getUserInterests
	addInterestToUser
)

var queryMap map[int]Query

func init() {
	queryMap = map[int]Query{}

	queryMap[createIfNotExists] = Query{
		SQL: `INSERT INTO interests (` + "`name`" + `) VALUES (?)
				ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id);`,
		Timeout: 10 * time.Second,
	}

	queryMap[listInterests] = Query{
		SQL: `SELECT id
			, name
			FROM interests
			ORDER BY name`,
		Timeout: 10 * time.Second,
	}

	queryMap[getUserInterests] = Query{
		SQL: `SELECT ui.interest_id
				, i.name
				FROM user_interests ui
				LEFT JOIN interests i on ui.interest_id = i.id
				WHERE ui.user_id = ?
				ORDER BY name`,
		Timeout: 10 * time.Second,
	}

	queryMap[addInterestToUser] = Query{
		SQL:     `INSERT INTO user_interests (` + "`user_id`, `interest_id`" + `) VALUES (?, ?)`,
		Timeout: 10 * time.Second,
	}
}
