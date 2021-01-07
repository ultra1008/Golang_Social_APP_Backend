package user

import (
	"github.com/niklod/highload-social-network/user/city"
)

type User struct {
	ID        int
	FirstName string
	Lastname  string
	Age       int
	Sex       string
	City      city.City
	Login     string
	Password  string
}

func (u *User) Sanitize() {
	u.Password = ""
}
