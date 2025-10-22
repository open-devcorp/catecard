package web

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/handlers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// In-memory request log buffer (dev)
type RequestLog struct {
	Time    time.Time `json:"time"`
	Method  string    `json:"method"`
	Path    string    `json:"path"`
	Status  int       `json:"status"`
	Latency int64     `json:"latency_ms"`
}

var (
	logsMu  sync.Mutex
	logsBuf []RequestLog
	maxLogs = 200
)

func pushLog(entry RequestLog) {
	logsMu.Lock()
	defer logsMu.Unlock()
	logsBuf = append([]RequestLog{entry}, logsBuf...)
	if len(logsBuf) > maxLogs {
		logsBuf = logsBuf[:maxLogs]
	}
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logsMu.Lock()
	copyLogs := make([]RequestLog, len(logsBuf))
	copy(copyLogs, logsBuf)
	logsMu.Unlock()
	enc := json.NewEncoder(w)
	_ = enc.Encode(copyLogs)
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, r)
			dur := time.Since(start)

			sc := statusColor(sw.status)
			mc := methodColor(r.Method)

			logger.Printf(
				"%s %s %s %s %s",
				c("[HTTP]", colorWhite),
				c(fmt.Sprintf("%3d", sw.status), sc),
				c(r.Method, mc),
				c(r.URL.Path, colorWhite),
				c(fmt.Sprintf("%v", dur), colorCyan),
			)
			// push to in-memory buffer when DEV enabled
			dev := os.Getenv("DEV")
			if dev == "1" || dev == "true" || dev == "TRUE" {
				pushLog(RequestLog{
					Time:    start,
					Method:  r.Method,
					Path:    r.URL.Path,
					Status:  sw.status,
					Latency: int64(dur / time.Millisecond),
				})
			}
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func setupRouter(logger *log.Logger, qrHandler handlers.QrHandler, authHandler handlers.AuthHandler, groupHandler handlers.GroupHandler, catechumenHandler handlers.CatechumenHandler) *mux.Router {
	r := mux.NewRouter()
	// Simple middleware: print only status code and endpoint path
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, req)
			// Print minimal info: colored status, colored method and endpoint
			logger.Printf("%s %s %s",
				c(fmt.Sprintf("%3d", sw.status), statusColor(sw.status)),
				c(req.Method, methodColor(req.Method)),
				req.URL.Path,
			)
		})
	})
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

	// Expose logs endpoint for dev UI
	r.HandleFunc("/__logs", func(w http.ResponseWriter, r *http.Request) {
		logsHandler(w, r)
	}).Methods("GET")

	//////////////////////////VIEWS//////////////////////////

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
	r.HandleFunc("/scanner", handlers.Scanner).Methods("GET")

	//////////////////////////APIS//////////////////////////////

	//AUTH
	r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", handlers.DeleteSession).Methods("POST")

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

	r.HandleFunc("/add-catechumen", func(w http.ResponseWriter, r *http.Request) {
		catechumenHandler.AddCatechumen(w, r)
	}).Methods("POST")

	r.HandleFunc("/catechumens", func(w http.ResponseWriter, r *http.Request) {
		catechumenHandler.GetAllCatechumens(w, r)
	}).Methods("GET")

	r.HandleFunc("/catechumen/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		catechumenHandler.GetCatechumenById(id, w, r)
	}).Methods("GET")

	r.HandleFunc("/catechumen/{id}", func(w http.ResponseWriter, r *http.Request) {
		catechumenHandler.UpdateCatechumen(w, r)
	}).Methods("PUT")

	return r
}
