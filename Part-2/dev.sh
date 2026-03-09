#!/bin/bash

# Script para desarrollo con hot-reload

echo "🔧 Iniciando modo desarrollo..."
echo ""

# Verificar que Docker esté corriendo
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker no está corriendo."
    exit 1
fi

echo "✅ Docker está corriendo"
echo ""

# Parar contenedores existentes
echo "🛑 Deteniendo contenedores existentes..."
docker-compose -f docker-compose.dev.yml down

echo ""
echo "🏗️  Construyendo imágenes Docker (modo desarrollo)..."
docker-compose -f docker-compose.dev.yml build

echo ""
echo "🚀 Iniciando servicios en modo desarrollo..."
docker-compose -f docker-compose.dev.yml up

# El comando anterior bloquea, así que el mensaje de abajo solo aparecerá cuando se detenga
echo ""
echo "🛑 Servicios detenidos"
