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

// StartWebServer starts the HTTP server using values from cfg (DB path, port, etc.).
func StartWebServer(cfg *config.Config) error {
	logger := log.New(os.Stdout, "[CATECART] ", log.LstdFlags|log.Lshortfile)

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

	// rutas templates

	tmplPath := "./pkg/infrastructure/web/templates"
	if err := handlers.LoadTemplates(tmplPath); err != nil {
		logger.Fatalf("Failed to load templates: %v", err)
	}

	//REPOSITORIES
	authRepo := repositories.NewUserRepository(logger, db)
	groupRepo := repositories.NewGroupRepository(logger, db)
	qrRepo := repositories.NewQrRepository(logger, db)

	//USECASES
	authUC := usecases.NewAuthUseCase(logger, authRepo)
	groupUC := usecases.NewGroupUsecase(logger, groupRepo)
	qrUC := usecases.NewQrUsecase(logger, qrRepo)

	//HANDLERS

	authHandler := handlers.NewAuthenticationHandler(logger, authUC, tmplPath)

	groupHandler := handlers.NewGroupHandler(logger, groupUC, tmplPath)
	qrHandler := handlers.NewQrHandler(logger, qrUC, tmplPath)

	//Router
	r := setupRouter(qrHandler, authHandler, groupHandler)
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
