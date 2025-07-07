# 🚀 Sitio Web Profesional con CMS Integrado

Un sitio web completo y moderno con sistema de gestión de contenido (CMS) construido con **Go**, **Fiber**, **GORM** y **PostgreSQL**. Diseñado para estudios fotográficos, empresas de servicios o cualquier negocio que necesite una presencia web profesional con gestión dinámica de contenido.

## 📋 Descripción del Proyecto

Este proyecto implementa una solución completa de sitio web corporativo con las siguientes características principales:

- **Frontend dinámico** con diseño moderno y responsivo
- **CMS integrado** para gestión completa del contenido
- **API RESTful** para comunicación entre componentes
- **Base de datos relacional** con PostgreSQL
- **Despliegue containerizado** con Docker
- **Optimización SEO** y rendimiento

## ✅ Lo que ya está implementado

### 🎨 Frontend Completo

#### **Páginas Implementadas:**

- ✅ **Página de Inicio** - Con slider dinámico, categorías destacadas y videos
- ✅ **Página de Productos** - Vista por categorías con productos
- ✅ **Página de Categoría** - Vista específica de cada categoría
- ✅ **Página de Contacto** - Formulario funcional e información de contacto
- ✅ **Página de Ubicaciones** - Múltiples ubicaciones con mapas
- ✅ **Catálogo Completo** - Vista con filtros, búsqueda y paginación

#### **Características del Frontend:**

- ✅ **Diseño Responsivo** - Adaptado para móviles, tablets y desktop
- ✅ **Slider Automático** - Carrusel de imágenes con controles y autoplay
- ✅ **Búsqueda en Tiempo Real** - Filtrado dinámico de productos
- ✅ **Navegación Intuitiva** - Menú hamburguesa para móviles
- ✅ **Efectos Visuales** - Hover effects, transiciones suaves
- ✅ **Optimización SEO** - Meta tags, estructura semántica
- ✅ **Accesibilidad** - Navegación por teclado, screen readers

#### **Tecnologías Frontend:**

- ✅ **HTML5** - Estructura semántica
- ✅ **CSS3** - Estilos modernos con Flexbox y Grid
- ✅ **JavaScript ES6+** - Funcionalidades interactivas
- ✅ **Font Awesome** - Iconografía profesional
- ✅ **Google Fonts** - Tipografía Inter

### 🔧 Backend Implementado

#### **API RESTful Completa:**

- ✅ **Endpoints Públicos** - Configuración, slides, categorías, productos, contactos
- ✅ **Endpoints CMS** - CRUD completo para todas las entidades
- ✅ **Validaciones** - Validación de datos con go-playground/validator
- ✅ **Manejo de Errores** - Sistema centralizado de errores
- ✅ **CORS** - Configurado para desarrollo y producción

#### **Modelos de Datos:**

- ✅ **SiteConfig** - Configuración dinámica del sitio
- ✅ **Slides** - Carrusel de imágenes con orden y estado
- ✅ **Categories** - Categorías de productos con slug único
- ✅ **Products** - Productos con imágenes, precios y características
- ✅ **ContactInfo** - Información de contacto tipificada
- ✅ **Users** - Sistema de usuarios con roles extendidos

#### **Características Backend:**

- ✅ **GORM** - ORM con migraciones automáticas
- ✅ **Fiber** - Framework web de alto rendimiento
- ✅ **Middleware** - Cache, logging, CORS
- ✅ **Templates** - Renderizado de vistas HTML
- ✅ **Pool de Conexiones** - Optimización de base de datos

### 🔐 Sistema de Seguridad Avanzado

#### **Sistema de Roles y Permisos:**

- ✅ **4 Roles Implementados** - super_admin, admin, editor, viewer
- ✅ **Sistema de Permisos Granular** - Control detallado por funcionalidad
- ✅ **Middleware de Autorización** - Verificación automática de permisos
- ✅ **Validación de Roles** - Verificación en tiempo real
- ✅ **Gestión de Usuarios** - CRUD completo con roles

#### **Rate Limiting Inteligente:**

