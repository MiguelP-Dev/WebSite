package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// LogConfig configuración para el sistema de logging
type LogConfig struct {
	Format     string // Formato del log (json, text, custom)
	Output     string // Archivo de salida (stdout, stderr, archivo)
	TimeFormat string // Formato de tiempo
	TimeZone   string // Zona horaria
	Colors     bool   // Si usar colores en la salida
	IP         bool   // Si incluir IP del cliente
	UserAgent  bool   // Si incluir User-Agent
	Referer    bool   // Si incluir Referer
	Latency    bool   // Si incluir latencia
	Status     bool   // Si incluir código de estado
	Method     bool   // Si incluir método HTTP
	Path       bool   // Si incluir ruta
	Query      bool   // Si incluir query parameters
	Body       bool   // Si incluir body del request
	Headers    bool   // Si incluir headers
}

// ========================================
// CONFIGURACIONES PREDEFINIDAS DE LOGGING
// ========================================

// LogConfigs contiene configuraciones predefinidas de logging
var LogConfigs = struct {
	Development LogConfig // Para desarrollo local
	Production  LogConfig // Para producción
	API         LogConfig // Para APIs
	Security    LogConfig // Para logs de seguridad
	Minimal     LogConfig // Configuración mínima
}{
	// Configuración para desarrollo: logs detallados con colores
	Development: LogConfig{
		Format:     "json",
		Output:     "stdout",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Colors:     true,
		IP:         true,
		UserAgent:  true,
		Referer:    true,
		Latency:    true,
		Status:     true,
		Method:     true,
		Path:       true,
		Query:      true,
		Body:       false, // No incluir body por defecto en desarrollo
		Headers:    false, // No incluir headers por defecto
	},
	// Configuración para producción: logs estructurados
	Production: LogConfig{
		Format:     "json",
		Output:     "stdout",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		TimeZone:   "UTC",
		Colors:     false,
		IP:         true,
		UserAgent:  true,
		Referer:    true,
		Latency:    true,
		Status:     true,
		Method:     true,
		Path:       true,
		Query:      false, // No incluir query en producción por seguridad
		Body:       false, // No incluir body en producción por seguridad
		Headers:    false, // No incluir headers en producción por seguridad
	},
	// Configuración para APIs: logs específicos para APIs
	API: LogConfig{
		Format:     "json",
		Output:     "stdout",
		TimeFormat: "2006-01-02T15:04:05.000Z",
		TimeZone:   "UTC",
		Colors:     false,
		IP:         true,
		UserAgent:  true,
		Referer:    false,
		Latency:    true,
		Status:     true,
		Method:     true,
		Path:       true,
		Query:      true,
		Body:       false,
		Headers:    false,
	},
	// Configuración para logs de seguridad: muy detallada
	Security: LogConfig{
		Format:     "json",
		Output:     "security.log",
		TimeFormat: "2006-01-02T15:04:05.000Z",
		TimeZone:   "UTC",
		Colors:     false,
		IP:         true,
		UserAgent:  true,
		Referer:    true,
		Latency:    true,
		Status:     true,
		Method:     true,
		Path:       true,
		Query:      true,
		Body:       true, // Incluir body para análisis de seguridad
		Headers:    true, // Incluir headers para análisis de seguridad
	},
	// Configuración mínima: solo información esencial
	Minimal: LogConfig{
		Format:     "text",
		Output:     "stdout",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
		Colors:     true,
		IP:         false,
		UserAgent:  false,
		Referer:    false,
		Latency:    false,
		Status:     true,
		Method:     true,
		Path:       true,
		Query:      false,
		Body:       false,
		Headers:    false,
	},
}

// ========================================
// FUNCIONES DE CONFIGURACIÓN DE LOGGING
// ========================================

// LogWithConfig crea middleware de logging con configuración personalizada
// Parámetros:
//   - config: Configuración de logging personalizada
//
// Retorna: Middleware de logging configurado
func LogWithConfig(config LogConfig) fiber.Handler {
	// Configurar opciones de logger
	loggerConfig := logger.Config{
		Format:     generateLogFormat(config),
		TimeFormat: config.TimeFormat,
		TimeZone:   config.TimeZone,
		Output:     getLogOutput(config.Output),
	}

	return logger.New(loggerConfig)
}

// generateLogFormat genera el formato de log basado en la configuración
// Parámetros:
//   - config: Configuración de logging
//
// Retorna: String con el formato de log
func generateLogFormat(config LogConfig) string {
	if config.Format == "json" {
		return generateJSONFormat(config)
	}
	return generateTextFormat(config)
}

// generateJSONFormat genera formato JSON para logging
// Parámetros:
//   - config: Configuración de logging
//
// Retorna: String con formato JSON
func generateJSONFormat(config LogConfig) string {
	format := `{"timestamp":"${time}","level":"${status}","method":"${method}","path":"${path}"`

	if config.IP {
		format += `,"ip":"${ip}"`
	}
	if config.UserAgent {
		format += `,"user_agent":"${ua}"`
	}
	if config.Referer {
		format += `,"referer":"${referer}"`
	}
	if config.Latency {
		format += `,"latency":"${latency}"`
	}
	if config.Query {
		format += `,"query":"${query}"`
	}
	if config.Body {
		format += `,"body":"${body}"`
	}
	if config.Headers {
		format += `,"headers":"${headers}"`
	}

	format += `}`
	return format
}

