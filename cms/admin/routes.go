package admin

import (
	"strconv"
	"website/backend/middleware"
	"website/backend/models"
	"website/backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	admin := app.Group("/admin")

	// Auth routes with strict rate limiting
	authGroup := admin.Group("/auth")
	authGroup.Use(middleware.RateLimitAuth())
	authGroup.Post("/login", loginHandler(db))
	authGroup.Post("/register", registerHandler(db))

	// Protected routes
	protected := admin.Use(middleware.AuthMiddleware(db))

	// Site Config - Solo admin y super_admin
	protected.Get("/config", middleware.IsEditor(), getSiteConfig(db))
	protected.Put("/config", middleware.IsAdmin(), updateSiteConfig(db))

	// Slides Management - Editor y superior
	protected.Get("/slides", middleware.IsEditor(), getSlides(db))
	protected.Post("/slides", middleware.IsEditor(), createSlide(db))
	protected.Put("/slides/:id", middleware.IsEditor(), updateSlide(db))
	protected.Delete("/slides/:id", middleware.IsAdmin(), deleteSlide(db))

	// Categories Management - Editor y superior
	protected.Get("/categories", middleware.IsEditor(), getCategories(db))
	protected.Post("/categories", middleware.IsEditor(), createCategory(db))
	protected.Put("/categories/:id", middleware.IsEditor(), updateCategory(db))
	protected.Delete("/categories/:id", middleware.IsAdmin(), deleteCategory(db))

	// Products Management - Editor y superior
	protected.Get("/products", middleware.IsEditor(), getProducts(db))
	protected.Post("/products", middleware.IsEditor(), createProduct(db))
	protected.Put("/products/:id", middleware.IsEditor(), updateProduct(db))
	protected.Delete("/products/:id", middleware.IsAdmin(), deleteProduct(db))

	// Contact Info - Editor y superior
	protected.Get("/contacts", middleware.IsEditor(), getContacts(db))
	protected.Post("/contacts", middleware.IsEditor(), createContact(db))
	protected.Put("/contacts/:id", middleware.IsEditor(), updateContact(db))
	protected.Delete("/contacts/:id", middleware.IsAdmin(), deleteContact(db))

	// Users Management - Solo admin y super_admin
	protected.Get("/users", middleware.IsAdmin(), getUsers(db))
	protected.Post("/users", middleware.IsAdmin(), createUser(db))
	protected.Put("/users/:id", middleware.IsAdmin(), updateUser(db))
	protected.Delete("/users/:id", middleware.IsSuperAdmin(), deleteUser(db))

	// System Management - Solo super_admin
	protected.Get("/system/status", middleware.IsSuperAdmin(), getSystemStatus(db))
	protected.Post("/system/backup", middleware.IsSuperAdmin(), createBackup(db))
}

// Auth Handlers
func loginHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginData struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		if err := utils.ParseAndValidate(c, &loginData); err != nil {
			return err
		}

		var user models.User
		if err := db.Where("username = ? OR email = ?", loginData.Username, loginData.Username).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Credenciales inválidas"})
		}

		if !user.Active {
			return c.Status(401).JSON(fiber.Map{"error": "Usuario inactivo"})
		}

		if !middleware.CheckPassword(loginData.Password, user.Password) {
			return c.Status(401).JSON(fiber.Map{"error": "Credenciales inválidas"})
		}

		token, err := middleware.GenerateToken(user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al generar token"})
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user":  user,
		})
	}
}

func registerHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		if err := utils.ParseAndValidate(c, &user); err != nil {
			return err
		}

		// Hash password
		hashedPassword, err := middleware.HashPassword(user.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al procesar contraseña"})
		}
		user.Password = hashedPassword

		// Generate slug for username if needed
		if user.Username == "" {
			user.Username = utils.GenerateSlug(user.Email)
		}

		if err := db.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al crear usuario"})
		}

		return c.JSON(fiber.Map{"message": "Usuario creado exitosamente", "user": user})
	}
}

