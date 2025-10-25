package entities

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
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

func NewUser(username, email, password string, role Role) User {

	return User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     role,
	}

}

func FakeUser() User {
	return User{
		ID:       1,
		Username: "rodrigo@devcorp.pe",
		Email:    "rodrigo@devcorp.pe",
		Password: "123456",
		Role:     ADMIN,
	}
}
