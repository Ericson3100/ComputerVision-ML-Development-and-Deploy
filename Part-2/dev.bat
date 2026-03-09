@echo off
REM Script para desarrollo con hot-reload (Windows)

echo 🔧 Iniciando modo desarrollo...
echo.

REM Verificar que Docker esté corriendo
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Docker no está corriendo.
    pause
    exit /b 1
)

echo ✅ Docker está corriendo
echo.

REM Parar contenedores existentes
echo 🛑 Deteniendo contenedores existentes...
docker-compose -f docker-compose.dev.yml down

echo.
echo 🏗️  Construyendo imágenes Docker (modo desarrollo)...
docker-compose -f docker-compose.dev.yml build

echo.
echo 🚀 Iniciando servicios en modo desarrollo...
echo.
docker-compose -f docker-compose.dev.yml up

REM El comando anterior bloquea, así que el mensaje de abajo solo aparecerá cuando se detenga
echo.
echo 🛑 Servicios detenidos

pause
