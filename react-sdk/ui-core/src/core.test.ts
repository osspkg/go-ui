import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";
import type { UIValue } from "./schema.js";
import { resolvePath } from "./path.js";
import { resolveValue, type ExpressionLimits } from "./expression.js";
import { createStateRuntime } from "./runtime.js";
import { createRequestTracker } from "./request.js";
import { parseManifest, parseView, validateView, SchemaValidationError } from "./validate.js";

function fixture(name: string): unknown {
  return JSON.parse(readFileSync(new URL(`../../../../fixtures/${name}`, import.meta.url), "utf8")) as unknown;
}

const scope = {
  state: { status: "active", name: "Ada" },
  context: { workspace: { id: "w-1" } },
  sources: { users: { status: "success" as const, data: { items: [{ id: 1 }] } } },
  event: { value: "draft" },
  result: { id: 42 },
};

describe("ui core", () => {
  it("rejects prototype traversal", () => {
    expect(() => resolvePath({}, "constructor.name")).toThrow("unsafe path");
  });

  it("resolves references and expressions", () => {
    expect(resolveValue({ $state: "name" }, scope)).toBe("Ada");
    expect(resolveValue({ $source: "users.items" }, scope)).toEqual([{ id: 1 }]);
    expect(resolveValue({ $eq: [{ $state: "status" }, "active"] }, scope)).toBe(true);
  });

  it.each([
    ["idle", false],
    ["loading", true],
    ["success", false],
    ["error", false],
  ] as const)("maps source loading marker for %s", (status, expected) => {
    const sourceScope = {
      ...scope,
      sources: {
        users: {
          status,
          data: { items: [{ id: 1 }] },
          ...(status === "error" ? { error: { code: "FAILED", message: "failed" } } : {}),
        },
      },
    };
    expect(resolveValue({ $source: "users.$loading" }, sourceScope)).toBe(expected);
  });

  it("maps source errors and data paths", () => {
    const error = { code: "FAILED", message: "failed" };
    const sourceScope = { ...scope, sources: { users: { status: "error" as const, error, data: { items: [] } } } };
    expect(resolveValue({ $source: "users.$error" }, sourceScope)).toEqual(error);
    expect(resolveValue({ $source: "users.items" }, sourceScope)).toEqual([]);
  });

  it("accepts only the latest request for a source", () => {
    const tracker = createRequestTracker();
    const first = tracker.next("users");
    const second = tracker.next("users");

    expect(tracker.isCurrent("users", first)).toBe(false);
    expect(tracker.isCurrent("users", second)).toBe(true);
    expect(tracker.isCurrent("orders", second)).toBe(false);
  });

  it("rejects invalid layouts", () => {
    expect(() => validateView({
      protocolVersion: "1.0",
      id: "broken",
      regions: { "top-header": [], "left-panel": [], content: [{ id: "n", component: "card", layout: { row: 1, cols: 8, offset: 5 } }], "right-panel": [], bottom: [] },
    })).toThrow(SchemaValidationError);
  });

  it("parses manifests and views only after validation", () => {
    expect(parseManifest({ protocolVersion: "1.0", plugin: { id: "users", title: "Users", version: "1.0" }, views: [] }).plugin.id).toBe("users");
    expect(parseView({
      protocolVersion: "1.0",
      id: "users.list",
      regions: { "top-header": [], "left-panel": [], content: [], "right-panel": [], bottom: [] },
    }).id).toBe("users.list");
    expect(() => parseView({ protocolVersion: "1.0", id: "broken", regions: { content: [] } })).toThrow(SchemaValidationError);
  });

  it("rejects references to unknown sources", () => {
    expect(() => validateView({
      protocolVersion: "1.0",
      id: "users.list",
      regions: { "top-header": [], "left-panel": [], content: [{ id: "table", component: "table", props: { rows: { $source: "missing.items" } } }], "right-panel": [], bottom: [] },
    })).toThrow(SchemaValidationError);
  });

  it("accepts and rejects the shared Go fixtures identically", () => {
    expect(parseView(fixture("valid-view.json")).id).toBe("users.list");
    expect(() => parseView(fixture("invalid-layout.json"))).toThrow(SchemaValidationError);
  });

  it("enforces expression collection, object, and string limits", () => {
    const limits: ExpressionLimits = { maxDepth: 4, maxArrayItems: 2, maxObjectKeys: 1, maxStringLength: 4 };
    expect(() => resolveValue([1, 2, 3], scope, 0, limits)).toThrow("expression array limit exceeded");
    expect(() => resolveValue({ first: 1, second: 2 }, scope, 0, limits)).toThrow("expression object limit exceeded");
    expect(() => resolveValue("12345", scope, 0, limits)).toThrow("expression string limit exceeded");
    expect(() => resolveValue({ $concat: ["123", "45"] }, scope, 0, limits)).toThrow("expression string limit exceeded");
  });

  it.each([
    ["$eq", { $eq: [1, 1] }, true],
    ["$ne", { $ne: [1, 2] }, true],
    ["$gt", { $gt: [2, 1] }, true],
    ["$gte", { $gte: [2, 2] }, true],
    ["$lt", { $lt: [1, 2] }, true],
    ["$lte", { $lte: [2, 2] }, true],
    ["$and", { $and: [true, true] }, true],
    ["$or", { $or: [false, true] }, true],
    ["$not", { $not: [false] }, true],
    ["$in", { $in: ["active", ["draft", "active"]] }, true],
    ["$exists", { $exists: [{ $state: "name" }] }, true],
    ["$concat", { $concat: ["Ada", " Lovelace"] }, "Ada Lovelace"],
    ["$coalesce", { $coalesce: [{ $state: "missing" }, "fallback"] }, "fallback"],
    ["$if", { $if: [true, "yes", "no"] }, "yes"],
  ] as const)("evaluates %s", (_operator, expression, expected) => {
    expect(resolveValue(expression as UIValue, scope)).toBe(expected);
  });

  it("keeps schema, path, and expression decoders bounded for hostile inputs", () => {
    let seed = 0x12345678;
    const next = (): number => {
      seed = (seed * 1664525 + 1013904223) >>> 0;
      return seed;
    };
    for (let index = 0; index < 256; index += 1) {
      const token = next().toString(16);
      expect(() => resolvePath({ value: token }, index % 2 === 0 ? `value.${token}` : "__proto__.x")).not.toThrow(TypeError);
      expect(() => resolveValue({ $concat: [token, { $state: "missing" }] }, scope)).not.toThrow(TypeError);
      expect(() => parseView(index % 2 === 0 ? { protocolVersion: "bad", id: token } : { [token]: next() })).not.toThrow(TypeError);
    }
  });

  it.each([
    ["set", (runtime: ReturnType<typeof createStateRuntime>) => runtime.set("form.name", "Grace")],
    ["merge", (runtime: ReturnType<typeof createStateRuntime>) => runtime.merge("form", { valid: true })],
  ])("supports local state action: %s", (_name, apply) => {
    const initial = { form: { name: "Ada", valid: false } };
    const runtime = createStateRuntime(initial);

    apply(runtime);

    expect(runtime.get()).not.toBe(initial);
    expect(initial).toEqual({ form: { name: "Ada", valid: false } });
    expect(runtime.get()).toEqual(_name === "set"
      ? { form: { name: "Grace", valid: false } }
      : { form: { name: "Ada", valid: true } });
  });
});
