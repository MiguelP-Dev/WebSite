package controllers

import (
	"strings"
	"time"
	"website/backend/middleware"
	"website/backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Public API with rate limiting
	api := app.Group("/api")

	// Rate limiting específico para APIs
	api.Use(middleware.RateLimitAPI())

	api.Get("/config", getSiteConfig(db))
	api.Get("/slides", getActiveSlides(db))
	api.Get("/categories", getActiveCategories(db))
	api.Get("/products", getActiveProducts(db))
	api.Get("/contacts", getActiveContacts(db))

	// Web Pages with relaxed rate limiting
	app.Get("/", homeHandler(db))
	app.Get("/productos", productsHandler(db))
	app.Get("/productos/:category", categoryHandler(db))
	app.Get("/contacto", contactHandler(db))
	app.Get("/ubicaciones", locationsHandler(db))
	app.Get("/catalogo", catalogHandler(db))
}

// API Handlers
func getSiteConfig(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var configs []models.SiteConfig
		if err := db.Find(&configs).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener configuración"})
		}

		// Convert to map for easier frontend consumption
		configMap := make(map[string]string)
		for _, config := range configs {
			configMap[config.Key] = config.Value
		}

		return c.JSON(configMap)
	}
}

func getActiveSlides(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var slides []models.Slide
		if err := db.Where("active = ?", true).Order("order ASC").Find(&slides).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener slides"})
		}
		return c.JSON(slides)
	}
}

func getActiveCategories(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var categories []models.Category
		if err := db.Where("active = ?", true).Order("order ASC").Find(&categories).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener categorías"})
		}
		return c.JSON(categories)
	}
}

func getActiveProducts(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var products []models.Product
		query := db.Where("active = ?", true).Preload("Category")

		// Filter by category if provided
		if categoryID := c.Query("category_id"); categoryID != "" {
			query = query.Where("category_id = ?", categoryID)
		}

		if err := query.Order("created_at DESC").Find(&products).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener productos"})
		}
		return c.JSON(products)
	}
}

func getActiveContacts(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contacts []models.ContactInfo
		if err := db.Where("active = ?", true).Order("order ASC").Find(&contacts).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al obtener contactos"})
		}
		return c.JSON(contacts)
	}
}

// Helper function to get common data for all pages
func getCommonData(db *gorm.DB) fiber.Map {
	var siteConfig models.SiteConfig
	var contacts []models.ContactInfo

	db.Where("key = ?", "site_name").First(&siteConfig)
	siteName := siteConfig.Value
	if siteName == "" {
		siteName = "Mi Sitio Web"
	}

	db.Where("key = ?", "site_description").First(&siteConfig)
	siteDescription := siteConfig.Value
	if siteDescription == "" {
		siteDescription = "Descripción del sitio web"
	}

	db.Where("active = ?", true).Order("order ASC").Find(&contacts)

	return fiber.Map{
		"SiteName":        siteName,
		"SiteDescription": siteDescription,
		"CurrentYear":     time.Now().Year(),
		"Contacts":        contacts,
	}
}

// Web Page Handlers
func homeHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var slides []models.Slide
		var categories []models.Category
		var siteConfig models.SiteConfig

		// Get active slides
		db.Where("active = ?", true).Order("order ASC").Find(&slides)

		// Get active categories with some products
		db.Preload("Products", "active = ?", true).Where("active = ?", true).Order("order ASC").Limit(6).Find(&categories)

		// Get site config
		db.Where("key = ?", "home_videos").First(&siteConfig)

		// Prepare data for template
		commonData := getCommonData(db)
		data := fiber.Map{
			"Title":       "Inicio",
			"CurrentPage": "home",
			"Slides":      slides,
			"Categories":  categories,
			"Videos":      strings.Split(siteConfig.Value, ","),
		}

		// Merge common data
		for k, v := range commonData {
			data[k] = v
		}

		return c.Render("pages/index", data)
	}
}

func productsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var categories []models.Category
		if err := db.Preload("Products", "active = ?", true).Where("active = ?", true).Order("order ASC").Find(&categories).Error; err != nil {
			return c.Status(500).SendString("Error interno del servidor")
		}

		commonData := getCommonData(db)
		data := fiber.Map{
			"Title":       "Productos",
			"CurrentPage": "products",
			"Categories":  categories,
		}

		// Merge common data
		for k, v := range commonData {
			data[k] = v
		}

		return c.Render("pages/products", data)
	}
}

func categoryHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		categorySlug := c.Params("category")

		var category models.Category
		if err := db.Preload("Products", "active = ?", true).Where("slug = ? AND active = ?", categorySlug, true).First(&category).Error; err != nil {
			return c.Status(404).SendString("Categoría no encontrada")
		}

		commonData := getCommonData(db)
		data := fiber.Map{
			"Title":       category.Name,
			"CurrentPage": "category",
			"Category":    category,
		}

		// Merge common data
		for k, v := range commonData {
			data[k] = v
		}

		return c.Render("pages/category", data)
	}
}

func contactHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contacts []models.ContactInfo
		if err := db.Where("active = ?", true).Order("order ASC").Find(&contacts).Error; err != nil {
			return c.Status(500).SendString("Error interno del servidor")
		}

		commonData := getCommonData(db)
		data := fiber.Map{
			"Title":       "Contacto",
			"CurrentPage": "contact",
			"Contacts":    contacts,
		}

		// Merge common data
		for k, v := range commonData {
			data[k] = v
		}

		return c.Render("pages/contact", data)
	}
}

func locationsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var siteConfig models.SiteConfig
		db.Where("key = ?", "locations").First(&siteConfig)

		commonData := getCommonData(db)
		data := fiber.Map{
			"Title":       "Ubicaciones",
			"CurrentPage": "locations",
			"Locations":   strings.Split(siteConfig.Value, ","),
		}

		// Merge common data
		for k, v := range commonData {
			data[k] = v
		}

		return c.Render("pages/locations", data)
	}
}

func catalogHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var products []models.Product
		if err := db.Preload("Category").Where("active = ?", true).Order("created_at DESC").Find(&products).Error; err != nil {
			return c.Status(500).SendString("Error interno del servidor")
		}

		commonData := getCommonData(db)
		data := fiber.Map{
			"Title":       "Catálogo",
			"CurrentPage": "catalog",
			"Products":    products,
		}

		// Merge common data
		for k, v := range commonData {
			data[k] = v
		}

		return c.Render("pages/catalog", data)
	}
}
