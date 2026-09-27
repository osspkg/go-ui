# go-ui protocol reference

This file is a self-contained working reference for protocol 1.0. The protocol is JSON and independent of Go, React, and MCP. The Go module import path is `go.osspkg.com/ui`; SDK packages are `go.osspkg.com/ui/go-sdk/ui` and `go.osspkg.com/ui/go-sdk/mcp`. The npm packages are `@osspkg/ui-core`, `@osspkg/ui-transport`, and `@osspkg/ui-react`.

## Object model

### Manifest

```json
{
  "protocolVersion": "1.0",
  "plugin": { "id": "profile", "title": "Profile", "version": "1.0.0" },
  "views": [{ "id": "profile.edit", "title": "Edit profile", "schema": "ui://views/profile.edit" }],
  "requiredComponents": ["card"]
}
```

`plugin.id`, `plugin.title`, `plugin.version`, and the `views` array are required. A view descriptor has `id`, `title`, and `schema`; `route` is optional. `requiredComponents` is optional.

### View

A view contains `protocolVersion`, `id`, and all five `regions`. Optional members are `revision`, `title`, `state`, `sources`, and `actions`.

Region keys are exactly:

- `top-header`
- `left-panel`
- `content`
- `right-panel`
- `bottom`

Each region is an array of nodes.

### Node and layout

A node has `id` and `component`. Optional members are `layout`, `props`, `events`, `children`, `slots`, and `when`.

```json
{
  "id": "profile-card",
  "component": "card",
  "layout": { "row": 1, "cols": 8, "offset": 2 },
  "children": [
    { "id": "heading", "component": "h1", "layout": { "row": 1, "cols": 12 }, "props": { "text": "Profile" } },
    { "id": "name", "component": "input", "layout": { "row": 2, "cols": 8, "offset": 2 }, "props": { "placeholder": "Ada Lovelace" } }
  ]
}
```

A layout is a 12-column placement: `row` identifies its grid row; `cols` is the width; `offset` is the number of preceding columns. Children render in their parent's own grid. Nodes without a layout use the available row/width defaults.

### Values and references

`UIValue` is a JSON scalar, array, object, reference, or expression. A reference is a single-key object:

- `{ "$state": "path" }`
- `{ "$context": "path" }`
- `{ "$source": "sourceName.path" }`
- `{ "$event": "path" }`
- `{ "$result": "path" }`

Expressions use one operator key, for example `{ "$eq": [{ "$state": "saved" }, true] }` or `{ "$concat": ["Hello, ", { "$state": "name" }] }`. Preserve safe-path checks and rejection of prototype-sensitive keys such as `__proto__`, `prototype`, and `constructor`.

### Sources, actions, and effects

A source describes a host tool operation:

```json
{ "type": "tool", "tool": "profile.read", "policy": "on-mount", "cache": { "ttl": 30 } }
```

`policy` is `manual` or `on-mount`; optional `input`, `cache.ttl`, and `refreshOn` configure its behavior.

Action types: `tool`, `set-state`, `merge-state`, `refresh-source`, `invalidate`. Common fields include `tool`, `input`, `path`, `value`, `source`, and `effects`.

Effect types: `set-state`, `merge-state`, `invalidate`, `refresh-source`, `refresh-view`, `patch-view`, `navigate`, `toast`, `dialog`, and `close-dialog`. Effects can use `path`, `source`, `value`, `to`, `variant`, and `message` depending on type.

An event handler can name one `action` or a list of `steps`. Each step has `type: "action"`, an action name, and optional input. A view condition `when` takes a `UIValue` expression.

## Package responsibilities

- `@osspkg/ui-core`: protocol types, validation, value/expression resolution, request tracking, and RPC method/request types.
- `@osspkg/ui-transport`: `RPCTransport`, JSON-RPC primitives, `HTTPTransport`, and `WebSocketTransport`.
- `@osspkg/ui-react`: component registry, default HTML/shadcn component definitions, `UIProvider`, `PluginView`, `ViewRenderer`, hooks, and runtime actions/effects.
- Go `ui` package: protocol types, builders, validation, values, and shared RPC method names/types.
- Go `mcp` package: advertises `osspkg.ui` capability and registers `ui://manifest` plus each view schema as MCP resources. It does not implement the host's UI RPC handlers.

Shared RPC method strings are `ui.get`, `data.call`, and `ui.action`. The host implements these methods for remote/runtime use. `PluginView` requests `{ plugin?, view, context? }` through `ui.get`; source and tool actions use `data.call` and `ui.action` respectively. Do not assume registering MCP resources creates these RPC handlers.

## Compatibility checklist

For a wire or runtime behavior change, review all applicable layers:

1. Go protocol structs, builders, validation, and JSON serialization.
2. TypeScript types, validation, resolution, and React runtime/rendering.
3. JSON Schemas for manifest, view, node, layout, actions, and effects.
4. Valid/invalid fixtures and relevant Go/TypeScript tests.
5. Go and React examples, including nested layout placement.

JSON Schema validation covers structural constraints; SDK validators also enforce semantic checks such as references, safe paths, supported actions, and limits. Keep both layers consistent.
