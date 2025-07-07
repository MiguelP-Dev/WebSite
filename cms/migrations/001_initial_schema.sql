-- Migration: 001_initial_schema.sql
-- Description: Initial database schema creation

-- Create site_configs table
CREATE TABLE IF NOT EXISTS site_configs (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create slides table
CREATE TABLE IF NOT EXISTS slides (
    id SERIAL PRIMARY KEY,
    image_url VARCHAR(500) NOT NULL,
    title VARCHAR(200) NOT NULL,
    subtitle VARCHAR(500),
    link VARCHAR(500),
    "order" INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    "order" INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    image_urls TEXT[],
    features TEXT,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create contact_infos table
CREATE TABLE IF NOT EXISTS contact_infos (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL CHECK (type IN ('whatsapp', 'facebook', 'instagram', 'email', 'phone', 'website')),
    value VARCHAR(200) NOT NULL,
    icon VARCHAR(50),
    "order" INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'admin' CHECK (role IN ('admin', 'editor')),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_slides_active_order ON slides(active, "order");
CREATE INDEX IF NOT EXISTS idx_categories_active_order ON categories(active, "order");
CREATE INDEX IF NOT EXISTS idx_products_category_active ON products(category_id, active);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_contact_infos_active_order ON contact_infos(active, "order");
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Insert default admin user (password: admin123)
INSERT INTO users (username, email, password, role, active) 
VALUES ('admin', 'admin@example.com', '$2a$14$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', true)
ON CONFLICT (username) DO NOTHING;

-- Insert default site configurations
INSERT INTO site_configs (key, value) VALUES
('site_name', 'Mi Sitio Web'),
('site_description', 'Descripción del sitio web'),
('home_videos', 'https://www.youtube.com/watch?v=example1,https://www.youtube.com/watch?v=example2'),
('locations', 'Ubicación 1,Ubicación 2,Ubicación 3'),
('contact_email', 'contacto@example.com'),
('contact_phone', '+1234567890')
ON CONFLICT (key) DO NOTHING; 