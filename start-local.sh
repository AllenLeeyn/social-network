#!/bin/bash

# Exit if any command fails
set -e

# Set env variable for frontend
export NEXT_PUBLIC_BACKEND_URL=http://localhost:8080

echo "Starting Go backend..."
cd backEnd
go run main.go &

BACKEND_PID=$!

echo "Starting Next.js frontend..."
cd ../frontend
npm install
npm install react-toastify
npm install react-icons
npm run dev &

FRONTEND_PID=$!

# Trap Ctrl+C (SIGINT) to clean up
trap "echo 'Stopping servers...'; kill $BACKEND_PID $FRONTEND_PID; exit" SIGINT

# Wait for background processes
wait
