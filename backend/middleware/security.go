package middleware

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SecurityConfig configuración de seguridad
type SecurityConfig struct {
	EnableHTTPS     bool
	EnableHSTS      bool
	EnableCSP       bool
	EnableXSS       bool
	EnableFrameDeny bool
	EnableNoSniff   bool
	EnableReferrer  bool
}

// SecurityMiddleware aplica configuraciones de seguridad
func SecurityMiddleware(config SecurityConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// HTTPS redirect
		if config.EnableHTTPS && c.Protocol() == "http" {
			host := c.Hostname()
			path := c.Path()
			query := c.Context().QueryArgs().String()

			redirectURL := fmt.Sprintf("https://%s%s", host, path)
			if query != "" {
				redirectURL += "?" + query
			}

			return c.Redirect(redirectURL, 301)
		}

		// Security headers
		if config.EnableHSTS {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		if config.EnableCSP {
			c.Set("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; "+
					"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://cdn.jsdelivr.net; "+
					"font-src 'self' https://fonts.gstatic.com https://cdn.jsdelivr.net; "+
					"img-src 'self' data: https:; "+
					"connect-src 'self' https:; "+
					"frame-src 'self'; "+
					"object-src 'none'; "+
					"base-uri 'self'; "+
					"form-action 'self'")
		}

		if config.EnableXSS {
			c.Set("X-XSS-Protection", "1; mode=block")
		}

		if config.EnableFrameDeny {
			c.Set("X-Frame-Options", "DENY")
		}

		if config.EnableNoSniff {
			c.Set("X-Content-Type-Options", "nosniff")
		}

		if config.EnableReferrer {
			c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		}

		// Additional security headers
		c.Set("X-Permitted-Cross-Domain-Policies", "none")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		return c.Next()
	}
}

// DefaultSecurityConfig configuración de seguridad por defecto
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		EnableHTTPS:     getEnvBool("ENABLE_HTTPS", false),
		EnableHSTS:      getEnvBool("ENABLE_HSTS", true),
		EnableCSP:       getEnvBool("ENABLE_CSP", true),
		EnableXSS:       getEnvBool("ENABLE_XSS", true),
		EnableFrameDeny: getEnvBool("ENABLE_FRAME_DENY", true),
		EnableNoSniff:   getEnvBool("ENABLE_NO_SNIFF", true),
		EnableReferrer:  getEnvBool("ENABLE_REFERRER", true),
	}
}

// getEnvBool obtiene una variable de entorno como booleano
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true" || value == "1"
}

// TLSServerConfig configuración del servidor TLS
func TLSServerConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256, tls.CurveP384},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

// SecurityHeadersMiddleware aplica solo headers de seguridad
func SecurityHeadersMiddleware() fiber.Handler {
	return SecurityMiddleware(DefaultSecurityConfig())
}

// HTTPSRedirectMiddleware solo redirige HTTP a HTTPS
func HTTPSRedirectMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Protocol() == "http" {
			host := c.Hostname()
			path := c.Path()
			query := c.Context().QueryArgs().String()

			redirectURL := fmt.Sprintf("https://%s%s", host, path)
			if query != "" {
				redirectURL += "?" + query
			}

			return c.Redirect(redirectURL, 301)
		}
		return c.Next()
	}
}
