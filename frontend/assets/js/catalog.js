/**
 * Catálogo JavaScript - Funcionalidades específicas para el catálogo
 */

class CatalogManager {
    constructor() {
        this.products = [];
        this.filteredProducts = [];
        this.currentPage = 1;
        this.productsPerPage = 12;
        this.categories = [];
        this.init();
    }

    init() {
        this.loadProducts();
        this.setupEventListeners();
        this.setupSearch();
        this.setupFilters();
        this.setupPagination();
    }

    async loadProducts() {
        try {
            const response = await fetch('/api/products');
            if (response.ok) {
                this.products = await response.json();
                this.filteredProducts = [...this.products];
                this.extractCategories();
                this.renderProducts();
                this.updatePagination();
            }
        } catch (error) {
            console.error('Error cargando productos:', error);
            this.showNotification('Error al cargar los productos', 'error');
        }
    }

    extractCategories() {
        const categorySet = new Set();
        this.products.forEach(product => {
            if (product.Category && product.Category.Name) {
                categorySet.add(product.Category.Name);
            }
        });
        this.categories = Array.from(categorySet);
        this.populateCategoryFilter();
    }

    populateCategoryFilter() {
        const categoryFilter = document.getElementById('categoryFilter');
        if (!categoryFilter) return;

        // Limpiar opciones existentes excepto la primera
        categoryFilter.innerHTML = '<option value="">Todas las categorías</option>';
        
        this.categories.forEach(category => {
            const option = document.createElement('option');
            option.value = category;
            option.textContent = category;
            categoryFilter.appendChild(option);
        });
    }

