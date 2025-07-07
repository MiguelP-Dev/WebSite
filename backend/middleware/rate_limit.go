package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RateLimitConfig configuración para el sistema de rate limiting
type RateLimitConfig struct {
	MaxRequests            int                       // Número máximo de requests permitidos
	Window                 time.Duration             // Ventana de tiempo para el límite
	KeyGenerator           func(c *fiber.Ctx) string // Función para generar la clave única
	SkipSuccessfulRequests bool                      // Si saltar requests exitosos
	SkipFailedRequests     bool                      // Si saltar requests fallidos
}

// RateLimiter implementa el sistema de rate limiting personalizado
type RateLimiter struct {
	requests map[string][]time.Time // Mapa de claves a timestamps de requests
	mutex    sync.RWMutex           // Mutex para acceso concurrente seguro
	config   RateLimitConfig        // Configuración del rate limiter
}

// NewRateLimiter crea una nueva instancia de rate limiter
// Parámetros:
//   - config: Configuración del rate limiter
//
// Retorna: Instancia configurada del rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	// Crear nueva instancia
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		config:   config,
	}

	// Iniciar limpieza automática en background
	go rl.cleanup()

	return rl
}

// cleanup limpia requests antiguos periódicamente para evitar fugas de memoria
func (rl *RateLimiter) cleanup() {
	// Crear ticker que se ejecuta cada minuto
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// Ejecutar limpieza en cada tick
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.config.Window)

		// Limpiar requests antiguos de cada clave
		for key, requests := range rl.requests {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}
			// Si no hay requests válidos, eliminar la clave
			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

// Limit middleware para aplicar rate limiting
// Retorna: Middleware de Fiber que aplica rate limiting
func (rl *RateLimiter) Limit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generar clave única para este request
		key := rl.config.KeyGenerator(c)

		rl.mutex.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.config.Window)

		// Filtrar requests antiguos
		var validRequests []time.Time
		if requests, exists := rl.requests[key]; exists {
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}
		}

		// Verificar si se excedió el límite
		if len(validRequests) >= rl.config.MaxRequests {
			rl.mutex.Unlock()

			// Calcular tiempo de espera
			oldestRequest := validRequests[0]
			retryAfter := oldestRequest.Add(rl.config.Window).Sub(now)

			// Configurar headers de rate limiting
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.MaxRequests))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", oldestRequest.Add(rl.config.Window).Unix()))
			c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))

			// Devolver error de rate limiting
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "demasiadas solicitudes, intente nuevamente más tarde",
				"code":        "RATE_LIMIT_EXCEEDED",
				"retry_after": retryAfter.Seconds(),
				"limit":       rl.config.MaxRequests,
				"window":      rl.config.Window.String(),
			})
		}

		// Agregar request actual a la lista
		validRequests = append(validRequests, now)
		rl.requests[key] = validRequests
		rl.mutex.Unlock()

		// Configurar headers de rate limiting
		remaining := rl.config.MaxRequests - len(validRequests)
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.MaxRequests))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(rl.config.Window).Unix()))

		// Continuar con el siguiente middleware o handler
		return c.Next()
	}
}

// ========================================
// GENERADORES DE CLAVES PREDEFINIDOS
// ========================================

// KeyGenerators contiene funciones predefinidas para generar claves de rate limiting
var KeyGenerators = struct {
	IP         func(c *fiber.Ctx) string // Por dirección IP
	UserID     func(c *fiber.Ctx) string // Por ID de usuario
	Username   func(c *fiber.Ctx) string // Por nombre de usuario
	Endpoint   func(c *fiber.Ctx) string // Por endpoint
	IPEndpoint func(c *fiber.Ctx) string // Por IP y endpoint combinados
}{
	// Generar clave por dirección IP
	IP: func(c *fiber.Ctx) string {
		return c.IP()
	},
	// Generar clave por ID de usuario (si está autenticado)
	UserID: func(c *fiber.Ctx) string {
		if userID := c.Locals("userID"); userID != nil {
			return fmt.Sprintf("user:%d", userID)
		}
		return c.IP()
	},
	// Generar clave por nombre de usuario (si está autenticado)
	Username: func(c *fiber.Ctx) string {
		if user := c.Locals("user"); user != nil {
			if u, ok := user.(interface{ GetUsername() string }); ok {
				return fmt.Sprintf("username:%s", u.GetUsername())
			}
		}
		return c.IP()
	},
	// Generar clave por endpoint
	Endpoint: func(c *fiber.Ctx) string {
		return fmt.Sprintf("endpoint:%s", c.Path())
	},
	// Generar clave por IP y endpoint combinados
	IPEndpoint: func(c *fiber.Ctx) string {
		return fmt.Sprintf("%s:%s", c.IP(), c.Path())
	},
}

