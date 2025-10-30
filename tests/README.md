# Facter-OSS Tests

This directory contains comprehensive tests for the facter-oss project.

## Test Organization

### Unit Tests
Unit tests are located alongside the source code in each package directory with the `_test.go` suffix.

### Integration Tests
Integration tests are located in `tests/integration/` and test the full workflow of the agent including:
- Agent run with various configurations
- Store operations and incremental updates
- Different output formats (JSON, Protocol Buffers)
- Performance profiling mode
- Error handling (corrupted store, etc.)

### End-to-End (E2E) Tests
E2E tests are located in `tests/e2e/` and test the complete CLI workflow:
- CLI execution with config files
- Output generation in different formats
- Multiple runs and store updates
- Debug mode
- Error handling for invalid configurations

## Running Tests

### Run All Tests
```bash
make test
```

### Run Unit Tests Only
```bash
go test ./pkg/...
```

### Run Integration Tests
```bash
go test -v ./tests/integration/...
```

### Run E2E Tests
```bash
go test -v ./tests/e2e/...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Generate Coverage Report
```bash
make test
# This generates coverage.out and cover.html in the tests directory
```

### Skip Long-Running Tests
```bash
go test -short ./...
```

## Test Coverage Summary

Current test coverage by package:

| Package | Coverage |
|---------|----------|
| cmd | 70.0% |
| pkg/agent | 84.6% |
| pkg/agent/collect/users | 100.0% |
| pkg/agent/collect/ssh | 95.2% |
| pkg/agent/store | 92.6% |
| pkg/agent/collectors/initSystem | 92.1% |
| pkg/agent/collectors/process | 91.2% |
| pkg/agent/collect/packages | 90.0% |
| pkg/agent/collectors/ssh | 89.9% |
| pkg/agent/inventory | 82.2% |
| pkg/agent/collect/systemservices | 81.8% |
| pkg/agent/collect/process | 83.3% |
| pkg/utils | 80.0% |
| pkg/agent/collectors/network | 80.5% |

Areas needing improvement:
- pkg/agent/collectors/firewall: 0.0%
- pkg/agent/collect/vulnerability: 6.7%
- pkg/agent/collect/networks: 28.4%
- pkg/agent/collect/platform: 34.4%
- pkg/agent/sink: 35.2%

## Test Types

### Unit Tests
- Test individual functions and methods
- Use mocks and stubs where appropriate
- Fast execution
- No external dependencies

### Integration Tests
- Test multiple components working together
- Test actual agent workflows
- Use temporary directories and files
- Clean up after execution

### E2E Tests
- Test the complete system from CLI to output
- Use the compiled binary
- Test real-world scenarios
- Validate actual output files

## Writing Tests

### Best Practices
1. Use descriptive test names that explain what is being tested
2. Follow the Arrange-Act-Assert pattern
3. Test both success and failure cases
4. Clean up resources (files, directories) after tests
5. Use `t.TempDir()` for temporary directories (automatically cleaned up)
6. Skip network-dependent tests in sandboxed environments
7. Use `-short` flag for quick test runs during development

### Example Test Structure
```go
func TestFeatureName(t *testing.T) {
    // Arrange - setup test data and dependencies
    tmpDir := t.TempDir()
    opts := setupTestOptions(tmpDir)
    
    // Act - execute the code under test
    result, err := functionUnderTest(opts)
    
    // Assert - verify the results
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedValue, result.Value)
}
```

## CI/CD Integration

Tests are run automatically in CI/CD pipelines:
- All unit tests run on every commit
- Integration tests run on pull requests
- E2E tests run before releases
- Coverage reports are generated and tracked

## Troubleshooting

### Tests Fail in CI but Pass Locally
- Check for hardcoded paths or assumptions about the environment
- Verify tests clean up after themselves
- Ensure tests don't depend on external services

### Slow Tests
- Use `-short` flag to skip long-running tests during development
- Consider mocking external dependencies
- Run specific test packages instead of all tests

### Coverage Not Increasing
- Check for error handling paths not covered
- Look for edge cases not tested
- Review generated coverage report HTML for specific uncovered lines

## Contributing

When adding new features:
1. Write unit tests for new functions
2. Add integration tests for new workflows
3. Add E2E tests for new CLI functionality
4. Ensure coverage doesn't decrease
5. Update this README if test structure changes