- ✅ **5 Configuraciones Diferentes** - Strict, Moderate, Relaxed, API, Auth
- ✅ **Rate Limiting por IP** - Protección contra ataques DDoS
- ✅ **Rate Limiting por Usuario** - Control por usuario autenticado
- ✅ **Rate Limiting por Endpoint** - Control específico por ruta
- ✅ **Headers Estándar** - X-RateLimit-Limit, X-RateLimit-Remaining, etc.
- ✅ **Limpieza Automática** - Gestión automática de memoria

#### **HTTPS y Seguridad TLS:**

- ✅ **Soporte HTTPS Completo** - Configuración TLS 1.2+
- ✅ **Certificados SSL** - Soporte para certificados autofirmados y CA
- ✅ **Headers de Seguridad** - HSTS, CSP, XSS Protection, etc.
- ✅ **Redirección HTTP a HTTPS** - Automática y configurable
- ✅ **Configuración TLS Robusta** - Cipher suites seguros, curvas modernas

#### **Headers de Seguridad Implementados:**

- ✅ **Strict-Transport-Security** - Forzar HTTPS
- ✅ **Content-Security-Policy** - Protección XSS y inyección
- ✅ **X-XSS-Protection** - Protección adicional XSS
- ✅ **X-Frame-Options** - Prevenir clickjacking
- ✅ **X-Content-Type-Options** - Prevenir MIME sniffing
- ✅ **Referrer-Policy** - Control de referrer
- ✅ **X-Permitted-Cross-Domain-Policies** - Control cross-domain
- ✅ **Permissions-Policy** - Control de permisos del navegador

### 🎛️ CMS (Panel de Administración)

#### **Interfaz de Administración:**

- ✅ **Dashboard** - Vista general con estadísticas
- ✅ **Gestión de Slides** - CRUD completo del carrusel
- ✅ **Gestión de Categorías** - Crear, editar, eliminar categorías
- ✅ **Gestión de Productos** - Administración completa de productos
- ✅ **Gestión de Contactos** - Información de contacto dinámica
- ✅ **Gestión de Usuarios** - Sistema de usuarios del CMS
- ✅ **Configuración del Sitio** - Ajustes generales
- ✅ **Estadísticas del Sistema** - Métricas y reportes
- ✅ **Sistema de Backups** - Exportación de datos

#### **Características del CMS:**

- ✅ **Autenticación JWT** - Sistema seguro de login
- ✅ **Interfaz Moderna** - Bootstrap 5 con iconos
- ✅ **Validaciones en Tiempo Real** - Feedback inmediato
- ✅ **Gestión de Archivos** - Subida y gestión de imágenes
- ✅ **Ordenamiento** - Control de orden de elementos
- ✅ **Estados Activo/Inactivo** - Control de visibilidad
- ✅ **Rate Limiting Estricto** - Protección del panel admin
- ✅ **Logs de Seguridad** - Registro de actividades

### 🗄️ Base de Datos

#### **Estructura Implementada:**

- ✅ **PostgreSQL** - Base de datos relacional robusta
- ✅ **Migraciones Automáticas** - GORM auto-migrate
- ✅ **Relaciones** - Foreign keys y joins optimizados
- ✅ **Índices** - Optimización de consultas
- ✅ **Validaciones** - Constraints a nivel de base de datos

### 🐳 Infraestructura

#### **Docker y Contenedores:**

- ✅ **Dockerfile** - Configuración para backend
- ✅ **Dockerfile.cms** - Configuración para CMS
- ✅ **docker-compose.yml** - Orquestación de servicios
- ✅ **Nginx** - Proxy inverso configurado
- ✅ **PostgreSQL** - Base de datos containerizada

### 🔧 Herramientas de Desarrollo

#### **Scripts y Utilidades:**

- ✅ **Generador de Certificados SSL** - Script para desarrollo
- ✅ **Configuración de Variables** - Archivo .env.example completo
- ✅ **Documentación de Seguridad** - Guía completa de seguridad
- ✅ **Logs Mejorados** - Incluye IP y métricas de seguridad

## 🚧 Lo que resta por implementar

### 🔐 Autenticación y Seguridad

- ✅ **Sistema de Roles** - ✅ IMPLEMENTADO: Roles y permisos para usuarios del CMS
- ✅ **Rate Limiting** - ✅ IMPLEMENTADO: Protección contra ataques
- ✅ **HTTPS** - ✅ IMPLEMENTADO: Certificados SSL/TLS
- ⏳ **Autenticación OAuth** - Login con Google, Facebook, etc.
- ⏳ **2FA** - Autenticación de dos factores
- ⏳ **Audit Logs** - Registro detallado de actividades de seguridad

