import { useEffect, useState, type CSSProperties, type ReactElement, type ReactNode } from "react";
import { uiRPCMethods } from "@osspkg/ui-core";
import type { EventHandler, RegionName, UIGetViewResponse, UINode, ViewSchema } from "@osspkg/ui-core";
import { validateView } from "@osspkg/ui-core";
import { useUIRuntime, useUIViewRuntime, ViewRuntimeProvider } from "./runtime.js";

export interface PluginViewProps {
  plugin?: string;
  view: string;
  context?: unknown;
  schema?: ViewSchema;
}

export function PluginView(props: PluginViewProps): ReactElement {
  const { transport } = useUIRuntime();
  const [schema, setSchema] = useState<ViewSchema | undefined>(props.schema);
  const [error, setError] = useState<Error | undefined>();
  useEffect(() => {
    if (props.schema) {
      setSchema(props.schema);
      setError(undefined);
      return;
    }
    setSchema(undefined);
    setError(undefined);
    if (!transport || !props.plugin) return;
    let active = true;
    void transport.call<UIGetViewResponse>(uiRPCMethods.getView, { plugin: props.plugin, view: props.view, context: props.context }).then((response) => {
      if (active) setSchema(response.schema);
    }).catch((reason: unknown) => { if (active) setError(reason instanceof Error ? reason : new Error("view loading failed")); });
    return () => { active = false; };
  }, [props.context, props.plugin, props.schema, props.view, transport]);
  if (error) return <div role="alert">{error.message}</div>;
  if (!schema) return <div role="status">Loading…</div>;
  return <ViewRenderer schema={schema} {...(props.plugin ? { plugin: props.plugin } : {})} context={props.context ?? {}} />;
}

export function ViewRenderer({ schema, plugin, context = {} }: { schema: ViewSchema; plugin?: string; context?: unknown }): ReactElement {
  validateView(schema);
  const runtimeKey = `${plugin ?? ""}:${schema.id}:${schema.revision ?? ""}`;
  return <ViewRuntimeProvider key={runtimeKey} schema={schema} {...(plugin ? { plugin } : {})} context={context}><ViewRegions /></ViewRuntimeProvider>;
}

function ViewRegions(): ReactElement {
  const runtime = useUIViewRuntime();
  const regions: RegionName[] = ["top-header", "left-panel", "content", "right-panel", "bottom"];
  return <>
    {regions.map((region) => <section key={region} data-ui-region={region} className="grid grid-cols-12 gap-4">{runtime.schema.regions[region].map((node) => <NodeRenderer key={node.id} node={node} />)}</section>)}
  </>;
}

function NodeRenderer({ node }: { node: UINode }): ReactNode {
  const runtime = useUIViewRuntime();
  if (node.when !== undefined && !runtime.resolve(node.when)) return null;
  const definition = runtime.components.get(node.component);
  if (!definition) return <div role="alert" data-unsupported-component={node.component}>Unsupported component: {node.component}</div>;
  const rawProps = Object.fromEntries(Object.entries(node.props ?? {}).map(([key, value]) => [key, runtime.resolve(value)]));
  const props = validateProps(rawProps, definition.allowedProps, definition.propsSchema);
  const events = bindEvents(node.events, definition.events, runtime);
  const children = definition.children === false || !node.children?.length ? undefined : <div data-ui-children={node.id} className="grid grid-cols-12 gap-4">{node.children.map((child) => <NodeRenderer key={child.id} node={child} />)}</div>;
  const slots = renderSlots(node.slots, definition.slots);
  const Component = definition.component as (props: Record<string, unknown>) => ReactElement;
  const layout = node.layout;
  const layoutStyle = {
    "--ui-row": layout?.row ?? "auto",
    "--ui-cols": layout?.cols ?? 12,
    "--ui-start": (layout?.offset ?? 0) + 1,
  } as CSSProperties;
  return <div data-ui-node={node.id} data-ui-row={layout?.row} data-ui-cols={layout?.cols} data-ui-offset={layout?.offset} className="col-span-[var(--ui-cols)] col-start-[var(--ui-start)] row-start-[var(--ui-row)]" style={layoutStyle}><Component {...props} {...events} {...(slots ? { slots } : {})}>{children}</Component></div>;
}

function renderSlots(slots: Record<string, UINode[]> | undefined, allowedSlots: readonly string[] | undefined): Record<string, ReactNode[]> | undefined {
  if (!slots) return undefined;
  if (!allowedSlots) throw new Error("component slots are not declared");
  const rendered: Record<string, ReactNode[]> = {};
  for (const [name, nodes] of Object.entries(slots)) {
    if (!allowedSlots.includes(name)) throw new Error(`slot ${name} is not allowed`);
    rendered[name] = nodes.map((child) => <NodeRenderer key={child.id} node={child} />);
  }
  return rendered;
}

function validateProps(rawProps: Record<string, unknown>, allowedProps: readonly string[] | undefined, schema: ((props: Record<string, unknown>) => unknown) | undefined): Record<string, unknown> {
  for (const key of Object.keys(rawProps)) {
    if (forbiddenPluginProps.has(key) || key.startsWith("on")) throw new Error(`prop ${key} is forbidden`);
  }
  if (!allowedProps && !schema && Object.keys(rawProps).length > 0) throw new Error("component props are not declared");
  if (allowedProps) for (const key of Object.keys(rawProps)) if (!allowedProps.includes(key)) throw new Error(`prop ${key} is not allowed`);
  return schema ? schema(rawProps) as Record<string, unknown> : rawProps;
}

const forbiddenPluginProps = new Set(["dangerouslySetInnerHTML", "style", "className", "ref", "key"]);

function bindEvents(handlers: Record<string, EventHandler> | undefined, allowedEvents: readonly string[] | undefined, runtime: ReturnType<typeof useUIViewRuntime>): Record<string, (event: unknown) => Promise<void>> {
  const entries = Object.entries(handlers ?? {});
  if (entries.length > 0 && !allowedEvents) throw new Error("component events are not declared");
  return Object.fromEntries(entries.map(([name, handler]) => {
    if (!/^[a-z][a-z0-9-]*$/.test(name) || !allowedEvents?.includes(name)) throw new Error(`event ${name} is not allowed`);
    return [eventProp(name), createEventHandler(handler, runtime)];
  }));
}

function eventProp(name: string): string {
  return `on${name.slice(0, 1).toUpperCase()}${name.slice(1)}`;
}

function createEventHandler(handler: EventHandler, runtime: ReturnType<typeof useUIViewRuntime>): (event: unknown) => Promise<void> {
  return async (event) => {
    if (handler.action) await runtime.runAction(handler.action, event);
    for (const step of handler.steps ?? []) if (step.action) await runtime.runAction(step.action, event);
  };
}
