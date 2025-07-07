package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RoleMiddleware verifica que el usuario tenga el rol requerido
func RoleMiddleware(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "Acceso denegado: rol no encontrado",
			})
		}

		role := userRole.(string)

		// Verificar si el rol del usuario está en los roles requeridos
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error":          "Acceso denegado: permisos insuficientes",
			"required_roles": requiredRoles,
			"user_role":      role,
		})
	}
}

// PermissionMiddleware verifica permisos específicos
func PermissionMiddleware(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "Acceso denegado: rol no encontrado",
			})
		}

		role := userRole.(string)

		// Obtener permisos del rol
		permissions := getRolePermissions(role)

		// Verificar si el usuario tiene todos los permisos requeridos
		for _, requiredPermission := range requiredPermissions {
			if !hasPermission(permissions, requiredPermission) {
				return c.Status(403).JSON(fiber.Map{
					"error":               "Acceso denegado: permiso requerido",
					"required_permission": requiredPermission,
					"user_permissions":    permissions,
				})
			}
		}

		return c.Next()
	}
}

// getRolePermissions retorna los permisos para un rol específico
func getRolePermissions(role string) []string {
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

	if perms, exists := permissions[role]; exists {
		return perms
	}

	return []string{}
}

// hasPermission verifica si un usuario tiene un permiso específico
func hasPermission(userPermissions []string, requiredPermission string) bool {
	for _, permission := range userPermissions {
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

// IsSuperAdmin verifica si el usuario es super administrador
func IsSuperAdmin() fiber.Handler {
	return RoleMiddleware("super_admin")
}

// IsAdmin verifica si el usuario es administrador o superior
func IsAdmin() fiber.Handler {
	return RoleMiddleware("super_admin", "admin")
}

// IsEditor verifica si el usuario es editor o superior
func IsEditor() fiber.Handler {
	return RoleMiddleware("super_admin", "admin", "editor")
}

// CanManageUsers verifica si el usuario puede gestionar usuarios
func CanManageUsers() fiber.Handler {
	return PermissionMiddleware("users:read", "users:write")
}

// CanManageContent verifica si el usuario puede gestionar contenido
func CanManageContent() fiber.Handler {
	return PermissionMiddleware("slides:write", "categories:write", "products:write")
}

// CanManageSystem verifica si el usuario puede gestionar el sistema
func CanManageSystem() fiber.Handler {
	return PermissionMiddleware("system:admin")
}
