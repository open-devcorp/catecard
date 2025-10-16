package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"errors"
	"log"
)

type AuthUseCase interface {
	Login(input LoginStruct) (*entities.User, error)
	SignUp(input SignupStruct) (*entities.User, error)
}

type authUseCase struct {
	log      *log.Logger
	userRepo repositories.UserRepository
}

type SignupStruct struct {
	// User info
	Username string        `json:"username"`
	Password string        `json:"password"`
	Role     entities.Role `json:"role"`
}
type LoginStruct struct {
	//User info
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserRepository interface {
	GetAll() ([]*entities.User, error)
	GetUser(username string, password string) *entities.User
	SaveUser(user *entities.User) (*entities.User, error)
}

func NewAuthUseCase(logger *log.Logger, r repositories.UserRepository) AuthUseCase {
	return &authUseCase{logger, r}
}

// SignUp implements AuthUseCase.
func (uc *authUseCase) SignUp(input SignupStruct) (*entities.User, error) {

	user := entities.NewUser(input.Username, input.Password, input.Role)

	savedUser, err := uc.userRepo.SaveUser(&user)
	if err != nil {
		uc.log.Printf("Error saving user: %v", err)
		return nil, err
	}

	uc.log.Printf("USER SAVED: %v", savedUser)
	return savedUser, nil
}

// UserLogin implements AuthUseCase.
func (uc *authUseCase) Login(input LoginStruct) (*entities.User, error) {
	user := uc.userRepo.GetUser(input.Username, input.Password)
	if user == nil {
		uc.log.Printf("Login failed for username: %s", input.Username)
		return nil, errors.New("invalid credentials")
	}

	uc.log.Printf("User logged in successfully: %s", user.Username)
	return user, nil
}
