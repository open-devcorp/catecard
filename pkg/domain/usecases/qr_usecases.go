package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"fmt"
	"log"
)

type QrUseCase interface {
	Add(User *entities.User, qr *entities.Qr) error
	GetAll() ([]*entities.Qr, error)
	GetById(id int) (*entities.Qr, error)
	DeleteById(id int) error
	ClaimQr(qrId int) (*entities.Qr, error)
}
type qrUseCase struct {
	logger *log.Logger
	qrRepo repositories.QrRepository
}

// ClaimQr implements QrUseCase.
func (q *qrUseCase) ClaimQr(qrId int) (*entities.Qr, error) {

	if qrId == 0 {
		return nil, fmt.Errorf("QR ID cannot be 0")
	}
	qr, err := q.qrRepo.GetById(qrId)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving QR: %v", err)
	}
	if qr.Forum != 0 {
		qr.Forum -= 1
	}
	err = q.qrRepo.Update(qr)
	if err != nil {
		return qr, err
	}

	return qr, nil

}

func (q *qrUseCase) Add(User *entities.User, qr *entities.Qr) error {

	if User == nil {
		return fmt.Errorf("unauthenticated: user required to create QR")
	}

	if User.Role != entities.CATECHIST {
		return fmt.Errorf("forbidden: only catechist can create QR")
	}

	if qr.Forum == 0 {
		return fmt.Errorf("Forum cannot be empty")
	}

	if qr.GroupId == 0 {
		return fmt.Errorf("GroupId cannot be empty")
	}

	return q.qrRepo.Add(qr)

}

// DeleteById implements QrUseCase.
func (q *qrUseCase) DeleteById(id int) error {
	panic("unimplemented")
}

// GetAll implements QrUseCase.
func (q *qrUseCase) GetAll() ([]*entities.Qr, error) {
	return q.qrRepo.GetAll()
}

// GetById implements QrUseCase.
func (q *qrUseCase) GetById(id int) (*entities.Qr, error) {
	return q.qrRepo.GetById(id)
}

func NewQrUsecase(logger *log.Logger, qrRepo repositories.QrRepository) QrUseCase {
	return &qrUseCase{
		logger: logger,
		qrRepo: qrRepo,
	}
}
