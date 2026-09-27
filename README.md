# go-ui

[![Go Reference](https://pkg.go.dev/badge/go.osspkg.com/ui/go-sdk/ui.svg)](https://pkg.go.dev/go.osspkg.com/ui/go-sdk/ui)

`go-ui` is a declarative UI protocol and SDK for describing plugin interfaces as JSON schemas. The Go SDK builds and validates manifests and views; the React SDK validates those schemas and renders them with registered HTML and shadcn components. The UI model is independent of the transport, with an MCP adapter available in `go-sdk/mcp`.

## Features

- Declarative manifests, views, regions, nodes, layouts, data sources, actions, and effects.
- Go builders for composing and validating UI schemas.
- React packages for schema validation, transport, and rendering.
- Default HTML and shadcn component registries, with support for custom components.
- Child layout using `row`, `cols`, and `offset` in a 12-column grid.
- JSON Schemas for protocol objects in [`schemas/`](schemas/).

## Repository layout

| Path | Description |
| --- | --- |
| [`go-sdk/ui`](go-sdk/ui) | Go schema types, builders, validation, and component definitions |
| [`go-sdk/mcp`](go-sdk/mcp) | Adapter that registers the UI manifest and views as MCP resources |
| [`react-sdk/ui-core`](react-sdk/ui-core) | TypeScript protocol types, validation, and runtime helpers |
| [`react-sdk/ui-transport`](react-sdk/ui-transport) | HTTP, WebSocket, and RPC transport helpers |
| [`react-sdk/ui-react`](react-sdk/ui-react) | React renderer, registries, and default components |
| [`example/go`](example/go) | Go examples for an HTML form and dashboard |
| [`example/react`](example/react) | Vite app with interactive React examples |
| [`schemas`](schemas) | JSON Schema documents for protocol entities |

## Requirements

- Go 1.26.8 or newer.
- Node.js and pnpm 12.6.0, as declared by the workspace package manifest.

## Go usage

The root Go module is `go.osspkg.com/ui`. Add it to a Go module with:

```bash
go get go.osspkg.com/ui
```

This example builds a profile view with a card and individually positioned child nodes:

```go
package main

import "go.osspkg.com/ui/go-sdk/ui"

func newProfileApp() (*ui.App, error) {
    app := ui.New(
        ui.AppID("profile"),
        ui.AppTitle("Profile"),
        ui.AppVersion("1.0.0"),
    )
    view := ui.View(
        "profile.edit",
        ui.Title("Edit profile"),
        ui.Content(
            ui.HTML(ui.HTMLArticle, "profile-card").Row(1).Cols(8).Offset(2).Children(
                ui.HTML(ui.HTMLH1, "heading").Row(1).Cols(12).Prop("text", "Edit profile"),
                ui.HTML(ui.HTMLInput, "name").Row(2).Cols(8).Offset(2).Prop("placeholder", "Ada Lovelace"),
            ),
        ),
    )
    if err := app.AddView(view); err != nil {
        return nil, err
    }
    return app, nil
}
```

To expose the manifest and views as MCP resources, register the app with the MCP adapter:

```go
import (
    goMCP "go.osspkg.com/mcp"
    "go.osspkg.com/ui/go-sdk/ui"
    uiMCP "go.osspkg.com/ui/go-sdk/mcp"
)

func registerUI(app *ui.App) error {
    server, err := goMCP.New(
        goMCP.ServerInfo{Name: "profile", Version: "1.0.0"},
        uiMCP.Capabilities(),
    )
    if err != nil {
        return err
    }
    return uiMCP.Register(server, app)
}
```

The adapter publishes resources; a host that uses the React runtime must separately handle `ui.get`, `data.call`, and `ui.action`. See the [dashboard example](example/go/dashboard/dashboard.go) for sources, actions, and effects.

## React usage

Install the packages used by your application:

```bash
pnpm add @osspkg/ui-core @osspkg/ui-react
```

Register the built-in component sets, load the package stylesheet, and render a `ViewSchema`:

```tsx
import {
  createComponentRegistry,
  registerDefaultHTMLComponents,
  registerDefaultShadcnComponents,
  UIProvider,
  ViewRenderer,
} from "@osspkg/ui-react";
import type { ViewSchema } from "@osspkg/ui-core";
import "@osspkg/ui-react/styles.css";

const components = createComponentRegistry();
registerDefaultHTMLComponents(components);
registerDefaultShadcnComponents(components);

const schema: ViewSchema = {
  protocolVersion: "1.0",
  id: "welcome",
  title: "Welcome",
  regions: {
    "top-header": [],
    "left-panel": [],
    "right-panel": [],
    bottom: [],
    content: [
      {
        id: "welcome-card",
        component: "card",
        layout: { row: 1, cols: 8, offset: 2 },
        children: [
          { id: "title", component: "h1", layout: { row: 1, cols: 12 }, props: { text: "Welcome" } },
          { id: "body", component: "p", layout: { row: 2, cols: 8, offset: 2 }, props: { text: "Rendered from a schema." } },
        ],
      },
    ],
  },
};

export function WelcomeView() {
  return (
    <UIProvider components={components}>
      <ViewRenderer schema={schema} />
    </UIProvider>
  );
}
```

Nodes with children create their own 12-column grid. Set a child's `layout` to control its row, width (`cols`), and horizontal offset; nodes without a layout use the available row. The [interactive example](example/react/README.md) demonstrates forms, default registries, and nested layouts.

## Development

Install the JavaScript workspace dependencies and build the publishable React packages:

```bash
make npm-install
make npm-build
```

Run the interactive example with:

```bash
make examples-dev
```

Check and format the React SDK and example sources with:

```bash
make npm-lint
make npm-format-check
make npm-format
```

`make npm-lint` applies the Airbnb JavaScript/React rules with TypeScript and React Hooks compatibility adjustments. `make npm-format` applies Prettier; use `make npm-format-check` in CI or before committing to verify formatting without changing files. These targets are separate from `make lint`, which runs the Go linter.

The Makefile also provides `make tests`, `make lint`, `make build`, and `make ci` for Go project workflows. `make ci` runs the `pre-commit` prerequisites, including dependency/tool setup and generated license work; consult the [Makefile](Makefile) before running it.

## Packages

The JavaScript workspace contains three publishable packages, currently version `0.1.0`:

- [`@osspkg/ui-core`](react-sdk/ui-core/package.json)
- [`@osspkg/ui-transport`](react-sdk/ui-transport/package.json)
- [`@osspkg/ui-react`](react-sdk/ui-react/package.json)

`make npm-publish` builds and publishes these packages to the public npm registry. Run it only as part of an intentional release.

## License

This project is licensed under the [BSD 3-Clause License](LICENSE).
