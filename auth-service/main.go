package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RuslanVNikolov/one-stop-shop/backend/auth-service/internal/config"
	"github.com/RuslanVNikolov/one-stop-shop/backend/auth-service/internal/database"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	log.Println("✓ Configuration loaded")

	db := database.Connect(cfg.DatabaseURL)

	sqlDB, err := database.GetDB(db)
	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}
	defer sqlDB.Close()

	database.Migrate(db)

	router := setupRouter(db)

	log.Printf("🚀 Auth Service starting on http://localhost%s\n", cfg.Port)
	log.Printf("📊 Environment: %s\n", cfg.Environment)
	log.Fatal(http.ListenAndServe(cfg.Port, router))
}

func setupRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/health/db", func(w http.ResponseWriter, r *http.Request) {
		databaseHealthCheck(w, r, db)
	}).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()

	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", notImplementedHandler).Methods("POST")
	authRouter.HandleFunc("/login", notImplementedHandler).Methods("POST")
	authRouter.HandleFunc("/refresh", notImplementedHandler).Methods("POST")
	authRouter.HandleFunc("/logout", notImplementedHandler).Methods("POST")
	authRouter.HandleFunc("/me", notImplementedHandler).Methods("GET")

	return router
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","service":"auth-service"}`)
}

func databaseHealthCheck(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"Failed to get database instance"}`)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"Database ping failed"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"Database connection healthy"}`)
}

func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error":"This endpoint is not yet implemented"}`)
}
