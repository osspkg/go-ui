import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import type { ReactNode } from "react";
import type { ViewSchema } from "@osspkg/ui-core";
import { createComponentRegistry } from "./registry.js";
import { UIProvider } from "./runtime.js";
import { ViewRenderer } from "./renderer.js";

function emptyRegions() {
  return { "top-header": [], "left-panel": [], content: [], "right-panel": [], bottom: [] };
}

describe("static renderer", () => {
  it("renders registered components, children, and grid metadata", () => {
    const registry = createComponentRegistry();
    registry.register("card", {
      component: (props: Record<string, unknown>) => <article>{String(props.title)}{props.children as ReactNode}</article>,
      allowedProps: ["title"],
    });
    registry.register("text", {
      component: (props: Record<string, unknown>) => <span>{String(props.value)}</span>,
      allowedProps: ["value"],
      children: false,
    });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "renderer",
      regions: {
        ...emptyRegions(),
        content: [{
          id: "card",
          component: "card",
          layout: { row: 1, cols: 8, offset: 2 },
          props: { title: "Users" },
          children: [{ id: "text", component: "text", props: { value: "Ada" } }],
        }],
      },
    };

    const markup = renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>);

    expect(markup).toContain("Users");
    expect(markup).toContain("Ada");
    expect(markup).toContain('data-ui-cols="8"');
    expect(markup).toContain('data-ui-offset="2"');
    expect(markup).toContain("grid grid-cols-12 gap-4");
    expect(markup).toContain("col-span-[var(--ui-cols)]");
    expect(markup).toContain("--ui-row:1;--ui-cols:8;--ui-start:3");
  });

  it("renders declared slots and rejects undeclared slots", () => {
    const registry = createComponentRegistry();
    registry.register("card", {
      component: (props: Record<string, unknown>) => <article>{(props.slots as Record<string, ReactNode[]>).footer}</article>,
      slots: ["footer"],
    });
    registry.register("text", { component: (props: Record<string, unknown>) => <span>{String(props.value)}</span>, allowedProps: ["value"] });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "slots",
      regions: { ...emptyRegions(), content: [{ id: "card", component: "card", slots: { footer: [{ id: "footer", component: "text", props: { value: "Footer" } }] } }] },
    };

    expect(renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>)).toContain("Footer");

    schema.regions.content[0]!.slots = { other: [] };
    expect(() => renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>)).toThrow("slot other is not allowed");
  });

  it("renders a safe fallback for unsupported components", () => {
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "unsupported",
      regions: { ...emptyRegions(), content: [{ id: "missing", component: "missing-component" }] },
    };

    const markup = renderToStaticMarkup(<UIProvider components={createComponentRegistry()}><ViewRenderer schema={schema} /></UIProvider>);

    expect(markup).toContain("Unsupported component: missing-component");
    expect(markup).toContain('data-unsupported-component="missing-component"');
  });

  it("rejects props without a host declaration", () => {
    const registry = createComponentRegistry();
    registry.register("text", { component: () => <span>text</span> });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "props",
      regions: { ...emptyRegions(), content: [{ id: "text", component: "text", props: { value: "unsafe" } }] },
    };

    expect(() => renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>)).toThrow("component props are not declared");
  });

  it("rejects events without a host allowlist", () => {
    const registry = createComponentRegistry();
    registry.register("button", { component: () => <button type="button">button</button> });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "events",
      actions: { save: { type: "set-state", path: "saved", value: true } },
      regions: { ...emptyRegions(), content: [{ id: "button", component: "button", events: { click: { action: "save" } } }] },
    };

    expect(() => renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>)).toThrow("component events are not declared");
  });

  it("rejects plugin styling and raw React props", () => {
    const registry = createComponentRegistry();
    registry.register("card", { component: () => <article />, allowedProps: ["style"] });
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "unsafe-props",
      regions: { ...emptyRegions(), content: [{ id: "card", component: "card", props: { style: { color: "red" } } }] },
    };

    expect(() => renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>)).toThrow("prop style is forbidden");
  });
});
