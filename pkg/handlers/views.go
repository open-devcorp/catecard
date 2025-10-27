package handlers

import (
	"catecard/pkg/infrastructure/repositories"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var templates *template.Template

func LoadTemplates(dir string) error {

	var err error
	templates, err = template.ParseGlob(dir + "/*.html")
	if err != nil {
		return err
	}
	return nil
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Could not get working dir:%v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return

	}
	basePath := filepath.Join(wd, "pkg", "infrastructure", "web", "templates", "base.html")
	viewPath := filepath.Join(wd, "pkg", "infrastructure", "web", "templates", name) //html

	t, err := template.ParseFiles(basePath, viewPath)
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "Error cargar plantilla"+err.Error(), http.StatusInternalServerError)
		return
	}

	baseName := filepath.Base(basePath)
	// Ensure data is a map so we can inject a ShowNav flag for templates.
	var dataMap map[string]interface{}
	if data == nil {
		dataMap = map[string]interface{}{}
		data = dataMap
	} else {
		if dm, ok := data.(map[string]interface{}); ok {
			dataMap = dm
		} else {
			// wrap non-map data so existing templates that expect a map still work
			dataMap = map[string]interface{}{"Payload": data}
			data = dataMap
		}
	}

	// By default show nav; hide it for the login view
	if _, exists := dataMap["ShowNav"]; !exists {
		if name == "login.html" {
			dataMap["ShowNav"] = false
		} else {
			dataMap["ShowNav"] = true
		}
	}

	if err := t.ExecuteTemplate(w, baseName, data); err != nil {
		http.Error(w, "Error al ejecutar plantilla"+err.Error(), http.StatusInternalServerError)
	}

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "signup.html", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "login.html", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "home.html", nil)
}

var GroupRepo repositories.GroupRepository // Debe ser inicializado en el main/web.go

func Catechist(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
	var group interface{} = nil
	if user != nil && GroupRepo != nil {
		groupId, err := GroupRepo.GetByCatechistsId(user.ID)
		if err == nil && groupId != 0 {
			g, err := GroupRepo.GetById(groupId)
			if err == nil && g != nil {
				group = g
			}
		}
	}
	RenderTemplate(w, "catechist.html", map[string]interface{}{
		"User":  user,
		"Group": group,
	})
}

func Denied(w http.ResponseWriter, r *http.Request) {
	// read optional message/data from query params and pass to template
	msg := r.URL.Query().Get("message")
	data := r.URL.Query().Get("data")
	payload := map[string]string{"Message": msg, "Data": data}
	RenderTemplate(w, "denied.html", payload)
}

func Success(w http.ResponseWriter, r *http.Request) {
	// read optional message/data from query params and pass to template
	msg := r.URL.Query().Get("message")
	data := r.URL.Query().Get("data")
	payload := map[string]string{"Message": msg, "Data": data}
	RenderTemplate(w, "success.html", payload)
}

func Scanner(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "scanner.html", nil)
}

func QrList(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "qr.html", nil)
}
