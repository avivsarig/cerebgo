# Cerebgo - Functional Markdown Task Management

Cerebgo is a Go-based task management system that processes markdown files.
I designed it to process markdown files in private repos - originally for my Obsidian notes, while the code itself is available at a public repository.
It is a learning project to learn Go, and practice pseudo-functional programming to expand my knowledge.

## Philosophy and Intended Use

This is a personal take on `Getting Things Done` and `Building a Second Brain`.

The system helps organize information and tasks through several key components:

- Notes serve as an Inbox - a place to quickly capture information and action items before they slip away - ideally it should be cleared as often as possible

- Daily journals are used to record events and link them to tasks, people, and archived records for better context and retrieval

- Tasks represent actions that need to be done - from simple actions like "Turn the light off" to complex projects like "Research before buying a car". They are handled differently based on their nature:

  - Tasks with content are marked as projects and moved to archives when done, preserving the knowledge
  - Simple tasks without content remain briefly after completion, then are cleaned up automatically

- Archives store collected knowledge with tags for easy retrieval

- Lists (name pending) keep track of media to consume - articles, movies, and books

## Current Features

The system processes task files in two main places:

- Active Tasks - These need attention and action
- Completed Tasks - These are done and waiting to be either archived or cleaned up

Currently, the system helps by:

- Updating `Do Date` (when a task is planned to be worked on) to not be in the past
- Converting standalone tasks to projects when content is added
- Cleaning up completed tasks based on their type and age

## Approach and Architecture

### Repository Separation

This system is meant to be run on two separated Github Repositories:

Private Repository

- This repository stores all markdown files in the structure mentioned in [Philosophy and Intended Use](##Philosophy-and-Intended-Use)

- The repository should stay private to protect your files, assuming this system is used for personal use

- Run the Github Actions Workflows from this repository to process your files in a secure manner

Public Repository (this one)

- Contains the processing logic in executable scripts, including:
  - Task management engine
  - Supporting modules for parsing, process and file management
  - Docker packaging and Github Actions workflows

#### Security Considerations

The repository separation design provides several security benefits:

- Data Privacy:
  - Personal information remains in your private repository
  - Processing logic is public but can't access your data without proper credentials
- Access Control
  - The public repository contains no sensitive information
  - Private repository access is limited to authorized users
  - GitHub Actions use restricted tokens
- Configuration Security
  - Sensitive paths and settings are stored in the private repository
  - Environment variables are used for runtime configuration
  - No hardcoded secrets in the codebase
- Process Isolation
  - Docker containers run with limited permissions
  - File system access is restricted to configured paths
  - No network access required during processing

### Data Flow

The process:

1. Tasks are created as markdown files with YAML frontmatter
2. The system processes these files using configurable rules
3. Task states and metadata get updated automatically
4. Changes are saved back to the files

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Git
- GitHub account
- Docker (optional)

### Setting Up Private Repository

1. Create a new private GitHub repository
2. Create the following directory structure:

```
.
├── config/
│   └── config.yaml
├── tasks/
│   └── completed/
├── journal/
│   └── completed/
├── archives/
└── lists/
```

3. Create the configuration file at "/config/config.yaml"

```
paths:
  base:
    inbox: inbox.md
    tasks: /tasks
    journal: /journal
    lists: /lists
    people: /people
    archives: /archives

  subdirs:
    tasks:
      completed: ${paths.base.tasks}/completed

    journal:
      completed: ${paths.base.journal}/completed

# User Journals
journals:
  - name: belle
  - name: pure

# System Settings
settings:
  retention:
    empty_task: 30
    project_before_archive: 7

  patterns:
    date_format: "YYYY-MM-DD"
    file_format: "*-${date_format}"
```

4. Set up GitHub Actions:
   - Create `.github/workflows/process.yml`

```
name: Process Tasks

on:
  push:
    branches: [ main ]
  schedule:
    - cron: '0 _/6 _ \* \*' # Run every 6 hours
  workflow_dispatch:

jobs:
process:
runs-on: ubuntu-latest
steps: - uses: actions/checkout@v4

      - name: Process Tasks
        uses: docker://ghcr.io/avivsarig/cerebgo:latest
        env:
          CONFIG_PATH: ${{ secrets.CONFIG_PATH }}

      - name: Commit changes
        run: |
          git config --global user.name 'GitHub Action'
          git config --global user.email 'action@github.com'
          git add .
          git diff --staged --quiet || git commit -m "Auto-process tasks"
          git push
```

- Add the following secure secrets:
  - `CEREBGO_TOKEN`: GitHub token with repo access
  - `CONFIG_PATH`: Path to your config directory

### Running Locally

```bash
# Clone both repositories
git clone https://github.com/yourusername/cerebgo-private private
git clone https://github.com/avivsarig/cerebgo cerebgo

# Build Cerebgo
cd cerebgo
go build -o cerebgo ./cmd/main/main.go

# Run with your config
export CONFIG_PATH=/path/to/private/config
./cerebgo
```

### Docker Deployment

```bash
# Build image
docker build -t cerebgo .

# Run with mounted config
docker run -v /path/to/private/config:/app/config cerebgo
```

## Configuration

The system uses a YAML configuration file (config.yaml) that defines:

- Directory structures for tasks and journals
- File naming patterns
- Retention policies for completed tasks
- System-wide settings

Example configuration:

```yaml
paths:
  base:
    tasks: /tasks
    journal: /journal

settings:
  retention:
    empty_task: 30 # days to keep completed non-project tasks
    project_before_archive: 7 # days to keep completed projects
```

## Project Roadmap

### Soon

- Support repetitive tasks
- Daily reports and planning
- Journal and Notes processing for action items and archive records extraction
- Better testing

### Eventually

- LLM integration for:
  - Recommendation of relevant archive records in newly created tasks
  - Periodical Journal summaries
  - Better reports
  - Tasks creation for Lists

## License

This code is under the MIT License - see the LICENSE file for the details.
