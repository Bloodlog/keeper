package server

import (
	"context"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/handler"
	"keeper/internal/handler/web"
	"keeper/internal/repository"
	"keeper/internal/service"
	"keeper/internal/store"
	"log"
	"os/signal"
	"syscall"

	chi "github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

// Init Application.
func run(cfg *config.MainServerConfig) error {
	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(rootCtx)

	context.AfterFunc(ctx, func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the service")
	})

	// Init logger
	l, err := initLogger(rootCtx)
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	// Init DB
	database, err := store.NewDB(ctx, cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	// Init repositories
	userRepo := repository.NewUserRepository(database.Pool)
	vaultRepo := repository.NewVaultRepository(database.Pool)
	accessRepo := repository.NewAccessRepository(database.Pool)

	// Init services
	jwtService := service.NewJwtService(cfg.Security)
	authService := service.NewAuthService(database.Pool, userRepo, accessRepo, jwtService, cfg.Security, l)
	vaultService := service.NewVaultService(vaultRepo, cfg.Security.DataEncryptionKey)
	if cfg.Storage.StorageType == "file_storage" {
		_, err := store.NewFileStorage(ctx, cfg.FileStorageConfig)
		if err != nil {
			return fmt.Errorf("file storage error: %w", err)
		}
	}

	// Init handlers
	// WEB handlers.
	staticHandler := web.NewStaticPageHandler(l)
	fileHandler := web.NewFileServerHandler(l, cfg)
	downloadHandler := web.NewDownloadHandler(l, cfg)

	router := chi.NewRouter()
	router.Handle("/downloads/*", fileHandler.FileServerHandler(ctx))
	router.Get("/download", downloadHandler.DownloadHandler())
	router.NotFound(staticHandler.NotFoundHandler(context.Background()))

	// GRPC handlers.
	authHandler := handler.NewAuthHandler(l, authService)
	vaultHandler := handler.NewVaultHandler(l, vaultService)

	// Start HTTP server
	initHTTPServer(ctx, g, cfg, router, l)

	// Start Grpc Server
	initGRPCServer(ctx, g, cfg, l, authHandler, vaultHandler, jwtService)

	err = g.Wait()
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