// Site Config Handlers
func getSiteConfig(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var configs []models.SiteConfig
		if err := db.Find(&configs).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener configuración"})
		}
		return c.JSON(configs)
	}
}

func updateSiteConfig(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var config models.SiteConfig
		if err := utils.ParseAndValidate(c, &config); err != nil {
			return err
		}

		if err := db.Save(&config).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar configuración"})
		}

		return c.JSON(config)
	}
}

// Slides Handlers
func getSlides(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var slides []models.Slide
		if err := db.Order("order ASC").Find(&slides).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener slides"})
		}
		return c.JSON(slides)
	}
}

func createSlide(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var slide models.Slide
		if err := utils.ParseAndValidate(c, &slide); err != nil {
			return err
		}

		if err := db.Create(&slide).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al crear slide"})
		}

		return c.JSON(slide)
	}
}

func updateSlide(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		var slide models.Slide
		if err := utils.ParseAndValidate(c, &slide); err != nil {
			return err
		}

		slide.ID = uint(id)
		if err := db.Save(&slide).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar slide"})
		}

		return c.JSON(slide)
	}
}

func deleteSlide(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		if err := db.Delete(&models.Slide{}, id).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar slide"})
		}

		return c.JSON(fiber.Map{"message": "Slide eliminado exitosamente"})
	}
}

// Categories Handlers
func getCategories(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var categories []models.Category
		if err := db.Order("order ASC").Find(&categories).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener categorías"})
		}
		return c.JSON(categories)
	}
}

func createCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var category models.Category
		if err := utils.ParseAndValidate(c, &category); err != nil {
			return err
		}

		// Generate slug if not provided
		if category.Slug == "" {
			category.Slug = utils.GenerateSlug(category.Name)
		}

		if err := db.Create(&category).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al crear categoría"})
		}

		return c.JSON(category)
	}
}

func updateCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		var category models.Category
		if err := utils.ParseAndValidate(c, &category); err != nil {
			return err
		}

		category.ID = uint(id)
		if err := db.Save(&category).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar categoría"})
		}

		return c.JSON(category)
	}
}

func deleteCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		if err := db.Delete(&models.Category{}, id).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar categoría"})
		}

		return c.JSON(fiber.Map{"message": "Categoría eliminada exitosamente"})
	}
}

// Products Handlers
func getProducts(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var products []models.Product
		query := db.Preload("Category")

		if categoryID := c.Query("category_id"); categoryID != "" {
			query = query.Where("category_id = ?", categoryID)
		}

		if err := query.Order("created_at DESC").Find(&products).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener productos"})
		}
		return c.JSON(products)
	}
}

func createProduct(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var product models.Product
		if err := utils.ParseAndValidate(c, &product); err != nil {
			return err
		}

		// Generate slug if not provided
		if product.Slug == "" {
			product.Slug = utils.GenerateSlug(product.Name)
		}

		if err := db.Create(&product).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al crear producto"})
		}

		return c.JSON(product)
	}
}

func updateProduct(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		var product models.Product
		if err := utils.ParseAndValidate(c, &product); err != nil {
			return err
		}

		product.ID = uint(id)
		if err := db.Save(&product).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar producto"})
		}

		return c.JSON(product)
	}
}

func deleteProduct(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		if err := db.Delete(&models.Product{}, id).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar producto"})
		}

		return c.JSON(fiber.Map{"message": "Producto eliminado exitosamente"})
	}
}

// Contacts Handlers
func getContacts(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contacts []models.ContactInfo
		if err := db.Order("order ASC").Find(&contacts).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener contactos"})
		}
		return c.JSON(contacts)
	}
}

func createContact(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contact models.ContactInfo
		if err := utils.ParseAndValidate(c, &contact); err != nil {
			return err
		}

		if err := db.Create(&contact).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al crear contacto"})
		}

		return c.JSON(contact)
	}
}

func updateContact(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		var contact models.ContactInfo
		if err := utils.ParseAndValidate(c, &contact); err != nil {
			return err
		}

		contact.ID = uint(id)
		if err := db.Save(&contact).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar contacto"})
		}

		return c.JSON(contact)
	}
}

