#!/bin/bash

# Script for development with hot-reload

echo "Starting development mode..."
echo ""

# Verify Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running."
    exit 1
fi

echo "Docker is running"
echo ""

# Stop existing containers
echo "Stopping existing containers..."
docker-compose -f docker-compose.dev.yml down

echo ""
echo "Building Docker images (development mode)... "
docker-compose -f docker-compose.dev.yml build

echo ""
echo "Starting services in development mode..."
docker-compose -f docker-compose.dev.yml up

# The above command blocks, so the message below will only appear when stopped
echo ""
echo "Services stopped"
