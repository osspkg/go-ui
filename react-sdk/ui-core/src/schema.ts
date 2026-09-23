export const protocolVersion = "1.0" as const;

export type RegionName = "top-header" | "left-panel" | "content" | "right-panel" | "bottom";
export type SourcePolicy = "manual" | "on-mount";
export type ValueReference = "$state" | "$context" | "$source" | "$event" | "$result";
export type ExpressionOperator = "$eq" | "$ne" | "$gt" | "$gte" | "$lt" | "$lte" | "$and" | "$or" | "$not" | "$in" | "$exists" | "$concat" | "$coalesce" | "$if";
export type ActionType = "tool" | "set-state" | "merge-state" | "refresh-source" | "invalidate";
export type EffectType = "set-state" | "merge-state" | "invalidate" | "refresh-source" | "refresh-view" | "patch-view" | "navigate" | "toast" | "dialog" | "close-dialog";

export interface Manifest {
  protocolVersion: string;
  plugin: { id: string; title: string; version: string; description?: string };
  views: ViewDescriptor[];
  requiredComponents?: string[];
}

export interface ViewDescriptor {
  id: string;
  title: string;
  route?: string;
  schema: string;
}

export interface ViewSchema {
  protocolVersion: string;
  id: string;
  revision?: string;
  title?: string;
  state?: Record<string, unknown>;
  sources?: Record<string, DataSource>;
  actions?: Record<string, ActionDefinition>;
  regions: Record<RegionName, UINode[]>;
}

export interface GridLayout {
  row: number;
  cols?: number;
  offset?: number;
}

export interface UINode {
  id: string;
  component: string;
  layout?: GridLayout;
  props?: Record<string, UIValue>;
  events?: Record<string, EventHandler>;
  children?: UINode[];
  slots?: Record<string, UINode[]>;
  when?: UIValue;
}

export interface DataSource {
  type: "tool";
  tool: string;
  input?: Record<string, UIValue>;
  policy?: SourcePolicy;
  cache?: { ttl?: number };
  refreshOn?: string[];
}

export interface EventHandler {
  action?: string;
  steps?: ActionStep[];
}

export interface ActionStep {
  type: "action";
  action?: string;
  input?: Record<string, UIValue>;
}

export interface ActionDefinition {
  type: ActionType;
  tool?: string;
  input?: Record<string, UIValue>;
  path?: string;
  value?: UIValue;
  source?: string;
  to?: string;
  variant?: string;
  message?: string;
  effects?: Effect[];
}

export interface Effect {
  type: EffectType;
  path?: string;
  source?: string;
  value?: UIValue;
  baseRevision?: string;
  revision?: string;
  to?: string;
  variant?: string;
  message?: string;
}

export type JSONValue = string | number | boolean | null | readonly JSONValue[] | { [key: string]: JSONValue };
export type UIReference = { $state: string } | { $context: string } | { $source: string } | { $event: string } | { $result: string };
export type UIExpression = { [operator in ExpressionOperator]?: readonly UIValue[] };
export type UIValue = JSONValue | UIReference | UIExpression;

export interface ResolveScope {
  state: unknown;
  context: unknown;
  sources: Record<string, SourceRuntimeState>;
  event?: unknown;
  result?: unknown;
}

export interface SourceRuntimeState<T = unknown> {
  status: "idle" | "loading" | "success" | "error";
  data?: T;
  error?: RuntimeError;
  updatedAt?: number;
  requestId?: string;
}

export interface RuntimeError {
  code: string;
  message: string;
  details?: unknown;
  retryable?: boolean;
}

export const regionNames: readonly RegionName[] = ["top-header", "left-panel", "content", "right-panel", "bottom"];
