import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState, type Dispatch, type PropsWithChildren, type ReactElement, type SetStateAction } from "react";
import { createRequestTracker, resolveValue, resolvePath, setPath, type Effect, type ResolveScope, type SourceRuntimeState, type UIValue, type ViewSchema } from "@osspkg/ui-core";
import type { RPCTransport } from "@osspkg/ui-transport";
import { ComponentRegistry } from "./registry.js";

export interface UIEffectHandlers {
  toast?(message: string, variant?: string): void | Promise<void>;
  navigate?(to: string): void | Promise<void>;
  dialog?(effect: Effect): void | Promise<void>;
  closeDialog?(effect: Effect): void | Promise<void>;
}

export interface UIRuntime {
  components: ComponentRegistry;
  transport?: RPCTransport;
  effects?: UIEffectHandlers;
}

export interface ViewRuntime extends UIRuntime {
  schema: ViewSchema;
  plugin?: string;
  context: unknown;
  state: unknown;
  sources: Record<string, SourceRuntimeState>;
  setState(path: string, value: unknown): void;
  mergeState(path: string, value: Record<string, unknown>): void;
  resolve(value: unknown, event?: unknown, result?: unknown): unknown;
  runAction(name: string, event?: unknown): Promise<unknown>;
  refreshSource(name: string): Promise<void>;
}

export interface UISource<T = unknown> extends SourceRuntimeState<T> {
  refresh(): Promise<void>;
}

const runtimeContext = createContext<UIRuntime | undefined>(undefined);
const viewRuntimeContext = createContext<ViewRuntime | undefined>(undefined);

export function UIProvider({ components, transport, effects, children }: PropsWithChildren<UIRuntime>): ReactElement {
  const value = useMemo(() => ({ components, ...(transport ? { transport } : {}), ...(effects ? { effects } : {}) }), [components, effects, transport]);
  return <runtimeContext.Provider value={value}>{children}</runtimeContext.Provider>;
}

export function useUIRuntime(): UIRuntime {
  const runtime = useContext(runtimeContext);
  if (!runtime) throw new Error("UIProvider is required");
  return runtime;
}

export function useUIViewRuntime(): ViewRuntime {
  const runtime = useContext(viewRuntimeContext);
  if (!runtime) throw new Error("ViewRuntime is required");
  return runtime;
}

export function useUISource<T = unknown>(name: string): UISource<T> {
  const runtime = useUIViewRuntime();
  const source = runtime.sources[name] as SourceRuntimeState<T> | undefined;
  const refresh = useCallback(() => runtime.refreshSource(name), [name, runtime]);
  return useMemo(() => ({ ...(source ?? { status: "idle" as const }), refresh }), [refresh, source]);
}

export function useUIAction(name: string): (event?: unknown) => Promise<unknown> {
  const runtime = useUIViewRuntime();
  return useCallback((event?: unknown) => runtime.runAction(name, event), [name, runtime]);
}

