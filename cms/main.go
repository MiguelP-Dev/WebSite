package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"
	"website/backend/middleware"
	"website/backend/models"
	"website/cms/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type CMS struct {
	App *fiber.App
	DB  *gorm.DB
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("CMS_DB_HOST"),
		os.Getenv("CMS_DB_USER"),
		os.Getenv("CMS_DB_PASSWORD"),
		os.Getenv("CMS_DB_NAME"),
		os.Getenv("CMS_DB_PORT"),
	)

	// Configure GORM logger
	gormLogger := gormlogger.Default.LogMode(gormlogger.Info)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get underlying sql.DB object to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate models
	db.AutoMigrate(
		&models.SiteConfig{},
		&models.Slide{},
		&models.Category{},
		&models.Product{},
		&models.ContactInfo{},
		&models.User{},
	)

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Security middleware (debe ir primero)
	app.Use(middleware.SecurityHeadersMiddleware())

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset",
	}))

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path} - IP: ${ip}\n",
	}))

	// Rate limiting estricto para CMS
	app.Use(middleware.RateLimitStrict())

	// Setup admin routes
	admin.SetupRoutes(app, db)

	// Get port from environment or use default
	port := os.Getenv("CMS_PORT")
	if port == "" {
		port = "4000"
	}

	// Check if HTTPS is enabled
	enableHTTPS := os.Getenv("ENABLE_HTTPS") == "true"
	certFile := os.Getenv("SSL_CERT_FILE")
	keyFile := os.Getenv("SSL_KEY_FILE")

	if enableHTTPS && certFile != "" && keyFile != "" {
		// Configure TLS
		tlsConfig := middleware.TLSServerConfig()

		// Create custom listener with TLS
		ln, err := tls.Listen("tcp", ":"+port, tlsConfig)
		if err != nil {
			log.Fatal("Failed to create TLS listener:", err)
		}

		log.Printf("HTTPS CMS starting on port %s", port)
		log.Fatal(app.Listener(ln))
	} else {
		// Start HTTP server
		log.Printf("HTTP CMS starting on port %s", port)
		log.Fatal(app.Listen(":" + port))
	}
}
