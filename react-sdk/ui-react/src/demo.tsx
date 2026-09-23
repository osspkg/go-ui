import { useMemo, type ReactElement, type ReactNode } from "react";
import type { RPCTransport } from "@osspkg/ui-transport";
import { createComponentRegistry, type ComponentRegistry } from "./registry.js";
import { PluginView } from "./renderer.js";
import { UIProvider } from "./runtime.js";

export interface UsersDemoProps {
  transport?: RPCTransport;
}

export function createUsersDemoRegistry(): ComponentRegistry {
  const registry = createComponentRegistry();
  registry.register("heading", { component: ({ text }: { text?: unknown }) => <h1>{String(text ?? "")}</h1>, allowedProps: ["text"] });
  registry.register("text", { component: ({ text }: { text?: unknown }) => <p>{String(text ?? "")}</p>, allowedProps: ["text"] });
  registry.register("button", { component: ({ text, children }: { text?: unknown; children?: ReactNode }) => <button type="button">{String(text ?? "")}{children}</button>, allowedProps: ["text"], events: ["click"] });
  registry.register("text-field", { component: ({ value }: { value?: unknown }) => <input value={String(value ?? "")} readOnly />, allowedProps: ["value"] });
  registry.register("form", { component: ({ children }: { children?: ReactNode }) => <form>{children}</form> });
  registry.register("table", { component: UsersTable, allowedProps: ["rows", "loading", "error"], children: false });
  return registry;
}

function UsersTable({ rows, loading, error }: { rows?: unknown; loading?: unknown; error?: unknown }): ReactElement {
  if (loading) return <div role="status">Loading users…</div>;
  if (error) return <div role="alert">Unable to load users</div>;
  const items = Array.isArray(rows) ? rows : [];
  return <table><tbody>{items.map((row, index) => <tr key={index}><td>{String((row as Record<string, unknown>).name ?? "")}</td></tr>)}</tbody></table>;
}

export function UsersDemo({ transport }: UsersDemoProps): ReactElement {
  const components = useMemo(() => createUsersDemoRegistry(), []);
  return <UIProvider components={components} {...(transport ? { transport } : {})}><PluginView plugin="users" view="users.list" context={{ workspaceId: "demo" }} /></UIProvider>;
}
