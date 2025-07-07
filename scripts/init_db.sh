#!/bin/bash

# Script para inicializar la base de datos
# Uso: ./scripts/init_db.sh

set -e

# Cargar variables de entorno
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "Error: Archivo .env no encontrado"
    exit 1
fi

# Variables de conexión
DB_HOST=${DB_HOST:-localhost}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-website_db}
DB_PORT=${DB_PORT:-5432}

echo "🔧 Inicializando base de datos..."

# Verificar si PostgreSQL está corriendo
if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER > /dev/null 2>&1; then
    echo "❌ Error: PostgreSQL no está corriendo en $DB_HOST:$DB_PORT"
    echo "Por favor, asegúrate de que PostgreSQL esté iniciado"
    exit 1
fi

# Crear base de datos si no existe
echo "📦 Creando base de datos '$DB_NAME' si no existe..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || echo "Base de datos ya existe"

# Ejecutar migraciones
echo "🔄 Ejecutando migraciones..."
for migration in cms/migrations/*.sql; do
    if [ -f "$migration" ]; then
        echo "Ejecutando: $migration"
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
    fi
done

echo "✅ Base de datos inicializada correctamente!"
echo ""
echo "📋 Información de conexión:"
echo "   Host: $DB_HOST"
echo "   Puerto: $DB_PORT"
echo "   Base de datos: $DB_NAME"
echo "   Usuario: $DB_USER"
echo ""
echo "🔑 Usuario admin por defecto:"
echo "   Username: admin"
echo "   Password: admin123"
echo "   Email: admin@example.com"
echo ""
echo "🚀 Puedes ahora ejecutar:"
echo "   go run backend/main.go    # Para el backend"
echo "   go run cms/main.go        # Para el CMS" 