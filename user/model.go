package user

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	Lastname  string `db:"last_name"`
	Age       int    `db:"age"`
	Sex       string `db:"sex"`
	City      string `db:"city"`
}
