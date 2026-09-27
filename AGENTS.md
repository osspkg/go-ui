# Repository instructions

## Project map

- `go-sdk/ui` contains the Go protocol types, builders, and validation.
- `go-sdk/mcp` adapts UI manifests and views to MCP resources.
- `react-sdk/ui-core`, `ui-transport`, and `ui-react` contain the TypeScript SDK packages.
- `example/go` and `example/react` are working examples; `schemas/` contains protocol JSON Schemas.

Keep the protocol transport-independent. When changing protocol fields or behavior, check the Go schema and validation, React schema and validation, JSON Schemas, fixtures, and relevant examples for consistency.

## Commands

Run commands from the repository root:

| Command | Purpose |
| --- | --- |
| `make npm-install` | Install the pnpm workspace from the frozen lockfile |
| `make npm-build` | Build the core, transport, and React npm packages |
| `make examples-dev` | Start the Vite React examples |
| `make tests` | Run the Go test workflow through goppy |
| `make lint` | Run the Go lint workflow through goppy |
| `make build` | Run the Go build workflow through goppy |
| `make ci` | Run `pre-commit` (tool setup, license, lint, tests, and build) |

For focused TypeScript workflows, use `pnpm --filter <workspace-package> typecheck`, `test`, or `build` where that script exists. Package names and scripts are defined in each workspace `package.json`.

`make install` installs `goppy@latest` and runs `goppy setup-lib`; `make ci` includes this setup through `pre-commit`. These workflows can update local/generated files, so inspect `git status` and the Makefile before using them when preserving a worktree matters.

`make npm-publish` publishes all three `@osspkg` packages publicly. Never run a publish command unless the user explicitly asked for a release.

## Change guidance

- Preserve the Go and TypeScript representations of the protocol in sync, including layout behavior for nested nodes.
- Keep JSON Schema documents and valid/invalid fixtures aligned with runtime validation.
- Prefer extending the existing SDK APIs and examples over introducing duplicate abstractions.
- Do not commit, publish, or run broad mutating workflows unless the user requests it.
