import { describe, expect, it } from "vitest";
import { createComponentRegistry } from "./registry.js";

describe("ComponentRegistry", () => {
  it("registers and resolves host components", () => {
    const registry = createComponentRegistry();
    const component = () => null;
    registry.register("text", { component });
    expect(registry.has("text")).toBe(true);
    expect(registry.get("text")?.component).toBe(component);
  });

  it("rejects duplicate logical component names", () => {
    const registry = createComponentRegistry();
    registry.register("text", { component: () => null });
    expect(() => registry.register("text", { component: () => null })).toThrow("already registered");
  });
});
