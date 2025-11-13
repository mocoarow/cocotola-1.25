Before handing work over:
- `gofmt` modified Go files (or rely on `goimports`).
- Run `go test ./...` (or `bazelisk test //...`) once tests exist.
- If BUILD/MODULE files changed, rerun `task gazelle` and `task mod-tidy` when relevant.
- Smoke-test binaries with `task run-hello-world`.