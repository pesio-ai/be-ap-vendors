package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pesio-ai/be-go-common/auth"
	"github.com/pesio-ai/be-go-common/config"
	"github.com/pesio-ai/be-go-common/database"
	"github.com/pesio-ai/be-go-common/logger"
	"github.com/pesio-ai/be-go-common/middleware"
	pb "github.com/pesio-ai/be-go-proto/gen/go/ap"
	identitypb "github.com/pesio-ai/be-go-proto/gen/go/platform"
	"github.com/pesio-ai/be-vendors-service/internal/handler"
	"github.com/pesio-ai/be-vendors-service/internal/repository"
	"github.com/pesio-ai/be-vendors-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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

	// Connect to identity service for authentication
	identityGrpcAddr := getEnv("IDENTITY_GRPC_URL", "localhost:9081")
	identityConn, err := grpc.NewClient(identityGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to identity service")
	}
	defer identityConn.Close()

	identityClient := identitypb.NewIdentityServiceClient(identityConn)
	log.Info().Str("identity_grpc", identityGrpcAddr).Msg("Identity service client initialized")

	// Setup HTTP handler
	httpHandler := handler.NewHTTPHandler(vendorService, log)

	// Setup gRPC handler
	grpcHandler := handler.NewGRPCHandler(vendorService, log)
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
	mux.HandleFunc("/api/v1/vendors/activate", httpHandler.ActivateVendor)
	mux.HandleFunc("/api/v1/vendors/deactivate", httpHandler.DeactivateVendor)
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

	// Setup gRPC server with auth interceptor
	grpcPort := 9084 // gRPC port (9000 + service number)

	// Create auth interceptor
	authInterceptor := auth.NewInterceptor(identityClient, log)

	// Create gRPC server with auth interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryServerInterceptor()),
	)
	pb.RegisterVendorsServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal().Err(err).Int("port", grpcPort).Msg("Failed to create gRPC listener")
	}

	go func() {
		log.Info().Int("port", grpcPort).Msg("Starting gRPC server")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Error().Err(err).Msg("gRPC server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown failed")
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	log.Info().Msg("Servers stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
