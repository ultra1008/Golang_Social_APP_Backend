package city

type City struct {
	ID            int    `db:"id"`
	Name          string `db:"city_name"`
	CreatedByUser bool   `db:"created_by_user"`
}
