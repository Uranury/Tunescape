# Contributing to Tunescape

## Branching

We follow GitFlow. The two permanent branches are `main` and `dev`.

- `main` is production-ready code only. Never push directly to it.
- `dev` is the integration branch. All features merge here first.

Branch off `dev` for any new work:

```
feature/spotify-oauth
fix/token-expiry
refactor/service-x
chore/update-dependencies
docs/update-readme
```

Never push directly to `main` or `dev`. All changes go through a Merge Request.

## Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/). Every commit message should be prefixed with a type:

| Prefix          | When to use |
|-----------------|---|
| `feature/feat:` | Adding a new feature |
| `fix:`          | Fixing a bug |
| `chore:`        | Maintenance, config, CI/CD, dependencies |
| `docs:`         | Documentation only |
| `refactor:`     | Code change that isn't a fix or feature |
| `test:`         | Adding or updating tests |

Examples:

```
feat: add Spotify OAuth2 callback handler
fix: refresh token not rotating on reuse
chore: add golangci-lint to CI pipeline
refactor: split jwtService and refreshTokenService
```

Keep the subject line short and in the imperative mood — "add handler" not "added handler".

## Merge Requests

Before opening an MR make sure your branch is up to date with `dev`:

```bash
git fetch origin
git rebase origin/dev
```

Every MR should have a short description explaining what the change does and why. For anything non-obvious — a transaction, a locking strategy, a schema decision — explain the reasoning. Your teammates learn from it.

MRs require at least one approval before merging. The pipeline must pass — lint, tests, and build all green.

Keep MRs focused. One concern per MR. A 1000-line MR that touches auth, analytics, and CI at once is hard to review and hard to revert if something goes wrong.

## Project Structure

```
.
├── cmd/
│   └── api/              # Entrypoint (main.go)
├── internal/
│   ├── app/              # HTTP server and routing
│   ├── auth/             # JWT and refresh token logic
│   └── infra/            # Infrastructure wiring (DB, Redis, config)
├── pkg/
│   ├── config/           # Config loading
│   └── database/         # DB init, migrations, transaction support
├── migrations/           # SQL migration files
```

When adding a new feature, follow the existing layering: handler → service → repository. Handlers should not touch the database directly. Services should not build HTTP responses.

## Code Style

We use `golangci-lint` to enforce code quality. Run it locally before pushing:

```bash
golangci-lint run ./...
```

A few rules we follow regardless of the linter:

- Always wrap errors with context using `fmt.Errorf("what failed: %w", err)`
- Define sentinel errors as package-level variables (`var ErrNotFound = errors.New(...)`) rather than inline strings
- Use `context.Context` as the first argument in any function that touches the DB or makes network calls
- Keep constructors in the form `NewX(deps) X`

## Running Locally

See the [README](README.md) for setup instructions.