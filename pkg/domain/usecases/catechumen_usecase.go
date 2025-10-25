package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"fmt"
	"log"
)

type CatechumenUseCase interface {
	Add(User *entities.User, catechumen *entities.Catechumen) (*entities.Catechumen, *entities.Qr, error)
	Update(User *entities.User, catechumen *entities.Catechumen) (*entities.Catechumen, error)
	GetAll(User *entities.User) ([]*entities.Catechumen, error)
	GetById(User *entities.User, id int) (*entities.Catechumen, error)
}

type catechumenUseCase struct {
	logger         *log.Logger
	catechumenRepo repositories.CatechumenRepository
	groupRepo      repositories.GroupRepository
	qrRepo         repositories.QrRepository
}

func NewCatechumenUsecase(logger *log.Logger, catechumenRepo repositories.CatechumenRepository, groupRepo repositories.GroupRepository, qrRepo repositories.QrRepository) CatechumenUseCase {
	return &catechumenUseCase{logger, catechumenRepo, groupRepo, qrRepo}
}

func (c *catechumenUseCase) Add(User *entities.User, catechumen *entities.Catechumen) (*entities.Catechumen, *entities.Qr, error) {
	if User == nil {
		return nil, nil, fmt.Errorf("Unauthorized: User is nil")
	}
	if User.Role != entities.CATECHIST {
		return nil, nil, fmt.Errorf("Unauthorized: User is not a catechist")
	}
	if catechumen.FullName == "" {
		return nil, nil, fmt.Errorf("Catechumen name cannot be empty")
	}
	if catechumen.Age == "" {
		return nil, nil, fmt.Errorf("Catechumen age cannot be empty")
	}

	groupId, err := c.groupRepo.GetByCatechistsId(User.ID)

	if err != nil {

		return nil, nil, fmt.Errorf("Error fetching group by catechist ID: %v", err)
	}

	if groupId == 0 {
		return nil, nil, fmt.Errorf("No group found for catechist ID: %d", User.ID)
	}

	catechumen.GroupId = groupId

	qr := entities.NewQr(3, groupId)

	qrId, qrErr := c.qrRepo.Add(qr)
	if qrErr != nil {
		c.logger.Printf("Error adding QR: %v", qrErr)
		return nil, nil, qrErr
	}
	catechumen.QrId = qrId

	cateId, err := c.catechumenRepo.Add(catechumen)
	if err != nil {
		c.logger.Printf("Error adding catechumen: %v", err)
		return nil, nil, err
	}
	catechumen.ID = cateId

	return catechumen, qr, nil
}

func (c *catechumenUseCase) Update(User *entities.User, catechumen *entities.Catechumen) (*entities.Catechumen, error) {
	if User == nil {
		return nil, fmt.Errorf("Unauthorized: User is nil")
	}
	if User.Role != entities.CATECHIST {
		return nil, fmt.Errorf("Unauthorized: User is not a catechist")
	}
	if User.ID == catechumen.ID {
		return nil, fmt.Errorf("Unauthorized: Cannot update own catechumen record")
	}

	if catechumen.FullName == "" {
		return nil, fmt.Errorf("Group name cannot be empty")
	}
	if catechumen.Age == "" {
		return nil, fmt.Errorf("Catechumen age cannot be empty")
	}
	if catechumen.GroupId == 0 {
		return nil, fmt.Errorf("Group ID cannot be 0")
	}

	updatedCatechumen, err := c.catechumenRepo.Update(catechumen)
	if err != nil {
		c.logger.Printf("Error updating catechumen: %v", err)
		return nil, err
	}

	return updatedCatechumen, nil
}

func (c *catechumenUseCase) GetAll(user *entities.User) ([]*entities.Catechumen, error) {
	if user == nil {
		return nil, fmt.Errorf("unauthorized: user is nil")
	}

	if user.Role != entities.CATECHIST {
		return nil, fmt.Errorf("unauthorized: user is not a catechist")
	}

	catechumens, err := c.catechumenRepo.GetAll()
	if err != nil {
		c.logger.Printf("Error getting all catechumens: %v", err)
		return nil, err
	}

	catechumensByUser := make([]*entities.Catechumen, 0)

	for _, catechumen := range catechumens {
		g, err := c.groupRepo.GetById(catechumen.GroupId)
		if err != nil {
			c.logger.Printf("Error getting group for catechumen %d: %v", catechumen.ID, err)
			return nil, err
		}

		if g.CatechistId == user.ID {
			catechumen.User = user
			catechumen.Group = g
			catechumensByUser = append(catechumensByUser, catechumen)
		}
	}

	return catechumensByUser, nil
}

func (c *catechumenUseCase) GetById(user *entities.User, id int) (*entities.Catechumen, error) {
	if user == nil {
		return nil, fmt.Errorf("unauthorized: user is nil")
	}
	if user.Role != entities.CATECHIST {
		return nil, fmt.Errorf("unauthorized: user is not a catechist")
	}
	if id == 0 {
		return nil, fmt.Errorf("invalid catechumen ID")
	}

	catechumen, err := c.catechumenRepo.GetById(id)
	if err != nil {
		c.logger.Printf("Error getting catechumen by ID: %v", err)
		return nil, err
	}

	group, err := c.groupRepo.GetById(catechumen.GroupId)
	if err != nil {
		c.logger.Printf("Error getting group for catechumen %d: %v", id, err)
		return nil, err
	}

	if group.CatechistId != user.ID {
		return nil, fmt.Errorf("unauthorized: cannot access this catechumen")
	}

	catechumen.User = user
	catechumen.Group = group

	return catechumen, nil
}
