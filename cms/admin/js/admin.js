// CMS Admin JavaScript
class CMSAdmin {
    constructor() {
        this.token = localStorage.getItem('admin_token');
        this.currentUser = JSON.parse(localStorage.getItem('admin_user') || '{}');
        this.init();
    }
    
    init() {
        this.checkAuth();
        this.setupEventListeners();
        this.loadDashboard();
    }
    
    checkAuth() {
        if (!this.token) {
            this.showLoginForm();
        } else {
            this.showAdminPanel();
        }
    }
    
    showLoginForm() {
        document.body.innerHTML = `
            <div class="login-container">
                <div class="login-card">
                    <h2>Iniciar Sesión</h2>
                    <form id="loginForm">
                        <div class="form-group">
                            <label for="username">Usuario o Email</label>
                            <input type="text" id="username" name="username" required>
                        </div>
                        <div class="form-group">
                            <label for="password">Contraseña</label>
                            <input type="password" id="password" name="password" required>
                        </div>
                        <button type="submit" class="btn btn-primary">Iniciar Sesión</button>
                    </form>
                </div>
            </div>
        `;
        
        document.getElementById('loginForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.login();
        });
    }
    
    async login() {
        const formData = new FormData(document.getElementById('loginForm'));
        const data = {
            username: formData.get('username'),
            password: formData.get('password')
        };
        
        try {
            const response = await fetch('/admin/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
            
            const result = await response.json();
            
            if (response.ok) {
                this.token = result.token;
                this.currentUser = result.user;
                localStorage.setItem('admin_token', this.token);
                localStorage.setItem('admin_user', JSON.stringify(this.currentUser));
                this.showAdminPanel();
            } else {
                this.showNotification(result.error, 'error');
            }
        } catch (error) {
            this.showNotification('Error de conexión', 'error');
        }
    }
    
    showAdminPanel() {
        document.body.innerHTML = `
            <div class="container-fluid">
                <div class="row">
                    <!-- Sidebar -->
                    <nav class="col-md-3 col-lg-2 d-md-block bg-dark sidebar">
                        <div class="position-sticky pt-3">
                            <h4 class="text-white px-3">Admin Panel</h4>
                            <ul class="nav flex-column">
                                <li class="nav-item">
                                    <a class="nav-link active" href="#dashboard" data-bs-toggle="tab">
                                        <i class="bi bi-speedometer2"></i> Dashboard
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link" href="#slides" data-bs-toggle="tab">
                                        <i class="bi bi-images"></i> Slides
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link" href="#categories" data-bs-toggle="tab">
                                        <i class="bi bi-collection"></i> Categorías
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link" href="#products" data-bs-toggle="tab">
                                        <i class="bi bi-box-seam"></i> Productos
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link" href="#contacts" data-bs-toggle="tab">
                                        <i class="bi bi-telephone"></i> Contacto
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link" href="#users" data-bs-toggle="tab">
                                        <i class="bi bi-people"></i> Usuarios
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link" href="#config" data-bs-toggle="tab">
                                        <i class="bi bi-gear"></i> Configuración
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link" href="#" id="logoutBtn">
                                        <i class="bi bi-box-arrow-right"></i> Cerrar Sesión
                                    </a>
                                </li>
                            </ul>
                        </div>
                    </nav>

                    <!-- Main Content -->
                    <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                        <div class="tab-content">
                            <div class="tab-pane fade show active" id="dashboard">
                                <div class="d-flex justify-content-between my-3">
                                    <h2>Dashboard</h2>
                                    <span class="text-muted">Bienvenido, ${this.currentUser.username}</span>
                                </div>
                                <div id="dashboardContent"></div>
                            </div>
                            
                            <div class="tab-pane fade" id="slides">
                                <div class="d-flex justify-content-between my-3">
                                    <h2>Slides del Carrusel</h2>
                                    <button class="btn btn-primary" id="addSlide">
                                        <i class="bi bi-plus"></i> Nuevo Slide
                                    </button>
                                </div>
                                <div id="slidesList"></div>
                            </div>
                            
                            <div class="tab-pane fade" id="categories">
                                <div class="d-flex justify-content-between my-3">
                                    <h2>Categorías</h2>
                                    <button class="btn btn-primary" id="addCategory">
                                        <i class="bi bi-plus"></i> Nueva Categoría
                                    </button>
                                </div>
                                <div id="categoriesList"></div>
                            </div>
                            
                            <div class="tab-pane fade" id="products">
                                <div class="d-flex justify-content-between my-3">
                                    <h2>Productos</h2>
                                    <button class="btn btn-primary" id="addProduct">
                                        <i class="bi bi-plus"></i> Nuevo Producto
                                    </button>
                                </div>
                                <div id="productsList"></div>
                            </div>
                            
                            <div class="tab-pane fade" id="contacts">
                                <div class="d-flex justify-content-between my-3">
                                    <h2>Información de Contacto</h2>
                                    <button class="btn btn-primary" id="addContact">
                                        <i class="bi bi-plus"></i> Nuevo Contacto
                                    </button>
                                </div>
                                <div id="contactsList"></div>
                            </div>
                            
                            <div class="tab-pane fade" id="users">
                                <div class="d-flex justify-content-between my-3">
                                    <h2>Usuarios</h2>
                                    <button class="btn btn-primary" id="addUser">
                                        <i class="bi bi-plus"></i> Nuevo Usuario
                                    </button>
                                </div>
                                <div id="usersList"></div>
                            </div>
                            
                            <div class="tab-pane fade" id="config">
                                <div class="d-flex justify-content-between my-3">
                                    <h2>Configuración del Sitio</h2>
                                </div>
                                <div id="configContent"></div>
                            </div>
                        </div>
                    </main>
                </div>
            </div>
        `;
        
        this.setupEventListeners();
        this.loadDashboard();
    }
    
    setupEventListeners() {
        // Tab navigation
        document.querySelectorAll('.nav-link[data-bs-toggle="tab"]').forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const target = e.target.getAttribute('href').substring(1);
                this.loadTabContent(target);
            });
        });
        
        // Logout
        document.getElementById('logoutBtn')?.addEventListener('click', (e) => {
            e.preventDefault();
            this.logout();
        });
        
        // CRUD buttons
        document.getElementById('addSlide')?.addEventListener('click', () => this.showSlideForm());
        document.getElementById('addCategory')?.addEventListener('click', () => this.showCategoryForm());
        document.getElementById('addProduct')?.addEventListener('click', () => this.showProductForm());
        document.getElementById('addContact')?.addEventListener('click', () => this.showContactForm());
        document.getElementById('addUser')?.addEventListener('click', () => this.showUserForm());
    }
    
    async loadTabContent(tab) {
        switch(tab) {
            case 'dashboard':
                this.loadDashboard();
                break;
            case 'slides':
                this.loadSlides();
                break;
            case 'categories':
                this.loadCategories();
                break;
            case 'products':
                this.loadProducts();
                break;
            case 'contacts':
                this.loadContacts();
                break;
            case 'users':
                this.loadUsers();
                break;
            case 'config':
                this.loadConfig();
                break;
        }
    }
    
    async loadDashboard() {
        try {
            const [slides, categories, products, contacts] = await Promise.all([
                this.apiCall('/admin/slides'),
                this.apiCall('/admin/categories'),
                this.apiCall('/admin/products'),
                this.apiCall('/admin/contacts')
            ]);
            
            const dashboardContent = document.getElementById('dashboardContent');
            dashboardContent.innerHTML = `
                <div class="row">
                    <div class="col-md-3">
                        <div class="card text-white bg-primary">
                            <div class="card-body">
                                <h5 class="card-title">Slides</h5>
                                <p class="card-text display-6">${slides.length}</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-3">
                        <div class="card text-white bg-success">
                            <div class="card-body">
                                <h5 class="card-title">Categorías</h5>
                                <p class="card-text display-6">${categories.length}</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-3">
                        <div class="card text-white bg-warning">
                            <div class="card-body">
                                <h5 class="card-title">Productos</h5>
                                <p class="card-text display-6">${products.length}</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-3">
                        <div class="card text-white bg-info">
                            <div class="card-body">
                                <h5 class="card-title">Contactos</h5>
                                <p class="card-text display-6">${contacts.length}</p>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        } catch (error) {
            this.showNotification('Error al cargar dashboard', 'error');
        }
    }
    
    async loadSlides() {
        try {
            const slides = await this.apiCall('/admin/slides');
            this.renderSlidesTable(slides);
        } catch (error) {
            this.showNotification('Error al cargar slides', 'error');
        }
    }
    
    renderSlidesTable(slides) {
        const slidesList = document.getElementById('slidesList');
        
        if (slides.length === 0) {
            slidesList.innerHTML = '<div class="alert alert-info">No hay slides configurados</div>';
            return;
        }
        
        const table = `
            <table class="table table-striped">
                <thead>
                    <tr>
                        <th>Imagen</th>
                        <th>Título</th>
                        <th>Subtítulo</th>
                        <th>Orden</th>
                        <th>Estado</th>
                        <th>Acciones</th>
                    </tr>
                </thead>
                <tbody>
                    ${slides.map(slide => `
                        <tr>
                            <td><img src="${slide.image_url}" style="max-width: 100px; height: 60px; object-fit: cover;"></td>
                            <td>${slide.title || 'Sin título'}</td>
                            <td>${slide.subtitle || '-'}</td>
                            <td>${slide.order}</td>
                            <td>${slide.active ? '<span class="badge bg-success">Activo</span>' : '<span class="badge bg-secondary">Inactivo</span>'}</td>
                            <td>
                                <button class="btn btn-sm btn-outline-primary edit-slide" data-id="${slide.id}">
                                    <i class="bi bi-pencil"></i>
                                </button>
                                <button class="btn btn-sm btn-outline-danger delete-slide" data-id="${slide.id}">
                                    <i class="bi bi-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
        
        slidesList.innerHTML = table;
        
        // Add event listeners
        slidesList.querySelectorAll('.edit-slide').forEach(btn => {
            btn.addEventListener('click', () => this.editSlide(btn.getAttribute('data-id')));
        });
        
        slidesList.querySelectorAll('.delete-slide').forEach(btn => {
            btn.addEventListener('click', () => this.deleteSlide(btn.getAttribute('data-id')));
        });
    }
    
    showSlideForm(slide = null) {
        const modal = `
            <div class="modal fade" id="slideModal" tabindex="-1">
                <div class="modal-dialog">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title">${slide ? 'Editar' : 'Nuevo'} Slide</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                        </div>
                        <div class="modal-body">
                            <form id="slideForm">
                                <div class="mb-3">
                                    <label for="image_url" class="form-label">URL de la Imagen</label>
                                    <input type="url" class="form-control" id="image_url" name="image_url" value="${slide?.image_url || ''}" required>
                                </div>
                                <div class="mb-3">
                                    <label for="title" class="form-label">Título</label>
                                    <input type="text" class="form-control" id="title" name="title" value="${slide?.title || ''}" required>
                                </div>
                                <div class="mb-3">
                                    <label for="subtitle" class="form-label">Subtítulo</label>
                                    <input type="text" class="form-control" id="subtitle" name="subtitle" value="${slide?.subtitle || ''}">
                                </div>
                                <div class="mb-3">
                                    <label for="link" class="form-label">Enlace</label>
                                    <input type="url" class="form-control" id="link" name="link" value="${slide?.link || ''}">
                                </div>
                                <div class="mb-3">
                                    <label for="order" class="form-label">Orden</label>
                                    <input type="number" class="form-control" id="order" name="order" value="${slide?.order || 0}">
                                </div>
                                <div class="mb-3">
                                    <div class="form-check">
                                        <input class="form-check-input" type="checkbox" id="active" name="active" ${slide?.active ? 'checked' : ''}>
                                        <label class="form-check-label" for="active">Activo</label>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
                            <button type="button" class="btn btn-primary" id="saveSlide">Guardar</button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        document.body.insertAdjacentHTML('beforeend', modal);
        
        const modalElement = document.getElementById('slideModal');
        const modalInstance = new bootstrap.Modal(modalElement);
        modalInstance.show();
        
        document.getElementById('saveSlide').addEventListener('click', () => {
            this.saveSlide(slide?.id);
        });
        
        modalElement.addEventListener('hidden.bs.modal', () => {
            modalElement.remove();
        });
    }
    
    async saveSlide(id = null) {
        const form = document.getElementById('slideForm');
        const formData = new FormData(form);
        
        const data = {
            image_url: formData.get('image_url'),
            title: formData.get('title'),
            subtitle: formData.get('subtitle'),
            link: formData.get('link'),
            order: parseInt(formData.get('order')),
            active: formData.get('active') === 'on'
        };
        
        try {
            const url = id ? `/admin/slides/${id}` : '/admin/slides';
            const method = id ? 'PUT' : 'POST';
            
            await this.apiCall(url, method, data);
            
            this.showNotification(`Slide ${id ? 'actualizado' : 'creado'} exitosamente`, 'success');
            bootstrap.Modal.getInstance(document.getElementById('slideModal')).hide();
            this.loadSlides();
        } catch (error) {
            this.showNotification('Error al guardar slide', 'error');
        }
    }
    
    async deleteSlide(id) {
        if (confirm('¿Estás seguro de que quieres eliminar este slide?')) {
            try {
                await this.apiCall(`/admin/slides/${id}`, 'DELETE');
                this.showNotification('Slide eliminado exitosamente', 'success');
                this.loadSlides();
            } catch (error) {
                this.showNotification('Error al eliminar slide', 'error');
            }
        }
    }
    
    async editSlide(id) {
        try {
            const slide = await this.apiCall(`/admin/slides/${id}`);
            this.showSlideForm(slide);
        } catch (error) {
            this.showNotification('Error al cargar slide', 'error');
        }
    }
    
    // Similar methods for other entities...
    async loadCategories() {
        try {
            const categories = await this.apiCall('/admin/categories');
            this.renderCategoriesTable(categories);
        } catch (error) {
            this.showNotification('Error al cargar categorías', 'error');
        }
    }
    
    async loadProducts() {
        try {
            const products = await this.apiCall('/admin/products');
            this.renderProductsTable(products);
        } catch (error) {
            this.showNotification('Error al cargar productos', 'error');
        }
    }
    
    async loadContacts() {
        try {
            const contacts = await this.apiCall('/admin/contacts');
            this.renderContactsTable(contacts);
        } catch (error) {
            this.showNotification('Error al cargar contactos', 'error');
        }
    }
    
    async loadUsers() {
        try {
            const users = await this.apiCall('/admin/users');
            this.renderUsersTable(users);
        } catch (error) {
            this.showNotification('Error al cargar usuarios', 'error');
        }
    }
    
    async loadConfig() {
        try {
            const config = await this.apiCall('/admin/config');
            this.renderConfigForm(config);
        } catch (error) {
            this.showNotification('Error al cargar configuración', 'error');
        }
    }
    
    // API helper
    async apiCall(endpoint, method = 'GET', data = null) {
        const options = {
            method,
            headers: {
                'Authorization': `Bearer ${this.token}`,
                'Content-Type': 'application/json'
            }
        };
        
        if (data && method !== 'GET') {
            options.body = JSON.stringify(data);
        }
        
        const response = await fetch(endpoint, options);
        
        if (response.status === 401) {
            this.logout();
            return;
        }
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Error de API');
        }
        
        return await response.json();
    }
    
    logout() {
        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_user');
        this.token = null;
        this.currentUser = {};
        this.showLoginForm();
    }
    
    showNotification(message, type = 'info') {
        const alertClass = {
            success: 'alert-success',
            error: 'alert-danger',
            warning: 'alert-warning',
            info: 'alert-info'
        }[type];
        
        const notification = `
            <div class="alert ${alertClass} alert-dismissible fade show" role="alert">
                ${message}
                <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
            </div>
        `;
        
        const container = document.querySelector('.col-md-9');
        container.insertAdjacentHTML('afterbegin', notification);
        
        // Auto remove after 5 seconds
        setTimeout(() => {
            const alert = container.querySelector('.alert');
            if (alert) {
                alert.remove();
            }
        }, 5000);
    }
}

// Initialize CMS Admin
document.addEventListener('DOMContentLoaded', function() {
    new CMSAdmin();
});