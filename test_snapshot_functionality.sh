#!/bin/bash

# Test script to verify snapshot functionality
set -e

echo "=== Testing Snapshot Functionality ==="

# Build the application
echo "1. Building memscope..."
go build -o memscope cmd/memscope/main.go

# Test snapshot commands help
echo "2. Testing snapshot command structure..."
./memscope snapshot --help
echo

# Test snapshot subcommands help
echo "3. Testing snapshot subcommands..."
./memscope snapshot create --help
./memscope snapshot list --help
./memscope snapshot compare --help
./memscope snapshot delete --help
echo

# Test list with no snapshots
echo "4. Testing snapshot list (empty)..."
./memscope snapshot list
echo

# Run comprehensive tests
echo "5. Running all snapshot tests..."
go test ./internal/snapshot -v
echo

# Run CLI tests
echo "6. Running CLI tests..."
go test ./internal/cli -v
echo

# Run all project tests
echo "7. Running all project tests..."
go test ./...
echo

echo "=== All Tests Passed Successfully! ==="
echo "Snapshot functionality is fully implemented and working."