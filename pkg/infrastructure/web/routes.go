package web

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/handlers"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func setupRouter(qrHandler handlers.QrHandler, authHandler handlers.AuthHandler, productHandler handlers.ProductHandler, groupHandler handlers.GroupHandler) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

	//////////////////////////VIEWS//////////////////////////
	r.HandleFunc("/add-product", handlers.AddProductView).Methods("GET")
	r.HandleFunc("/products", productHandler.GetAllProducts).Methods("GET")

	////// HOME REDIRECTION //////
	r.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if user.Role == entities.CATECHIST {
			handlers.Catechist(w, r)
			return
		}
		handlers.Home(w, r)

	}).Methods("GET")

	//AUTH
	r.HandleFunc("/signup", handlers.SignUp).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("GET")

	//////////////////////////APIS//////////////////////////////

	//AUTH
	r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	//////////// ADMIN /////////////
	r.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		authHandler.GetUserById(user, id, w, r)
	}).Methods("GET")
	//Catechists
	r.HandleFunc("/all-catechists", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.GetAllCatechists(user, w, r)
	}).Methods("GET")

	/////Groups
	r.HandleFunc("/all-groups", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		groupHandler.GetAllGroups(user, w, r)
	}).Methods("GET")
	r.HandleFunc("/group/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		groupHandler.GetGroupById(user, id, w, r)
	}).Methods("GET")

	r.HandleFunc("/delete-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		authHandler.DeleteUserById(user, id, w, r)
	}).Methods("DELETE")
	////////////////////////////////////////////////////////////
	//Scanners
	r.HandleFunc("/all-scanners", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.GetAllScanners(user, w, r)
	}).Methods("GET")

	//CATECHIST
	r.HandleFunc("/add-catechist", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.CreateAccounts(user, w, r)
	}).Methods("POST")

	///////////////////////////////GROUPS//////////////////////////////////////////////////
	r.HandleFunc("/add-group", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		groupHandler.AddGroup(user, w, r)
	}).Methods("POST")

	// DELETE GROUP
	r.HandleFunc("/delete-group/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		groupHandler.DeleteGroupById(user, id, w, r)
	}).Methods("DELETE")

	// EDIT GROUP
	r.HandleFunc("/edit-group/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		vars := mux.Vars(r)
		idStr := vars["id"]
		_, _ = strconv.Atoi(idStr)

		groupHandler.EditGroup(user, w, r)
	}).Methods("PUT")

	////////SCARNER///////////
	r.HandleFunc("/all-scanners", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.GetAllScanners(user, w, r)
	}).Methods("GET")

	//CATECHIST
	r.HandleFunc("/add-qr", func(w http.ResponseWriter, r *http.Request) {
		qrHandler.AddQr(w, r)
	}).Methods("POST")
	r.HandleFunc("/all-qr", func(w http.ResponseWriter, r *http.Request) {
		qrHandler.GetAllQrs(w, r)
	}).Methods("GET")
	r.HandleFunc("/qr/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		qrHandler.GetQrById(id, w, r)
	}).Methods("GET")

	r.HandleFunc("/qr/{id}/claim", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		qrHandler.ClaimQr(id, w, r)
	}).Methods("POST")

	return r
}
