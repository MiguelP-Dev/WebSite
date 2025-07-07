# 🔒 Guía de Seguridad

Esta guía describe las medidas de seguridad implementadas en el sitio web y cómo configurarlas correctamente.

## 📋 Medidas de Seguridad Implementadas

### 🔐 Sistema de Roles y Permisos

#### Roles Disponibles

1. **super_admin** - Acceso completo al sistema
   - Gestión de usuarios (crear, editar, eliminar)
   - Gestión de contenido completo
   - Configuración del sistema
   - Acceso a estadísticas y backups

2. **admin** - Administrador general
   - Gestión de usuarios (crear, editar)
   - Gestión de contenido completo
   - Configuración del sitio
   - No puede eliminar usuarios

3. **editor** - Editor de contenido
   - Crear y editar contenido
   - No puede eliminar contenido
   - No puede gestionar usuarios
   - Solo lectura de configuración

4. **viewer** - Solo lectura
   - Solo puede ver contenido
   - No puede hacer modificaciones

#### Permisos por Rol

```go
// Permisos de super_admin
"users:read", "users:write", "users:delete",
"slides:read", "slides:write", "slides:delete",
"categories:read", "categories:write", "categories:delete",
"products:read", "products:write", "products:delete",
"contacts:read", "contacts:write", "contacts:delete",
"config:read", "config:write", "system:admin"

// Permisos de admin
"users:read", "users:write",
"slides:read", "slides:write", "slides:delete",
"categories:read", "categories:write", "categories:delete",
"products:read", "products:write", "products:delete",
"contacts:read", "contacts:write", "contacts:delete",
"config:read", "config:write"

// Permisos de editor
"slides:read", "slides:write",
"categories:read", "categories:write",
"products:read", "products:write",
"contacts:read", "contacts:write",
"config:read"

// Permisos de viewer
"slides:read",
"categories:read",
"products:read",
"contacts:read",
"config:read"
```

### 🚦 Rate Limiting

#### Configuraciones Disponibles

1. **RateLimitStrict** - Para endpoints críticos
   - 10 requests por minuto por IP
   - Aplicado al CMS completo

2. **RateLimitModerate** - Para uso general
   - 60 requests por minuto por IP
   - Aplicado al backend principal

3. **RateLimitRelaxed** - Para contenido público
   - 300 requests por minuto por IP
   - Para páginas web

4. **RateLimitAPI** - Para APIs
   - 1000 requests por hora por IP
   - Para endpoints de API

5. **RateLimitAuth** - Para autenticación
   - 5 requests por 15 minutos por IP
   - Para login y registro

#### Headers de Rate Limiting

El sistema incluye headers estándar de rate limiting:

```plaintext
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1640995200
Retry-After: 30
```

### 🔒 HTTPS y Seguridad TLS

#### Configuración TLS

- **TLS 1.2+** requerido
- **Cipher suites** seguros habilitados
- **Curvas elípticas** modernas (X25519, P-256, P-384)
- **HSTS** habilitado por defecto

#### Headers de Seguridad

1. **Strict-Transport-Security**

   ```plaintext
   Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
   ```

2. **Content-Security-Policy**

   ```plaintext
   Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https:; connect-src 'self' https:; frame-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'
   ```

3. **X-XSS-Protection**

   ```plaintext
   X-XSS-Protection: 1; mode=block
   ```

4. **X-Frame-Options**

   ```plaintext
   X-Frame-Options: DENY
   ```

5. **X-Content-Type-Options**

   ```plaintext
   X-Content-Type-Options: nosniff
   ```

6. **Referrer-Policy**

   ```plaintext
   Referrer-Policy: strict-origin-when-cross-origin
   ```

7. **X-Permitted-Cross-Domain-Policies**

   ```plaintext
   X-Permitted-Cross-Domain-Policies: none
   ```

8. **Permissions-Policy**

   ```plaintext
   Permissions-Policy: geolocation=(), microphone=(), camera=()
   ```

## ⚙️ Configuración

### Variables de Entorno de Seguridad

```env
# Habilitar HTTPS
ENABLE_HTTPS=true

# Headers de seguridad
ENABLE_HSTS=true
ENABLE_CSP=true
ENABLE_XSS=true
ENABLE_FRAME_DENY=true
ENABLE_NO_SNIFF=true
ENABLE_REFERRER=true

# Certificados SSL
SSL_CERT_FILE=/path/to/certificate.crt
SSL_KEY_FILE=/path/to/private.key

# Rate limiting
RATE_LIMIT_STRICT_MAX=10
RATE_LIMIT_STRICT_WINDOW=1m
RATE_LIMIT_MODERATE_MAX=60
RATE_LIMIT_MODERATE_WINDOW=1m
RATE_LIMIT_RELAXED_MAX=300
RATE_LIMIT_RELAXED_WINDOW=1m
RATE_LIMIT_API_MAX=1000
RATE_LIMIT_API_WINDOW=1h
RATE_LIMIT_AUTH_MAX=5
RATE_LIMIT_AUTH_WINDOW=15m
```

### Generar Certificados SSL de Desarrollo

```bash
# Generar certificados autofirmados
./scripts/generate_ssl.sh

# Configurar variables de entorno
ENABLE_HTTPS=true
SSL_CERT_FILE=ssl/certificate.crt
SSL_KEY_FILE=ssl/private.key
```

## 🛡️ Mejores Prácticas

### 1. Gestión de Usuarios

- **Cambiar contraseñas** regularmente
- **Usar contraseñas fuertes** (mínimo 8 caracteres, mayúsculas, minúsculas, números, símbolos)
- **Habilitar 2FA** cuando esté disponible
- **Revisar logs** de acceso regularmente

### 2. Configuración de Producción

- **Usar certificados SSL** de una CA confiable (Let's Encrypt, DigiCert, etc.)
- **Configurar firewall** para limitar acceso
- **Habilitar logs** de seguridad
- **Monitorear** intentos de acceso fallidos

### 3. Desarrollo

- **No usar credenciales** de producción en desarrollo
- **Usar certificados** autofirmados solo para desarrollo
- **Revisar logs** de errores regularmente
- **Probar** todas las funcionalidades de seguridad

### 4. Mantenimiento

- **Actualizar dependencias** regularmente
- **Revisar logs** de seguridad
- **Hacer backups** regulares
- **Monitorear** métricas de rate limiting

## 🚨 Respuesta a Incidentes

### 1. Detección

- **Monitorear logs** de acceso
- **Revisar métricas** de rate limiting
- **Verificar** intentos de acceso fallidos
- **Revisar** cambios en archivos críticos

### 2. Contención

- **Bloquear IPs** maliciosas
- **Cambiar contraseñas** si es necesario
- **Deshabilitar** cuentas comprometidas
- **Revisar** permisos de usuarios

### 3. Recuperación

- **Restaurar** desde backups si es necesario
- **Actualizar** credenciales
- **Revisar** logs de seguridad
- **Documentar** el incidente

### 4. Prevención

- **Implementar** medidas adicionales
- **Actualizar** políticas de seguridad
- **Capacitar** al equipo
- **Revisar** configuraciones

## 📞 Contacto de Seguridad

Si encuentras una vulnerabilidad de seguridad:

1. **No la divulgues** públicamente
2. **Reporta** el problema al equipo de desarrollo
3. **Proporciona** detalles específicos
4. **Espera** respuesta antes de tomar acciones adicionales

---

**Nota**: Esta guía debe actualizarse regularmente con las mejores prácticas de seguridad y las nuevas amenazas identificadas.
