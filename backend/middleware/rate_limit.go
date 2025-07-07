package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RateLimitConfig configuración para rate limiting
type RateLimitConfig struct {
	MaxRequests            int                       // Número máximo de requests
	Window                 time.Duration             // Ventana de tiempo
	KeyGenerator           func(c *fiber.Ctx) string // Función para generar la clave
	SkipSuccessfulRequests bool                      // Si saltar requests exitosos
	SkipFailedRequests     bool                      // Si saltar requests fallidos
}

// RateLimiter implementa rate limiting personalizado
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	config   RateLimitConfig
}

// NewRateLimiter crea una nueva instancia de rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		config:   config,
	}

	// Limpiar requests antiguos periódicamente
	go rl.cleanup()

	return rl
}

// cleanup limpia requests antiguos cada minuto
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.config.Window)

		for key, requests := range rl.requests {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}
			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

// Limit middleware para rate limiting
func (rl *RateLimiter) Limit() fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		// Verificar límite
		if len(validRequests) >= rl.config.MaxRequests {
			rl.mutex.Unlock()

			// Calcular tiempo de espera
			oldestRequest := validRequests[0]
			retryAfter := oldestRequest.Add(rl.config.Window).Sub(now)

			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.MaxRequests))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", oldestRequest.Add(rl.config.Window).Unix()))
			c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))

			return c.Status(429).JSON(fiber.Map{
				"error":       "Demasiadas solicitudes",
				"retry_after": retryAfter.Seconds(),
				"limit":       rl.config.MaxRequests,
				"window":      rl.config.Window.String(),
			})
		}

		// Agregar request actual
		validRequests = append(validRequests, now)
		rl.requests[key] = validRequests
		rl.mutex.Unlock()

		// Agregar headers de rate limit
		remaining := rl.config.MaxRequests - len(validRequests)
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.MaxRequests))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(rl.config.Window).Unix()))

		return c.Next()
	}
}

// KeyGenerators funciones predefinidas para generar claves de rate limiting
var KeyGenerators = struct {
	IP         func(c *fiber.Ctx) string
	UserID     func(c *fiber.Ctx) string
	Username   func(c *fiber.Ctx) string
	Endpoint   func(c *fiber.Ctx) string
	IPEndpoint func(c *fiber.Ctx) string
}{
	IP: func(c *fiber.Ctx) string {
		return c.IP()
	},
	UserID: func(c *fiber.Ctx) string {
		if userID := c.Locals("userID"); userID != nil {
			return fmt.Sprintf("user:%d", userID)
		}
		return c.IP()
	},
	Username: func(c *fiber.Ctx) string {
		if user := c.Locals("user"); user != nil {
			if u, ok := user.(interface{ GetUsername() string }); ok {
				return fmt.Sprintf("username:%s", u.GetUsername())
			}
		}
		return c.IP()
	},
	Endpoint: func(c *fiber.Ctx) string {
		return fmt.Sprintf("endpoint:%s", c.Path())
	},
	IPEndpoint: func(c *fiber.Ctx) string {
		return fmt.Sprintf("%s:%s", c.IP(), c.Path())
	},
}

// RateLimitByIP limita por IP
func RateLimitByIP(maxRequests int, window time.Duration) fiber.Handler {
	config := RateLimitConfig{
		MaxRequests:  maxRequests,
		Window:       window,
		KeyGenerator: KeyGenerators.IP,
	}
	return NewRateLimiter(config).Limit()
}

// RateLimitByUser limita por usuario
func RateLimitByUser(maxRequests int, window time.Duration) fiber.Handler {
	config := RateLimitConfig{
		MaxRequests:  maxRequests,
		Window:       window,
		KeyGenerator: KeyGenerators.UserID,
	}
	return NewRateLimiter(config).Limit()
}

// RateLimitByEndpoint limita por endpoint
func RateLimitByEndpoint(maxRequests int, window time.Duration) fiber.Handler {
	config := RateLimitConfig{
		MaxRequests:  maxRequests,
		Window:       window,
		KeyGenerator: KeyGenerators.Endpoint,
	}
	return NewRateLimiter(config).Limit()
}

// RateLimitByIPAndEndpoint limita por IP y endpoint
func RateLimitByIPAndEndpoint(maxRequests int, window time.Duration) fiber.Handler {
	config := RateLimitConfig{
		MaxRequests:  maxRequests,
		Window:       window,
		KeyGenerator: KeyGenerators.IPEndpoint,
	}
	return NewRateLimiter(config).Limit()
}

// Configuraciones predefinidas de rate limiting
var RateLimitConfigs = struct {
	Strict   RateLimitConfig
	Moderate RateLimitConfig
	Relaxed  RateLimitConfig
	API      RateLimitConfig
	Auth     RateLimitConfig
}{
	Strict: RateLimitConfig{
		MaxRequests:  10,
		Window:       1 * time.Minute,
		KeyGenerator: KeyGenerators.IP,
	},
	Moderate: RateLimitConfig{
		MaxRequests:  60,
		Window:       1 * time.Minute,
		KeyGenerator: KeyGenerators.IP,
	},
	Relaxed: RateLimitConfig{
		MaxRequests:  300,
		Window:       1 * time.Minute,
		KeyGenerator: KeyGenerators.IP,
	},
	API: RateLimitConfig{
		MaxRequests:  1000,
		Window:       1 * time.Hour,
		KeyGenerator: KeyGenerators.IP,
	},
	Auth: RateLimitConfig{
		MaxRequests:  5,
		Window:       15 * time.Minute,
		KeyGenerator: KeyGenerators.IP,
	},
}

// RateLimitStrict rate limiting estricto
func RateLimitStrict() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.Strict).Limit()
}

// RateLimitModerate rate limiting moderado
func RateLimitModerate() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.Moderate).Limit()
}

// RateLimitRelaxed rate limiting relajado
func RateLimitRelaxed() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.Relaxed).Limit()
}

// RateLimitAPI rate limiting para APIs
func RateLimitAPI() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.API).Limit()
}

// RateLimitAuth rate limiting para autenticación
func RateLimitAuth() fiber.Handler {
	return NewRateLimiter(RateLimitConfigs.Auth).Limit()
}
