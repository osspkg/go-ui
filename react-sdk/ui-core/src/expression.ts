import { resolvePath } from "./path.js";
import type { ResolveScope, UIValue } from "./schema.js";

const operators = new Set(["$eq", "$ne", "$gt", "$gte", "$lt", "$lte", "$and", "$or", "$not", "$in", "$exists", "$concat", "$coalesce", "$if"]);

export interface ExpressionLimits {
  maxDepth: number;
  maxArrayItems: number;
  maxObjectKeys: number;
  maxStringLength: number;
}

export const defaultExpressionLimits: ExpressionLimits = {
  maxDepth: 32,
  maxArrayItems: 1000,
  maxObjectKeys: 256,
  maxStringLength: 16 * 1024,
};

export function resolveValue(value: UIValue, scope: ResolveScope, depth = 0, limits: ExpressionLimits = defaultExpressionLimits): unknown {
  if (depth > limits.maxDepth) throw new Error("expression depth exceeded");
  if (typeof value === "string" && value.length > limits.maxStringLength) throw new Error("expression string limit exceeded");
  if (Array.isArray(value)) {
    if (value.length > limits.maxArrayItems) throw new Error("expression array limit exceeded");
    return value.map((item) => resolveValue(item, scope, depth + 1, limits));
  }
  if (!value || typeof value !== "object") return value;
  const entries = Object.entries(value);
  if (entries.length > limits.maxObjectKeys) throw new Error("expression object limit exceeded");
  if (entries.length === 1) {
    const [key, raw] = entries[0]!;
    if (key === "$state" || key === "$context" || key === "$event" || key === "$result") {
      const root = key === "$state" ? scope.state : key === "$context" ? scope.context : key === "$event" ? scope.event : scope.result;
      return resolvePath(root, String(raw));
    }
    if (key === "$source") return resolveSource(String(raw), scope);
    if (operators.has(key)) return evaluateExpression(key, raw, scope, depth + 1, limits);
  }
  return Object.fromEntries(entries.map(([key, nested]) => [key, resolveValue(nested, scope, depth + 1, limits)]));
}

function resolveSource(path: string, scope: ResolveScope): unknown {
  const separator = path.indexOf(".");
  const name = separator < 0 ? path : path.slice(0, separator);
  const rest = separator < 0 ? "" : path.slice(separator + 1);
  const source = scope.sources[name];
  if (!source) return undefined;
  if (rest === "$loading") return source.status === "loading";
  if (rest === "$error") return source.error;
  if (rest === "$data" || rest === "") return source.data;
  if (rest === "$updatedAt") return source.updatedAt;
  return resolvePath(source.data, rest);
}

function evaluateExpression(operator: string, raw: unknown, scope: ResolveScope, depth: number, limits: ExpressionLimits): unknown {
  if (!Array.isArray(raw)) throw new Error(`operator ${operator} requires an array`);
  if (raw.length > limits.maxArrayItems) throw new Error("expression array limit exceeded");
  const args = raw.map((item) => resolveValue(item, scope, depth, limits));
  switch (operator) {
    case "$eq": return args[0] === args[1];
    case "$ne": return args[0] !== args[1];
    case "$gt": return comparable(args[0]) > comparable(args[1]);
    case "$gte": return comparable(args[0]) >= comparable(args[1]);
    case "$lt": return comparable(args[0]) < comparable(args[1]);
    case "$lte": return comparable(args[0]) <= comparable(args[1]);
    case "$and": return args.every(Boolean);
    case "$or": return args.some(Boolean);
    case "$not": return !args[0];
    case "$in": return Array.isArray(args[1]) && args[1].includes(args[0]);
    case "$exists": return args[0] !== undefined && args[0] !== null;
    case "$concat": {
      const result = args.map((item) => String(item ?? "")).join("");
      if (result.length > limits.maxStringLength) throw new Error("expression string limit exceeded");
      return result;
    }
    case "$coalesce": return args.find((item) => item !== undefined && item !== null);
    case "$if": return args[0] ? args[1] : args[2];
    default: throw new Error(`unsupported expression operator ${operator}`);
  }
}

function comparable(value: unknown): string | number {
  if (typeof value === "string" || typeof value === "number") return value;
  throw new Error("comparison requires a string or number");
}
