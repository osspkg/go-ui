import { describe, expect, it } from "vitest";
import { createUsersDemoRegistry } from "./demo.js";

describe("users demo registry", () => {
  it("contains host-owned CRUD primitives", () => {
    const registry = createUsersDemoRegistry();
    expect(registry.has("table")).toBe(true);
    expect(registry.has("form")).toBe(true);
    expect(registry.has("text-field")).toBe(true);
    expect(registry.get("table")?.allowedProps).toEqual(["rows", "loading", "error"]);
  });
});
