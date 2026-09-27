# go-ui code examples

These examples are included here so the skill does not depend on documentation outside this directory. Confirm details against the current project API when using them after a protocol or SDK change.

## Compose and validate a Go view

```go
package profile

import "go.osspkg.com/ui/go-sdk/ui"

func NewApp() (*ui.App, error) {
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
                ui.HTML(ui.HTMLButton, "save").Row(3).Cols(4).Offset(4).Prop("text", "Save"),
            ),
        ),
    )
    if err := app.AddView(view); err != nil {
        return nil, err
    }
    return app, nil
}
```

`App.AddView` validates the schema and stores a copy. Node builders support `.Row(n)`, `.Cols(n)`, `.Offset(n)`, `.Prop(key, value)`, `.Children(...)`, and `.On(event, ui.ActionRef(name))`. View options include `ui.Title`, `ui.State`, `ui.Source`, `ui.Content`, `ui.TopHeader`, `ui.LeftPanel`, `ui.RightPanel`, and `ui.Bottom`.

A local state action can be declared as:

```go
ui.View("profile.edit").Action(
    "profile.submit",
    ui.SetState("submitted", true).Then(ui.ToastSuccess("Profile saved")),
)
```

For a tool source, use `ui.Source("metrics", ui.Tool("metrics.summary").OnMount().TTL(30))`. Use the typed builders rather than hand-constructing action JSON when the builder supports the operation.

## Register UI resources on an MCP server

```go
import (
    goMCP "go.osspkg.com/mcp"
    uiMCP "go.osspkg.com/ui/go-sdk/mcp"
)

server, err := goMCP.New(
    goMCP.ServerInfo{Name: "profile", Version: "1.0.0"},
    uiMCP.Capabilities(),
)
if err != nil {
    return err
}
if err := uiMCP.Register(server, app); err != nil {
    return err
}
```

Registration publishes the manifest and view resources. Implement `ui.get`, `data.call`, and `ui.action` in the host if the React runtime needs RPC calls.

## Render a static schema in React

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
  regions: {
    "top-header": [],
    "left-panel": [],
    "right-panel": [],
    bottom: [],
    content: [
      {
        id: "card",
        component: "card",
        layout: { row: 1, cols: 8, offset: 2 },
        children: [
          { id: "title", component: "h1", layout: { row: 1, cols: 12 }, props: { text: "Welcome" } },
          { id: "body", component: "p", layout: { row: 2, cols: 8, offset: 2 }, props: { text: "A schema rendered by React." } },
        ],
      },
    ],
  },
};

export function WelcomeView() {
  return <UIProvider components={components}><ViewRenderer schema={schema} /></UIProvider>;
}
```

Default component registries provide dependency-free HTML and shadcn-name renderers. A host can register its own implementation first; the default registration functions leave existing registry entries intact. Each custom component definition declares allowed props, events, slots, and whether it accepts children.

## Connect a host RPC transport

```tsx
import { PluginView, UIProvider, createComponentRegistry, registerDefaultHTMLComponents, registerDefaultShadcnComponents } from "@osspkg/ui-react";
import { HTTPTransport } from "@osspkg/ui-transport";

const components = createComponentRegistry();
registerDefaultHTMLComponents(components);
registerDefaultShadcnComponents(components);
const transport = new HTTPTransport("/rpc");

export function RemoteProfile() {
  return (
    <UIProvider components={components} transport={transport}>
      <PluginView plugin="profile" view="profile.edit" />
    </UIProvider>
  );
}
```

`HTTPTransport` sends JSON-RPC calls to its endpoint. `WebSocketTransport` supports calls, subscriptions, and optional reconnect handling. Both implement `RPCTransport`. Close long-lived transports when their owning application lifecycle ends.

The server must return `{ "schema": <ViewSchema> }` for `ui.get`, accept source requests for `data.call`, and handle declared tool actions for `ui.action`. React local actions such as `set-state` execute in the runtime; host calls are required for tool-backed behavior.

## Reference project patterns

For an HTML form, use component names such as `article`, `form`, `label`, `input`, and `button`, bind the form `submit` event to an action, and use `when` to conditionally show a result. For data dashboards, use a source with `on-mount`, bind displayed values to `$source` references, and refresh sources from action effects. In both cases, place children with their own layouts when they should not simply stack at default width.
