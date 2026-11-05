package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/config"
	"portfolio/api/internal/database"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/handler"
	"portfolio/api/internal/middleware"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/service"
	"portfolio/api/internal/storage"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Configuration error", "error", err)
		os.Exit(1)
	}

	// 2. Initialize database connection pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Database connection pool failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("PostgreSQL connection pool established")

	// 3. Initialize shared infrastructure
	rateLimiter := middleware.NewMemoryRateLimiter()

	storage.SetGlobalStorageConfig(cfg.AppBaseURL, cfg.StorageDriver)

	var storageProvider storage.StorageProvider
	if cfg.StorageDriver == "supabase" {
		storageProvider = storage.NewSupabaseStorageProvider(
			cfg.SupabaseURL,
			cfg.SupabaseBucket,
			cfg.SupabaseKey,
			cfg.AppBaseURL,
		)
		slog.Info("Storage provider initialized: Supabase Storage",
			"url", cfg.SupabaseURL,
			"bucket", cfg.SupabaseBucket,
		)
	} else {
		localStorage, err := storage.NewLocalStorageProvider(cfg.StorageLocalDir, cfg.AppBaseURL, cfg.StorageDriver)
		if err != nil {
			slog.Error("Storage provider initialization failed", "error", err)
			os.Exit(1)
		}
		storageProvider = localStorage
		slog.Info("Storage provider initialized: Local Storage", "dir", cfg.StorageLocalDir)
	}

	// 4. Initialize repositories
	userRepo        := repository.NewPostgresUserRepository(pool)
	profileRepo     := repository.NewPostgresProfileRepository(pool)
	categoryRepo    := repository.NewPostgresCategoryRepository(pool)
	skillRepo       := repository.NewPostgresSkillRepository(pool)
	projectRepo     := repository.NewPostgresProjectRepository(pool)
	experienceRepo  := repository.NewPostgresExperienceRepository(pool)
	certificateRepo := repository.NewPostgresCertificateRepository(pool)
	postRepo        := repository.NewPostgresPostRepository(pool)
	messageRepo     := repository.NewPostgresMessageRepository(pool)

	// 5. Initialize services (dependency order matters)
	authSvc        := service.NewAuthService(userRepo, cfg)
	categorySvc    := service.NewCategoryService(categoryRepo)
	skillSvc       := service.NewSkillService(skillRepo)
	profileSvc     := service.NewProfileService(profileRepo)
	projectSvc     := service.NewProjectService(projectRepo, categorySvc, skillSvc)
	experienceSvc  := service.NewExperienceService(experienceRepo)
	certificateSvc := service.NewCertificateService(certificateRepo)
	postSvc        := service.NewPostService(postRepo, categorySvc)
	messageSvc     := service.NewMessageService(messageRepo, rateLimiter, cfg.JWTSecret)

	// 6. Initialize handlers
	healthHandler      := handler.NewHealthHandler(pool)
	authHandler        := handler.NewAuthHandler(authSvc, cfg, rateLimiter)
	profileHandler     := handler.NewProfileHandler(profileSvc)
	categoryHandler    := handler.NewCategoryHandler(categorySvc)
	skillHandler       := handler.NewSkillHandler(skillSvc)
	projectHandler     := handler.NewProjectHandler(projectSvc)
	experienceHandler  := handler.NewExperienceHandler(experienceSvc)
	certificateHandler := handler.NewCertificateHandler(certificateSvc)
	postHandler        := handler.NewPostHandler(postSvc)
	messageHandler     := handler.NewMessageHandler(messageSvc)
	uploadHandler      := handler.NewUploadHandler(storageProvider)

	// 7. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName:      "Portfolio & CMS API v1",
		BodyLimit:    2 * 1024 * 1024, // 2 MB (matches upload limit)
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(domain.NewErrorResponse(
				"HTTP_ERROR",
				err.Error(),
				nil,
			))
		},
	})

	// 8. Apply Middleware Stack (order per PRD Bagian 4.5)
	app.Use(middleware.Recover())
	app.Use(middleware.RequestIDAndLogger())
	app.Use(middleware.SecurityHeaders(cfg.CookieSecure))
	app.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	app.Use(middleware.GlobalRateLimit(rateLimiter, 100, time.Minute))

	// 9. Serve uploaded static files
	app.Static("/uploads", cfg.StorageLocalDir, fiber.Static{
		Browse: false,
	})

	// 10. Register API Routes
	api := app.Group("/api/v1")

	// -- Health Check --
	api.Get("/health", healthHandler.Check)

	// -- Auth Routes --
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", middleware.RequireAuth(authSvc), authHandler.Me)

	// -- Public Content Routes --
	api.Get("/profile", profileHandler.Get)
	api.Get("/categories", categoryHandler.List)
	api.Get("/skills", skillHandler.ListPublic)
	api.Get("/projects", projectHandler.ListPublic)
	api.Get("/projects/:slug", projectHandler.GetBySlugPublic)
	api.Get("/experiences", experienceHandler.ListPublic)
	api.Get("/certificates", certificateHandler.ListPublic)
	api.Get("/posts", postHandler.ListPublic)
	api.Get("/posts/:slug", postHandler.GetBySlugPublic)
	api.Post("/contact", messageHandler.CreateMessage)

	// -- Admin Protected Group (RequireAuth + ValidateCSRF on mutations) --
	admin := api.Group("/admin", middleware.RequireAuth(authSvc), middleware.ValidateCSRF())

	// Admin stats
	admin.Get("/stats", func(c *fiber.Ctx) error {
		return c.JSON(domain.NewSuccessResponse(fiber.Map{
			"message": "Admin area accessed successfully",
		}, nil))
	})

	// Profile (admin)
	admin.Put("/profile", profileHandler.Update)

	// Categories (admin)
	admin.Get("/categories", categoryHandler.List)
	admin.Get("/categories/:id", categoryHandler.GetByID)
	admin.Post("/categories", categoryHandler.Create)
	admin.Put("/categories/:id", categoryHandler.Update)
	admin.Delete("/categories/:id", categoryHandler.Delete)

	// Skills (admin)
	admin.Get("/skills", skillHandler.ListAdmin)
	admin.Get("/skills/:id", skillHandler.GetByID)
	admin.Post("/skills", skillHandler.Create)
	admin.Put("/skills/:id", skillHandler.Update)
	admin.Delete("/skills/:id", skillHandler.Delete)

	// Projects (admin)
	admin.Get("/projects", projectHandler.ListAdmin)
	admin.Get("/projects/:id", projectHandler.GetByID)
	admin.Post("/projects", projectHandler.Create)
	admin.Put("/projects/:id", projectHandler.Update)
	admin.Put("/projects/:id/skills", projectHandler.UpdateSkills)
	admin.Delete("/projects/:id", projectHandler.Delete)

	// Experiences (admin)
	admin.Get("/experiences", experienceHandler.ListAdmin)
	admin.Get("/experiences/:id", experienceHandler.GetByID)
	admin.Post("/experiences", experienceHandler.Create)
	admin.Put("/experiences/:id", experienceHandler.Update)
	admin.Delete("/experiences/:id", experienceHandler.Delete)

	// Certificates (admin)
	admin.Get("/certificates", certificateHandler.ListAdmin)
	admin.Get("/certificates/:id", certificateHandler.GetByID)
	admin.Post("/certificates", certificateHandler.Create)
	admin.Put("/certificates/:id", certificateHandler.Update)
	admin.Delete("/certificates/:id", certificateHandler.Delete)

	// Posts (admin)
	admin.Get("/posts", postHandler.ListAdmin)
	admin.Get("/posts/:id", postHandler.GetByID)
	admin.Post("/posts", postHandler.Create)
	admin.Put("/posts/:id", postHandler.Update)
	admin.Delete("/posts/:id", postHandler.Delete)

	// Messages (admin)
	admin.Get("/messages", messageHandler.ListMessages)
	admin.Get("/messages/:id", messageHandler.GetMessageByID)
	admin.Patch("/messages/:id/status", messageHandler.UpdateStatus)
	admin.Delete("/messages/:id", messageHandler.DeleteMessage)

	// Uploads (admin + CSRF)
	admin.Post("/uploads/avatar", uploadHandler.UploadAvatar)
	admin.Post("/uploads/resume", uploadHandler.UploadResume)
	admin.Post("/uploads/project", uploadHandler.UploadProjectThumbnail)
	admin.Post("/uploads/certificate", uploadHandler.UploadCertificateFile)
	admin.Post("/uploads/post", uploadHandler.UploadPostCover)
	admin.Post("/uploads/skill", uploadHandler.UploadSkillIcon)
	admin.Delete("/uploads", uploadHandler.DeleteFile)

	// 11. Start Server with Graceful Shutdown
	go func() {
		addr := fmt.Sprintf(":%d", cfg.AppPort)
		slog.Info("Portfolio API server listening", "address", addr, "env", cfg.AppEnv)
		if err := app.Listen(addr); err != nil {
			slog.Info("Server stopped listening", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		slog.Error("Server forced shutdown", "error", err)
	}
	slog.Info("Server exited cleanly")
}
