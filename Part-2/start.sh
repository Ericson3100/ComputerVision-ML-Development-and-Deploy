#!/bin/bash

# Script to start the application in production mode

echo "Starting Image Classifier..."
echo ""

# Verify Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running."
    echo "Please start Docker Desktop and try again."
    exit 1
fi

echo "Docker is running"
echo ""

# Stop existing containers
echo "Stopping existing containers..."
docker-compose down

echo ""
echo "Building Docker images..."
docker-compose build --no-cache

echo ""
echo "Starting services..."
docker-compose up -d

echo ""
echo "Waiting for services to be ready..."
sleep 5

# Verify services are running
if docker-compose ps | grep -q "Up"; then
    echo ""
    echo "Application started successfully!"
    echo ""
    echo "Available services:"
    echo "   - Frontend: http://localhost:3000"
    echo "   - Backend:  http://localhost:8080"
    echo "   - Health:   http://localhost:8080/api/health"
    echo ""
    echo "To view logs: docker-compose logs -f"
    echo "To stop:      docker-compose down"
else
    echo ""
    echo "Error: Services did not start correctly"
    echo "Check logs with: docker-compose logs"
    exit 1
fi
