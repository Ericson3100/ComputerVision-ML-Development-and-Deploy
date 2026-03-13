@echo off
REM Script to start the application in production mode (Windows)

echo Starting Image Classifier...
echo.

REM Verify Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo Error: Docker is not running.
    echo Please start Docker Desktop and try again.
    pause
    exit /b 1
)

echo Docker is running
echo.

REM Stop existing containers
echo Stopping existing containers...
docker-compose down

echo.
echo Building Docker images...
docker-compose build --no-cache

echo.
echo Starting services...
docker-compose up -d

echo.
echo Waiting for services to be ready...
timeout /t 5 /nobreak >nul

echo.
echo Application started successfully!
echo.
echo Available services:
echo    - Frontend: http://localhost:3000
echo    - Backend:  http://localhost:8080
echo    - Health:   http://localhost:8080/api/health
echo.
echo To view logs: docker-compose logs -f
echo To stop:      docker-compose down
echo.

pause
