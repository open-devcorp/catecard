package handlers

import (
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
	if err := t.ExecuteTemplate(w, baseName, data); err != nil {
		http.Error(w, "Error al ejecutar plantilla"+err.Error(), http.StatusInternalServerError)
	}

}
func AddProductView(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "add_product.html", nil)
}

func ProductsView(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "products.html", nil)
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

func Catechist(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "catechist.html", nil)
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
