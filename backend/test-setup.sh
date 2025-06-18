#!/bin/bash

# Test script for Strudel-nvim setup

echo "🧪 Testing Strudel-nvim Setup"
echo "============================="

# Test 1: Check Go installation
echo "1. Checking Go installation..."
if command -v go &> /dev/null; then
    echo "   ✅ Go is installed: $(go version)"
else
    echo "   ❌ Go is not installed"
    exit 1
fi

# Test 2: Check if backend dependencies are installed
echo "2. Checking backend dependencies..."
if [ -f "go.sum" ]; then
    echo "   ✅ Go dependencies are installed"
else
    echo "   📦 Installing Go dependencies..."
    go mod tidy
fi

# Test 3: Test backend compilation
echo "3. Testing backend compilation..."
if go build -o test-backend .; then
    echo "   ✅ Backend compiles successfully"
    rm -f test-backend
else
    echo "   ❌ Backend compilation failed"
    exit 1
fi

# Test 4: Check Lua dependencies
echo "4. Checking Lua dependencies..."
cd ..
if luarocks --local list | grep -q "copas\|luasocket\|luabitop"; then
    echo "   ✅ Lua WebSocket dependencies are installed"
else
    echo "   ⚠️  Some Lua dependencies may be missing"
    echo "   Install with: luarocks --local install copas lua-websockets"
fi

# Test 5: Check port availability
echo "5. Checking port 8080 availability..."
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "   ⚠️  Port 8080 is in use. You may need to stop existing services."
else
    echo "   ✅ Port 8080 is available"
fi

echo ""
echo "🎉 Setup test complete!"
echo ""
echo "Next steps:"
echo "1. Start the backend: ./start-backend.sh"
echo "2. Open Neovim and run: :StrudelStart"
echo "3. Try evaluating: sound \"bd sn hh sn\""
echo ""

