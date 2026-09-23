import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import type { ViewSchema } from "@osspkg/ui-core";
import { createComponentRegistry } from "./registry.js";
import { htmlTagNames, registerDefaultHTMLComponents } from "./html.js";
import { registerDefaultShadcnComponents, shadcnComponentNames } from "./shadcn.js";
import { UIProvider } from "./runtime.js";
import { ViewRenderer } from "./renderer.js";

describe("default component registrations", () => {
  it("registers every HTML tag and shadcn component without overwriting a host component", () => {
    const registry = createComponentRegistry();
    const hostButton = () => null;
    registry.register("button", { component: hostButton });

    registerDefaultHTMLComponents(registry);
    registerDefaultShadcnComponents(registry);

    for (const name of htmlTagNames) expect(registry.has(name)).toBe(true);
    for (const name of shadcnComponentNames) expect(registry.has(name)).toBe(true);
    expect(registry.get("button")?.component).toBe(hostButton);
    expect(registry.get("a")?.events).toContain("click");
    expect(registry.get("data-table")?.allowedProps).toContain("rows");
  });

  it("renders safe native elements and the baseline data table", () => {
    const registry = createComponentRegistry();
    registerDefaultHTMLComponents(registry);
    registerDefaultShadcnComponents(registry);
    const schema: ViewSchema = {
      protocolVersion: "1.0",
      id: "defaults",
      regions: {
        "top-header": [], "left-panel": [], "right-panel": [], bottom: [],
        content: [
          { id: "link", component: "a", props: { href: "https://example.com", text: "Example" } },
          { id: "profile-form", component: "form", props: { method: "post", enctype: "multipart/form-data", novalidate: true, autocomplete: "off" } },
          { id: "email-label", component: "label", props: { for: "email", text: "Email" } },
          { id: "email", component: "input", props: { id: "email", type: "email", autocomplete: "email", inputmode: "email", enterkeyhint: "next", maxlength: 255, spellcheck: false } },
          { id: "notes", component: "textarea", props: { dirname: "notes.dir", wrap: "soft" } },
          { id: "submit", component: "button", props: { type: "submit", text: "Save" } },
          { id: "users", component: "data-table", props: { rows: [{ name: "Ada" }] } },
        ],
      },
    };

    const markup = renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>);
    expect(markup).toContain('href="https://example.com"');
    expect(markup).toContain('method="post"');
    expect(markup).toContain('encType="multipart/form-data"');
    expect(markup).toContain('for="email"');
    expect(markup).toContain('autoComplete="email"');
    expect(markup).toContain('inputMode="email"');
    expect(markup).toContain('enterKeyHint="next"');
    expect(markup).toContain('maxLength="255"');
    expect(markup).toContain('dirname="notes.dir"');
    expect(markup).toContain('wrap="soft"');
    expect(markup).toContain('type="submit"');
    expect(markup).toContain("Ada");

    schema.regions.content[0]!.props = { href: "javascript:alert(1)", text: "unsafe" };
    expect(() => renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>)).toThrow("href must be a safe URL");

    schema.regions.content[0]!.component = "form";
    schema.regions.content[0]!.props = { action: "https://attacker.example" };
    expect(() => renderToStaticMarkup(<UIProvider components={registry}><ViewRenderer schema={schema} /></UIProvider>)).toThrow("prop action is not allowed");
  });
});
