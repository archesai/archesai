# Arches Studio UI Design

## Overview

The Arches Studio is a web-based IDE for building full-stack applications from OpenAPI specifications. It runs on port 3000 and manages generated applications that run on separate ports (8080+).

## Two-Process Architecture

### Process 1: Arches Studio (Port 3000)

```bash
arches studio
```

This runs:

- **Studio API** - Project management, code generation triggers
- **Studio UI** - Web interface for building apps
- **File watchers** - Monitor for OpenAPI changes
- **Generation orchestration** - Triggers codegen when needed

### Process 2: Generated App (Port 8080+)

```bash
arches run <project-name>
# or from project directory
arches dev
```

This runs:

- **Generated backend API** - The actual app being built
- **Generated frontend** - The app's UI
- **Hot reload** - Restarts when regeneration happens

## Studio UI Layout

```
┌─────────────────────────────────────────────────────┐
│  Arches Studio                                [□ ×] │
├──────────┬──────────────────────────────────────────┤
│          │                                          │
│ Projects │  Schema Editor                           │
│ -------- │  ┌────────────────────────────────────┐  │
│          │  │ openapi: 3.1.0                    │  │
│ my-app   │  │ info:                             │  │
│ todo-api │  │   title: My App                   │  │
│ blog     │  │ paths:                            │  │
│          │  │   /users:                         │  │
├──────────┤  │     get: ...                      │  │
│          │  └────────────────────────────────────┘  │
│ Views    │                                          │
│ -----    │  [Generate] [Validate] [AI Assist]       │
│          │                                          │
│ • Schema │  ┌────────────────────────────────────┐  │
│ • Handlers│ │ Preview (iframe of localhost:8080) │  │
│ • Preview │  │                                    │  │
│ • Deploy  │  │  [Your app running here]          │  │
│ • Logs    │  │                                    │  │
│          │  └────────────────────────────────────┘  │
└──────────┴──────────────────────────────────────────┘
```

## Navigation Structure

### Left Sidebar - Projects

- List of all projects
- Active project highlighted
- Quick actions (new, delete, duplicate)

### Left Sidebar - Views

For the selected project:

- **Schema** - OpenAPI editor with syntax highlighting
- **Handlers** - Custom business logic editor
- **Preview** - Live preview of running app
- **Deploy** - Deployment configuration and triggers
- **Logs** - Development server logs

### Main Content Area

Split view option:

- Top: Editor (Schema or Handlers)
- Bottom: Live Preview (iframe)

Or single view:

- Full editor
- Full preview
- Full logs

### Top Bar Actions

- Generate - Trigger code generation
- Validate - Check OpenAPI validity
- AI Assist - Natural language to OpenAPI
- Start/Stop Dev Server
- View at localhost:8080 (external)

## Interaction Flow

1. **Select Project** - Click project in sidebar
2. **Edit Schema** - Modify OpenAPI in editor
3. **Auto-Generate** - On save, triggers generation
4. **Preview Updates** - iframe refreshes with new app
5. **Edit Handlers** - Switch to handlers view for custom logic
6. **Test Changes** - See results immediately in preview

## Technical Implementation

### Frontend Stack

- React 19 with TypeScript
- Monaco Editor for code editing
- WebSocket for hot reload notifications
- TanStack Router for navigation
- Iframe for app preview

### Communication

```
Studio UI ←→ Studio API ←→ File System ←→ Generated App
         WebSocket      Watch & Generate    Auto-restart
```

### File Watching

- Watch `projects/*/api/openapi.yaml`
- Watch `projects/*/handlers/*`
- Trigger regeneration on changes
- Notify UI via WebSocket

## Responsive Design

### Desktop (1920px+)

- Full three-column layout
- Sidebar: 250px
- Editor: 50% of remaining
- Preview: 50% of remaining

### Laptop (1366px-1919px)

- Collapsible project list
- Editor: 60%
- Preview: 40%

### Tablet (768px-1365px)

- Tab-based switching between editor and preview
- Hamburger menu for project list

### Mobile (< 768px)

- Stack views vertically
- Bottom sheet for preview
- Swipe between views

## Theme Support

### Light Mode (Default)

- White background
- Gray-900 text
- Blue-600 accents
- Gray-100 borders

### Dark Mode

- Gray-900 background
- Gray-100 text
- Blue-400 accents
- Gray-800 borders

## Keyboard Shortcuts

- `Cmd/Ctrl + S` - Save and regenerate
- `Cmd/Ctrl + P` - Quick project switch
- `Cmd/Ctrl + K` - Command palette
- `Cmd/Ctrl + Shift + P` - Toggle preview
- `Cmd/Ctrl + B` - Toggle sidebar
- `Cmd/Ctrl + Enter` - Generate code

## Future Enhancements

1. **Split Panes** - Multiple files open simultaneously
2. **Git Integration** - Version control in UI
3. **Collaboration** - Real-time multi-user editing
4. **Templates Gallery** - Pre-built app templates
5. **Marketplace** - Share handlers and schemas
6. **Cloud Sync** - Projects backed up to cloud
7. **AI Chat** - Integrated chat for assistance
8. **Debugging** - Integrated debugger for handlers
