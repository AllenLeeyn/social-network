#!/bin/bash

echo "Starting Docker containers..."
docker-compose build --no-cache
docker-compose up
