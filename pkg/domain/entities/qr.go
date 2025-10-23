package entities

type Qr struct {
	ID      int `json:"id"`
	Forum   int `json:"forum"`
	GroupId int `json:"group_id"`
	Count   int `json:"count"`
}

func NewQr(forum int, groupId int) *Qr {
	return &Qr{
		Forum:   forum,
		GroupId: groupId,
		Count:   0,
	}
}

func FakeQr() Qr {
	return Qr{
		ID:      1,
		Forum:   2,
		GroupId: 1,
		Count:   5,
	}
}
