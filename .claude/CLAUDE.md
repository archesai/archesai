# Arches Assistant Guide

## What is Arches?

@../README.md

## Project Layout

**IMPORTANT: For token-efficient navigation, use the XML structure:**

- **XML Format (70% fewer tokens)**: @../docs/architecture/project-structure.xml
  The XML format provides structured metadata (file counts, types, purposes) and is much more efficient for finding specific directories and understanding project organization.

## Dev Commands

See Makefile Commands

@../docs/guides/makefile-commands.md

## Project Conventions

- **Generate first, code second** - Define in OpenAPI/SQL before implementing
- **Use generated types** - Don't create manual type definitions

## Task Master AI Instructions

**Import Task Master's development workflow commands and guidelines, treat as if import is in the main CLAUDE.md file.**

@./.taskmaster/CLAUDE.md

## Tips

- **Build fails**: `make generate && make lint`
- **Type errors**: Check generated files are up to date
- **Directory moving**: Do not CD into other directories. You should ideally do everything through the Makefile
- **DO NOT SWITCH DIRECTORIES, STAY IN THE ROOT AT ALL TIMES**
- **Do not create your own mocks** - Always try to use mockery and generate from an interface
- **We have done this many times in this project**
- **DO NOT KEEP DEPRECATED OR LEGACY CODE** - Always implement latest patterns
- **Improve test coverage as much as possible**
- **You should only ever use mocks from mockery** that will be found in `mocks_test.go`
- **If you need to get a mocked interface from another package**, alias the interface in your local package and add it to `.mockery.yaml`

### Mockery Guidelines

- **ALWAYS USE MOCKERY FOR GETTING MOCKS** - Never create mocked services or repositories manually
- **Run `go tool -modfile=tools.mod mockery`**
- **We are running Mockery v3**
- **Mockery config is `.mockery.yaml`**

DO NOT UNDER ANY CIRCUMSTANCE HARD CODE VALUES IN TESTS OR IN TEMPLATES OR IN ANYTHING ELSE
FOR THE SAKE OF HANDLING SPECIAL CASES. ALWAYS DO IT THE CORRECT WAY.

DO NOT ADD TODOS IN THE CODE. IF YOU NEED TO CREATE A TASK, CREATE IT IN THE PROJECT MANAGEMENT TOOL.

DO NOT EVER EVER USER INTERFACE{} OR ANY.
