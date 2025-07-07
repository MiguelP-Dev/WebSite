package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSConfig configuración para el middleware de CORS
type CORSConfig struct {
	AllowOrigins     []string // Lista de orígenes permitidos
	AllowMethods     []string // Métodos HTTP permitidos
	AllowHeaders     []string // Headers permitidos
	AllowCredentials bool     // Si permitir credenciales
	ExposeHeaders    []string // Headers expuestos al cliente
	MaxAge           int      // Tiempo máximo de cache en segundos
}

// ========================================
// CONFIGURACIONES PREDEFINIDAS DE CORS
// ========================================

// CORSConfigs contiene configuraciones predefinidas de CORS
var CORSConfigs = struct {
	Development CORSConfig // Para desarrollo local
	Production  CORSConfig // Para producción
	API         CORSConfig // Para APIs públicas
	Restricted  CORSConfig // Para APIs restringidas
}{
	// Configuración para desarrollo: permite todos los orígenes
	Development: CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		MaxAge:           86400, // 24 horas
	},
	// Configuración para producción: orígenes específicos
	Production: CORSConfig{
		AllowOrigins:     []string{"https://tudominio.com", "https://www.tudominio.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		MaxAge:           86400, // 24 horas
	},
	// Configuración para APIs públicas: más permisiva
	API: CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-API-Key"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count", "X-RateLimit-Limit", "X-RateLimit-Remaining"},
		MaxAge:           3600, // 1 hora
	},
	// Configuración restringida: solo orígenes específicos
	Restricted: CORSConfig{
		AllowOrigins:     []string{"https://admin.tudominio.com", "https://dashboard.tudominio.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           1800, // 30 minutos
	},
}

// ========================================
// FUNCIONES DE CONFIGURACIÓN DE CORS
// ========================================

// CORSWithConfig crea middleware de CORS con configuración personalizada
// Parámetros:
//   - config: Configuración de CORS personalizada
//
// Retorna: Middleware de CORS configurado
func CORSWithConfig(config CORSConfig) fiber.Handler {
	// Configurar opciones de CORS
	corsConfig := cors.Config{
		AllowOrigins:     config.AllowOrigins[0],                 // Fiber espera un string, no un slice
		AllowMethods:     strings.Join(config.AllowMethods, ","), // Convertir slice a string
		AllowHeaders:     strings.Join(config.AllowHeaders, ","), // Convertir slice a string
		AllowCredentials: config.AllowCredentials,
		ExposeHeaders:    strings.Join(config.ExposeHeaders, ","), // Convertir slice a string
		MaxAge:           config.MaxAge,
	}

	// Si se permite cualquier origen, usar "*"
	if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
		corsConfig.AllowOrigins = "*"
	} else {
		// Para múltiples orígenes, usar el primero como principal
		// Nota: Fiber tiene limitaciones con múltiples orígenes
		corsConfig.AllowOrigins = config.AllowOrigins[0]
	}

	return cors.New(corsConfig)
}

// CORSDevelopment aplica configuración de CORS para desarrollo
// Retorna: Middleware de CORS configurado para desarrollo
func CORSDevelopment() fiber.Handler {
	return CORSWithConfig(CORSConfigs.Development)
}

// CORSProduction aplica configuración de CORS para producción
// Retorna: Middleware de CORS configurado para producción
func CORSProduction() fiber.Handler {
	return CORSWithConfig(CORSConfigs.Production)
}

// CORSAPI aplica configuración de CORS para APIs públicas
// Retorna: Middleware de CORS configurado para APIs
func CORSAPI() fiber.Handler {
	return CORSWithConfig(CORSConfigs.API)
}

// CORSRestricted aplica configuración de CORS restringida
// Retorna: Middleware de CORS configurado de forma restringida
func CORSRestricted() fiber.Handler {
	return CORSWithConfig(CORSConfigs.Restricted)
}

// ========================================
// FUNCIONES DE CONFIGURACIÓN AVANZADA
// ========================================

// CORSWithOrigins crea middleware de CORS con orígenes específicos
// Parámetros:
//   - origins: Lista de orígenes permitidos
//
// Retorna: Middleware de CORS con orígenes específicos
func CORSWithOrigins(origins ...string) fiber.Handler {
	config := CORSConfigs.Development // Usar configuración base
	config.AllowOrigins = origins
	return CORSWithConfig(config)
}

// CORSWithMethods crea middleware de CORS con métodos específicos
// Parámetros:
//   - methods: Lista de métodos HTTP permitidos
//
// Retorna: Middleware de CORS con métodos específicos
func CORSWithMethods(methods ...string) fiber.Handler {
	config := CORSConfigs.Development // Usar configuración base
	config.AllowMethods = methods
	return CORSWithConfig(config)
}

// CORSWithHeaders crea middleware de CORS con headers específicos
// Parámetros:
//   - headers: Lista de headers permitidos
//
// Retorna: Middleware de CORS con headers específicos
func CORSWithHeaders(headers ...string) fiber.Handler {
	config := CORSConfigs.Development // Usar configuración base
	config.AllowHeaders = headers
	return CORSWithConfig(config)
}

// CORSWithCredentials crea middleware de CORS que permite credenciales
// Parámetros:
//   - allowCredentials: Si permitir credenciales
//
// Retorna: Middleware de CORS con configuración de credenciales
func CORSWithCredentials(allowCredentials bool) fiber.Handler {
	config := CORSConfigs.Development // Usar configuración base
	config.AllowCredentials = allowCredentials
	return CORSWithConfig(config)
}

// ========================================
// MIDDLEWARE DE CORS PREDEFINIDO
// ========================================

// CORS aplica configuración de CORS por defecto (desarrollo)
// Retorna: Middleware de CORS con configuración de desarrollo
func CORS() fiber.Handler {
	return CORSDevelopment()
}
