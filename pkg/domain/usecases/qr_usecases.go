package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"errors"
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
	logger         *log.Logger
	qrRepo         repositories.QrRepository
	catechumenRepo repositories.CatechumenRepository
	authRepo       repositories.UserRepository
	groupRepo      repositories.GroupRepository
}

// ClaimQr implements QrUseCase.
var ErrQrFull = errors.New("max participants reached") // ← sentinela

// ClaimQr implements QrUseCase.
func (q *qrUseCase) ClaimQr(qrId int) (*entities.Qr, error) {
	if qrId == 0 {
		return nil, fmt.Errorf("QR ID cannot be 0")
	}

	// 1) Trae el QR
	qr, err := q.qrRepo.GetById(qrId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving QR: %w", err)
	}
	if qr == nil {
		return nil, fmt.Errorf("QR not found for ID: %d", qrId)
	}

	// 2) Nil-safety para el contenedor
	if qr.Catechumen == nil {
		qr.Catechumen = &entities.Catechumen{}
	}

	// 3) Trae catecúmeno (idealmente por FK guardada en el QR)
	var catechumenID int
	switch {
	case qr.Catechumen.ID != 0:
		catechumenID = qr.Catechumen.ID
	case qr.Catechumen.ID != 0:
		catechumenID = qr.Catechumen.ID
	default:
		// ⚠️ Evita asumir que catechumenID == qrId si tu modelo NO lo define así.
		catechumenID = qrId
	}

	catechum, err := q.catechumenRepo.GetById(catechumenID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Catechumen: %w", err)
	}
	if catechum == nil {
		return nil, fmt.Errorf("no catechumen found for ID: %d", catechumenID)
	}
	log.Printf("[ClaimQr] catechumen: %#v", catechum)

	// 4) Grupo y catequista
	group, err := q.groupRepo.GetById(catechum.GroupId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Group: %w", err)
	}
	if group == nil {
		return nil, fmt.Errorf("no group found for Group ID: %d", catechum.GroupId)
	}
	log.Printf("[ClaimQr] group: %#v", group)

	catechist, err := q.authRepo.GetById(group.CatechistId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving User: %w", err)
	}
	if catechist == nil {
		return nil, fmt.Errorf("no user found for User ID: %d", group.CatechistId)
	}
	log.Printf("[ClaimQr] catechist: %#v", catechist)

	// 5) Reglas de negocio: validar cupo ANTES de decrementar
	// Si 'Forum' representa INVITADOS RESTANTES, negar solo cuando ya no quede (<= 0).
	// Si además tienes un límite total (ej. qr.Limit), valida también Count < Limit.
	// 6) Poblamos relaciones (usa los objetos que acabas de cargar)
	qr.Catechumen.ID = catechum.ID
	qr.Catechumen.FullName = catechum.FullName // ajusta a tus campos
	qr.Catechumen.GroupId = catechum.GroupId
	qr.Catechumen.Group = group
	qr.Catechumen.User = catechist
	if qr.Forum <= 0 {
		return qr, ErrQrFull
	}
	// Ejemplo adicional si existe 'Limit':
	// if qr.Limit > 0 && qr.Count >= qr.Limit {
	//     return qr, ErrQrFull
	// }

	// 7) Actualiza contadores (después de validar)
	qr.Forum -= 1     // ahora sí decrementa
	if qr.Forum < 0 { // por seguridad, no permitas negativos
		qr.Forum = 0
	}
	qr.Count += 1

	if err := q.qrRepo.Update(qr); err != nil {
		return qr, fmt.Errorf("error updating QR: %w", err)
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

func NewQrUsecase(logger *log.Logger, qrRepo repositories.QrRepository, catechum repositories.CatechumenRepository, authRepo repositories.UserRepository, groupRepo repositories.GroupRepository) QrUseCase {
	return &qrUseCase{
		logger:         logger,
		qrRepo:         qrRepo,
		catechumenRepo: catechum,
		authRepo:       authRepo,
		groupRepo:      groupRepo,
	}
}
