import { regionNames, type Manifest, type RuntimeError, type ViewSchema } from "./schema.js";

export interface ValidationLimits {
  maxBytes: number;
  maxNodes: number;
  maxDepth: number;
  maxPropsPerNode: number;
  maxActions: number;
  maxSources: number;
  maxStrings: number;
  maxStringLength: number;
  maxExpressionDepth: number;
  maxArrayItems: number;
  maxObjectKeys: number;
  maxPathLength: number;
}

export const defaultLimits: ValidationLimits = {
  maxBytes: 1 << 20,
  maxNodes: 1000,
  maxDepth: 32,
  maxPropsPerNode: 64,
  maxActions: 128,
  maxSources: 64,
  maxStrings: 10000,
  maxStringLength: 16 * 1024,
  maxExpressionDepth: 32,
  maxArrayItems: 1000,
  maxObjectKeys: 256,
  maxPathLength: 1024,
};

export class SchemaValidationError extends Error {
  readonly runtimeError: RuntimeError;

  constructor(message: string, details?: unknown) {
    super(message);
    this.name = "SchemaValidationError";
    this.runtimeError = { code: "INVALID_UI_SCHEMA", message, ...(details === undefined ? {} : { details }) };
  }
}

export function parseManifest(input: unknown, limits: Partial<ValidationLimits> = {}): Manifest {
  if (!isRecord(input)) throw new SchemaValidationError("manifest must be an object");
  const manifest = input as Manifest;
  validateManifest(manifest, limits);
  return manifest;
}

export function parseView(input: unknown, limits: Partial<ValidationLimits> = {}): ViewSchema {
  if (!isRecord(input)) throw new SchemaValidationError("view must be an object");
  const view = input as ViewSchema;
  validateView(view, limits);
  return view;
}

export function validateManifest(manifest: Manifest, customLimits: Partial<ValidationLimits> = {}): void {
  const limits = withDefaults(customLimits);
  if (!isRecord(manifest) || manifest.protocolVersion !== "1.0") throw new SchemaValidationError("unsupported manifest protocol version");
  if (!isRecord(manifest.plugin) || !isSafeName(manifest.plugin.id) || !isText(manifest.plugin.title) || !isText(manifest.plugin.version)) {
    throw new SchemaValidationError("manifest plugin metadata is required");
  }
  if (!Array.isArray(manifest.views)) throw new SchemaValidationError("manifest views must be an array");
  const ids = new Set<string>();
  for (const view of manifest.views) {
    if (!isRecord(view) || !isSafeName(view.id) || !isText(view.title) || !isText(view.schema) || ids.has(view.id)) {
      throw new SchemaValidationError("manifest view is invalid");
    }
    ids.add(view.id);
  }
  if (manifest.requiredComponents !== undefined && (!Array.isArray(manifest.requiredComponents) || manifest.requiredComponents.some((name) => !isComponentName(name)))) {
    throw new SchemaValidationError("manifest required components are invalid");
  }
  validateSerializedLimits(manifest, limits);
}

export function validateView(view: ViewSchema, customLimits: Partial<ValidationLimits> = {}): void {
  const limits = withDefaults(customLimits);
  if (!isRecord(view) || view.protocolVersion !== "1.0") throw new SchemaValidationError("unsupported view protocol version");
  if (!isSafeName(view.id)) throw new SchemaValidationError("view id is required");
  if (view.title !== undefined && !isText(view.title)) throw new SchemaValidationError("view title is invalid");
  if (view.sources !== undefined && !isRecord(view.sources)) throw new SchemaValidationError("sources must be an object");
  if (view.actions !== undefined && !isRecord(view.actions)) throw new SchemaValidationError("actions must be an object");
  if (!isRecord(view.regions)) throw new SchemaValidationError("regions are required");
  validateSerializedLimits(view, limits);

  const actions = view.actions ?? {};
  const sources = view.sources ?? {};
  if (Object.keys(actions).length > limits.maxActions || Object.keys(sources).length > limits.maxSources) {
    throw new SchemaValidationError("source or action limit exceeded");
  }
  for (const [name, source] of Object.entries(sources)) validateSource(name, source, view, limits);
  for (const [name, action] of Object.entries(actions)) validateAction(name, action, view, limits);

  let nodes = 0;
  const nodeIDs = new Set<string>();
  for (const region of regionNames) {
    const nodesInRegion = view.regions[region];
    if (!Array.isArray(nodesInRegion)) throw new SchemaValidationError(`region ${region} must be an array`);
    for (const node of nodesInRegion) nodes = validateNode(node, region, 1, view, limits, nodeIDs, nodes);
  }
}

