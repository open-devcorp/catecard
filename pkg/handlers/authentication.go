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
	CreateCatechist(User *entities.User, w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	log      *log.Logger
	uc       usecases.AuthUseCase
	tmplPath string
}

// CreateCatechist implements AuthHandler.
func (a *authHandler) CreateCatechist(User *entities.User, w http.ResponseWriter, r *http.Request) {
	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}
	if User.Role != entities.ADMIN {
		log.Printf("ROLE WAS: %v", User)
		writeJSONError(w, http.StatusForbidden, "Invalid role value")
		return
	}

	signupData := usecases.SignupStruct{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Role:     entities.CATECHIST,
	}

	user, err := a.uc.CreateCatechist(User, signupData)
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

func NewAuthenticationHandler(logger *log.Logger, uc usecases.AuthUseCase, tmpl string) AuthHandler {
	return &authHandler{logger, uc, tmpl}
}
