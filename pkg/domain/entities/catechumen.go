package entities

type Catechumen struct {
	ID       int    `json:"id"`
	FullName string `json:"name"`
	Age      string `json:"age"`
	GroupId  int    `json:"group_id"`
}

func NewCatechumen(fullName string, age string, groupId int) Catechumen {
	return Catechumen{
		FullName: fullName,
		Age:      age,
		GroupId:  groupId,
	}
}

func FakeCatechumen() Catechumen {
	return Catechumen{
		ID:       1,
		FullName: "Italo D'Alessandro Luna Capuñay",
		Age:      "18",
		GroupId:  1,
	}
}
