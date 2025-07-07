package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RoleMiddleware verifica que el usuario tenga el rol requerido
// Parámetros:
//   - requiredRoles: Lista de roles permitidos
//
// Retorna: Middleware de Fiber que valida roles
func RoleMiddleware(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener el rol del usuario del contexto
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "rol de usuario no encontrado en el contexto",
				"code":  "ROLE_NOT_FOUND",
			})
		}

		// Convertir el rol a string
		role := userRole.(string)

		// Verificar si el rol del usuario está en los roles requeridos
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				return c.Next()
			}
		}

		// Si no tiene el rol requerido, devolver error
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":          "permisos insuficientes para acceder a este recurso",
			"code":           "INSUFFICIENT_PERMISSIONS",
			"required_roles": requiredRoles,
			"user_role":      role,
		})
	}
}

// PermissionMiddleware verifica permisos específicos del usuario
// Parámetros:
//   - requiredPermissions: Lista de permisos requeridos
//
// Retorna: Middleware de Fiber que valida permisos
func PermissionMiddleware(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener el rol del usuario del contexto
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "rol de usuario no encontrado en el contexto",
				"code":  "ROLE_NOT_FOUND",
			})
		}

		// Convertir el rol a string
		role := userRole.(string)

		// Obtener permisos del rol
		permissions := getRolePermissions(role)

		// Verificar si el usuario tiene todos los permisos requeridos
		for _, requiredPermission := range requiredPermissions {
			if !hasPermission(permissions, requiredPermission) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error":               "permiso requerido no encontrado",
					"code":                "PERMISSION_DENIED",
					"required_permission": requiredPermission,
					"user_permissions":    permissions,
				})
			}
		}

		// Si tiene todos los permisos, continuar
		return c.Next()
	}
}

// getRolePermissions retorna los permisos para un rol específico
// Parámetros:
//   - role: Rol del usuario
//
// Retorna: Lista de permisos asociados al rol
func getRolePermissions(role string) []string {
	// Mapeo de roles a permisos
	permissions := map[string][]string{
		"super_admin": {
			"users:read", "users:write", "users:delete",
			"slides:read", "slides:write", "slides:delete",
			"categories:read", "categories:write", "categories:delete",
			"products:read", "products:write", "products:delete",
			"contacts:read", "contacts:write", "contacts:delete",
			"config:read", "config:write",
			"system:admin",
		},
		"admin": {
			"users:read", "users:write",
			"slides:read", "slides:write", "slides:delete",
			"categories:read", "categories:write", "categories:delete",
			"products:read", "products:write", "products:delete",
			"contacts:read", "contacts:write", "contacts:delete",
			"config:read", "config:write",
		},
		"editor": {
			"slides:read", "slides:write",
			"categories:read", "categories:write",
			"products:read", "products:write",
			"contacts:read", "contacts:write",
			"config:read",
		},
		"viewer": {
			"slides:read",
			"categories:read",
			"products:read",
			"contacts:read",
			"config:read",
		},
	}

	// Retornar permisos del rol si existe, sino lista vacía
	if perms, exists := permissions[role]; exists {
		return perms
	}

	return []string{}
}

// hasPermission verifica si un usuario tiene un permiso específico
// Parámetros:
//   - userPermissions: Lista de permisos del usuario
//   - requiredPermission: Permiso requerido
//
// Retorna: true si tiene el permiso, false en caso contrario
func hasPermission(userPermissions []string, requiredPermission string) bool {
	for _, permission := range userPermissions {
		// Verificar permiso exacto
		if permission == requiredPermission {
			return true
		}
		// Verificar permisos wildcard (ej: "users:*" para "users:read")
		if strings.HasSuffix(permission, ":*") {
			basePermission := strings.TrimSuffix(permission, ":*")
			if strings.HasPrefix(requiredPermission, basePermission+":") {
				return true
			}
		}
	}
	return false
}

// ========================================
// MIDDLEWARES PREDEFINIDOS PARA ROLES
// ========================================

// IsSuperAdmin verifica si el usuario es super administrador
// Retorna: Middleware que solo permite super_admin
func IsSuperAdmin() fiber.Handler {
	return RoleMiddleware("super_admin")
}

// IsAdmin verifica si el usuario es administrador o superior
// Retorna: Middleware que permite super_admin y admin
func IsAdmin() fiber.Handler {
	return RoleMiddleware("super_admin", "admin")
}

// IsEditor verifica si el usuario es editor o superior
// Retorna: Middleware que permite super_admin, admin y editor
func IsEditor() fiber.Handler {
	return RoleMiddleware("super_admin", "admin", "editor")
}

// ========================================
// MIDDLEWARES PREDEFINIDOS PARA PERMISOS
// ========================================

// CanManageUsers verifica si el usuario puede gestionar usuarios
// Retorna: Middleware que requiere permisos de gestión de usuarios
func CanManageUsers() fiber.Handler {
	return PermissionMiddleware("users:read", "users:write")
}

// CanManageContent verifica si el usuario puede gestionar contenido
// Retorna: Middleware que requiere permisos de gestión de contenido
func CanManageContent() fiber.Handler {
	return PermissionMiddleware("slides:write", "categories:write", "products:write")
}

// CanManageSystem verifica si el usuario puede gestionar el sistema
// Retorna: Middleware que requiere permisos de administración del sistema
func CanManageSystem() fiber.Handler {
	return PermissionMiddleware("system:admin")
}
