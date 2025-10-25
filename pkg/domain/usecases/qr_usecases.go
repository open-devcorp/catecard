package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"errors"
	"fmt"
	"log"
)

type QrUseCase interface {
	GetAll() ([]*entities.Qr, error)
	GetById(id int) (*entities.Qr, error)
	ClaimQr(User *entities.User, qrId int) (*entities.Qr, error)
	GetAllScans(User *entities.User) ([]repositories.ScanAndCatechume, error)
}
type qrUseCase struct {
	logger         *log.Logger
	qrRepo         repositories.QrRepository
	catechumenRepo repositories.CatechumenRepository
	authRepo       repositories.UserRepository
	groupRepo      repositories.GroupRepository
	scanCateRepo   repositories.ScanCatechumenRepository
}

// GetAllScans implements QrUseCase.
func (q *qrUseCase) GetAllScans(User *entities.User) ([]repositories.ScanAndCatechume, error) {

	if User == nil {
		return nil, fmt.Errorf("unauthenticated: user required to get scans")
	}

	if User.Role != entities.ADMIN {
		return nil, fmt.Errorf("forbidden: only Admin")
	}

	scans, err := q.scanCateRepo.GetAll()
	if err != nil {
		q.logger.Printf("Error getting all scans: %v", err)
		return nil, err
	}

	return scans, nil

}

var ErrQrFull = errors.New("max participants reached")

func (q *qrUseCase) ClaimQr(User *entities.User, qrId int) (*entities.Qr, error) {

	if qrId == 0 {
		return nil, fmt.Errorf("QR ID cannot be 0")
	}
	//get qr
	qr, err := q.qrRepo.GetById(qrId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving QR: %w", err)
	}

	if qr.Catechumen == nil {
		qr.Catechumen = &entities.Catechumen{}
	}

	catechum, err := q.catechumenRepo.GetByQrId(qrId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Catechumen: %w", err)
	}

	// Get group and catechists
	group, err := q.groupRepo.GetById(catechum.GroupId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Group: %w", err)
	}

	catechist, err := q.authRepo.GetById(group.CatechistId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving User: %w", err)
	}

	//Set data
	qr.Catechumen.ID = catechum.ID
	qr.Catechumen.FullName = catechum.FullName
	qr.Catechumen.GroupId = catechum.GroupId
	qr.Catechumen.Group = group
	qr.Catechumen.User = catechist

	if qr.Forum <= 0 {

		return qr, ErrQrFull
	}

	// Prepare updates
	qr.Forum -= 1
	if qr.Forum < 0 {
		qr.Forum = 0 // hard guard
	}
	qr.Count += 1

	// Update QR
	if err := q.qrRepo.Update(qr); err != nil {
		return qr, fmt.Errorf("error updating QR: %w", err)
	}

	err = q.scanCateRepo.Add(entities.NewScanCatechumen(catechum.ID, User.ID))
	if err != nil {
		return qr, fmt.Errorf("error logging scan catechumen: %w", err)
	}

	return qr, nil
}

// func (q *qrUseCase) Add(User *entities.User, qr *entities.Qr) error {

// 	if User == nil {
// 		return fmt.Errorf("unauthenticated: user required to create QR")
// 	}

// 	if User.Role != entities.CATECHIST {
// 		return fmt.Errorf("forbidden: only catechist can create QR")
// 	}

// 	if qr.Forum == 0 {
// 		return fmt.Errorf("Forum cannot be empty")
// 	}

// 	if qr.GroupId == 0 {
// 		return fmt.Errorf("GroupId cannot be empty")
// 	}

// 	return q.qrRepo.Add(qr)

// }

// DeleteById implements QrUseCase.

// GetAll implements QrUseCase.
func (q *qrUseCase) GetAll() ([]*entities.Qr, error) {
	return q.qrRepo.GetAll()
}

// GetById implements QrUseCase.
func (q *qrUseCase) GetById(id int) (*entities.Qr, error) {
	return q.qrRepo.GetById(id)
}

func NewQrUsecase(logger *log.Logger, qrRepo repositories.QrRepository, catechum repositories.CatechumenRepository, authRepo repositories.UserRepository, groupRepo repositories.GroupRepository, scanCateRepo repositories.ScanCatechumenRepository) QrUseCase {
	return &qrUseCase{
		logger:         logger,
		qrRepo:         qrRepo,
		catechumenRepo: catechum,
		authRepo:       authRepo,
		groupRepo:      groupRepo,
		scanCateRepo:   scanCateRepo,
	}
}
