package web

import (
	"catecard/pkg/config"
	"catecard/pkg/domain/usecases"
	"catecard/pkg/handlers"
	"catecard/pkg/infrastructure/repositories"
	"database/sql"

	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
)

func c(s, col string) string { return col + s + colorReset }

func statusColor(code int) string {
	switch {
	case code >= 500:
		return colorRed
	case code >= 400:
		return colorYellow
	case code >= 300:
		return colorCyan
	default:
		return colorGreen
	}
}

func methodColor(m string) string {
	switch m {
	case "GET":
		return colorGreen
	case "POST":
		return colorBlue
	case "PUT", "PATCH":
		return colorMagenta
	case "DELETE":
		return colorRed
	default:
		return colorWhite
	}
}

func printBanner(logger *log.Logger, port, dbPath string) {
	line := strings.Repeat("─", 48)
	logger.Printf("%s", c("┌"+line+"┐", colorBlue))
	logger.Printf("%s  %s  %s", c("│", colorBlue), c("Server Running", colorMagenta), c("│", colorBlue))
	logger.Printf("%s  %s %s  %s", c("│", colorBlue), c("Port:", colorWhite), c(port, colorGreen), c("│", colorBlue))
	logger.Printf("%s  %s %s  %s", c("│", colorBlue), c("DB:  ", colorWhite), c(dbPath, colorGreen), c("│", colorBlue))
	logger.Printf("%s", c("└"+line+"┘", colorBlue))
}

// StartWebServer starts the HTTP server using values from cfg (DB path, port, etc.).
func StartWebServer(cfg *config.Config) error {
	logger := log.New(os.Stdout, c("[*] ", colorCyan), log.LstdFlags|log.Lshortfile)

	dbPath := cfg.DatabasePath
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Migrations
	migPath := "./db/migrations/tables.sql"
	migBytes, err := os.ReadFile(migPath)
	if err != nil {
		logger.Fatalf("Failed to read migration file %s: %v", migPath, err)
	}

	migSQL := string(migBytes)
	migSQL = replaceSerial(migSQL)
	migSQL = removeOnUpdate(migSQL)

	if err := applyMigrations(db, migSQL, logger); err != nil {
		logger.Fatalf("Failed to apply migrations: %v", err)
	}

	// Inicializar el store de sesiones para que use la misma BD
	handlers.InitSessionStore(db)
	printBanner(logger, cfg.ServerPort, dbPath)

	// rutas templates

	tmplPath := "./pkg/infrastructure/web/templates"
	if err := handlers.LoadTemplates(tmplPath); err != nil {
		logger.Fatalf("Failed to load templates: %v", err)
	}

	//REPOSITORIES
	authRepo := repositories.NewUserRepository(logger, db)
	groupRepo := repositories.NewGroupRepository(logger, db)
	qrRepo := repositories.NewQrRepository(logger, db)
	catechumenRepo := repositories.NewMockCatechumenRepository()

	//USECASES
	authUC := usecases.NewAuthUseCase(logger, authRepo)
	groupUC := usecases.NewGroupUsecase(logger, groupRepo)
	qrUC := usecases.NewQrUsecase(logger, qrRepo)
	catechumenUC := usecases.NewCatechumenUsecase(logger, catechumenRepo, qrRepo)

	//HANDLERS

	authHandler := handlers.NewAuthenticationHandler(logger, authUC, tmplPath)

	groupHandler := handlers.NewGroupHandler(logger, groupUC, tmplPath)

	catechumenHandler := handlers.NewCatechumenHandler(logger, catechumenUC, tmplPath)
	qrHandler := handlers.NewQrHandler(logger, qrUC, tmplPath)

	//Router
	r := setupRouter(logger, qrHandler, authHandler, groupHandler, catechumenHandler)
	logger.Printf("Server running on port: %s", cfg.ServerPort)
	logger.Printf("Database: %s", dbPath)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		logger.Fatal(err)
	}
	return nil
}

// SQLite helpers
func removeOnUpdate(s string) string {
	re := regexp.MustCompile(`ON UPDATE CURRENT_TIMESTAMP`)
	return re.ReplaceAllString(s, "")
}

func replaceSerial(s string) string {
	re := regexp.MustCompile(`(?i)\bSERIAL\b\s*PRIMARY\s*KEY`)
	s = re.ReplaceAllString(s, "INTEGER PRIMARY KEY AUTOINCREMENT")
	re2 := regexp.MustCompile(`(?i)\bSERIAL\b`)
	s = re2.ReplaceAllString(s, "INTEGER")
	return s
}

func applyMigrations(db *sql.DB, sqlContent string, logger *log.Logger) error {
	parts := strings.Split(sqlContent, ";")
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			logger.Printf("migration statement failed: %s; err: %v", stmt, err)
			return err
		}
	}
	return nil
}