function validateSource(name: string, source: unknown, view: ViewSchema, limits: ValidationLimits): void {
  if (!isSafeName(name) || !isRecord(source) || source.type !== "tool" || !isSafeName(source.tool)) throw new SchemaValidationError(`source ${name} is invalid`);
  if (source.policy !== undefined && source.policy !== "manual" && source.policy !== "on-mount") throw new SchemaValidationError(`source ${name} policy is invalid`);
  if (source.cache !== undefined && (!isRecord(source.cache) || (source.cache.ttl !== undefined && (!Number.isInteger(source.cache.ttl) || source.cache.ttl < 0)))) throw new SchemaValidationError(`source ${name} cache is invalid`);
  validateValues(source.input, view, limits);
  if (source.refreshOn !== undefined && (!Array.isArray(source.refreshOn) || source.refreshOn.some((action) => !view.actions?.[action]))) throw new SchemaValidationError(`source ${name} refresh action is invalid`);
}

function validateAction(name: string, action: unknown, view: ViewSchema, limits: ValidationLimits): void {
  if (!isSafeName(name) || !isRecord(action) || !["tool", "set-state", "merge-state", "refresh-source", "invalidate"].includes(action.type)) throw new SchemaValidationError(`action ${name} is invalid`);
  if (action.type === "tool" && !isSafeName(action.tool)) throw new SchemaValidationError(`tool action ${name} requires a tool`);
  if ((action.type === "set-state" || action.type === "merge-state") && !isSafePath(action.path, limits.maxPathLength)) throw new SchemaValidationError(`action ${name} requires a safe path`);
  if (action.path !== undefined && action.path !== "" && !isSafePath(action.path, limits.maxPathLength)) throw new SchemaValidationError(`action ${name} path is invalid`);
  if ((action.type === "refresh-source" || action.type === "invalidate") && !view.sources?.[action.source as string]) throw new SchemaValidationError(`action ${name} source is invalid`);
  validateValues(action.input, view, limits);
  if (action.value !== undefined) validateValue(action.value, view, limits, 0);
  if (action.effects !== undefined) {
    if (!Array.isArray(action.effects)) throw new SchemaValidationError(`action ${name} effects are invalid`);
    for (const effect of action.effects) validateEffect(effect, view, limits);
  }
}

function validateEffect(effect: unknown, view: ViewSchema, limits: ValidationLimits): void {
  if (!isRecord(effect) || !["set-state", "merge-state", "invalidate", "refresh-source", "refresh-view", "patch-view", "navigate", "toast", "dialog", "close-dialog"].includes(effect.type)) throw new SchemaValidationError("effect type is invalid");
  if ((effect.type === "set-state" || effect.type === "merge-state") && !isSafePath(effect.path, limits.maxPathLength)) throw new SchemaValidationError("state effect path is invalid");
  if (effect.path !== undefined && effect.path !== "" && !isSafePath(effect.path, limits.maxPathLength)) throw new SchemaValidationError("effect path is invalid");
  if ((effect.type === "invalidate" || effect.type === "refresh-source") && !view.sources?.[effect.source as string]) throw new SchemaValidationError("effect source is invalid");
  if (effect.type === "toast" && !isText(effect.message)) throw new SchemaValidationError("toast effect message is required");
  if (effect.type === "navigate" && !isText(effect.to)) throw new SchemaValidationError("navigate effect target is required");
  if (effect.value !== undefined) validateValue(effect.value, view, limits, 0);
}

