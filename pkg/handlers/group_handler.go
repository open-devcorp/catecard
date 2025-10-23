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
	GetGroupById(User *entities.User, id int, w http.ResponseWriter, r *http.Request)
	DeleteGroupById(User *entities.User, id int, w http.ResponseWriter, r *http.Request)
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

	err = g.uc.Add(User, &group)
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

// GetGroupById implements GroupHandler.
func (g *groupHandler) GetGroupById(User *entities.User, id int, w http.ResponseWriter, r *http.Request) {
	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}

	group, err := g.uc.GetById(User, id)
	if err != nil {
		g.log.Printf("Error getting group by ID: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving group")
		return
	}
	if group == nil {
		writeJSONError(w, http.StatusNotFound, "Group not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"status": "ok",
		"group":  group,
	}

	json.NewEncoder(w).Encode(resp)

}

// EditGroup implements GroupHandler.
func (g *groupHandler) EditGroup(User *entities.User, w http.ResponseWriter, r *http.Request) {
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

	idStr := r.FormValue("id")
	name := r.FormValue("name")
	catechistIdStr := r.FormValue("catechist_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		g.log.Printf("Error parsing group ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid group ID value")
		return
	}
	catechistId, err := strconv.Atoi(catechistIdStr)
	if err != nil {
		g.log.Printf("Error parsing catechist ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid catechist ID value")
		return
	}

	group := &entities.Group{
		ID:          id,
		Name:        name,
		CatechistId: catechistId,
	}

	updatedGroup, err := g.uc.Update(User, group)
	if err != nil {
		g.log.Printf("Error updating group: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error updating group")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"status":  "ok",
		"message": "Group updated",
		"group":   updatedGroup,
	}

	json.NewEncoder(w).Encode(resp)
}

// GetAllGroups implements GroupHandler.
func (g *groupHandler) GetAllGroups(User *entities.User, w http.ResponseWriter, r *http.Request) {
	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}

	groups, err := g.uc.GetAll(User)
	if err != nil {
		g.log.Printf("Error getting groups: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving groups")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"status": "ok",
		"groups": groups,
	}

	json.NewEncoder(w).Encode(resp)

}

func (g *groupHandler) DeleteGroupById(User *entities.User, id int, w http.ResponseWriter, r *http.Request) {
	if User == nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: no user in session")
		return
	}
	if User.Role != entities.ADMIN {
		log.Printf("ROLE WAS: %v", User)
		writeJSONError(w, http.StatusForbidden, "Invalid role value")
		return
	}

	err := g.uc.DeleteById(User, id)
	if err != nil {
		g.log.Printf("Error deleting group: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error deleting group")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Group deleted",
	})
}
