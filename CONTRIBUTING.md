# Contributing to kubectl-go-mcp-server

Thank you for your interest in contributing to kubectl-go-mcp-server! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.24 or later
- Docker (optional, for container testing)
- kubectl (for testing Kubernetes integration)
- Make

### Local Development

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/kubectl-go-mcp-server.git
   cd kubectl-go-mcp-server
   ```

2. **Install dependencies**:
   ```bash
   make deps
   ```

3. **Build the project**:
   ```bash
   make build
   ```

4. **Run tests**:
   ```bash
   make test
   ```

5. **Run with coverage**:
   ```bash
   make test-coverage
   ```

### VS Code Setup

We provide example VS Code configurations for a consistent development experience:

1. Copy the example configurations:
   ```bash
   cp .vscode/settings.json.example .vscode/settings.json
   cp .vscode/tasks.json.example .vscode/tasks.json
   cp .vscode/launch.json.example .vscode/launch.json
   ```

2. Install recommended VS Code extensions:
   - Go extension
   - golangci-lint extension

## Code Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting (run `make fmt`)
- Use `golangci-lint` for linting (run `make lint`)
- Write comprehensive tests for new functionality
- Maintain test coverage above 80%

### Git Workflow

1. **Fork the repository** and create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the coding standards

3. **Run the full test suite**:
   ```bash
   make check
   ```

4. **Commit your changes** with a descriptive commit message:
   ```bash
   git commit -m "feat: add new kubectl validation feature"
   ```

5. **Push to your fork** and create a pull request

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

## Testing Guidelines

### Unit Tests

- Write unit tests for all public functions
- Use table-driven tests for multiple test cases
- Place tests in the same package as the code being tested
- Use meaningful test names that describe the scenario

### Integration Tests

- Test real kubectl interactions when possible
- Use mocks for external dependencies in unit tests
- Ensure tests are deterministic and can run in parallel

### Test Coverage

- Aim for >80% test coverage
- Use `make test-coverage` to generate coverage reports
- Focus on testing edge cases and error conditions

## Documentation

### Code Documentation

- Add godoc comments for all exported functions and types
- Include examples in documentation where helpful
- Keep comments concise but informative

### README Updates

- Update the README when adding new features
- Include usage examples for new functionality
- Keep the feature list current

## Security Considerations

### Kubectl Command Safety

- All kubectl commands are validated before execution
- Interactive commands are blocked by default
- Destructive operations require explicit validation
- Input sanitization is mandatory

### Sensitive Information

- Never commit credentials or sensitive configuration
- Use example configuration files for documentation
- Sanitize logs and error messages

## Pull Request Process

1. **Ensure your code passes all checks**:
   ```bash
   make check
   ```

2. **Update documentation** as needed

3. **Add or update tests** for your changes

4. **Create a pull request** with:
   - Clear description of changes
   - Reference to any related issues
   - Screenshots or examples if applicable

5. **Address review feedback** promptly

6. **Ensure CI passes** before requesting final review

## Issue Reporting

### Bug Reports

Include the following information:

- Go version
- Operating system
- kubectl version
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error messages

### Feature Requests

Include:

- Clear description of the feature
- Use case and motivation
- Proposed implementation approach (if applicable)
- Potential breaking changes

## Release Process

The project uses [GoReleaser](https://goreleaser.com/) to automate the release process. Here's how to create a new release:

1. **Update the CHANGELOG.md**:
   - Add your changes to the `[Unreleased]` section
   - Follow the [Keep a Changelog](https://keepachangelog.com/) format

2. **Set the version**:
   ```bash
   export VERSION=X.Y.Z  # Replace with the actual version number
   ```

3. **Create the release**:
   ```bash
   make release VERSION=$VERSION
   ```

4. **Verify the release**:
   - Check the GitHub Actions workflow at the repository's Actions tab
   - Ensure all artifacts are properly uploaded to the GitHub release

The release workflow will:
- Build binaries for multiple platforms
- Create distribution packages (DEB, RPM)
- Generate checksums
- Update the changelog
- Create a GitHub release with release notes
- Upload all artifacts

### Testing Releases Locally

To test the release process without creating a Git tag or GitHub release:
```bash
make release-snapshot
```

This will generate all artifacts in the `dist/` directory without publishing them.

## Getting Help

- **Documentation**: Check the README and inline documentation
- **Issues**: Search existing issues before creating new ones
- **Discussions**: Use GitHub Discussions for questions and ideas

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## License

By contributing to kubectl-go-mcp-server, you agree that your contributions will be licensed under the same license as the project.
