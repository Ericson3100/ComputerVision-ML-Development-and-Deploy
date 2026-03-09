#!/bin/bash

# Script para iniciar la aplicación en modo producción

echo "🚀 Iniciando Clasificador de Imágenes..."
echo ""

# Verificar que Docker esté corriendo
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker no está corriendo."
    echo "Por favor, inicia Docker Desktop e intenta nuevamente."
    exit 1
fi

echo "✅ Docker está corriendo"
echo ""

# Parar contenedores existentes
echo "🛑 Deteniendo contenedores existentes..."
docker-compose down

echo ""
echo "🏗️  Construyendo imágenes Docker..."
docker-compose build --no-cache

echo ""
echo "🚀 Iniciando servicios..."
docker-compose up -d

echo ""
echo "⏳ Esperando a que los servicios estén listos..."
sleep 5

# Verificar que los servicios estén corriendo
if docker-compose ps | grep -q "Up"; then
    echo ""
    echo "✅ ¡Aplicación iniciada exitosamente!"
    echo ""
    echo "📡 Servicios disponibles:"
    echo "   - Frontend: http://localhost:3000"
    echo "   - Backend:  http://localhost:8080"
    echo "   - Health:   http://localhost:8080/api/health"
    echo ""
    echo "📋 Para ver los logs: docker-compose logs -f"
    echo "🛑 Para detener:     docker-compose down"
else
    echo ""
    echo "❌ Error: Los servicios no se iniciaron correctamente"
    echo "Revisa los logs con: docker-compose logs"
    exit 1
fi
