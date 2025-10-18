package entities

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CatechistId int    `json:"catechist_id"`
}

func NewGroup(name string, catechistId int) Group {
	return Group{
		Name:        name,
		CatechistId: catechistId,
	}
}

func FakeGroup() Group {
	return Group{
		ID:          1,
		Name:        "Grupo A",
		CatechistId: 1,
	}
}
