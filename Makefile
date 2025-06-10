# Test management for microservice-commons

.PHONY: test test-verbose test-coverage test-unit test-integration clean

# Run all tests
test:
	go test ./tests/... -v

# Run tests with verbose output
test-verbose:
	go test ./tests/... -v -race

# Generate test coverage report
test-coverage:
	go test ./tests/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run only unit tests (exclude integration tests)
test-unit:
	go test ./tests/... -short -v

# Test individual packages
test-config:
	go test ./tests/ -run TestConfig -v

test-server:
	go test ./tests/ -run TestServer -v

test-utils:
	go test ./tests/ -run TestUtils -v

test-middleware:
	go test ./tests/ -run TestMiddleware -v

test-responses:
	go test ./tests/ -run TestResponses -v

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html

# Benchmark tests
benchmark:
	go test ./tests/... -bench=. -benchmem

# Check test coverage percentage
coverage-check:
	@go test ./tests/... -coverprofile=coverage.out > /dev/null
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'