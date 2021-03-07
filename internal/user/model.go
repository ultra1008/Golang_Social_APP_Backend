package user

import (
	"github.com/niklod/highload-social-network/internal/user/city"
	"github.com/niklod/highload-social-network/internal/user/interest"
	"github.com/niklod/highload-social-network/internal/user/post"
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
	Posts     []post.Post
}

func (u *User) Sanitize() {
	u.Password = ""
}