export function ViewRuntimeProvider({ schema, plugin, context, children }: PropsWithChildren<{ schema: ViewSchema; plugin?: string; context: unknown }>): ReactElement {
  const base = useUIRuntime();
  const [state, setState] = useState<unknown>(schema.state ?? {});
  const [sources, setSources] = useState<Record<string, SourceRuntimeState>>(() => initialSourceStates(schema));
  const requestIDs = useRef(createRequestTracker());
  const controllers = useRef(new Map<string, AbortController>());

  useEffect(() => {
    requestIDs.current = createRequestTracker();
    setState(schema.state ?? {});
    setSources(initialSourceStates(schema));
  }, [schema]);

  const refreshSource = useCallback(async (name: string): Promise<void> => {
    const source = schema.sources?.[name];
    if (!source || !base.transport) return;
    controllers.current.get(name)?.abort();
    const controller = new AbortController();
    controllers.current.set(name, controller);
    const requestID = requestIDs.current.next(name);
    setSources((current) => ({ ...current, [name]: { status: "loading", requestId: String(requestID) } }));
    try {
      const input = resolveValue(source.input ?? {}, { state, context, sources });
      const data = await base.transport.call("data.call", { plugin, view: schema.id, source: name, operation: source.tool, input }, { signal: controller.signal });
      if (!requestIDs.current.isCurrent(name, requestID)) return;
      setSources((current) => ({ ...current, [name]: { status: "success", data, updatedAt: Date.now(), requestId: String(requestID) } }));
    } catch (error) {
      if (controller.signal.aborted || !requestIDs.current.isCurrent(name, requestID)) return;
      setSources((current) => ({ ...current, [name]: { status: "error", error: { code: "SOURCE_FAILED", message: error instanceof Error ? error.message : "source failed", retryable: true }, requestId: String(requestID) } }));
    } finally {
      if (controllers.current.get(name) === controller) controllers.current.delete(name);
    }
  }, [base.transport, context, plugin, schema, sources, state]);

  useEffect(() => {
    for (const [name, source] of Object.entries(schema.sources ?? {})) {
      if (source.policy === "on-mount") void refreshSource(name);
    }
    return () => { for (const controller of controllers.current.values()) controller.abort(); };
    // on-mount execution belongs to one schema lifecycle; source updates must not retrigger it.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [schema, base.transport, plugin]);

  useEffect(() => {
    const onReconnect = base.transport?.onReconnect;
    if (!onReconnect) return;
    return onReconnect(() => {
      for (const [name, source] of Object.entries(schema.sources ?? {})) {
        if (source.policy === "on-mount") void refreshSource(name);
      }
    });
  }, [base.transport, refreshSource, schema]);

  useEffect(() => {
    const transport = base.transport;
    if (!transport) return;
    return transport.subscribe((notification) => {
      const params = isRecord(notification.params) ? notification.params : {};
      const notificationPlugin = params.plugin;
      const notificationView = params.view;
      if (typeof notificationPlugin === "string" && notificationPlugin !== plugin) return;
      if (typeof notificationView === "string" && notificationView !== schema.id) return;
      if (notification.method === "data.invalidate" && typeof params.source === "string") {
        void refreshSource(params.source);
        return;
      }
      if (notification.method === "ui.invalidate" || notification.method === "plugin.changed") {
        for (const [name, source] of Object.entries(schema.sources ?? {})) {
          if (source.policy === "on-mount") void refreshSource(name);
        }
      }
    });
  }, [base.transport, plugin, refreshSource, schema]);

  const runtime = useMemo<ViewRuntime>(() => ({
    ...base,
    schema,
    ...(plugin ? { plugin } : {}),
    context,
    state,
    sources,
    setState: (path, value) => setState((current: unknown) => setPath(current, path, value)),
    mergeState: (path, value) => setState((current: unknown) => {
      const currentValue = resolvePath(current, path);
      const merged = currentValue && typeof currentValue === "object" ? { ...currentValue, ...value } : value;
      return setPath(current, path, merged);
    }),
    resolve: (value, event, result) => resolveValue(value as UIValue, { state, context, sources, event, result }),
    runAction: (name, event) => executeAction(name, schema, base.transport, base.effects, plugin, state, context, sources, setState, refreshSource, event),
    refreshSource,
  }), [base, context, plugin, refreshSource, schema, sources, state]);
  return <viewRuntimeContext.Provider value={runtime}>{children}</viewRuntimeContext.Provider>;
}

function initialSourceStates(schema: ViewSchema): Record<string, SourceRuntimeState> {
  return Object.fromEntries(Object.keys(schema.sources ?? {}).map((name) => [name, { status: "idle" as const }]));
}

async function executeAction(name: string, schema: ViewSchema, transport: RPCTransport | undefined, handlers: UIEffectHandlers | undefined, plugin: string | undefined, state: unknown, context: unknown, sources: Record<string, SourceRuntimeState>, setState: Dispatch<SetStateAction<unknown>>, refreshSource: (name: string) => Promise<void>, event?: unknown): Promise<unknown> {
  const action = schema.actions?.[name];
  if (!action) throw new Error(`action ${name} is not defined`);
  const scope: ResolveScope = { state, context, sources, event };
  if (action.type === "set-state") {
    setState((current: unknown) => setPath(current, action.path ?? "", resolveValue(action.value ?? null, scope)));
    for (const effect of action.effects ?? []) await executeEffect(effect, undefined, scope, handlers, setState, refreshSource);
    await refreshTriggeredSources(name, schema, refreshSource);
    return undefined;
  }
  if (action.type === "merge-state") {
    const current = resolvePath(state, action.path ?? "");
    const next = resolveValue(action.value ?? null, scope);
    setState((value: unknown) => setPath(value, action.path ?? "", { ...(current && typeof current === "object" ? current : {}), ...(next && typeof next === "object" ? next : {}) }));
    for (const effect of action.effects ?? []) await executeEffect(effect, undefined, scope, handlers, setState, refreshSource);
    await refreshTriggeredSources(name, schema, refreshSource);
    return undefined;
  }
  if (action.type === "refresh-source" || action.type === "invalidate") {
    if (action.source) await refreshSource(action.source);
    await refreshTriggeredSources(name, schema, refreshSource);
    return undefined;
  }
  if (action.type !== "tool" || !transport) throw new Error(`action ${name} requires a transport`);
  const input = resolveValue(action.input ?? {}, scope);
  const response = await transport.call("ui.action", { plugin, view: schema.id, action: name, input });
  for (const effect of action.effects ?? []) await executeEffect(effect, response, scope, handlers, setState, refreshSource);
  await refreshTriggeredSources(name, schema, refreshSource);
  return response;
}

async function refreshTriggeredSources(name: string, schema: ViewSchema, refreshSource: (name: string) => Promise<void>): Promise<void> {
  for (const [sourceName, source] of Object.entries(schema.sources ?? {})) {
    if (source.refreshOn?.includes(name)) await refreshSource(sourceName);
  }
}

async function executeEffect(effect: Effect, result: unknown, scope: ResolveScope, handlers: UIEffectHandlers | undefined, setState: Dispatch<SetStateAction<unknown>>, refreshSource: (name: string) => Promise<void>): Promise<void> {
  const effectScope = { ...scope, result };
  if (effect.type === "set-state" && effect.path) {
    const value = effect.value === undefined ? result : resolveValue(effect.value, effectScope);
    setState((current: unknown) => setPath(current, effect.path!, value));
  }
  if (effect.type === "merge-state" && effect.path) {
    const value = resolveValue(effect.value ?? {}, effectScope);
    setState((current: unknown) => {
      const existing = resolvePath(current, effect.path!);
      const merged = existing && typeof existing === "object" ? { ...existing, ...(value && typeof value === "object" ? value : {}) } : value;
      return setPath(current, effect.path!, merged);
    });
  }
  if ((effect.type === "invalidate" || effect.type === "refresh-source") && effect.source) await refreshSource(effect.source);
  if (effect.type === "toast" && effect.message && handlers?.toast) await handlers.toast(effect.message, effect.variant);
  if (effect.type === "navigate" && effect.to && handlers?.navigate) await handlers.navigate(effect.to);
  if (effect.type === "dialog" && handlers?.dialog) await handlers.dialog(effect);
  if (effect.type === "close-dialog" && handlers?.closeDialog) await handlers.closeDialog(effect);
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return !!value && typeof value === "object" && !Array.isArray(value);
}
