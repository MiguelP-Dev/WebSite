package middleware

import (
	"os"
	"strings"
	"time"

	"website/backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware verifica el token JWT
func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token de autorización requerido",
			})
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{
				"error": "Formato de token inválido",
			})
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token inválido",
			})
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token inválido",
			})
		}

		// Check if user still exists and is active
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Usuario no encontrado",
			})
		}

		if !user.Active {
			return c.Status(401).JSON(fiber.Map{
				"error": "Usuario inactivo",
			})
		}

		// Add user info to context
		c.Locals("user", user)
		c.Locals("userID", claims.UserID)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// GenerateToken genera un token JWT para un usuario
func GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hours

	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// HashPassword hashea una contraseña usando bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword verifica si una contraseña coincide con su hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
