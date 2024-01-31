#!/bin/bash

# Read port from .env file
HOST_PORT=$(grep -E '^PORT=' .env | cut -d '=' -f2)

# Build Docker image
docker build -t test-app-container --build-arg PORT=$HOST_PORT .

# Run Docker container
docker run -e PORT=$HOST_PORT -p 4000:$HOST_PORT test-app-container