func deleteContact(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		if err := db.Delete(&models.ContactInfo{}, id).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar contacto"})
		}

		return c.JSON(fiber.Map{"message": "Contacto eliminado exitosamente"})
	}
}

// Users Handlers
func getUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener usuarios"})
		}
		return c.JSON(users)
	}
}

func createUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		if err := utils.ParseAndValidate(c, &user); err != nil {
			return err
		}

		// Hash password
		hashedPassword, err := middleware.HashPassword(user.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al procesar contraseña"})
		}
		user.Password = hashedPassword

		if err := db.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al crear usuario"})
		}

		return c.JSON(user)
	}
}

func updateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		var user models.User
		if err := utils.ParseAndValidate(c, &user); err != nil {
			return err
		}

		user.ID = uint(id)

		// Hash password if provided
		if user.Password != "" {
			hashedPassword, err := middleware.HashPassword(user.Password)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Error al procesar contraseña"})
			}
			user.Password = hashedPassword
		}

		if err := db.Save(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar usuario"})
		}

		return c.JSON(user)
	}
}

func deleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		if err := db.Delete(&models.User{}, id).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar usuario"})
		}

		return c.JSON(fiber.Map{"message": "Usuario eliminado exitosamente"})
	}
}

// System Management Handlers
func getSystemStatus(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var stats struct {
			TotalUsers       int64 `json:"total_users"`
			ActiveUsers      int64 `json:"active_users"`
			TotalSlides      int64 `json:"total_slides"`
			ActiveSlides     int64 `json:"active_slides"`
			TotalCategories  int64 `json:"total_categories"`
			ActiveCategories int64 `json:"active_categories"`
			TotalProducts    int64 `json:"total_products"`
			ActiveProducts   int64 `json:"active_products"`
			TotalContacts    int64 `json:"total_contacts"`
			ActiveContacts   int64 `json:"active_contacts"`
		}

		// Get user statistics
		db.Model(&models.User{}).Count(&stats.TotalUsers)
		db.Model(&models.User{}).Where("active = ?", true).Count(&stats.ActiveUsers)

		// Get slide statistics
		db.Model(&models.Slide{}).Count(&stats.TotalSlides)
		db.Model(&models.Slide{}).Where("active = ?", true).Count(&stats.ActiveSlides)

		// Get category statistics
		db.Model(&models.Category{}).Count(&stats.TotalCategories)
		db.Model(&models.Category{}).Where("active = ?", true).Count(&stats.ActiveCategories)

		// Get product statistics
		db.Model(&models.Product{}).Count(&stats.TotalProducts)
		db.Model(&models.Product{}).Where("active = ?", true).Count(&stats.ActiveProducts)

		// Get contact statistics
		db.Model(&models.ContactInfo{}).Count(&stats.TotalContacts)
		db.Model(&models.ContactInfo{}).Where("active = ?", true).Count(&stats.ActiveContacts)

		return c.JSON(stats)
	}
}

func createBackup(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Esta es una implementación básica de backup
		// En producción, deberías usar herramientas como pg_dump

		var backup struct {
			Timestamp string                 `json:"timestamp"`
			Data      map[string]interface{} `json:"data"`
		}

		backup.Timestamp = "2024-01-01T00:00:00Z"
		backup.Data = make(map[string]interface{})

		// Exportar datos principales
		var slides []models.Slide
		var categories []models.Category
		var products []models.Product
		var contacts []models.ContactInfo
		var configs []models.SiteConfig

		db.Find(&slides)
		db.Find(&categories)
		db.Find(&products)
		db.Find(&contacts)
		db.Find(&configs)

		backup.Data["slides"] = slides
		backup.Data["categories"] = categories
		backup.Data["products"] = products
		backup.Data["contacts"] = contacts
		backup.Data["configs"] = configs

		return c.JSON(fiber.Map{
			"message": "Backup creado exitosamente",
			"backup":  backup,
		})
	}
}
