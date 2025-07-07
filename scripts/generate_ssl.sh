#!/bin/bash

# Script para generar certificados SSL de desarrollo
# Uso: ./scripts/generate_ssl.sh

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔐 Generando certificados SSL de desarrollo...${NC}"

# Crear directorio para certificados si no existe
mkdir -p ssl

# Generar certificado autofirmado
echo -e "${YELLOW}Generando certificado autofirmado...${NC}"
openssl req -x509 -newkey rsa:4096 -keyout ssl/private.key -out ssl/certificate.crt -days 365 -nodes -subj "/C=ES/ST=Madrid/L=Madrid/O=Development/OU=IT/CN=localhost"

# Verificar que los archivos se crearon
if [ -f "ssl/certificate.crt" ] && [ -f "ssl/private.key" ]; then
    echo -e "${GREEN}✅ Certificados generados exitosamente!${NC}"
    echo -e "${YELLOW}📁 Ubicación:${NC}"
    echo -e "   Certificado: ssl/certificate.crt"
    echo -e "   Clave privada: ssl/private.key"
    echo ""
    echo -e "${YELLOW}🔧 Para usar HTTPS, actualiza tu .env:${NC}"
    echo -e "   ENABLE_HTTPS=true"
    echo -e "   SSL_CERT_FILE=ssl/certificate.crt"
    echo -e "   SSL_KEY_FILE=ssl/private.key"
    echo ""
    echo -e "${YELLOW}⚠️  Nota: Este es un certificado autofirmado para desarrollo.${NC}"
    echo -e "   En producción, usa certificados de una CA confiable."
else
    echo -e "${RED}❌ Error al generar certificados${NC}"
    exit 1
fi 