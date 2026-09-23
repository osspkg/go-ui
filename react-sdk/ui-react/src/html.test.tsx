import { describe, expect, it } from "vitest";
import { createComponentRegistry } from "./registry.js";
import { htmlTagNames, isHTMLTagName, registerHTMLTags } from "./html.js";

describe("HTML tag catalog", () => {
  it("contains safe semantic, form, table, and media tags", () => {
    expect(htmlTagNames).toContain("article");
    expect(htmlTagNames).toContain("input");
    expect(htmlTagNames).toContain("table");
    expect(htmlTagNames).toContain("video");
    expect(isHTMLTagName("script")).toBe(false);
    expect(isHTMLTagName("iframe")).toBe(false);
  });

  it("registers host implementations only for catalog tags", () => {
    const registry = createComponentRegistry();
    registerHTMLTags(registry, { article: { component: () => null } });
    expect(registry.has("article")).toBe(true);
  });
});
