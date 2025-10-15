package entities

import (
	"goapp/pkg/domain/entities"
	"testing"
)

func TestProduct(t *testing.T) {

	product := entities.NewProduct("Paracetamol", 22.4)

	fakeProduct := entities.FakeProduct()

	if product.Name != fakeProduct.Name {
		t.Errorf("Expected name was to be %s", fakeProduct.Name)
	}

	if product.Price != fakeProduct.Price {
		t.Errorf("Expected price was to be %v", fakeProduct.Price)
	}

}
