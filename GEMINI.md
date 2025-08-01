# Gemini Development Guide

This document outlines the development standards and procedures for this project. Adhering to these guidelines will ensure the codebase is robust, maintainable, and consistent.

## 1. Version Control (Git)

A clean and descriptive version history is crucial.

### Commit Command

To avoid shell-specific quoting issues, we will use the following simple and reliable command format for all commits:

```bash
git commit -m "type: subject"
```

- The entire message is enclosed in double quotes.
- No complex characters or quotes will be used in the subject line to ensure compatibility.

### Commit Message Convention

We will follow the **Conventional Commits** specification. This creates an explicit and readable commit history. The format is:

`type(scope): subject`

- **type**: Must be one of the following:
    - `feat`: A new feature.
    - `fix`: A bug fix.
    - `docs`: Documentation only changes.
    - `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc).
    - `refactor`: A code change that neither fixes a bug nor adds a feature.
    - `test`: Adding missing tests or correcting existing tests.
    - `chore`: Changes to the build process or auxiliary tools and libraries.
- **scope** (optional): A noun describing the section of the codebase affected.
- **subject**: A concise description of the change in the imperative mood (e.g., "add," "change," not "added," "changed").

**Example:**
`git commit -m "feat: add websocket echo handler"`
`git commit -m "fix(server): correct COOP header value"`

## 2. Testing

We will use a **Test-Driven Development (TDD)** approach.

### Workflow
1.  **Red**: Write a new test that fails because the feature or fix is not yet implemented.
2.  **Green**: Write the minimum amount of code necessary to make the test pass.
3.  **Refactor**: Clean up the code while ensuring all tests still pass.

### Execution
- All tests will be written using Go's built-in `testing` package.
- Run all tests from the project root using:
  ```bash
  go test ./...
  ```

## 3. Code Quality & Style

Consistency is key to readability and maintainability.

### Formatting
All Go code **must** be formatted with `gofmt`. Before committing, run:
```bash
gofmt -w .
```

### Linting
To catch common errors and enforce idiomatic Go, we will use `golangci-lint`. Before committing, run:
```bash
golangci-lint run
```
*(Note: This may require installing the tool first: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)*

## 4. Overall Development Process

1.  **Write Tests**: Add or update tests for the planned change (`go test ./...` should show a failure).
2.  **Implement**: Write the code to satisfy the tests.
3.  **Verify**: Run all tests to ensure they pass (`go test ./...`).
4.  **Format & Lint**: Run `gofmt -w .` and `golangci-lint run`.
5.  **Commit**: Stage the changes (`git add .`) and commit using the Conventional Commits standard.