### 📁 Gestión de Archivos Avanzada

- ⏳ **Subida Múltiple** - Subida de varias imágenes simultáneamente
- ⏳ **Optimización Automática** - Compresión y redimensionado de imágenes
- ⏳ **CDN Integration** - Integración con CloudFront, Cloudflare
- ⏳ **Almacenamiento en la Nube** - S3, Google Cloud Storage
- ⏳ **Watermarks** - Marcas de agua automáticas

### 🛒 Funcionalidades E-commerce

- ⏳ **Carrito de Compras** - Sistema de carrito funcional
- ⏳ **Pasarela de Pagos** - Stripe, PayPal, MercadoPago
- ⏳ **Sistema de Pedidos** - Gestión completa de pedidos
- ⏳ **Inventario** - Control de stock
- ⏳ **Cupones y Descuentos** - Sistema de promociones

### 📊 Analítica y Reportes

- ⏳ **Google Analytics** - Integración con GA4
- ⏳ **Dashboard Analytics** - Métricas del sitio
- ⏳ **Reportes de Ventas** - Si se implementa e-commerce
- ⏳ **Tracking de Usuarios** - Comportamiento de visitantes
- ⏳ **Heatmaps** - Análisis de clics y scroll

### 📱 Funcionalidades Avanzadas

- ⏳ **PWA** - Progressive Web App
- ⏳ **Notificaciones Push** - Notificaciones del navegador
- ⏳ **Chat en Vivo** - Integración con WhatsApp, Telegram
- ⏳ **Blog/Noticias** - Sistema de publicaciones
- ⏳ **Newsletter** - Sistema de suscripciones
- ⏳ **Reservas/Citas** - Sistema de agendamiento

### 🔧 Mejoras Técnicas

- ⏳ **Cache Redis** - Cache de alto rendimiento
- ⏳ **Microservicios** - Arquitectura distribuida
- ⏳ **API GraphQL** - Alternativa a REST
- ⏳ **WebSockets** - Comunicación en tiempo real
- ⏳ **Testing** - Tests unitarios y de integración
- ⏳ **CI/CD** - Pipeline de despliegue automático

### 🌐 Internacionalización

- ⏳ **Multiidioma** - Soporte para múltiples idiomas
- ⏳ **Monedas** - Soporte para diferentes monedas
- ⏳ **Zonas Horarias** - Configuración por región
- ⏳ **Formatos Locales** - Fechas, números, etc.

## 🛠️ Tecnologías Utilizadas

### Backend

- **Go 1.24+** - Lenguaje de programación
- **Fiber** - Framework web de alto rendimiento
- **GORM** - ORM para Go
- **PostgreSQL** - Base de datos relacional
- **JWT** - Autenticación y autorización

### Frontend

- **HTML5** - Estructura semántica
- **CSS3** - Estilos modernos
- **JavaScript ES6+** - Funcionalidades interactivas
- **Bootstrap 5** - Framework CSS (CMS)
- **Font Awesome** - Iconografía

### DevOps

- **Docker** - Contenedores
- **Docker Compose** - Orquestación
- **Nginx** - Proxy inverso
- **Git** - Control de versiones

## 📦 Instalación y Despliegue

### Prerrequisitos

- Go 1.24+
- PostgreSQL 12+
- Docker & Docker Compose (recomendado)

### Instalación Local

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd website

# 2. Configurar variables de entorno
cp .env.example .env
# Editar .env con tus configuraciones

# 3. Instalar dependencias
go mod tidy

# 4. Inicializar base de datos
./scripts/init_db.sh

# 5. Ejecutar el proyecto
go run backend/main.go  # Puerto 3000
go run cms/main.go      # Puerto 4000
```

### Despliegue con Docker

```bash
# 1. Construir y levantar contenedores
docker-compose up --build -d

# 2. Acceder a los servicios
# Sitio web: http://localhost
# Panel admin: http://localhost:4000/admin
```

## 🔧 Configuración

### Variables de Entorno (.env)

```env
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=website_db
DB_PORT=5432

