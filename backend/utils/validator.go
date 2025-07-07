package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators
	validate.RegisterValidation("alphanum", validateAlphanum)
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) []string {
	var errors []string

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			param := err.Param()

			var message string
			switch tag {
			case "required":
				message = field + " es requerido"
			case "email":
				message = field + " debe ser un email válido"
			case "url":
				message = field + " debe ser una URL válida"
			case "min":
				message = field + " debe tener al menos " + param + " caracteres"
			case "max":
				message = field + " debe tener máximo " + param + " caracteres"
			case "gte":
				message = field + " debe ser mayor o igual a " + param
			case "oneof":
				message = field + " debe ser uno de: " + param
			case "alphanum":
				message = field + " debe contener solo letras y números"
			default:
				message = field + " no es válido"
			}

			errors = append(errors, message)
		}
	}

	return errors
}

// validateAlphanum custom validator for alphanumeric strings
func validateAlphanum(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}

// ParseAndValidate parses JSON body and validates struct
func ParseAndValidate(c *fiber.Ctx, s interface{}) error {
	if err := c.BodyParser(s); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Error al parsear el JSON: " + err.Error(),
		})
	}

	errors := ValidateStruct(s)
	if len(errors) > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Errores de validación",
			"details": errors,
		})
	}

	return nil
}

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove special characters except hyphens
	var result strings.Builder
	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	// Remove multiple consecutive hyphens
	slug := result.String()
	slug = strings.ReplaceAll(slug, "--", "-")
	slug = strings.Trim(slug, "-")

	return slug
}