    setupEventListeners() {
        // Event listeners para botones de productos
        document.addEventListener('click', (e) => {
            if (e.target.matches('[onclick*="showProductDetails"]')) {
                const productId = e.target.getAttribute('onclick').match(/'([^']+)'/)[1];
                this.showProductDetails(productId);
            }
            if (e.target.matches('[onclick*="addToInquiry"]')) {
                const productId = e.target.getAttribute('onclick').match(/'([^']+)'/)[1];
                this.addToInquiry(productId);
            }
        });

        // Event listeners para modales
        const modals = document.querySelectorAll('.modal');
        const closeButtons = document.querySelectorAll('.close');

        closeButtons.forEach(button => {
            button.addEventListener('click', () => {
                modals.forEach(modal => modal.style.display = 'none');
            });
        });

        window.addEventListener('click', (e) => {
            modals.forEach(modal => {
                if (e.target === modal) {
                    modal.style.display = 'none';
                }
            });
        });
    }

    setupSearch() {
        const searchInput = document.getElementById('searchInput');
        const searchBtn = document.getElementById('searchBtn');

        if (searchInput) {
            searchInput.addEventListener('input', (e) => {
                this.performSearch(e.target.value);
            });

            searchInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    this.performSearch(e.target.value);
                }
            });
        }

        if (searchBtn) {
            searchBtn.addEventListener('click', () => {
                const searchValue = searchInput ? searchInput.value : '';
                this.performSearch(searchValue);
            });
        }
    }

    setupFilters() {
        const categoryFilter = document.getElementById('categoryFilter');
        const sortFilter = document.getElementById('sortFilter');

        if (categoryFilter) {
            categoryFilter.addEventListener('change', () => {
                this.applyFilters();
            });
        }

        if (sortFilter) {
            sortFilter.addEventListener('change', () => {
                this.applyFilters();
            });
        }
    }

    setupPagination() {
        const prevBtn = document.getElementById('prevPage');
        const nextBtn = document.getElementById('nextPage');

        if (prevBtn) {
            prevBtn.addEventListener('click', () => {
                if (this.currentPage > 1) {
                    this.currentPage--;
                    this.renderProducts();
                    this.updatePagination();
                }
            });
        }

        if (nextBtn) {
            nextBtn.addEventListener('click', () => {
                const totalPages = Math.ceil(this.filteredProducts.length / this.productsPerPage);
                if (this.currentPage < totalPages) {
                    this.currentPage++;
                    this.renderProducts();
                    this.updatePagination();
                }
            });
        }
    }

    performSearch(searchTerm) {
        const term = searchTerm.toLowerCase().trim();
        
        this.filteredProducts = this.products.filter(product => {
            const nameMatch = product.Name.toLowerCase().includes(term);
            const descMatch = product.Description && product.Description.toLowerCase().includes(term);
            const categoryMatch = product.Category && product.Category.Name.toLowerCase().includes(term);
            
            return nameMatch || descMatch || categoryMatch;
        });

        this.currentPage = 1;
        this.renderProducts();
        this.updatePagination();
    }

    applyFilters() {
        const categoryFilter = document.getElementById('categoryFilter');
        const sortFilter = document.getElementById('sortFilter');
        
        let filtered = [...this.products];

        // Aplicar filtro de categoría
        if (categoryFilter && categoryFilter.value) {
            filtered = filtered.filter(product => 
                product.Category && product.Category.Name === categoryFilter.value
            );
        }

        // Aplicar ordenamiento
        if (sortFilter && sortFilter.value) {
            switch (sortFilter.value) {
                case 'name':
                    filtered.sort((a, b) => a.Name.localeCompare(b.Name));
                    break;
                case 'price_low':
                    filtered.sort((a, b) => (a.Price || 0) - (b.Price || 0));
                    break;
                case 'price_high':
                    filtered.sort((a, b) => (b.Price || 0) - (a.Price || 0));
                    break;
                case 'newest':
                    filtered.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt));
                    break;
            }
        }

        this.filteredProducts = filtered;
        this.currentPage = 1;
        this.renderProducts();
        this.updatePagination();
    }

    renderProducts() {
        const productsGrid = document.getElementById('productsGrid');
        if (!productsGrid) return;

        const startIndex = (this.currentPage - 1) * this.productsPerPage;
        const endIndex = startIndex + this.productsPerPage;
        const productsToShow = this.filteredProducts.slice(startIndex, endIndex);

        if (productsToShow.length === 0) {
            productsGrid.innerHTML = `
                <div class="no-products">
                    <h2>No se encontraron productos</h2>
                    <p>Intenta ajustar los filtros de búsqueda.</p>
                </div>
            `;
            return;
        }

        productsGrid.innerHTML = productsToShow.map(product => this.createProductCard(product)).join('');
    }

    createProductCard(product) {
        const categoryName = product.Category ? product.Category.Name : '';
        const features = product.Features ? product.Features.join('</li><li>') : '';
        
        return `
            <div class="product-card" data-category="${categoryName}" data-name="${product.Name}">
                ${product.ImageURL ? `
                <div class="product-image">
                    <img src="${product.ImageURL}" alt="${product.Name}" loading="lazy">
                    ${categoryName ? `<div class="product-category"><span>${categoryName}</span></div>` : ''}
                </div>
                ` : ''}
                <div class="product-info">
                    <h3>${product.Name}</h3>
                    ${product.Description ? `<p class="product-description">${product.Description}</p>` : ''}
                    ${product.Price ? `<div class="product-price"><span class="price">$${product.Price}</span></div>` : ''}
                    ${features ? `
                    <div class="product-features">
                        <ul><li>${features}</li></ul>
                    </div>
                    ` : ''}
                    <div class="product-actions">
                        <button class="btn btn-primary btn-sm" onclick="catalogManager.showProductDetails('${product.ID}')">
                            Ver detalles
                        </button>
                        <button class="btn btn-secondary btn-sm" onclick="catalogManager.addToInquiry('${product.ID}')">
                            Cotizar
                        </button>
                    </div>
                </div>
            </div>
        `;
    }

    updatePagination() {
        const totalPages = Math.ceil(this.filteredProducts.length / this.productsPerPage);
        const currentPageSpan = document.getElementById('currentPage');
        const totalPagesSpan = document.getElementById('totalPages');
        const prevBtn = document.getElementById('prevPage');
        const nextBtn = document.getElementById('nextPage');

        if (currentPageSpan) currentPageSpan.textContent = this.currentPage;
        if (totalPagesSpan) totalPagesSpan.textContent = totalPages;
        if (prevBtn) prevBtn.disabled = this.currentPage <= 1;
        if (nextBtn) nextBtn.disabled = this.currentPage >= totalPages;
    }

    async showProductDetails(productId) {
        const product = this.products.find(p => p.ID === productId);
        if (!product) return;

        const modal = document.getElementById('productModal');
        const content = document.getElementById('productModalContent');
        
        if (!modal || !content) return;

        const categoryName = product.Category ? product.Category.Name : '';
        const features = product.Features ? product.Features.map(f => `<li>${f}</li>`).join('') : '';
        const specs = product.Specifications ? product.Specifications.map(s => `<li><strong>${s.Key}:</strong> ${s.Value}</li>`).join('') : '';

        content.innerHTML = `
            <h2>${product.Name}</h2>
            ${product.ImageURL ? `<img src="${product.ImageURL}" alt="${product.Name}" style="max-width: 100%; height: auto; margin: 1rem 0;">` : ''}
            ${product.Description ? `<p>${product.Description}</p>` : ''}
            ${product.Price ? `<div class="product-price"><span class="price">$${product.Price}</span></div>` : ''}
            ${categoryName ? `<p><strong>Categoría:</strong> ${categoryName}</p>` : ''}
            ${features ? `
            <div class="product-features">
                <h4>Características:</h4>
                <ul>${features}</ul>
            </div>
            ` : ''}
            ${specs ? `
            <div class="product-specs">
                <h4>Especificaciones:</h4>
                <ul>${specs}</ul>
            </div>
            ` : ''}
            <div class="product-actions" style="margin-top: 2rem;">
                <button class="btn btn-primary" onclick="catalogManager.addToInquiry('${product.ID}')">
                    Solicitar cotización
                </button>
            </div>
        `;

        modal.style.display = 'block';
    }

    addToInquiry(productId) {
        const product = this.products.find(p => p.ID === productId);
        if (!product) return;

        const modal = document.getElementById('inquiryModal');
        if (!modal) return;

        // Agregar información del producto al formulario
        const messageField = document.getElementById('inquiryMessage');
        if (messageField) {
            const currentMessage = messageField.value;
            const productInfo = `\n\nProducto de interés: ${product.Name}`;
            messageField.value = currentMessage + productInfo;
        }

        modal.style.display = 'block';
    }

    showNotification(message, type = 'info') {
        // Usar la función de notificación del main.js si está disponible
        if (typeof showNotification === 'function') {
            showNotification(message, type);
        } else {
            alert(message);
        }
    }
}

// Inicializar el catálogo cuando el DOM esté listo
let catalogManager;
document.addEventListener('DOMContentLoaded', () => {
    catalogManager = new CatalogManager();
});

// Funciones globales para compatibilidad con onclick
window.showProductDetails = function(productId) {
    if (catalogManager) {
        catalogManager.showProductDetails(productId);
    }
};

window.addToInquiry = function(productId) {
    if (catalogManager) {
        catalogManager.addToInquiry(productId);
    }
}; 