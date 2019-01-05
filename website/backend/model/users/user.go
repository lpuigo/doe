package users

type User struct {
	Id       int
	Name     string
	Password string
}

func NewUser(name string) *User {
	return &User{Name: name}
}
