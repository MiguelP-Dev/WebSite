package middleware

import (
	"os"
	"strings"
	"time"

	"website/backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthConfig configuración para el middleware de autenticación
type AuthConfig struct {
	SecretKey     string   // Clave secreta para JWT
	TokenPrefix   string   // Prefijo del token (Bearer, etc.)
	HeaderName    string   // Nombre del header que contiene el token
	QueryParam    string   // Parámetro de query que puede contener el token
	CookieName    string   // Nombre de la cookie que puede contener el token
	RequiredRoles []string // Roles requeridos para acceder
	Optional      bool     // Si la autenticación es opcional
}

// ========================================
// CONFIGURACIONES PREDEFINIDAS DE AUTENTICACIÓN
// ========================================

// AuthConfigs contiene configuraciones predefinidas de autenticación
var AuthConfigs = struct {
	Required AuthConfig // Autenticación requerida
	Optional AuthConfig // Autenticación opcional
	Admin    AuthConfig // Solo administradores
	Editor   AuthConfig // Solo editores
	User     AuthConfig // Solo usuarios autenticados
}{
	// Configuración para autenticación requerida
	Required: AuthConfig{
		SecretKey:     "tu-clave-secreta-aqui",
		TokenPrefix:   "Bearer",
		HeaderName:    "Authorization",
		QueryParam:    "token",
		CookieName:    "auth_token",
		RequiredRoles: []string{},
		Optional:      false,
	},
	// Configuración para autenticación opcional
	Optional: AuthConfig{
		SecretKey:     "tu-clave-secreta-aqui",
		TokenPrefix:   "Bearer",
		HeaderName:    "Authorization",
		QueryParam:    "token",
		CookieName:    "auth_token",
		RequiredRoles: []string{},
		Optional:      true,
	},
	// Configuración para solo administradores
	Admin: AuthConfig{
		SecretKey:     "tu-clave-secreta-aqui",
		TokenPrefix:   "Bearer",
		HeaderName:    "Authorization",
		QueryParam:    "token",
		CookieName:    "auth_token",
		RequiredRoles: []string{"admin", "super_admin"},
		Optional:      false,
	},
	// Configuración para solo editores
	Editor: AuthConfig{
		SecretKey:     "tu-clave-secreta-aqui",
		TokenPrefix:   "Bearer",
		HeaderName:    "Authorization",
		QueryParam:    "token",
		CookieName:    "auth_token",
		RequiredRoles: []string{"editor", "admin", "super_admin"},
		Optional:      false,
	},
	// Configuración para usuarios autenticados
	User: AuthConfig{
		SecretKey:     "tu-clave-secreta-aqui",
		TokenPrefix:   "Bearer",
		HeaderName:    "Authorization",
		QueryParam:    "token",
		CookieName:    "auth_token",
		RequiredRoles: []string{},
		Optional:      false,
	},
}

// ========================================
// FUNCIONES DE AUTENTICACIÓN
// ========================================

// AuthWithConfig crea middleware de autenticación con configuración personalizada
// Parámetros:
//   - config: Configuración de autenticación personalizada
//
// Retorna: Middleware de autenticación configurado
func AuthWithConfig(config AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extraer token de diferentes fuentes
		token := extractToken(c, config)

		// Si no hay token y la autenticación es opcional, continuar
		if token == "" && config.Optional {
			return c.Next()
		}

		// Si no hay token y la autenticación es requerida, devolver error
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token de autenticación requerido",
				"code":  "AUTH_TOKEN_REQUIRED",
			})
		}

		// Validar y parsear el token JWT
		claims, err := validateToken(token, config.SecretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "token de autenticación inválido",
				"code":    "AUTH_TOKEN_INVALID",
				"details": err.Error(),
			})
		}

		// Verificar roles requeridos si se especifican
		if len(config.RequiredRoles) > 0 {
			if !hasRequiredRole(claims, config.RequiredRoles) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error":          "permisos insuficientes para acceder a este recurso",
					"code":           "INSUFFICIENT_PERMISSIONS",
					"required_roles": config.RequiredRoles,
				})
			}
		}

		// Almacenar información del usuario en el contexto
		c.Locals("user", claims)
		c.Locals("userID", claims["user_id"])
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		// Continuar con el siguiente middleware o handler
		return c.Next()
	}
}

