package web

import (
	"catecard/pkg/handlers"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func setupRouter(authHandler handlers.AuthHandler, productHandler handlers.ProductHandler, groupHandler handlers.GroupHandler) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

	//products
	r.HandleFunc("/add-product", productHandler.AddProduct).Methods("POST")
	r.HandleFunc("/add-product", handlers.AddProductView).Methods("GET")
	r.HandleFunc("/products", productHandler.GetAllProducts).Methods("GET")

	//AUTH
	r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/signup", handlers.SignUp).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("GET")

	////DASHBOARD///////////////
	r.HandleFunc("/all-catechists", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.GetAllCatechists(user, w, r)
	}).Methods("GET")

	r.HandleFunc("/all-groups", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		groupHandler.GetAllGroups(user, w, r)
	}).Methods("GET")

	r.HandleFunc("/all-scanners", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.GetAllScanners(user, w, r)
	}).Methods("GET")

	//GET USER BY ID
	r.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		authHandler.GetUserById(user, id, w, r)
	}).Methods("GET")

	//GET GROUP BY ID
	r.HandleFunc("/group/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		groupHandler.GetGroupById(user, id, w, r)
	}).Methods("GET")

	///////////////////////////////
	//CATECHIST
	r.HandleFunc("/add-catechist", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.CreateAccounts(user, w, r)
	}).Methods("POST")

	//GROUPS
	r.HandleFunc("/add-group", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		groupHandler.AddGroup(user, w, r)
	}).Methods("POST")

	// HOME
	r.HandleFunc("/home", handlers.Home).Methods("GET")

	return r
}