function validateNode(node: unknown, region: string, depth: number, view: ViewSchema, limits: ValidationLimits, ids: Set<string>, count: number): number {
  if (!isRecord(node) || !isSafeName(node.id) || !isComponentName(node.component) || ids.has(node.id)) throw new SchemaValidationError("node id or component is invalid");
  count += 1;
  if (count > limits.maxNodes || depth > limits.maxDepth) throw new SchemaValidationError("schema complexity limit exceeded");
  ids.add(node.id);
  if (node.layout !== undefined) validateLayout(node.layout, region);
  if (node.props !== undefined && (!isRecord(node.props) || Object.keys(node.props).length > limits.maxPropsPerNode)) throw new SchemaValidationError("node props are invalid");
  validateValues(node.props, view, limits);
  if (node.events !== undefined) {
    if (!isRecord(node.events)) throw new SchemaValidationError("node events are invalid");
    for (const [event, handler] of Object.entries(node.events)) {
      if (!isSafeName(event) || !isRecord(handler) || (!handler.action && (!Array.isArray(handler.steps) || handler.steps.length === 0))) throw new SchemaValidationError("event handler is invalid");
      if (handler.action && !view.actions?.[handler.action]) throw new SchemaValidationError("event action is not defined");
      if (handler.steps) for (const step of handler.steps) {
        if (!isRecord(step) || step.type !== "action" || !isSafeName(step.action) || !view.actions?.[step.action]) throw new SchemaValidationError("event step is invalid");
        validateValues(step.input, view, limits);
      }
    }
  }
  if (node.when !== undefined) validateValue(node.when, view, limits, 0);
  if (node.children !== undefined) {
    if (!Array.isArray(node.children)) throw new SchemaValidationError("children must be an array");
    for (const child of node.children) count = validateNode(child, region, depth + 1, view, limits, ids, count);
  }
  if (node.slots !== undefined) {
    if (!isRecord(node.slots)) throw new SchemaValidationError("slots must be an object");
    for (const children of Object.values(node.slots)) {
      if (!Array.isArray(children)) throw new SchemaValidationError("slot children must be an array");
      for (const child of children) count = validateNode(child, region, depth + 1, view, limits, ids, count);
    }
  }
  return count;
}

function validateLayout(layout: unknown, region: string): void {
  if (!isRecord(layout) || !Number.isInteger(layout.row) || layout.row < 1) throw new SchemaValidationError("layout row is invalid");
  if (region === "left-panel" || region === "right-panel") {
    if (layout.cols !== undefined || layout.offset !== undefined) throw new SchemaValidationError("side panel layout is invalid");
    return;
  }
  const cols = layout.cols ?? 0;
  const offset = layout.offset ?? 0;
  if (!Number.isInteger(cols) || cols < 1 || cols > 12 || !Number.isInteger(offset) || offset < 0 || offset > 11 || cols + offset > 12) throw new SchemaValidationError("grid layout is invalid");
}

function validateValues(values: unknown, view: ViewSchema, limits: ValidationLimits): void {
  if (values === undefined) return;
  if (!isRecord(values)) throw new SchemaValidationError("values must be an object");
  for (const value of Object.values(values)) validateValue(value, view, limits, 0);
}

function validateValue(value: unknown, view: ViewSchema, limits: ValidationLimits, depth: number): void {
  if (depth > limits.maxExpressionDepth) throw new SchemaValidationError("expression depth exceeded");
  if (value === null || typeof value === "boolean" || typeof value === "number") return;
  if (typeof value === "string") {
    if (value.length > limits.maxStringLength) throw new SchemaValidationError("value is too long");
    return;
  }
  if (Array.isArray(value)) {
    if (value.length > limits.maxArrayItems) throw new SchemaValidationError("array is too large");
    for (const item of value) validateValue(item, view, limits, depth + 1);
    return;
  }
  if (!isRecord(value)) throw new SchemaValidationError("value is not JSON-compatible");
  for (const key of Object.keys(value)) if (["__proto__", "prototype", "constructor"].includes(key)) throw new SchemaValidationError("unsafe value key");
  const entries = Object.entries(value);
  if (entries.length > limits.maxObjectKeys) throw new SchemaValidationError("expression object is too large");
  if (entries.length === 1) {
    const [key, raw] = entries[0]!;
    if (["$state", "$context", "$event", "$result"].includes(key)) {
      if (!isSafePath(raw, limits.maxPathLength)) throw new SchemaValidationError("reference path is unsafe");
      return;
    }
    if (key === "$source") {
      if (!isSafePath(raw, limits.maxPathLength)) throw new SchemaValidationError("source path is unsafe");
      const sourceName = sourceForPath(raw, view.sources);
      if (!sourceName) throw new SchemaValidationError("source reference is not defined");
      return;
    }
    const arity = expressionArity[key];
    if (arity) {
      if (!Array.isArray(raw) || raw.length < arity[0] || (arity[1] >= 0 && raw.length > arity[1])) throw new SchemaValidationError(`operator ${key} has invalid arity`);
      for (const item of raw) validateValue(item, view, limits, depth + 1);
      return;
    }
    if (key.startsWith("$")) throw new SchemaValidationError(`unsupported expression operator ${key}`);
  }
  for (const nested of Object.values(value)) validateValue(nested, view, limits, depth + 1);
}

