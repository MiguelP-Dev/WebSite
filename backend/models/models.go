package models

import (
	"time"
)

type SiteConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;not null" json:"key" validate:"required,min=1,max=100"`
	Value     string    `gorm:"type:text;not null" json:"value" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Slide struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ImageURL  string    `json:"image_url" validate:"required,url"`
	Title     string    `json:"title" validate:"required,min=1,max=200"`
	Subtitle  string    `json:"subtitle" validate:"max=500"`
	Link      string    `json:"link" validate:"omitempty,url"`
	Order     int       `json:"order" validate:"gte=0"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name" validate:"required,min=1,max=100"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug" validate:"required,min=1,max=100,alphanum"`
	Description string    `gorm:"type:text" json:"description" validate:"max=1000"`
	ImageURL    string    `json:"image_url" validate:"omitempty,url"`
	Order       int       `json:"order" validate:"gte=0"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `json:"category_id" validate:"required"`
	Name        string    `gorm:"not null" json:"name" validate:"required,min=1,max=200"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug" validate:"required,min=1,max=200,alphanum"`
	Description string    `gorm:"type:text" json:"description" validate:"max=2000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
	ImageURLs   []string  `gorm:"type:text[]" json:"image_urls" validate:"omitempty,dive,url"`
	Features    string    `gorm:"type:text" json:"features" validate:"max=1000"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

type ContactInfo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Type      string    `gorm:"not null" json:"type" validate:"required,oneof=whatsapp facebook instagram email phone website"`
	Value     string    `gorm:"not null" json:"value" validate:"required,min=1,max=200"`
	Icon      string    `json:"icon" validate:"max=50"`
	Order     int       `json:"order" validate:"gte=0"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User model for CMS authentication
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50,alphanum"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"-" validate:"required,min=6"`
	Role      string    `gorm:"not null;default:'viewer'" json:"role" validate:"required,oneof=super_admin admin editor viewer"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUsername retorna el nombre de usuario
func (u *User) GetUsername() string {
	return u.Username
}

// HasRole verifica si el usuario tiene un rol específico
func (u *User) HasRole(role string) bool {
	return u.Role == role
}

// HasAnyRole verifica si el usuario tiene alguno de los roles especificados
func (u *User) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if u.Role == role {
			return true
		}
	}
	return false
}

// IsSuperAdmin verifica si el usuario es super administrador
func (u *User) IsSuperAdmin() bool {
	return u.Role == "super_admin"
}

// IsAdmin verifica si el usuario es administrador o superior
func (u *User) IsAdmin() bool {
	return u.HasAnyRole("super_admin", "admin")
}

// IsEditor verifica si el usuario es editor o superior
func (u *User) IsEditor() bool {
	return u.HasAnyRole("super_admin", "admin", "editor")
}

// CanManageUsers verifica si el usuario puede gestionar usuarios
func (u *User) CanManageUsers() bool {
	return u.HasAnyRole("super_admin", "admin")
}

// CanManageContent verifica si el usuario puede gestionar contenido
func (u *User) CanManageContent() bool {
	return u.HasAnyRole("super_admin", "admin", "editor")
}

// CanManageSystem verifica si el usuario puede gestionar el sistema
func (u *User) CanManageSystem() bool {
	return u.Role == "super_admin"
}

// GetPermissions retorna los permisos del usuario
func (u *User) GetPermissions() []string {
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

	if perms, exists := permissions[u.Role]; exists {
		return perms
	}

	return []string{}
}

// HasPermission verifica si el usuario tiene un permiso específico
func (u *User) HasPermission(permission string) bool {
	permissions := u.GetPermissions()
	for _, perm := range permissions {
		if perm == permission {
			return true
		}
	}
	return false
}
