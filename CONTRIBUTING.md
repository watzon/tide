<!--
 Copyright (c) 2024 Christopher Watson
 
 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# Contributing to Tide

Thank you for your interest in contributing to Tide! We love your input! We want to make contributing to Tide as easy and transparent as possible, whether it's:

- 🐛 Reporting a bug
- 📝 Improving documentation
- ✨ Submitting new features
- 🧪 Adding tests
- 🎨 Improving code quality

## Development Process

We use GitHub to host code, track issues and feature requests, as well as accept pull requests.

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code follows the existing style
6. Issue that pull request!

## Testing

Testing is crucial for maintaining code quality. Our current code coverage can be viewed at [Codecov](https://app.codecov.io/gh/watzon/tide/). We welcome contributions that increase our test coverage!

### Running Tests

To run the test suite:

```sh
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run tests for a specific package
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./pkg/backend/terminal

# View coverage report
go tool cover -html=coverage.txt
```

### Writing Tests

We use Go's built-in testing framework. Here's a basic example of a test:

```go
func TestMyFunction(t *testing.T) {
	t.Run("description of test case", func(t *testing.T) {
		result := MyFunction()
		if result != expectedValue {
			t.Errorf("expected %v, got %v", expectedValue, result)
		}
	})
}
```


Some tips for writing good tests:

- Use table-driven tests for testing multiple cases
- Test edge cases and error conditions
- Use meaningful test names and descriptions
- Follow the existing test patterns in the codebase
- Add comments explaining complex test logic

## Code Style

- Follow standard Go formatting (use `gofmt`)
- Write clear, readable code with meaningful variable names
- Add comments for complex logic
- Follow existing patterns in the codebase

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Update any relevant documentation
3. Add tests for new functionality
4. Ensure all tests pass and there are no linting errors
5. The PR will be merged once you have the sign-off of at least one maintainer

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to open an issue with your question or reach out to the maintainers directly.

---

Thank you for contributing to Tide! 🌊