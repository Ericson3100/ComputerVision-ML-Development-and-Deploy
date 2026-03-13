@echo off
REM Script for development with hot-reload (Windows)

echo Starting development mode...
echo.

REM Verify Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo Error: Docker is not running.
    pause
    exit /b 1
)

echo Docker is running
echo.

REM Stop existing containers
echo Stopping existing containers...
docker-compose -f docker-compose.dev.yml down

echo.
echo Building Docker images (development mode)...
docker-compose -f docker-compose.dev.yml build

echo.
echo Starting services in development mode...
echo.
docker-compose -f docker-compose.dev.yml up

REM The above command blocks, so the message below will only appear when stopped
echo.
echo Services stopped

pause
