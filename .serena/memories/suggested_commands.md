## Testing & Quality
- `task test`: Run all tests with coverage (small/medium sizes), starts MySQL via Docker Compose, generates coverage.lcov
- `task test-remote`: Same as test but uses remote build cache
- `task clear-test`: Stop and remove test Docker containers
- `task check`: Run golangci-lint on all Go projects with 5-minute timeout

## Building & Running
- `task build`: Build all Go apps (cocotola-empty) as OCI images
- `task run-cocotola-app`: Run main application (starts init service + Jaeger)
- `task run-cocotola-auth`: Run authentication service (starts Jaeger)
- `task run-cocotola-init`: Run initialization service (starts Jaeger)
- `task run-cocotola-empty`: Run empty/template application
- `task run-hello-world`: Run sample hello-world binary
- `task run-third-party-library`: Run third-party library demo
- `task clear-local-infra`: Stop and remove local Jaeger containers

## Bazel & Dependencies
- `task gazelle`: Regenerate BUILD files via Gazelle with test size tags (small,medium,large)
- `task mod-tidy`: Run `bazelisk run @rules_go//go -- mod tidy` to keep go.mod files clean
- `task update-mod`: Update all Go project dependencies with `go get -u ./...` and tidy
- `task init`: Initialize/improve Go workspace references (rarely needed once set up)

## OpenAPI & Code Generation
- `task download-openapi-yaml`: Download openapi.yaml from Apidog to openapi/openapi.yaml
- `task generate-openapi`: Generate OpenAPI models using custom templates, copy to cocotola-auth/openapi/

## Container Management
- `task push TAG=<version>`: Push all Go app images to registry with specified tag
- `task kics`: Run Checkmarx KICS security scanner via Docker

## Remote Builds
- Add `--config=remote` or use `-remote` task variants (e.g., `task run-hello-world-remote`) for remote build cache