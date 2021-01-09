package user

import (
	"github.com/niklod/highload-social-network/user/city"
	"github.com/niklod/highload-social-network/user/interest"
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
	Friends   []User
	Interests []interest.Interest
}

func (u *User) Sanitize() {
	u.Password = ""
}
