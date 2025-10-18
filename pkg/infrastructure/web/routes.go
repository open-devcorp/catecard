package web

import (
	"catecard/pkg/handlers"
	"net/http"

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
