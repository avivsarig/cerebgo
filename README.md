# Markdown Task Management System

A functional Go implementation for managing tasks and notes using markdown files. The system separates processing logic (this repository) from content storage, enabling public collaboration while maintaining data privacy.

## Features

- Pure functional approach to markdown processing
- Automated task categorization and organization
- Journal entry processing with action item extraction
- Cross-referenced people mentions
- Reading list management
- Secure multi-repository architecture

## Architecture

The system operates across two repositories:

- Public (this repo): Contains processing logic and GitHub Actions
- Private: Stores markdown content and configuration

### Data Flow

1. Changes in private repository trigger GitHub Actions
2. Processing logic operates on markdown files
3. Content is automatically organized based on rules
4. Updates are committed back to private repository

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Make
- GitHub account (for Actions)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/task-management
cd task-management

# Install dependencies
make dev-deps
```

### Basic Setup

1. Set up your private repository following the [Private Repo Setup Guide](docs/private-repo-setup.md)
2. Configure GitHub Actions in your private repository
3. Add your first markdown files

### Configuration

Configure system behavior through your private repository:

- File organization rules
- Task categorization settings
- Processing schedules
- Security parameters

## Development

### Available Commands

```bash
make help     # Show all available commands
make all      # Run full check suite (fmt, vet, lint, test, build)
make test     # Run tests with race detection
make fmt      # Format code
make lint     # Run linter
make coverage # Generate test coverage report
```

### Development Workflow

```bash
make watch-test  # Run tests automatically on file changes
make docker-test # Run tests in Docker environment
```

## Design Philosophy

- Emphasis on immutable data structures
- Pure functions where possible
- Clear data transformation pipelines
- Strong type safety
- Comprehensive testing
- Security by design

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
