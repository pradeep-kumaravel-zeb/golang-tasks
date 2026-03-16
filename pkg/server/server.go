package server

import (
	"fmt"
	"net/http"
	"student-enrollment/internal/config"
	"student-enrollment/internal/handlers"
	"student-enrollment/internal/repository"
	"student-enrollment/internal/service"
	"student-enrollment/pkg/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	router *mux.Router
	logger *zap.Logger
	db     *gorm.DB
}

func InitializeApp(creds *config.DBCredentials) (*Application, error) {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Build database URL
	dbURL := config.BuildDBUrl(*creds)

	// Establish database connection
	db, err := database.EstablishPostgresConnection(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize layers
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	handler := handlers.NewHandler(svc)

	// Setup router
	router := mux.NewRouter()

	// Enrollment routes
	router.HandleFunc("/enrollments", handler.GetAllEnrollments).Methods("GET")
	router.HandleFunc("/enrollments/{id}", handler.GetEnrollmentByID).Methods("GET")
	router.HandleFunc("/enrollments", handler.CreateEnrollment).Methods("POST")
	router.HandleFunc("/enrollments/{id}", handler.UpdateEnrollment).Methods("PUT")
	router.HandleFunc("/enrollments/{id}", handler.PatchEnrollment).Methods("PATCH")
	router.HandleFunc("/enrollments/{id}", handler.DeleteEnrollment).Methods("DELETE")

	logger.Info("Application initialized successfully")

	return &Application{
		router: router,
		logger: logger,
		db:     db,
	}, nil
}

func (app *Application) Start() error {
	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(app.router)

	app.logger.Info("Starting server on :8080")
	return http.ListenAndServe(":8080", handler)
}
