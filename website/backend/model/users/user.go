package users

type User struct {
	Id          int
	Name        string
	Password    string
	Clients     []string
	Groups      []int
	Permissions map[string]bool
}

func NewUser(name string) *User {
	return &User{
		Id:          0,
		Name:        name,
		Password:    "",
		Clients:     []string{},
		Groups:      []int{},
		Permissions: make(map[string]bool),
	}
}

func (u *User) HasPermissionHR() bool {
	return u.Permissions["HR"]
}

// IsSeeingAllGroups returns true if user has unrestricted Groups access
func (u *User) IsSeeingAllGroups() bool {
	return len(u.Groups) == 0
}