const expressionArity: Record<string, [number, number]> = {
  $eq: [2, 2], $ne: [2, 2], $gt: [2, 2], $gte: [2, 2], $lt: [2, 2], $lte: [2, 2],
  $and: [1, -1], $or: [1, -1], $not: [1, 1], $in: [2, 2], $exists: [1, 1],
  $concat: [0, -1], $coalesce: [1, -1], $if: [3, 3],
};

function validateSerializedLimits(value: unknown, limits: ValidationLimits): void {
  let encoded: string;
  try {
    encoded = JSON.stringify(value);
  } catch {
    throw new SchemaValidationError("schema is not JSON serializable");
  }
  if (new TextEncoder().encode(encoded).length > limits.maxBytes) throw new SchemaValidationError("schema byte limit exceeded");
  let strings = 0;
  walkSerialized(value, 0, limits, () => { strings += 1; if (strings > limits.maxStrings) throw new SchemaValidationError("string count limit exceeded"); });
}

function walkSerialized(value: unknown, depth: number, limits: ValidationLimits, onString: () => void): void {
  if (depth > limits.maxDepth) throw new SchemaValidationError("schema depth limit exceeded");
  if (typeof value === "string") {
    if (value.length > limits.maxStringLength) throw new SchemaValidationError("string length limit exceeded");
    onString();
  } else if (Array.isArray(value)) {
    if (value.length > limits.maxArrayItems) throw new SchemaValidationError("array item limit exceeded");
    for (const item of value) walkSerialized(item, depth + 1, limits, onString);
  } else if (isRecord(value)) {
    if (Object.keys(value).length > limits.maxObjectKeys) throw new SchemaValidationError("object key limit exceeded");
    for (const [key, nested] of Object.entries(value)) {
      if (["__proto__", "prototype", "constructor"].includes(key)) throw new SchemaValidationError("unsafe value key");
      walkSerialized(key, depth + 1, limits, onString);
      walkSerialized(nested, depth + 1, limits, onString);
    }
  }
}

function withDefaults(custom: Partial<ValidationLimits>): ValidationLimits {
  return { ...defaultLimits, ...custom };
}

function isRecord(value: unknown): value is Record<string, any> {
  return !!value && typeof value === "object" && !Array.isArray(value);
}

function isText(value: unknown): value is string {
  return typeof value === "string";
}

function isSafeName(value: unknown): value is string {
  return typeof value === "string" && value.length > 0 && !/[\r\n]/.test(value);
}

function isSafePath(value: unknown, maxLength: number): value is string {
  return typeof value === "string" && value.length > 0 && value.length <= maxLength && value.split(".").every((segment) => segment.length > 0 && !["__proto__", "prototype", "constructor"].includes(segment));
}

function sourceForPath(path: string, sources: Record<string, unknown> | undefined): string | undefined {
  for (let name = path; name; ) {
    if (sources && Object.prototype.hasOwnProperty.call(sources, name)) return name;
    const index = name.lastIndexOf(".");
    if (index < 0) break;
    name = name.slice(0, index);
  }
  return undefined;
}

function isComponentName(value: unknown): value is string {
  return typeof value === "string" && /^[a-z][a-z0-9-]*$/.test(value);
}
