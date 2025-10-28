package handlers

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/domain/usecases"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	CreateAccounts(User *entities.User, w http.ResponseWriter, r *http.Request)
	GetAllCatechists(User *entities.User, w http.ResponseWriter, r *http.Request)
	GetAllScanners(User *entities.User, w http.ResponseWriter, r *http.Request)
	GetUserById(User *entities.User, id int, w http.ResponseWriter, r *http.Request)
	DeleteUserById(User *entities.User, id int, w http.ResponseWriter, r *http.Request)
	GetAllCatechistsWithoutGroup(User *entities.User, id int, w http.ResponseWriter, r *http.Request) ([]*entities.User, error)
}

type authHandler struct {
	log      *log.Logger
	uc       usecases.AuthUseCase
	tmplPath string
}

// GetAllCatechistsWithoutGroup  AuthHandler.
func (a *authHandler) GetAllCatechistsWithoutGroup(User *entities.User, id int, w http.ResponseWriter, r *http.Request) ([]*entities.User, error) {
	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return nil, nil
	}
	if User.Role != entities.ADMIN {
		writeJSONError(w, http.StatusForbidden, "Invalid role value")
		return nil, nil
	}

	catechists, err := a.uc.GetAllCatechistsWithoutGroup()
	if err != nil {
		a.log.Printf("Error retrieving catechists without group: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving catechists without group")
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(catechists)
	return catechists, nil

}

// CreateCatechist implements AuthHandler.
func (a *authHandler) CreateAccounts(User *entities.User, w http.ResponseWriter, r *http.Request) {
	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}
	if User.Role != entities.ADMIN {
		log.Printf("ROLE WAS: %v", User)
		writeJSONError(w, http.StatusForbidden, "Invalid role value")
		return
	}

	if err := r.ParseForm(); err != nil {
		a.log.Printf("Error parsing form: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	roleStr := r.FormValue("role")
	roleInt, err := strconv.Atoi(roleStr)
	if err != nil {
		a.log.Printf("Invalid role value: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid role value")
		return
	}

	signupData := usecases.SignupStruct{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Role:     entities.Role(roleInt),
	}

	user, err := a.uc.CreateAccounts(User, signupData)
	if err != nil {
		a.log.Printf("Registration failed: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Registration failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]interface{}{
		"status":  "ok",
		"message": "catechist created",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"password": user.Password,
			"role":     user.Role,
		},
	}

	json.NewEncoder(w).Encode(resp)

}

func (a *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	roleStr := r.FormValue("role")
	roleInt, err := strconv.Atoi(roleStr)
	if err != nil {
		a.log.Printf("Invalid role value: %v", err)
		http.Error(w, "Invalid role value", http.StatusBadRequest)
		return
	}

	signupData := usecases.SignupStruct{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
		Role:     entities.Role(roleInt),
	}

	user, err := a.uc.SignUp(signupData)
	if err != nil {
		a.log.Printf("Registration failed: %v", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	log.Printf("USER WAS ADDED %v", user)
	RenderTemplate(w, "login.html", nil)
}

func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	loginData := usecases.LoginStruct{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	user, err := a.uc.Login(loginData)
	if err != nil {
		a.log.Printf("Login failed: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// crear sesión
	if _, err := CreateSession(w, user); err != nil {
		a.log.Printf("Error creating session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	a.log.Printf("USER WAS LOGGED %v", user)
	http.Redirect(w, r, "/home", http.StatusSeeOther)

}

// Logout implements AuthHandler.
func (a *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// borrar sesión
	DeleteSession(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Logout successful"}
	json.NewEncoder(w).Encode(response)
}

// writeJSONError escribe un error JSON con estructura {status: "error", message: "..."}
func writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

func (a *authHandler) GetAllCatechists(User *entities.User, w http.ResponseWriter, r *http.Request) {

	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}
	if User.Role != entities.ADMIN {
		log.Printf("ROLE WAS: %v", User)
		writeJSONError(w, http.StatusForbidden, "Invalid role value")
		return
	}

	catechists, err := a.uc.GetAllAccountsByRole(User, entities.CATECHIST)
	if err != nil {
		a.log.Printf("Error retrieving catechists: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving catechists")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(catechists)

}

// GetAllScanners implements AuthHandler.
func (a *authHandler) GetAllScanners(User *entities.User, w http.ResponseWriter, r *http.Request) {

	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}
	if User.Role != entities.ADMIN {
		log.Printf("ROLE WAS: %v", User)
		writeJSONError(w, http.StatusForbidden, "Invalid role value")
		return
	}

	scanners, err := a.uc.GetAllAccountsByRole(User, entities.SCANNER)
	if err != nil {
		a.log.Printf("Error retrieving scanners: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving scanners")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scanners)

}

func (a *authHandler) GetUserById(User *entities.User, id int, w http.ResponseWriter, r *http.Request) {

	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}
	if User.Role != entities.ADMIN {
		log.Printf("ROLE WAS: %v", User)
		writeJSONError(w, http.StatusForbidden, "Invalid role value")
		return
	}

	user, err := a.uc.GetUserById(User, id)
	if err != nil {
		a.log.Printf("Error retrieving user: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func (a *authHandler) DeleteUserById(User *entities.User, id int, w http.ResponseWriter, r *http.Request) {

	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}
	if User.Role != entities.ADMIN {
		log.Printf("ROLE WAS: %v", User)
		writeJSONError(w, http.StatusForbidden, "Invalid role value")
		return
	}

	err := a.uc.DeleteUserById(User, id)
	if err != nil {
		a.log.Printf("Error deleting user: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error deleting user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "User deleted successfully",
	})

}

func NewAuthenticationHandler(logger *log.Logger, uc usecases.AuthUseCase, tmpl string) AuthHandler {
	return &authHandler{logger, uc, tmpl}
}
