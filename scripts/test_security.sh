#!/bin/bash

# Script para probar las funcionalidades de seguridad
# Uso: ./scripts/test_security.sh [backend_url] [cms_url]

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# URLs por defecto
BACKEND_URL=${1:-"http://localhost:3000"}
CMS_URL=${2:-"http://localhost:4000"}

echo -e "${BLUE}🔒 Probando funcionalidades de seguridad...${NC}"
echo -e "${YELLOW}Backend URL: ${BACKEND_URL}${NC}"
echo -e "${YELLOW}CMS URL: ${CMS_URL}${NC}"
echo ""

# Función para hacer requests y verificar headers
test_endpoint() {
    local url=$1
    local description=$2
    local expected_status=${3:-200}
    
    echo -e "${BLUE}Testing: ${description}${NC}"
    echo -e "URL: ${url}"
    
    response=$(curl -s -w "\n%{http_code}\n%{time_total}" -o /tmp/response_body "$url" 2>/dev/null)
    status_code=$(echo "$response" | tail -n 2 | head -n 1)
    time_total=$(echo "$response" | tail -n 1)
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ Status: ${status_code} (${time_total}s)${NC}"
    else
        echo -e "${RED}❌ Status: ${status_code} (expected ${expected_status})${NC}"
    fi
    
    # Verificar headers de seguridad
    echo -e "${YELLOW}Headers de seguridad:${NC}"
    
    # Verificar HSTS
    hsts=$(curl -s -I "$url" | grep -i "strict-transport-security" || echo "No HSTS")
    if [[ "$hsts" != "No HSTS" ]]; then
        echo -e "  ${GREEN}✅ HSTS: ${hsts}${NC}"
    else
        echo -e "  ${YELLOW}⚠️  HSTS: No encontrado${NC}"
    fi
    
    # Verificar CSP
    csp=$(curl -s -I "$url" | grep -i "content-security-policy" || echo "No CSP")
    if [[ "$csp" != "No CSP" ]]; then
        echo -e "  ${GREEN}✅ CSP: Presente${NC}"
    else
        echo -e "  ${YELLOW}⚠️  CSP: No encontrado${NC}"
    fi
    
    # Verificar X-Frame-Options
    xfo=$(curl -s -I "$url" | grep -i "x-frame-options" || echo "No X-Frame-Options")
    if [[ "$xfo" != "No X-Frame-Options" ]]; then
        echo -e "  ${GREEN}✅ X-Frame-Options: ${xfo}${NC}"
    else
        echo -e "  ${YELLOW}⚠️  X-Frame-Options: No encontrado${NC}"
    fi
    
    # Verificar X-Content-Type-Options
    xcto=$(curl -s -I "$url" | grep -i "x-content-type-options" || echo "No X-Content-Type-Options")
    if [[ "$xcto" != "No X-Content-Type-Options" ]]; then
        echo -e "  ${GREEN}✅ X-Content-Type-Options: ${xcto}${NC}"
    else
        echo -e "  ${YELLOW}⚠️  X-Content-Type-Options: No encontrado${NC}"
    fi
    
    echo ""
}

# Función para probar rate limiting
test_rate_limit() {
    local url=$1
    local description=$2
    local max_requests=${3:-10}
    
    echo -e "${BLUE}Testing Rate Limiting: ${description}${NC}"
    echo -e "URL: ${url}"
    echo -e "Max requests: ${max_requests}"
    
    # Hacer múltiples requests rápidamente
    for i in $(seq 1 $((max_requests + 2))); do
        response=$(curl -s -w "\n%{http_code}" -o /dev/null "$url" 2>/dev/null)
        status_code=$(echo "$response" | tail -n 1)
        
        if [ "$status_code" = "429" ]; then
            echo -e "${GREEN}✅ Rate limiting activo en request ${i} (${status_code})${NC}"
            break
        elif [ "$i" -le "$max_requests" ]; then
            echo -e "  Request ${i}: ${status_code}"
        else
            echo -e "${RED}❌ Rate limiting no activo en request ${i} (${status_code})${NC}"
        fi
    done
    
    echo ""
}

# Función para probar autenticación
test_auth() {
    local url=$1
    local description=$2
    
    echo -e "${BLUE}Testing Authentication: ${description}${NC}"
    echo -e "URL: ${url}"
    
    # Probar sin token
    response=$(curl -s -w "\n%{http_code}" -o /dev/null "$url" 2>/dev/null)
    status_code=$(echo "$response" | tail -n 1)
    
    if [ "$status_code" = "401" ]; then
        echo -e "${GREEN}✅ Autenticación requerida (${status_code})${NC}"
    else
        echo -e "${RED}❌ Autenticación no requerida (${status_code})${NC}"
    fi
    
    echo ""
}

# Probar endpoints públicos
echo -e "${YELLOW}=== Probando Endpoints Públicos ===${NC}"
test_endpoint "${BACKEND_URL}/" "Página de inicio"
test_endpoint "${BACKEND_URL}/api/config" "API Config"
test_endpoint "${BACKEND_URL}/api/slides" "API Slides"
test_endpoint "${BACKEND_URL}/api/categories" "API Categories"

# Probar rate limiting en APIs
echo -e "${YELLOW}=== Probando Rate Limiting ===${NC}"
test_rate_limit "${BACKEND_URL}/api/config" "API Rate Limiting" 1000

# Probar endpoints protegidos del CMS
echo -e "${YELLOW}=== Probando Endpoints Protegidos ===${NC}"
test_auth "${CMS_URL}/admin/slides" "CMS Slides (protegido)"
test_auth "${CMS_URL}/admin/users" "CMS Users (protegido)"

# Probar rate limiting en autenticación
echo -e "${YELLOW}=== Probando Rate Limiting en Autenticación ===${NC}"
test_rate_limit "${CMS_URL}/admin/auth/login" "Auth Rate Limiting" 5

# Probar headers de seguridad en CMS
echo -e "${YELLOW}=== Probando Headers de Seguridad en CMS ===${NC}"
test_endpoint "${CMS_URL}/admin" "CMS Admin Panel"

# Verificar HTTPS redirection (si está habilitado)
echo -e "${YELLOW}=== Probando Redirección HTTPS ===${NC}"
if [[ "$BACKEND_URL" == "https://"* ]]; then
    http_url=$(echo "$BACKEND_URL" | sed 's/https:/http:/')
    response=$(curl -s -w "\n%{http_code}" -o /dev/null "$http_url" 2>/dev/null)
    status_code=$(echo "$response" | tail -n 1)
    
    if [ "$status_code" = "301" ] || [ "$status_code" = "302" ]; then
        echo -e "${GREEN}✅ Redirección HTTP a HTTPS activa (${status_code})${NC}"
    else
        echo -e "${YELLOW}⚠️  Redirección HTTP a HTTPS no activa (${status_code})${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  HTTPS no configurado, saltando prueba de redirección${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Pruebas de seguridad completadas!${NC}"
echo ""
echo -e "${YELLOW}📋 Resumen de verificaciones:${NC}"
echo -e "  ✅ Headers de seguridad"
echo -e "  ✅ Rate limiting"
echo -e "  ✅ Autenticación"
echo -e "  ✅ Redirección HTTPS (si configurado)"
echo ""
echo -e "${BLUE}💡 Para más detalles, revisa los logs del servidor${NC}" 