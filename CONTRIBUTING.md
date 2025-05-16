# Contributing to E-commerce Monorepo

Thank you for your interest in contributing to our project! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Development Process

1. Fork the repository
2. Create a new branch for your feature/fix
3. Make your changes
4. Run tests and ensure they pass
5. Submit a pull request

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Update the CHANGELOG.md with your changes
3. The PR will be merged once you have the sign-off of at least one other developer

## Development Setup

1. Install dependencies:

   ```bash
   make install-deps
   ```

2. Set up development environment:

   ```bash
   make setup-dev
   ```

3. Run tests:

   ```bash
   make test
   ```

## Coding Standards

- Follow the Go code style guide for backend services
- Use ESLint and Prettier for frontend code
- Write meaningful commit messages
- Include tests for new features
- Update documentation as needed

## Testing

- Write unit tests for all new features
- Ensure all tests pass before submitting PR
- Include integration tests for API changes
- Update existing tests if needed

## Documentation

- Update README.md if needed
- Document new API endpoints
- Add comments to complex code
- Update CHANGELOG.md

## Questions?

Feel free to open an issue or contact the development team at [dev@example.com](mailto:dev@example.com)