// extractToken extrae el token de diferentes fuentes (header, query, cookie)
// Parámetros:
//   - c: Contexto de Fiber
//   - config: Configuración de autenticación
//
// Retorna: Token extraído o string vacío si no se encuentra
func extractToken(c *fiber.Ctx, config AuthConfig) string {
	// Intentar extraer del header Authorization
	if header := c.Get(config.HeaderName); header != "" {
		if strings.HasPrefix(header, config.TokenPrefix+" ") {
			return strings.TrimPrefix(header, config.TokenPrefix+" ")
		}
		return header
	}

	// Intentar extraer del parámetro de query
	if queryToken := c.Query(config.QueryParam); queryToken != "" {
		return queryToken
	}

	// Intentar extraer de la cookie
	if cookieToken := c.Cookies(config.CookieName); cookieToken != "" {
		return cookieToken
	}

	return ""
}

// validateToken valida y parsea un token JWT
// Parámetros:
//   - token: Token JWT a validar
//   - secretKey: Clave secreta para validar el token
//
// Retorna: Claims del token y error si ocurre alguno
func validateToken(token, secretKey string) (jwt.MapClaims, error) {
	// Parsear el token JWT
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Verificar que el token sea válido
	if !parsedToken.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// Extraer claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// hasRequiredRole verifica si el usuario tiene uno de los roles requeridos
// Parámetros:
//   - claims: Claims del token JWT
//   - requiredRoles: Lista de roles requeridos
//
// Retorna: true si el usuario tiene un rol requerido, false en caso contrario
func hasRequiredRole(claims jwt.MapClaims, requiredRoles []string) bool {
	// Obtener el rol del usuario desde los claims
	userRole, ok := claims["role"].(string)
	if !ok {
		return false
	}

	// Verificar si el rol del usuario está en la lista de roles requeridos
	for _, requiredRole := range requiredRoles {
		if userRole == requiredRole {
			return true
		}
	}

	return false
}

// ========================================
// MIDDLEWARES DE AUTENTICACIÓN PREDEFINIDOS
// ========================================

// AuthRequired aplica autenticación requerida
// Retorna: Middleware de autenticación requerida
func AuthRequired() fiber.Handler {
	return AuthWithConfig(AuthConfigs.Required)
}

// AuthOptional aplica autenticación opcional
// Retorna: Middleware de autenticación opcional
func AuthOptional() fiber.Handler {
	return AuthWithConfig(AuthConfigs.Optional)
}

// AuthAdmin aplica autenticación solo para administradores
// Retorna: Middleware de autenticación para administradores
func AuthAdmin() fiber.Handler {
	return AuthWithConfig(AuthConfigs.Admin)
}

// AuthEditor aplica autenticación para editores y administradores
// Retorna: Middleware de autenticación para editores
func AuthEditor() fiber.Handler {
	return AuthWithConfig(AuthConfigs.Editor)
}

// AuthUser aplica autenticación para usuarios autenticados
// Retorna: Middleware de autenticación para usuarios
func AuthUser() fiber.Handler {
	return AuthWithConfig(AuthConfigs.User)
}

// ========================================
// FUNCIONES DE CONFIGURACIÓN AVANZADA
// ========================================

// AuthWithSecretKey crea middleware de autenticación con clave secreta específica
// Parámetros:
//   - secretKey: Clave secreta para JWT
//
// Retorna: Middleware de autenticación con clave específica
func AuthWithSecretKey(secretKey string) fiber.Handler {
	config := AuthConfigs.Required // Usar configuración base
	config.SecretKey = secretKey
	return AuthWithConfig(config)
}

// AuthWithRoles crea middleware de autenticación con roles específicos
// Parámetros:
//   - roles: Lista de roles requeridos
//
// Retorna: Middleware de autenticación con roles específicos
func AuthWithRoles(roles ...string) fiber.Handler {
	config := AuthConfigs.Required // Usar configuración base
	config.RequiredRoles = roles
	return AuthWithConfig(config)
}

// AuthWithTokenSource crea middleware de autenticación con fuente de token específica
// Parámetros:
//   - headerName: Nombre del header para el token
//   - queryParam: Parámetro de query para el token
//   - cookieName: Nombre de la cookie para el token
//
// Retorna: Middleware de autenticación con fuente específica
func AuthWithTokenSource(headerName, queryParam, cookieName string) fiber.Handler {
	config := AuthConfigs.Required // Usar configuración base
	config.HeaderName = headerName
	config.QueryParam = queryParam
	config.CookieName = cookieName
	return AuthWithConfig(config)
}

// AuthWithPrefix crea middleware de autenticación con prefijo específico
// Parámetros:
//   - prefix: Prefijo del token (Bearer, etc.)
//
// Retorna: Middleware de autenticación con prefijo específico
func AuthWithPrefix(prefix string) fiber.Handler {
	config := AuthConfigs.Required // Usar configuración base
	config.TokenPrefix = prefix
	return AuthWithConfig(config)
}

// ========================================
// FUNCIONES DE UTILIDAD
// ========================================

// GetUserFromContext obtiene la información del usuario desde el contexto
// Parámetros:
//   - c: Contexto de Fiber
//
// Retorna: Claims del usuario o nil si no está autenticado
func GetUserFromContext(c *fiber.Ctx) jwt.MapClaims {
	if user := c.Locals("user"); user != nil {
		if claims, ok := user.(jwt.MapClaims); ok {
			return claims
		}
	}
	return nil
}

// GetUserIDFromContext obtiene el ID del usuario desde el contexto
// Parámetros:
//   - c: Contexto de Fiber
//
// Retorna: ID del usuario o nil si no está autenticado
func GetUserIDFromContext(c *fiber.Ctx) interface{} {
	return c.Locals("userID")
}

// GetUsernameFromContext obtiene el nombre de usuario desde el contexto
// Parámetros:
//   - c: Contexto de Fiber
//
// Retorna: Nombre de usuario o nil si no está autenticado
func GetUsernameFromContext(c *fiber.Ctx) interface{} {
	return c.Locals("username")
}

// GetUserRoleFromContext obtiene el rol del usuario desde el contexto
// Parámetros:
//   - c: Contexto de Fiber
//
// Retorna: Rol del usuario o nil si no está autenticado
func GetUserRoleFromContext(c *fiber.Ctx) interface{} {
	return c.Locals("role")
}

// IsAuthenticated verifica si el usuario está autenticado
// Parámetros:
//   - c: Contexto de Fiber
//
// Retorna: true si está autenticado, false en caso contrario
func IsAuthenticated(c *fiber.Ctx) bool {
	return c.Locals("user") != nil
}

// ========================================
// MIDDLEWARE DE AUTENTICACIÓN PREDEFINIDO
// ========================================

// Auth aplica autenticación por defecto (requerida)
// Retorna: Middleware de autenticación requerida
func Auth() fiber.Handler {
	return AuthRequired()
}

// GenerateToken genera un token JWT para un usuario
// Parámetros:
//   - user: Usuario para el cual generar el token
//
// Retorna: Token JWT firmado y error si ocurre alguno
func GenerateToken(user models.User) (string, error) {
	// Configurar tiempo de expiración (24 horas)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Crear las reclamaciones del token usando MapClaims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
		"sub":      user.Username,
	}

	// Crear el token con las reclamaciones
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// HashPassword hashea una contraseña usando bcrypt
// Parámetros:
//   - password: Contraseña en texto plano
//
// Retorna: Hash de la contraseña y error si ocurre
func HashPassword(password string) (string, error) {
	// Generar hash con costo 14 (balance entre seguridad y rendimiento)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword verifica si una contraseña coincide con su hash
// Parámetros:
//   - password: Contraseña en texto plano
//   - hash: Hash de la contraseña almacenado
//
// Retorna: true si la contraseña coincide, false en caso contrario
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
