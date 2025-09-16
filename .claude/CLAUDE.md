# Arches Assistant Guide

## What is Arches?

@../README.md

## Project Layout

@../docs/architecture/project-layout.md
@../docs/architecture/overview.md

## Dev Commands

See Makefile Commands

@../docs/guides/makefile-commands.md

## Project Conventions

- **Generate first, code second** - Define in OpenAPI/SQL before implementing
- **Use generated types** - Don't create manual type definitions

## Code Generation

After modifying:

- Run `make generate` to run all of the generators.

## Testing

```bash
make test
```

See more at @../docs/guides/testing.md

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
- **Run `go tool mockery`**
- **We are running Mockery v3**
- **Mockery config is `.mockery.yaml`**
