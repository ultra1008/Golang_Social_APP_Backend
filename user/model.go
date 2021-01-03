package user

import (
	"github.com/niklod/highload-social-network/user/city"
)

type User struct {
	ID        int
	FirstName string
	Lastname  string
	Age       int
	Sex       Sex
	City      city.City
	Login     string
	Password  string
}

type Sex struct {
	ID   int
	Name string
}
