package entities

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type Role int

const (
	ADMIN Role = iota
	CATECHIST
	CATECHUMEN
	SCANNER
)

func NewUser(username, password string, role Role) User {

	return User{
		Username: username,
		Password: password,
		Role:     role,
	}

}

func FakeUser() User {
	return User{
		Username: "rodrigo@devcorp.pe",
		Password: "123456",
		Role:     ADMIN,
	}
}
