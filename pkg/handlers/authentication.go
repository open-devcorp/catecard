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
}

type authHandler struct {
	log      *log.Logger
	uc       usecases.AuthUseCase
	tmplPath string
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

	log.Printf("USER WAS LOGGED %v", user)
	RenderTemplate(w, "home.html", nil)

}

// Logout implements AuthHandler.
func (a *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Logout successful",
	}

	json.NewEncoder(w).Encode(response)
}

func NewAuthenticationHandler(logger *log.Logger, uc usecases.AuthUseCase, tmpl string) AuthHandler {
	return &authHandler{logger, uc, tmpl}
}