// ========================================
// FUNCIONES DE RATE LIMITING POR TIPO
// ========================================

// RateLimitByIP limita requests por dirección IP
// Parámetros:
//   - maxRequests: Número máximo de requests
//   - window: Ventana de tiempo
//
// Retorna: Middleware de rate limiting por IP
func RateLimitByIP(maxRequests int, window time.Duration) fiber.Handler {
	config := RateLimitConfig{
		MaxRequests:  maxRequests,
		Window:       window,
		KeyGenerator: KeyGenerators.IP,
	}
	return NewRateLimiter(config).Limit()
}

// RateLimitByUser limita requests por usuario autenticado
// Parámetros:
//   - maxRequests: Número máximo de requests
//   - window: Ventana de tiempo
//
// Retorna: Middleware de rate limiting por usuario
func RateLimitByUser(maxRequests int, window time.Duration) fiber.Handler {
	config := RateLimitConfig{
		MaxRequests:  maxRequests,
		Window:       window,
		KeyGenerator: KeyGenerators.UserID,
	}
	return NewRateLimiter(config).Limit()
}

// RateLimitByEndpoint limita requests por endpoint específico
// Parámetros:
//   - maxRequests: Número máximo de requests
//   - window: Ventana de tiempo
//
// Retorna: Middleware de rate limiting por endpoint
func RateLimitByEndpoint(maxRequests int, window time.Duration) fiber.Handler {
	config := RateLimitConfig{
		MaxRequests:  maxRequests,
		Window:       window,
		KeyGenerator: KeyGenerators.Endpoint,
	}
	return NewRateLimiter(config).Limit()
}

// RateLimitByIPAndEndpoint limita requests por IP y endpoint combinados
// Parámetros:
//   - maxRequests: Número máximo de requests
//   - window: Ventana de tiempo
//
// Retorna: Middleware de rate limiting por IP y endpoint
func RateLimitByIPAndEndpoint(maxRequests int, window time.Duration) fiber.Handler {
	config := RateLimitConfig{
		MaxRequests:  maxRequests,
		Window:       window,
		KeyGenerator: KeyGenerators.IPEndpoint,
	}
	return NewRateLimiter(config).Limit()
}

// ========================================
// CONFIGURACIONES PREDEFINIDAS
// ========================================

// RateLimitConfigs contiene configuraciones predefinidas de rate limiting
var RateLimitConfigs = struct {
	Strict   RateLimitConfig // Muy restrictivo
	Moderate RateLimitConfig // Moderado
	Relaxed  RateLimitConfig // Relajado
	API      RateLimitConfig // Para APIs
	Auth     RateLimitConfig // Para autenticación
}{
	// Configuración estricta: 10 requests por minuto
	Strict: RateLimitConfig{
		MaxRequests:  10,
		Window:       1 * time.Minute,
		KeyGenerator: KeyGenerators.IP,
	},
	// Configuración moderada: 60 requests por minuto
	Moderate: RateLimitConfig{
		MaxRequests:  60,
		Window:       1 * time.Minute,
		KeyGenerator: KeyGenerators.IP,
	},
	// Configuración relajada: 300 requests por minuto
	Relaxed: RateLimitConfig{
		MaxRequests:  300,
		Window:       1 * time.Minute,
		KeyGenerator: KeyGenerators.IP,
	},
	// Configuración para APIs: 1000 requests por hora
	API: RateLimitConfig{
		MaxRequests:  1000,
		Window:       1 * time.Hour,
		KeyGenerator: KeyGenerators.IP,
	},
	// Configuración para autenticación: 5 requests por 15 minutos
	Auth: RateLimitConfig{
		MaxRequests:  5,
		Window:       15 * time.Minute,
		KeyGenerator: KeyGenerators.IP,
	},
}

// ========================================
// MIDDLEWARES PREDEFINIDOS
// ========================================

// RateLimitStrict aplica rate limiting estricto
// Retorna: Middleware con límite de 10 requests por minuto
func RateLimitStrict() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.Strict).Limit()
}

// RateLimitModerate aplica rate limiting moderado
// Retorna: Middleware con límite de 60 requests por minuto
func RateLimitModerate() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.Moderate).Limit()
}

// RateLimitRelaxed aplica rate limiting relajado
// Retorna: Middleware con límite de 300 requests por minuto
func RateLimitRelaxed() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.Relaxed).Limit()
}

// RateLimitAPI aplica rate limiting para APIs
// Retorna: Middleware con límite de 1000 requests por hora
func RateLimitAPI() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.API).Limit()
}

// RateLimitAuth aplica rate limiting para autenticación
// Retorna: Middleware con límite de 5 requests por 15 minutos
func RateLimitAuth() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.Auth).Limit()
}
