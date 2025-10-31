package entities

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
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

func NewUser(username, fullName, password string, role Role) User {

	return User{
		Username: username,
		FullName: fullName,
		Password: password,
		Role:     role,
	}

}

func FakeUser() User {
	return User{
		ID:       1,
		Username: "rodrigo@devcorp.pe",
		FullName: "Rodrigo Devcorp",
		Password: "123456",
		Role:     ADMIN,
	}
}
