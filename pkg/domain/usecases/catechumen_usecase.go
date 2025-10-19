package usecases

import (
	"catecard/pkg/domain/entities"
)

type CatechumenUseCase interface {
	Add(catechumen *entities.Catechumen) (*entities.Catechumen, error)
}
