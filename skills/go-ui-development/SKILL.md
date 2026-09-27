---
name: go-ui-development
description: Implement and integrate features in the go-ui declarative protocol, Go SDK, React SDK, and examples while keeping their behavior aligned.
---

# go-ui development

Use this skill for development involving the `go-ui` protocol, Go or React SDKs, MCP adapter, transports, or examples.

## Principles

- Treat the wire contract below as shared by Go and TypeScript. Preserve parity in field names, enum values, reference resolution, validation, limits, and layout behavior.
- Keep the UI schema independent of transport. MCP resource publication and host RPC handling are separate concerns.
- A node's `layout` positions that node. Each node with children owns a new 12-column grid, so each child chooses its own `row`, `cols`, and `offset` within that immediate parent.
- Keep validation and safe value/path handling at the schema boundary. React component props and events must pass through registry allowlists; never forward arbitrary plugin props.
- For implementation work, inspect the existing corresponding Go and TypeScript behavior before changing the contract. Update schemas, fixtures, examples, and behavior checks together when applicable.
- Read [the protocol reference](references/protocol-map.md) for the complete model and [the examples reference](references/examples.md) for Go, React, and host usage snippets.

## Development commands

Run from the repository root. The workspace uses Go 1.26.8 and pnpm 12.6.0.

- `make install`: install frozen pnpm dependencies, install `goppy@latest`, and run `goppy setup-lib`.
- `make lint`: run pnpm lint, rewrite files with `pnpm format`, check formatting, and run `goppy lint`.
- `make tests`: run pnpm tests and `goppy test`. In a clean checkout, run `make build` first because this target does not build npm packages.
- `make build`: build `@osspkg/ui-core`, `@osspkg/ui-transport`, and `@osspkg/ui-react` in order, run pnpm typecheck, and run `goppy build --arch=amd64`.
- `make license`: run `goppy license`.
- `make pre-commit`: run install, license, lint, build, and tests in sequence. There is no `make ci` target.
- `make examples-dev`: start the React example app.
- `pnpm --filter @osspkg/ui-core typecheck|test|build`, likewise `@osspkg/ui-transport` and `@osspkg/ui-react`, for focused package scripts.

Inspect the worktree before running `make install`, `make lint`, or `make pre-commit`: setup, formatting, lint, and build workflows can update local, generated, or formatted files. `make npm-publish` runs the full build, then publishes all three packages publicly. Do not run publishing unless the user explicitly requests a release.
