package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
)

type AuthUseCase interface {
	Login(input LoginStruct) (*entities.User, error)
	SignUp(input SignupStruct) (*entities.User, error)
	CreateAccounts(User *entities.User, input SignupStruct) (*entities.User, error)
	GetAllAccountsByRole(User *entities.User, role entities.Role) ([]*entities.User, error)
	GetUserById(User *entities.User, id int) (*entities.User, error)
	DeleteUserById(User *entities.User, id int) error
}

type authUseCase struct {
	log      *log.Logger
	userRepo repositories.UserRepository
}

type SignupStruct struct {
	// User info
	Username string        `json:"username"`
	Email    string        `json:"email"`
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

	user := entities.NewUser(input.Username, input.Email, input.Password, input.Role)

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

// CreateAccounts implements AuthUseCase.
func (uc *authUseCase) CreateAccounts(User *entities.User, input SignupStruct) (*entities.User, error) {

	if User.Role != entities.ADMIN {
		uc.log.Printf("only users with role Admin can create a new catechist")
		return nil, errors.New("only users with role Admin can create a new catechist")
	}

	// generate a secure random password
	const passLen = 12
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_+=<>?"
	b := make([]byte, passLen)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			uc.log.Printf("Error generating password: %v", err)
			return nil, err
		}
		b[i] = charset[n.Int64()]
	}
	randomPass := string(b)
	input.Password = randomPass

	user := entities.NewUser(input.Username, input.Email, input.Password, input.Role)

	catechist, err := uc.userRepo.SaveUser(&user)
	if err != nil {
		uc.log.Printf("Error saving user: %v", err)
		return nil, err
	}

	uc.log.Printf("Generated password for user %s: %s", user.Username, randomPass)

	return catechist, nil

}

func (uc *authUseCase) GetAllAccountsByRole(User *entities.User, role entities.Role) ([]*entities.User, error) {
	if User.Role != entities.ADMIN {
		uc.log.Printf("only users with role Admin can view catechists")
		return nil, errors.New("only users with role Admin can view catechists")
	}

	allUsers, err := uc.userRepo.GetAll()
	if err != nil {
		uc.log.Printf("Error retrieving users: %v", err)
		return nil, err
	}

	var filteredUsers []*entities.User
	for _, u := range allUsers {
		if u.Role == role {
			filteredUsers = append(filteredUsers, u)
		}
	}

	return filteredUsers, nil

}

func (uc *authUseCase) GetUserById(User *entities.User, id int) (*entities.User, error) {

	if User.Role != entities.ADMIN {
		uc.log.Printf("only users with role Admin can view user details")
		return nil, errors.New("only users with role Admin can view user details")
	}

	allUsers, err := uc.userRepo.GetAll()
	if err != nil {
		uc.log.Printf("Error retrieving users: %v", err)
		return nil, err
	}

	for _, u := range allUsers {
		if u.ID == id {
			return u, nil
		}
	}

	uc.log.Printf("User with ID %d not found", id)
	return nil, errors.New("user not found")

}

func (uc *authUseCase) DeleteUserById(User *entities.User, id int) error {

	if User.Role != entities.ADMIN {
		uc.log.Printf("only users with role Admin can delete user accounts")
		return errors.New("only users with role Admin can delete user accounts")
	}

	err := uc.userRepo.DeleteUserById(id)
	if err != nil {
		uc.log.Printf("Error deleting user with ID %d: %v", id, err)

		return err
	}

	uc.log.Printf("User with ID %d deleted successfully", id)
	return nil
}
