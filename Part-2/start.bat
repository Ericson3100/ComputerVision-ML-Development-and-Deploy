@echo off
REM Script para iniciar la aplicación en modo producción (Windows)

echo 🚀 Iniciando Clasificador de Imágenes...
echo.

REM Verificar que Docker esté corriendo
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Docker no está corriendo.
    echo Por favor, inicia Docker Desktop e intenta nuevamente.
    pause
    exit /b 1
)

echo ✅ Docker está corriendo
echo.

REM Parar contenedores existentes
echo 🛑 Deteniendo contenedores existentes...
docker-compose down

echo.
echo 🏗️  Construyendo imágenes Docker...
docker-compose build --no-cache

echo.
echo 🚀 Iniciando servicios...
docker-compose up -d

echo.
echo ⏳ Esperando a que los servicios estén listos...
timeout /t 5 /nobreak >nul

echo.
echo ✅ ¡Aplicación iniciada exitosamente!
echo.
echo 📡 Servicios disponibles:
echo    - Frontend: http://localhost:3000
echo    - Backend:  http://localhost:8080
echo    - Health:   http://localhost:8080/api/health
echo.
echo 📋 Para ver los logs: docker-compose logs -f
echo 🛑 Para detener:     docker-compose down
echo.

pause
