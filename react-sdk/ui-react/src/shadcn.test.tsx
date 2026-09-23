import { describe, expect, it } from "vitest";
import { createComponentRegistry } from "./registry.js";
import { isShadcnComponentName, registerShadcnComponents, shadcnComponentNames } from "./shadcn.js";

describe("shadcn component catalog", () => {
  it("contains the official top-level component names", () => {
    expect(shadcnComponentNames).toHaveLength(64);
    expect(isShadcnComponentName("alert-dialog")).toBe(true);
    expect(isShadcnComponentName("input-otp")).toBe(true);
    expect(isShadcnComponentName("not-a-component")).toBe(false);
  });

  it("registers only catalog components with host-owned definitions", () => {
    const registry = createComponentRegistry();
    registerShadcnComponents(registry, { button: { component: () => null, events: ["click"] } });
    expect(registry.has("button")).toBe(true);
  });
});
