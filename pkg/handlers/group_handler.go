package handlers

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/domain/usecases"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type GroupHandler interface {
	AddGroup(User *entities.User, w http.ResponseWriter, r *http.Request)
	EditGroup(User *entities.User, w http.ResponseWriter, r *http.Request)
	GetAllGroups(User *entities.User, w http.ResponseWriter, r *http.Request)
}

type groupHandler struct {
	log      *log.Logger
	uc       usecases.GroupUseCase
	tmplPath string
}

func NewGroupHandler(logger *log.Logger, uc usecases.GroupUseCase, tmplPath string) GroupHandler {
	return &groupHandler{logger, uc, tmplPath}
}

func (g *groupHandler) AddGroup(User *entities.User, w http.ResponseWriter, r *http.Request) {
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
		g.log.Printf("Error parsing form: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	name := r.FormValue("name")
	catechistIdStr := r.FormValue("catechist_id")

	catechistId, err := strconv.Atoi(catechistIdStr)
	if err != nil {
		g.log.Printf("Error parsing catechist ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid catechist ID value")
		return
	}

	group := entities.Group{
		Name:        name,
		CatechistId: catechistId,
	}

	err = g.uc.Add(&group)
	if err != nil {
		g.log.Printf("Error adding group: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error adding group")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]interface{}{
		"status":  "ok",
		"message": "Group created",
		"group":   group,
	}

	json.NewEncoder(w).Encode(resp)

}

// EditGroup implements GroupHandler.
func (g *groupHandler) EditGroup(User *entities.User, w http.ResponseWriter, r *http.Request) {
	// Implementation for editing a group
}

// GetAllGroups implements GroupHandler.
func (g *groupHandler) GetAllGroups(User *entities.User, w http.ResponseWriter, r *http.Request) {
	// Implementation for retrieving all groups
}
