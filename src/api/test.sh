#!/bin/bash
# Test runner script for Authgrid API

echo "Running Authgrid API Tests..."
echo "=============================="
echo ""

# Run tests with verbose output
go test -v ./...

# Check exit code
if [ $? -eq 0 ]; then
    echo ""
    echo "=============================="
    echo "✅ All tests passed!"
    echo "=============================="
else
    echo ""
    echo "=============================="
    echo "❌ Some tests failed"
    echo "=============================="
    exit 1
fi
