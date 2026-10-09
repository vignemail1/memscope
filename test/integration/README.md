# Memscope Integration Tests

This directory contains comprehensive integration tests for the memscope CLI tool.

## Test Categories

### End-to-End Workflow Testing
- Complete user workflow validation
- Cross-command integration testing
- Data flow validation between components
- Error propagation and recovery testing

### Performance and Load Testing  
- Memory usage profiling
- Response time benchmarks
- Large dataset handling validation
- Resource cleanup verification

### Cross-Platform Validation
- Windows-specific functionality verification
- Non-Windows graceful fallback validation
- Cross-platform CLI behavior consistency
- File path and permissions handling

### Regression Testing
- Automated test suite for all major functionality
- Backward compatibility validation
- Data format consistency checking
- Command interface stability verification

## Running Tests

```bash
# Run all integration tests
go test ./test/integration -v

# Run only fast tests
go test ./test/integration -v -short

# Run with coverage
go test ./test/integration -v -cover

# Run specific test category
go test ./test/integration -v -run TestMemscopeEndToEndWorkflow
```

## Test Data

Tests use synthetic test data to validate functionality without requiring actual hardware access.

## Platform Notes

- Windows-specific tests validate WMI integration and hardware access
- Non-Windows tests validate graceful fallbacks and error handling
- All tests should pass on any platform with appropriate messaging