package entities

import (
	"catecard/pkg/domain/entities"
	"testing"
)

func TestUser(t *testing.T) {

	fakeUser := entities.FakeUser()
	user := entities.NewUser(fakeUser.Username, fakeUser.Username, fakeUser.Password, fakeUser.Role)

	if user.Username != fakeUser.Username {
		t.Errorf("Expected Username was to be %s", fakeUser.Username)
	}
	if user.Password != fakeUser.Password {
		t.Errorf("Expected Password was to be %s", fakeUser.Password)
	}
	if user.Role != fakeUser.Role {
		t.Errorf("Expected Role was to be %d", fakeUser.Role)
	}

}
