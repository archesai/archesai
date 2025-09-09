# Hugo Documentation Site

The ArchesAI documentation site has been migrated from MkDocs to Hugo using the Stack theme for
better performance.

## Overview

- **Location**: `/website/` directory
- **Theme**: [Hugo Stack Theme](https://github.com/CaiJimmy/hugo-theme-stack) v3
- **Build Tool**: Hugo Extended (required for SCSS processing)
- **Output**: Static site generated in `/website/public/`

## Quick Start

### Prerequisites

Hugo Extended is automatically installed via the Makefile, but you can install manually:

```bash
go install -tags extended github.com/gohugoio/hugo@latest
```

### Build Documentation

```bash
# Build static site (installs Hugo if needed)
make build-hugo-docs

# Serve with Go server on port 3000
make run-hugo-docs

# Clean build artifacts
make clean-docs
```

### Local Development

For live reload during development:

```bash
cd website
hugo server -D
# Site available at http://localhost:1313
```

## Project Structure

```text
website/
├── content/          # Documentation content (Markdown files)
│   ├── architecture/ # System architecture docs
│   ├── features/     # Feature documentation
│   ├── guides/       # Developer guides
│   └── _index.md     # Homepage content
├── config/           # Hugo configuration
│   └── _default/
│       ├── config.toml   # Site configuration
│       └── params.toml   # Theme parameters
├── static/           # Static files (images, etc.)
├── public/           # Generated static site (build output)
├── serve.go          # Go static file server
└── go.mod           # Hugo module dependencies
```

## Adding Documentation

### Create New Pages

1. Add Markdown files to appropriate `content/` subdirectory
2. Include front matter:

```markdown
---
title: "Your Page Title"
date: 2025-09-08
draft: false
categories: ["guides"]
---

Your content here...
```

### Update Navigation

Navigation is automatically generated from content structure and front matter.

## Theme Features

The Stack theme provides:

- ✅ **Responsive Design** - Mobile-friendly layout
- ✅ **Dark Mode** - Automatic theme switching
- ✅ **Search** - Full-text search functionality
- ✅ **Table of Contents** - Auto-generated for long pages
- ✅ **Code Highlighting** - Syntax highlighting for code blocks
- ✅ **Categories/Tags** - Automatic content organization
- ✅ **Fast Performance** - Static site generation

## Deployment Options

The generated `website/public/` directory contains all static files and can be deployed to:

- **GitHub Pages** - Upload public/ contents
- **Netlify/Vercel** - Connect repository with build command `make build-hugo-docs`
- **Docker** - Serve with nginx/Apache
- **CDN** - Upload to any static hosting service

## Migration from MkDocs

The following content has been migrated from the original `/docs` directory:

- Architecture documentation → `/website/content/architecture/`
- Developer guides → `/website/content/guides/`
- Feature docs → `/website/content/features/`
- All Markdown files have Hugo front matter added

## Makefile Integration

New commands added to project Makefile:

| Command             | Description                      |
| ------------------- | -------------------------------- |
| `make install-hugo` | Install Hugo Extended            |
| `make build-docs`   | Build static documentation site  |
| `make run-docs`     | Serve documentation on port 3000 |
| `make clean-docs`   | Clean build artifacts            |

## Configuration

### Site Configuration (`config/_default/config.toml`)

- Site title, URL, language settings
- Hugo module configuration
- Pagination settings

### Theme Parameters (`config/_default/params.toml`)

- Sidebar configuration
- Widget settings
- Comment system (disabled)
- Color scheme settings

## Updating Theme

To update the Stack theme to latest version:

```bash
cd website
hugo mod get -u github.com/CaiJimmy/hugo-theme-stack/v3
hugo mod tidy
```

## Performance

Static site generation provides:

- **Fast Load Times** - No server-side processing
- **CDN Friendly** - All assets can be cached
- **SEO Optimized** - Static HTML with proper meta tags
- **Mobile Optimized** - Responsive design and fast rendering

## Support

- **Hugo Documentation**: [https://gohugo.io/documentation/](https://gohugo.io/documentation/)
- **Stack Theme Docs**: [https://stack.jimmycai.com/](https://stack.jimmycai.com/)
- **Stack Theme GitHub**:
  [https://github.com/CaiJimmy/hugo-theme-stack](https://github.com/CaiJimmy/hugo-theme-stack)
