# GitHub Copilot Instructions

This file contains specific instructions and reminders for GitHub Copilot when working on this repository.

IMPORTANT: Do NOT commit or push changes to the repository. Prepare patches or suggest edits and ask the user to commit the changes to the remote branch.

IMPORTANT: Do not use Python or Ruby. They are not installed on the system and will lead to errors.

## Dependency Updates

## Go Version Updates

When updating the Go version:
1. Update `go.mod` file
2. Update `.github/workflows/ci.yml` - change the `go-version` field
3. Run `go get -u ./...` to update all dependencies
4. Run `go mod tidy` to clean up dependencies
5. Verify build and core tests still pass

## GitHub Actions Workflow Updates

Always use the newest stable versions of GitHub Actions in workflows. When updating action versions, follow the step-by-step process below to avoid introducing CI breakage and to minimize unnecessary network requests or trial-and-error steps.

Step-by-step process for updating action versions (minimize premium/trial requests):
1. Search the workflow files under `.github/workflows/` for `uses:` entries that reference third-party actions.
2. For each action, visit the action's official repository or the Marketplace page to find the latest stable release tag. Prefer a pinned MAJOR tag (for example `actions/checkout@v4`). Avoid pinning to an exact patch (for example `@v4.0.0`) unless you need fully reproducible runs — major-only tags receive compatible non-breaking updates (minor/patch) and are easier to maintain. When a new major release is available, evaluate compatibility and bump the major tag if safe. A successful pipeline run is needed to complete a pull request. You can do network calls for that.
3. Update the workflow to the new tag. Example replacements:
	- `actions/checkout@v3` → `actions/checkout@v4`
	- `actions/setup-go@v4` → `actions/setup-go@v5`
	- `docker/build-push-action@v4` → `docker/build-push-action@v5`

Best practices and safety notes:
- Prefer pinned MAJOR tags (e.g., `@v4`) rather than floating tags like `@main` or `@master`.
- NEVER pin to an exact patch (e.g., `@v4.0.0`) unless you require byte-for-byte reproducibility of CI runs. If you do pin an exact patch, document the reason in the PR and include steps to test updates.
- Read the action's changelog for breaking changes before updating.

Recommended usage examples:
- Use `actions/checkout@v5` to get stable, non-breaking bug/security fixes within v5.
- Use `actions/checkout@v5.0.0` only if you need to lock the runner to that exact release (rare).

Examples of common actions in this repo (verify these when updating):
- `actions/checkout@v4`
- `actions/setup-go@v5`
- `codecov/codecov-action@v4`
- `docker/build-push-action@v5`
- `actions/github-script@v7`

Rollback guidance:
- If an updated action breaks CI in a way that's not immediately fixable, revert the workflow commit

## Code Quality

- Maintain existing code style and patterns
- Ensure all exported functions have proper documentation
- Add unit tests for new functionality
- Keep dependency updates atomic and well-documented

## Testing
- All tests should always pass
- Run `go test ./internal/app/document ./internal/app/server` to verify core functionality