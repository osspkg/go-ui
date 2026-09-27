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

- `make npm-install`: frozen pnpm workspace install.
- `make npm-build`: build core, transport, and React packages.
- `make examples-dev`: start the React example app.
- `make tests`, `make lint`, `make build`: Go workflows via goppy.
- `make ci`: setup, license, lint, tests, and build; setup/license steps may modify generated/local files.
- `pnpm --filter @osspkg/ui-core typecheck|test|build`, likewise `@osspkg/ui-transport` and `@osspkg/ui-react`, for focused package scripts.

Inspect the worktree before workflows that generate or update files. `make npm-publish` publishes all three packages publicly; do not run it unless the user explicitly requests a release.
