package entities

type Catechumen struct {
	ID       int    `json:"id"`
	FullName string `json:"name"`
	Age      string `json:"age"`
	GroupId  int    `json:"group_id"`
	User     *User  `json:"user,omitempty"`
	Group    *Group `json:"group,omitempty"`
}

func NewCatechumen(fullName string, age string) Catechumen {
	return Catechumen{
		FullName: fullName,
		Age:      age,
	}
}

func FakeCatechumen() Catechumen {
	return Catechumen{
		ID:       1,
		FullName: "Italo D'Alessandro Luna Capuñay",
		Age:      "18",
		GroupId:  1,
		User:     &User{},
	}
}