# CMS Database Configuration
CMS_DB_HOST=localhost
CMS_DB_USER=postgres
CMS_DB_PASSWORD=password
CMS_DB_NAME=website_db
CMS_DB_PORT=5432

# Server Configuration
BACKEND_PORT=3000
CMS_PORT=4000

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h
```

## 📚 API Documentation

### Endpoints Públicos

```
GET  /api/config          # Configuración del sitio
GET  /api/slides          # Slides activos
GET  /api/categories      # Categorías activas
GET  /api/products        # Productos activos
GET  /api/contacts        # Contactos activos
```

### Endpoints CMS (Protegidos)

```
POST /admin/login         # Autenticación
POST /admin/register      # Registro de usuarios

# Configuración
GET  /admin/config        # Obtener configuración
PUT  /admin/config        # Actualizar configuración

# Slides
GET    /admin/slides      # Listar slides
POST   /admin/slides      # Crear slide
PUT    /admin/slides/:id  # Actualizar slide
DELETE /admin/slides/:id  # Eliminar slide

# Categorías
GET    /admin/categories      # Listar categorías
POST   /admin/categories      # Crear categoría
PUT    /admin/categories/:id  # Actualizar categoría
DELETE /admin/categories/:id  # Eliminar categoría

# Productos
GET    /admin/products        # Listar productos
POST   /admin/products        # Crear producto
PUT    /admin/products/:id    # Actualizar producto
DELETE /admin/products/:id    # Eliminar producto

# Contactos
GET    /admin/contacts        # Listar contactos
POST   /admin/contacts        # Crear contacto
PUT    /admin/contacts/:id    # Actualizar contacto
DELETE /admin/contacts/:id    # Eliminar contacto

# Usuarios
GET    /admin/users           # Listar usuarios
POST   /admin/users           # Crear usuario
PUT    /admin/users/:id       # Actualizar usuario
DELETE /admin/users/:id       # Eliminar usuario
```

## 📁 Estructura del Proyecto

```
├── backend/
│   ├── controllers/     # Controladores de la API
│   ├── models/         # Modelos de datos
│   ├── middleware/     # Middleware (auth, cache, etc.)
│   ├── utils/          # Utilidades (validación, etc.)
│   └── main.go         # Servidor principal
├── cms/
│   ├── admin/          # Controladores del CMS
│   │   ├── js/         # JavaScript del panel admin
│   │   └── css/        # Estilos del panel admin
│   ├── migrations/     # Migraciones de BD
│   └── main.go         # Servidor CMS
├── frontend/
│   ├── templates/      # Templates HTML
│   │   ├── layout.html # Layout principal
│   │   └── pages/      # Páginas específicas
│   ├── assets/         # Assets estáticos
│   │   ├── css/        # Estilos CSS
│   │   ├── js/         # JavaScript
│   │   └── images/     # Imágenes
│   └── views/          # Vistas adicionales
├── scripts/            # Scripts de utilidad
├── docker/             # Configuración Docker
│   ├── nginx/          # Configuración Nginx
│   └── postgres/       # Configuración PostgreSQL
├── .env                # Variables de entorno
├── docker-compose.yml  # Orquestación Docker
├── Dockerfile          # Imagen del backend
├── Dockerfile.cms      # Imagen del CMS
└── README.MD           # Este archivo
```

## 🚀 Características Destacadas

### ✨ Diseño Moderno

- Interfaz limpia y profesional
- Diseño responsivo para todos los dispositivos
- Animaciones suaves y efectos visuales
- Tipografía moderna con Google Fonts

### ⚡ Alto Rendimiento

- Caché inteligente para optimizar velocidad
- Lazy loading de imágenes
- Compresión de assets
- Base de datos optimizada

### 🔒 Seguridad

- Autenticación JWT segura
- Validación de datos en múltiples capas
- Protección CSRF
- Sanitización de inputs

### 📱 Experiencia de Usuario

- Navegación intuitiva
- Búsqueda en tiempo real
- Filtros avanzados
- Paginación optimizada

## 🤝 Contribuir

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 📞 Soporte

Si tienes alguna pregunta o problema:

- 📧 Abrir un issue en el repositorio
- 📖 Revisar la documentación
- 💬 Contactar al equipo de desarrollo

---

**Desarrollado con ❤️ usando Go, Fiber y tecnologías modernas**
