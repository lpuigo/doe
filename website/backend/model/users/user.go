package users

type User struct {
	Id          int
	Name        string
	Password    string
	Clients     []string
	Permissions map[string]bool
}

func NewUser(name string) *User {
	return &User{
		Id:          0,
		Name:        name,
		Password:    "",
		Clients:     []string{},
		Permissions: make(map[string]bool),
	}
}

func (u *User) HasPermissionHR() bool {
	return u.Permissions["HR"]
}
