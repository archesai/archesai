# Arches CLI Reference

Complete reference for all `archesai` CLI commands.

## Installation

```bash
# Install from source
go install github.com/archesai/archesai/cmd/archesai@latest

# Or build locally
make build
```

The binary will be available as `archesai`.

## Global Flags

These flags are available for all commands:

- `--config` - Config file path (default: `.archesai.yaml`)
- `-v, --verbose` - Enable verbose output
- `--pretty` - Enable pretty logging output

## Commands

### `archesai`

The base command for the Arches platform.

```bash
archesai [flags]
archesai [command]
```

Running `archesai` without a subcommand starts the API server.

---

### `archesai dev`

Run development server with hot reload. This command runs both the API server
(with air for hot reload) and the Platform UI (with Vite) concurrently.

```bash
archesai dev [flags]
```

**Flags:**

- `--tui` - Enable TUI mode for interactive log viewing

**Example:**

```bash
# Run development server with standard output
archesai dev

# Run with interactive TUI for better log management
archesai dev --tui
```

**What it does:**

1. Starts the API server on port 3001 with hot reload
2. Starts the Platform UI on port 3000 with Vite
3. Watches for file changes and automatically restarts

---

### `archesai generate`

Generate code from specifications. This is the core command for code generation.

#### `archesai generate openapi`

Generate complete application code from an OpenAPI specification.

```bash
archesai generate openapi [spec-path] [flags]
```

**Arguments:**

- `spec-path` - Path to the OpenAPI specification file (YAML or JSON)

**Required Flags:**

- `--output` - Output directory for generated code

**Optional Flags:**

- `--bundle` - Bundle OpenAPI spec into single file instead of generating code
- `--orval-fix` - Apply fixes for Orval compatibility (only with --bundle)

**Example:**

```bash
# Generate full application from OpenAPI spec
archesai generate openapi api/openapi.yaml --output ./generated

# Bundle a multi-file OpenAPI spec into a single file
archesai generate openapi api/openapi.yaml --bundle --output api/bundled.yaml

# Bundle with Orval compatibility fixes
archesai generate openapi api/openapi.yaml --bundle --orval-fix --output api/bundled.yaml
```

**Generated Files:**

- `models/` - Data models and types
- `repositories/` - Database access layer
- `controllers/` - HTTP controllers
- `handlers/` - Business logic handlers
- `events/` - Event definitions
- `client/` - TypeScript/JavaScript client SDK
- `database/` - SQL migrations and schema
- `bootstrap/` - Application initialization code

#### `archesai generate jsonschema`

Generate Go structs from JSON Schema files.

```bash
archesai generate jsonschema [spec-path] [flags]
```

**Arguments:**

- `spec-path` - Path to the JSON Schema file

**Required Flags:**

- `--output` - Output file path for generated Go code

**Example:**

```bash
# Generate Go structs from JSON Schema
archesai generate jsonschema schemas/user.json --output models/user.go
```

---

### `archesai config`

Manage Arches configuration files.

#### `archesai config show`

Display the current configuration.

```bash
archesai config show [flags]
```

**Flags:**

- `-o, --output` - Output format: `yaml`, `json`, `tui` (default: `yaml`)

**Example:**

```bash
# Show config as YAML
archesai config show

# Show config as JSON
archesai config show -o json

# Show config in interactive TUI
archesai config show -o tui
```

#### `archesai config validate`

Validate configuration file for errors.

```bash
archesai config validate
```

**Example:**

```bash
# Validate the current configuration
archesai config validate
```

#### `archesai config init`

Create a default configuration file.

```bash
archesai config init [path]
```

**Arguments:**

- `path` - Path for config file (default: `config.yaml`)

**Example:**

```bash
# Create default config.yaml
archesai config init

# Create config at specific path
archesai config init myapp.yaml
```

#### `archesai config env`

Display all Arches-related environment variables.

```bash
archesai config env
```

**Example:**

```bash
# Show all environment variables
archesai config env
```

---

### `archesai version`

Print version information including version number, commit hash, and build date.

```bash
archesai version
```

**Example Output:**

```plaintext
Arches CLI
Version: v1.0.0
Commit: abc123def
Built: 2024-11-22T10:00:00Z
```

---

### `archesai completion`

Generate shell completion scripts for easier CLI usage.

```bash
archesai completion [shell]
```

**Supported Shells:**

- `bash`
- `zsh`
- `fish`
- `powershell`

**Installation Examples:**

**Bash:**

```bash
# Add to ~/.bashrc
echo 'source <(archesai completion bash)' >> ~/.bashrc
```

**Zsh:**

```bash
# Add to ~/.zshrc
echo 'source <(archesai completion zsh)' >> ~/.zshrc
```

**Fish:**

```bash
# Add to config
archesai completion fish > ~/.config/fish/completions/archesai.fish
```

**PowerShell:**

```powershell
# Add to profile
archesai completion powershell | Out-String | Invoke-Expression
```

---

## Planned Commands

These commands are planned but not yet implemented. See [ROADMAP.md](ROADMAP.md) for details.

### `archesai new` (Coming Soon)

Create a new Arches project from templates.

```bash
# Planned syntax
archesai new [project-name] [flags]
```

### `archesai init` (Coming Soon)

Initialize Arches in an existing project.

```bash
# Planned syntax
archesai init [flags]
```

### `archesai deploy` (Coming Soon)

Deploy application to cloud platforms.

```bash
# Planned syntax
archesai deploy [platform] [flags]
```

---

## Configuration File

Arches uses a YAML configuration file (`.archesai.yaml` by default).

### Example Configuration

```yaml
# .archesai.yaml
app:
  name: my-app
  version: 1.0.0

server:
  host: localhost
  port: 3001

database:
  driver: postgres
  host: localhost
  port: 5432
  name: archesdb
  user: postgres

redis:
  host: localhost
  port: 6379

auth:
  jwt_secret: your-secret-key
  token_expiry: 24h
```

### Environment Variables

All configuration values can be overridden with environment variables:

```bash
# Override server port
export ARCHES_SERVER_PORT=8080

# Override database host
export ARCHES_DATABASE_HOST=db.example.com

# Override Redis connection
export ARCHES_REDIS_HOST=redis.example.com
```

---

## Common Workflows

### Starting a New Project

```bash
# 1. Create your OpenAPI specification
vim api/openapi.yaml

# 2. Generate the application code
archesai generate openapi api/openapi.yaml --output ./app

# 3. Start development server
archesai dev --tui
```

### Regenerating After Schema Changes

```bash
# Update your OpenAPI spec
vim api/openapi.yaml

# Regenerate code
archesai generate openapi api/openapi.yaml --output ./app

# Changes are automatically picked up by dev server
```

### Preparing for Production

```bash
# 1. Validate configuration
archesai config validate

# 2. Build production binaries
make build

# 3. Run production server
archesai --config production.yaml
```

---

## Troubleshooting

### Command Not Found

If `archesai` is not found:

```bash
# Ensure it's built
make build

# Or install globally
go install github.com/archesai/archesai/cmd/archesai@latest
```

### Port Already in Use

If ports 3000 or 3001 are in use:

```bash
# Change ports in config
export ARCHES_SERVER_PORT=8001
export ARCHES_PLATFORM_PORT=8000
```

### Database Connection Issues

```bash
# Check database is running
docker-compose up -d postgres

# Verify connection settings
archesai config show
```

---

## See Also

- [Getting Started](getting-started.md) - Quick start guide
- [Development Guide](guides/development.md) - Development workflow
- [Code Generation](guides/code-generation.md) - How code generation works
- [Makefile Commands](guides/makefile-commands.md) - Development commands
