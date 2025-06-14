#!/bin/bash

# Exit if any command fails
set -e

# Kill both processes on exit or error
cleanup() {
    echo "Stopping servers..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true
    exit
}

trap cleanup SIGINT SIGTERM ERR EXIT

# Set env variable for frontend
export NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
export NEXT_PUBLIC_BACKEND_WS_URL=ws://localhost:8080

echo "Starting Go backend..."
cd backEnd
go run main.go &

BACKEND_PID=$!
cd ../frontend

echo "Starting Next.js frontend..."
npm install
npm run build
npm start &

FRONTEND_PID=$!

# Wait for background processes
wait
