import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import type { ViewSchema } from "@osspkg/ui-core";
import { createComponentRegistry } from "./registry.js";
import { ViewRenderer } from "./renderer.js";
import { UIProvider, useUISource, type UISource } from "./runtime.js";
import type { RPCTransport } from "@osspkg/ui-transport";

describe("view runtime bindings", () => {
  it("resolves state and context values before rendering and binds events", () => {
    let received: Record<string, unknown> | undefined;
    const registry = createComponentRegistry();
    registry.register("probe", {
      component: (props: Record<string, unknown>) => {
        received = props;
        return <output>{String(props.name)}:{String(props.workspace)}</output>;
      },
      allowedProps: ["name", "workspace"],
      events: ["click"],
    });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "bindings",
      state: { name: "Ada" },
      actions: { submit: { type: "set-state", path: "name", value: "Grace" } },
      regions: {
        "top-header": [],
        "left-panel": [],
        content: [{
          id: "probe",
          component: "probe",
          props: { name: { $state: "name" }, workspace: { $context: "workspace.id" } },
          events: { click: { action: "submit" } },
        }],
        "right-panel": [],
        bottom: [],
      },
    };

    const markup = renderToStaticMarkup(
      <UIProvider components={registry}>
        <ViewRenderer schema={schema} context={{ workspace: { id: "w-1" } }} />
      </UIProvider>,
    );

    expect(markup).toContain("Ada:w-1");
    expect(received?.name).toBe("Ada");
    expect(received?.workspace).toBe("w-1");
    expect(received?.onClick).toEqual(expect.any(Function));
  });

  it("does not render nodes whose when condition is false", () => {
    const registry = createComponentRegistry();
    registry.register("probe", { component: () => <span>visible</span> });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "conditions",
      state: { enabled: false },
      regions: {
        "top-header": [],
        "left-panel": [],
        content: [{ id: "conditional", component: "probe", when: { $state: "enabled" } }],
        "right-panel": [],
        bottom: [],
      },
    };

    const markup = renderToStaticMarkup(
      <UIProvider components={registry}>
        <ViewRenderer schema={schema} />
      </UIProvider>,
    );

    expect(markup).not.toContain("visible");
  });

  it("exposes idle source state and manual refresh", () => {
    let source: UISource | undefined;
    const registry = createComponentRegistry();
    registry.register("source-probe", {
      component: () => {
        source = useUISource("users");
        return <output>{source.status}</output>;
      },
    });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "source",
      sources: { users: { type: "tool", tool: "users.list", policy: "manual" } },
      regions: {
        "top-header": [],
        "left-panel": [],
        content: [{ id: "source-probe", component: "source-probe" }],
        "right-panel": [],
        bottom: [],
      },
    };

    const markup = renderToStaticMarkup(
      <UIProvider components={registry}>
        <ViewRenderer schema={schema} />
      </UIProvider>,
    );

    expect(markup).toContain("idle");
    expect(source?.refresh).toEqual(expect.any(Function));
  });

  it("executes server actions and host-controlled effects", async () => {
    let received: Record<string, unknown> | undefined;
    const calls: Array<{ method: string; params: unknown }> = [];
    const transport: RPCTransport = {
      async call<T>(method: string, params: unknown): Promise<T> {
        calls.push({ method, params });
        return { id: "u-1" } as T;
      },
      subscribe: () => () => undefined,
      close: () => undefined,
    };
    const toasts: string[] = [];
    const refreshedViews: string[] = [];
    const patchedViews: string[] = [];
    const registry = createComponentRegistry();
    registry.register("button", {
      component: (props: Record<string, unknown>) => {
        received = props;
        return <button type="button">save</button>;
      },
      events: ["click"],
    });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "actions",
      actions: {
        save: {
          type: "tool",
          tool: "users.save",
          input: { name: { $event: "value" } },
          effects: [{ type: "toast", variant: "success", message: "Saved" }, { type: "refresh-view", revision: "rev-2" }, { type: "patch-view", baseRevision: "rev-1", revision: "rev-2" }],
        },
      },
      sources: { users: { type: "tool", tool: "users.list", refreshOn: ["save"] } },
      regions: {
        "top-header": [],
        "left-panel": [],
        content: [{ id: "save", component: "button", events: { click: { action: "save" } } }],
        "right-panel": [],
        bottom: [],
      },
    };

    renderToStaticMarkup(
      <UIProvider components={registry} transport={transport} effects={{ toast: (message) => { toasts.push(message); }, refreshView: (effect) => { refreshedViews.push(effect.type); }, patchView: (effect) => { patchedViews.push(effect.type); } }}>
        <ViewRenderer schema={schema} plugin="users" />
      </UIProvider>,
    );
    await (received?.onClick as (event: unknown) => Promise<void>)({ value: "Ada" });

    expect(calls).toEqual([
      { method: "ui.action", params: { plugin: "users", view: "actions", action: "save", input: { name: "Ada" } } },
      { method: "data.call", params: { plugin: "users", view: "actions", source: "users", operation: "users.list", input: {} } },
    ]);
    expect(toasts).toEqual(["Saved"]);
    expect(refreshedViews).toEqual(["refresh-view"]);
    expect(patchedViews).toEqual(["patch-view"]);
  });
});
