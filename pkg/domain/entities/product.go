package entities

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func NewProduct(name string, price float32) Product {

	return Product{
		Name:  name,
		Price: price,
	}
}

func FakeProduct() Product {
	return Product{
		ID:    1,
		Name:  "Paracetamol",
		Price: 22.4,
	}
}