// generateTextFormat genera formato de texto para logging
// Parámetros:
//   - config: Configuración de logging
//
// Retorna: String con formato de texto
func generateTextFormat(config LogConfig) string {
	format := "${time} | ${status} | ${method} | ${path}"

	if config.IP {
		format += ` | ${ip}`
	}
	if config.UserAgent {
		format += ` | ${ua}`
	}
	if config.Referer {
		format += ` | ${referer}`
	}
	if config.Latency {
		format += ` | ${latency}`
	}
	if config.Query {
		format += ` | ${query}`
	}
	if config.Body {
		format += ` | ${body}`
	}
	if config.Headers {
		format += ` | ${headers}`
	}

	return format
}

// getLogOutput obtiene el writer de salida para los logs
// Parámetros:
//   - output: Tipo de salida (stdout, stderr, archivo)
//
// Retorna: Writer para la salida de logs
func getLogOutput(output string) *os.File {
	switch output {
	case "stderr":
		return os.Stderr
	case "security.log":
		// Abrir archivo de log de seguridad
		file, err := os.OpenFile("security.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Si no se puede abrir el archivo, usar stderr
			fmt.Fprintf(os.Stderr, "error al abrir archivo de log de seguridad: %v\n", err)
			return os.Stderr
		}
		return file
	default:
		return os.Stdout
	}
}

// ========================================
// MIDDLEWARES DE LOGGING PREDEFINIDOS
// ========================================

// LogDevelopment aplica configuración de logging para desarrollo
// Retorna: Middleware de logging configurado para desarrollo
func LogDevelopment() fiber.Handler {
	return LogWithConfig(LogConfigs.Development)
}

// LogProduction aplica configuración de logging para producción
// Retorna: Middleware de logging configurado para producción
func LogProduction() fiber.Handler {
	return LogWithConfig(LogConfigs.Production)
}

// LogAPI aplica configuración de logging para APIs
// Retorna: Middleware de logging configurado para APIs
func LogAPI() fiber.Handler {
	return LogWithConfig(LogConfigs.API)
}

// LogSecurity aplica configuración de logging para seguridad
// Retorna: Middleware de logging configurado para seguridad
func LogSecurity() fiber.Handler {
	return LogWithConfig(LogConfigs.Security)
}

// LogMinimal aplica configuración de logging mínima
// Retorna: Middleware de logging con configuración mínima
func LogMinimal() fiber.Handler {
	return LogWithConfig(LogConfigs.Minimal)
}

// ========================================
// FUNCIONES DE CONFIGURACIÓN AVANZADA
// ========================================

// LogWithFormat crea middleware de logging con formato específico
// Parámetros:
//   - format: Formato del log (json, text)
//
// Retorna: Middleware de logging con formato específico
func LogWithFormat(format string) fiber.Handler {
	config := LogConfigs.Development // Usar configuración base
	config.Format = format
	return LogWithConfig(config)
}

// LogWithOutput crea middleware de logging con salida específica
// Parámetros:
//   - output: Tipo de salida (stdout, stderr, archivo)
//
// Retorna: Middleware de logging con salida específica
func LogWithOutput(output string) fiber.Handler {
	config := LogConfigs.Development // Usar configuración base
	config.Output = output
	return LogWithConfig(config)
}

// LogWithTimeFormat crea middleware de logging con formato de tiempo específico
// Parámetros:
//   - timeFormat: Formato de tiempo
//
// Retorna: Middleware de logging con formato de tiempo específico
func LogWithTimeFormat(timeFormat string) fiber.Handler {
	config := LogConfigs.Development // Usar configuración base
	config.TimeFormat = timeFormat
	return LogWithConfig(config)
}

// LogWithFields crea middleware de logging con campos específicos
// Parámetros:
//   - fields: Campos a incluir en el log
//
// Retorna: Middleware de logging con campos específicos
func LogWithFields(fields map[string]bool) fiber.Handler {
	config := LogConfigs.Development // Usar configuración base

	// Aplicar campos específicos
	if ip, ok := fields["ip"]; ok {
		config.IP = ip
	}
	if ua, ok := fields["user_agent"]; ok {
		config.UserAgent = ua
	}
	if referer, ok := fields["referer"]; ok {
		config.Referer = referer
	}
	if latency, ok := fields["latency"]; ok {
		config.Latency = latency
	}
	if query, ok := fields["query"]; ok {
		config.Query = query
	}
	if body, ok := fields["body"]; ok {
		config.Body = body
	}
	if headers, ok := fields["headers"]; ok {
		config.Headers = headers
	}

	return LogWithConfig(config)
}

// ========================================
// MIDDLEWARE DE LOGGING PREDEFINIDO
// ========================================

// Log aplica configuración de logging por defecto (desarrollo)
// Retorna: Middleware de logging con configuración de desarrollo
func Log() fiber.Handler {
	return LogDevelopment()
}
