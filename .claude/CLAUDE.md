# Arches Assistant Guide

## What is Arches?

@../README.md

## Rules

- **Generate first, code second** - Define in OpenAPI/SQL before implementing
- **Type errors**: Check generated files are up to date
- **Directory moving**: Do not CD into other directories. You should ideally do everything through the Makefile
- **DO NOT SWITCH DIRECTORIES, STAY IN THE ROOT AT ALL TIMES**
- **Do not create your own mocks** - Always try to use mockery and generate from an interface
- **We have done this many times in this project**
- **DO NOT KEEP DEPRECATED OR LEGACY CODE** - Always implement latest patterns
- **Improve test coverage as much as possible**
- **Do not under any circumstance hard code values in tests or in templates or in anything else for the sake of handling special cases. Always do it the correct way.**
- **Do not add todos in the code. If you need to create a task, create it in the project management tool.**
- **Do not ever ever user interface{} or any.**
- **Use generated types** - Don't create manual type definitions
- **Never manually update a .gen.go file**
