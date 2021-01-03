package user

type repository interface {
	Create(user *User) (*User, error)
	List() ([]User, error)
	GetByID(id int) (*User, error)
	GetByLogin(login string) (*User, error)
}
