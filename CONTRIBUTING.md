# Contributing to Poodle Go SDK

Thank you for your interest in contributing to the Poodle Go SDK! We welcome contributions from the community.

## Development Setup

### Requirements

- Go 1.20 or higher
- Git

### Setup

1. Fork the repository
2. Clone your fork:

   ```bash
   git clone https://github.com/yourusername/poodle-go.git
   cd poodle-go
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

## Development Workflow

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run tests with verbose output
go test -v ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

### Building Examples

```bash
cd examples
go build ./...
```

### Making Changes

1. Create a feature branch:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes
3. Write or update tests
4. Ensure all tests pass and code is properly formatted
5. Commit your changes with a descriptive message
6. Push to your fork
7. Create a pull request

## Code Standards

### Go Standards

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Maintain Go 1.20+ compatibility

### Testing Standards

- Write unit tests for all new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Test both success and failure scenarios
- Include edge cases in your tests

### Documentation Standards

- Document all exported functions, types, and methods
- Include code examples in documentation
- Update README.md if adding new features
- Write clear commit messages

## Project Structure

```
poodle-go/
├── client.go          # Main client implementation
├── config.go          # Configuration management
├── email.go           # Email model and validation
├── response.go        # Response model
├── errors.go          # Error types and handling
├── http.go            # HTTP client implementation
├── *_test.go          # Test files
├── examples/          # Usage examples
├── .github/           # CI/CD configuration
└── README.md          # Main documentation
```

## Pull Request Guidelines

### Before Submitting

- Ensure all tests pass
- Run `go fmt ./...` and `go vet ./...`
- Update documentation if needed
- Add tests for new functionality
- Check that examples still work

### Pull Request Description

Please include:

- A clear description of the changes
- The motivation for the changes
- Any breaking changes
- Steps to test the changes
- Screenshots (if applicable)

### Review Process

1. All pull requests require review
2. CI checks must pass
3. Code coverage should not decrease
4. Documentation must be updated for new features

## Reporting Issues

### Bug Reports

Please include:

- Go version
- Operating system
- Steps to reproduce
- Expected behavior
- Actual behavior
- Code examples (if applicable)

### Feature Requests

Please include:

- Use case description
- Proposed API (if applicable)
- Examples of how it would be used
- Any alternatives considered

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## License

By contributing to Poodle Go SDK, you agree that your contributions will be licensed under the MIT License.

## Getting Help

- Check existing issues and documentation
- Ask questions in GitHub Discussions
- Contact the maintainers

Thank you for contributing!
