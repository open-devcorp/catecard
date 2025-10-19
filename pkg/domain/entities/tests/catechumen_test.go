package entities

import (
	"catecard/pkg/domain/entities"
	"testing"
)

func TestCatechumen(t *testing.T) {
	catechumen := entities.NewCatechumen("Italo D'Alessandro Luna Capuñay", "18", 1)
	fakeCatechumen := entities.FakeCatechumen()

	if catechumen.FullName != fakeCatechumen.FullName {
		t.Errorf("Expected fullname was to be %s", fakeCatechumen.FullName)
	}

	if catechumen.Age != fakeCatechumen.Age {
		t.Errorf("Expected age was to be %s", fakeCatechumen.Age)
	}

}
