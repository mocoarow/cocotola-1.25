Before handing work over:
- `gofmt` modified Go files (or rely on `goimports`).
- Run `task check` to verify code quality with golangci-lint.
- Run `task test` to verify all tests pass (small/medium sizes) and generate coverage report.
- If BUILD/MODULE files changed, rerun `task gazelle` with build tags and `task mod-tidy`.
- If OpenAPI spec changed, run `task generate-openapi` to regenerate models.
- Ensure commit messages follow Conventional Commits format (enforced by pre-commit hooks).
- Smoke-test applications: `task run-cocotola-app`, `task run-cocotola-auth`, or `task run-hello-world`.
- For Docker images, verify build with `task build` before pushing.
- Clean up test/local infrastructure with `task clear-test` and `task clear-local-infra` when done.