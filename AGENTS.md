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
| `make install` | Install frozen pnpm workspace dependencies, install `goppy@latest`, and run `goppy setup-lib` |
| `make npm-build` | Build core, transport, and React packages in dependency order |
| `make lint` | Run pnpm lint and formatting, then `goppy lint` |
| `make tests` | Run pnpm tests, then `goppy test` |
| `make build` | Build the three npm packages, run pnpm typecheck, then `goppy build --arch=amd64` |
| `make license` | Run `goppy license` |
| `make pre-commit` | Run install, license, lint, tests, and build in sequence |
| `make examples-dev` | Start the Vite React examples |
| `make npm-publish` | Build and publish the three `@osspkg` packages publicly |

For focused TypeScript workflows, use `pnpm --filter <workspace-package> typecheck`, `test`, or `build` where that script exists. Package names and scripts are defined in each workspace `package.json`.

`make install` installs `goppy@latest` and runs `goppy setup-lib`. `make lint` runs `pnpm format`, which rewrites files before checking formatting, and the Go lint workflow may also update generated or formatted files. Inspect `git status` before running these workflows when preserving a worktree matters. `make pre-commit` includes install, license, lint, tests, and build; there is no `make ci` target.

The `tests` and `build` targets depend on `npm-build` so package exports exist in clean checkouts. `npm-publish` also builds first and publishes publicly. Never run a publish command unless the user explicitly asked for a release.

## Change guidance

- Preserve the Go and TypeScript representations of the protocol in sync, including layout behavior for nested nodes.
- Keep JSON Schema documents and valid/invalid fixtures aligned with runtime validation.
- Prefer extending the existing SDK APIs and examples over introducing duplicate abstractions.
- Do not commit, publish, or run broad mutating workflows unless the user requests it.
