# ArchesAI Documentation

This directory contains the comprehensive documentation for the ArchesAI project, built using MkDocs Material.

## Setup

### Prerequisites

- Python 3.8 or later
- pip (Python package installer)

### Installation

1. Install MkDocs and dependencies:

   ```bash
   pip install -r requirements.txt
   ```

2. Serve the documentation locally:

   ```bash
   mkdocs serve
   ```

3. Open [http://localhost:8000](http://localhost:8000) in your browser.

### Building

To build the static site:

```bash
mkdocs build
```

The generated site will be in the `site/` directory.

## Structure

```text
docs/
├── mkdocs.yml              # MkDocs configuration
├── requirements.txt        # Python dependencies
├── index.md               # Homepage
├── DEVELOPMENT.md         # Development guide
├── api-reference/         # API documentation
├── architecture/          # Architecture docs
├── development/           # Development guides
├── deployment/            # Deployment guides
├── features/              # Feature documentation
├── troubleshooting/       # Troubleshooting guides
├── security/              # Security documentation
└── performance/           # Performance guides
```

## Features

- **Material Design**: Modern, responsive theme
- **Search**: Fast client-side search functionality
- **Mermaid Diagrams**: Built-in diagram rendering
- **Code Highlighting**: Syntax highlighting for multiple languages
- **Dark/Light Mode**: Toggle between themes
- **Versioning**: Support for multiple documentation versions
- **Git Integration**: Show last modified dates

## Contributing

When adding new documentation:

1. Follow the existing structure and navigation in `mkdocs.yml`
2. Use the markdownlint configuration in `.markdownlint.json`
3. Include the page in the navigation structure
4. Test locally with `mkdocs serve`

## Deployment

Documentation is automatically built and deployed via GitHub Actions when changes are pushed to the main branch.
