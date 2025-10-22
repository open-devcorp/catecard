package entities

type Qr struct {
	ID      int `json:"id"`
	Forum   int `json:"forum"`
	GroupId int `json:"group_id"`
}

func NewQr(forum int, groupId int) *Qr {
	return &Qr{
		Forum:   forum,
		GroupId: groupId,
	}
}

func FakeQr() Qr {
	return Qr{
		ID:      1,
		Forum:   2,
		GroupId: 1,
	}
}
