# GitHub Copilot Instructions for Go Projects

## General Go Best Practices
- Write idiomatic Go code. Follow standard Go conventions for naming, formatting, and project structure.
- Emphasize clarity, readability, and simplicity in the code.
- Ensure proper error handling using Go's `error` type. Check errors where they occur and return them to the caller or handle them appropriately.
- Use `golangci-lint` for formatting and linting code according to the project's `.golangci.yml` configuration. Ensure linting passes before considering work complete.
- Write clear and concise Go documentation comments for public APIs.

## Testing (with Testify)
- Generate unit tests for new functions and packages.
- Prefer table-driven tests for covering multiple scenarios concisely.
- Utilize the `testify` suite (especially `testify/assert` and `testify/require`) for more expressive assertions in tests.
- Encourage the use of standard Go testing packages in conjunction with `testify`.
- When suggesting tests, consider edge cases, error conditions, and typical usage.

## HTTP Routing (Chi)
- When working with HTTP services, use the `Chi` router.
- Define routes clearly and group related routes, possibly using Chi's mounting capabilities.
- Utilize Chi's middleware for common concerns like logging, authentication, request ID injection, and recovery.
- Ensure path parameters are correctly parsed and validated.

## Templating (Standard Go html/template)
- For HTML templating, use the standard Go `html/template` package to ensure context-aware escaping and help prevent XSS vulnerabilities.
- Organize templates logically, perhaps in a dedicated directory structure.
- Pass data to templates using structs for clarity and type safety.

## Logging (slog)
- Use the structured logging package `slog` for all application logging.
- Log meaningful, structured information, including relevant context (e.g., request IDs, user IDs, operation names).
- Use appropriate log levels (e.g., `slog.Debug`, `slog.Info`, `slog.Warn`, `slog.Error`).
- Avoid logging sensitive information directly; if necessary, ensure it is properly masked or redacted.
- When logging errors, include the error itself and any relevant contextual attributes.

## Concurrency
- When suggesting concurrent code, ensure it is safe. Highlight potential race conditions or deadlocks if applicable.
- Use channels for communication between goroutines where appropriate.
- Consider using `errgroup` for managing groups of goroutines that can return errors.

## Dependencies and Packages
- Manage dependencies using Go modules (`go.mod` and `go.sum`).
- Strive for well-defined, modular packages with clear responsibilities.
- Avoid circular dependencies between packages.

## Security
- Sanitize all user inputs, especially when dealing with data that will be used in queries, displayed in HTML, or passed to external systems.
- Be mindful of potential security vulnerabilities (e.g., SQL injection, XSS, command injection) and suggest secure coding practices.
- When using `html/template`, ensure proper context is used for variables to prevent XSS.
- Use `gosec` (often integrated via `golangci-lint`) to identify potential security issues.

## Performance
- While clarity is key, suggest efficient algorithms and data structures where appropriate without sacrificing readability.
- Point out potential performance bottlenecks if they are apparent.

## Code Style & Structure
- Keep functions and methods reasonably short and focused on a single responsibility.
- Use meaningful variable and function names.
- Avoid global variables when possible; prefer passing dependencies explicitly (dependency injection).