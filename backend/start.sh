#!/bin/bash

# Strudel Backend Startup Script

echo "🎵 Starting Strudel Backend Server..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Install dependencies if needed
if [ ! -f "go.sum" ]; then
    echo "📦 Installing dependencies..."
    go mod tidy
fi

# Check if port 8080 is already in use
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  Port 8080 is already in use. Please stop the existing service."
    echo "   You can find the process with: lsof -i :8080"
    exit 1
fi

# Start the server
echo "🚀 Launching server on http://localhost:8080"
echo "📡 WebSocket endpoint: ws://localhost:8080/ws"
echo "💡 Press Ctrl+C to stop the server"
echo ""

go run .

