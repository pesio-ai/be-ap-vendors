package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pesio-ai/be-go-common/config"
	"github.com/pesio-ai/be-go-common/database"
	"github.com/pesio-ai/be-go-common/logger"
	"github.com/pesio-ai/be-go-common/middleware"
	"github.com/pesio-ai/be-vendors-service/internal/handler"
	"github.com/pesio-ai/be-vendors-service/internal/repository"
	"github.com/pesio-ai/be-vendors-service/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(logger.Config{
		Level:       os.Getenv("LOG_LEVEL"),
		Environment: cfg.Service.Environment,
		ServiceName: cfg.Service.Name,
		Version:     cfg.Service.Version,
	})

	log.Info().
		Str("service", cfg.Service.Name).
		Str("version", cfg.Service.Version).
		Str("environment", cfg.Service.Environment).
		Msg("Starting Vendors Service (AP-1)")

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database
	db, err := database.New(ctx, database.Config{
		Host:        cfg.Database.Host,
		Port:        cfg.Database.Port,
		User:        cfg.Database.User,
		Password:    cfg.Database.Password,
		Database:    cfg.Database.Database,
		SSLMode:     cfg.Database.SSLMode,
		MaxConns:    cfg.Database.MaxConns,
		MinConns:    cfg.Database.MinConns,
		MaxConnTime: cfg.Database.MaxConnTime,
		MaxIdleTime: cfg.Database.MaxIdleTime,
		HealthCheck: cfg.Database.HealthCheck,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()
	log.Info().Msg("Database connection established")

	// Initialize repositories
	vendorRepo := repository.NewVendorRepository(db)

	// Initialize services
	vendorService := service.NewVendorService(vendorRepo, log)

	// Setup HTTP routes
	httpHandler := handler.NewHTTPHandler(vendorService, log)
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Vendor routes
	mux.HandleFunc("/api/v1/vendors", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			httpHandler.ListVendors(w, r)
		case http.MethodPost:
			httpHandler.CreateVendor(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/vendors/get", httpHandler.GetVendor)
	mux.HandleFunc("/api/v1/vendors/code", httpHandler.GetVendorByCode)
	mux.HandleFunc("/api/v1/vendors/update", httpHandler.UpdateVendor)
	mux.HandleFunc("/api/v1/vendors/delete", httpHandler.DeleteVendor)
	mux.HandleFunc("/api/v1/vendors/validate", httpHandler.ValidateVendor)

	// Vendor contact routes
	mux.HandleFunc("/api/v1/vendors/contacts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			httpHandler.GetVendorContacts(w, r)
		case http.MethodPost:
			httpHandler.AddVendorContact(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Payment terms routes
	mux.HandleFunc("/api/v1/payment-terms", httpHandler.GetPaymentTerms)

	// Vendor balance routes
	mux.HandleFunc("/api/v1/vendors/balance", httpHandler.UpdateBalance)

	// Apply middleware
	var h http.Handler = mux
	h = middleware.RequestID(h)
	h = middleware.Logger(&log.Logger)(h)
	h = middleware.Recovery(&log.Logger)(h)
	h = middleware.CORS([]string{"*"})(h)
	h = middleware.Timeout(30 * time.Second)(h)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      h,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info().Int("port", cfg.Server.Port).Msg("Starting HTTP server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown failed")
	}

	log.Info().Msg("Server stopped")
}
