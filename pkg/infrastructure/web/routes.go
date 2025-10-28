package web

import (
	"bufio"
	"catecard/pkg/domain/entities"
	"catecard/pkg/handlers"
	"encoding/json"
	"fmt"
	"log"
	"net"
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

// ---------- Logging middleware (único) ----------

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

// ---------- statusWriter blindado ----------

type statusWriter struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (w *statusWriter) WriteHeader(code int) {
	if w.wrote {
		// Evita segundo WriteHeader (causa del “superfluous...”)
		return
	}
	w.status = code
	w.wrote = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		// Fija 200 en el primer Write si nadie fijó status
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

// Soporte opcional de interfaces comunes (SSE/WebSockets/HTTP2)
func (w *statusWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacker not supported")
}

func (w *statusWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := w.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

// ---------- Router setup ----------

func setupRouter(
	logger *log.Logger,
	qrHandler handlers.QrHandler,
	authHandler handlers.AuthHandler,
	groupHandler handlers.GroupHandler,
	catechumenHandler handlers.CatechumenHandler,
) *mux.Router {

	r := mux.NewRouter()

	// Usa SOLO el middleware de logging formal (evita duplicar)
	r.Use(LoggingMiddleware(logger))

	////////////////////////////////////////////////////////
	// 📂 Archivos públicos y logs
	////////////////////////////////////////////////////////
	r.PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

	r.HandleFunc("/__logs", func(w http.ResponseWriter, r *http.Request) {
		logsHandler(w, r)
	}).Methods("GET")

	////////////////////////////////////////////////////////
	// 🌐 VISTAS (HTML)
	////////////////////////////////////////////////////////
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
		if user.Role == entities.SCANNER {
			handlers.Scanner(w, r)
			return
		}
		handlers.Home(w, r)
	}).Methods("GET")

	r.HandleFunc("/signup", handlers.SignUp).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("GET")
	r.HandleFunc("/scanner", handlers.Scanner).Methods("GET")
	r.HandleFunc("/dene", handlers.Denied)
	r.HandleFunc("/all-qr-list", handlers.QrList)

	////////////////////////////////////////////////////////
	// 🔐 AUTENTICACIÓN
	////////////////////////////////////////////////////////
	r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", handlers.DeleteSession).Methods("POST")

	////////////////////////////////////////////////////////
	// 👤 USUARIOS / ADMIN
	////////////////////////////////////////////////////////
	r.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		authHandler.GetUserById(user, id, w, r)
	}).Methods("GET")

	r.HandleFunc("/all-catechists", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.GetAllCatechists(user, w, r)
	}).Methods("GET")

	r.HandleFunc("/delete-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		authHandler.DeleteUserById(user, id, w, r)
	}).Methods("DELETE")

	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.CreateAccounts(user, w, r)
	}).Methods("POST")

	r.HandleFunc("/all-scanners", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.GetAllScanners(user, w, r)
	}).Methods("GET")

	r.HandleFunc("/admin/scans", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		qrHandler.GetAllScans(user, w, r)
	}).Methods("GET")

	r.HandleFunc("/catechists-without-group", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		authHandler.GetAllCatechistsWithoutGroup(user, 0, w, r)
	}).Methods("GET")

	////////////////////////////////////////////////////////
	// 👥 GRUPOS
	////////////////////////////////////////////////////////
	r.HandleFunc("/all-groups", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		groupHandler.GetAllGroups(user, w, r)
	}).Methods("GET")

	r.HandleFunc("/admin/groups/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		groupHandler.GetGroup(user, id, w, r)
	}).Methods("GET")

	r.HandleFunc("/group/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		groupHandler.GetGroupById(user, id, w, r)
	}).Methods("GET")

	r.HandleFunc("/add-group", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		groupHandler.AddGroup(user, w, r)
	}).Methods("POST")

	r.HandleFunc("/edit-group/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		groupHandler.EditGroup(user, w, r)
	}).Methods("PUT")

	r.HandleFunc("/delete-group/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		groupHandler.DeleteGroupById(user, id, w, r)
	}).Methods("DELETE")

	////////////////////////////////////////////////////////
	// 📱 QR / SCANNERS
	////////////////////////////////////////////////////////
	// r.HandleFunc("/add-qr", func(w http.ResponseWriter, r *http.Request) {
	// 	qrHandler.AddQr(w, r)
	// }).Methods("POST")

	r.HandleFunc("/all-qr", func(w http.ResponseWriter, r *http.Request) {
		qrHandler.GetAllQrs(w, r)
	}).Methods("GET")

	r.HandleFunc("/qr/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		qrHandler.GetQrById(id, w, r)
	}).Methods("GET")

	r.HandleFunc("/qr/{id}/claim", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		user := handlers.GetUserFromRequest(r)
		qrHandler.ClaimQr(user, id, w, r)
	}).Methods("POST")

	////////////////////////////////////////////////////////
	// 🙋 CATECÚMENOS
	////////////////////////////////////////////////////////
	r.HandleFunc("/add-catechumen", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		catechumenHandler.AddCatechumen(user, w, r)
	}).Methods("POST")

	r.HandleFunc("/catechumens", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		catechumenHandler.GetAllCatechumens(user, w, r)
	}).Methods("GET")

	r.HandleFunc("/catechumen/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		catechumenHandler.GetCatechumenById(user, id, w, r)
	}).Methods("GET")

	r.HandleFunc("/catechumen/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		catechumenHandler.UpdateCatechumen(user, w, r)
	}).Methods("PUT")
	r.HandleFunc("/delete-catechumen/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromRequest(r)
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		catechumenHandler.DeleteCatechumenById(user, id, w, r)
	}).Methods("DELETE")

	return r
}
