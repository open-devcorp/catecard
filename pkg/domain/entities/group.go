package entities

type Group struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	CatechistId      int    `json:"catechist_id"`
	User             *User  `json:"user,omitempty"`
	LimitCatechumens int    `json:"limit_catechumens"`
}

func NewGroup(name string, catechistId int, limiteCatechumen int) Group {
	return Group{
		Name:             name,
		CatechistId:      catechistId,
		LimitCatechumens: limiteCatechumen,
	}
}

func FakeGroup() Group {
	return Group{
		ID:          1,
		Name:        "Grupo A",
		CatechistId: 1,
	}
}
