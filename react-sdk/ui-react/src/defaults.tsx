import { createElement, type ComponentType, type ReactNode } from "react";
import type { ComponentDefinition } from "./registry.js";

type Props = Record<string, unknown>;

const globalProps = [
  "id", "title", "hidden", "dir", "lang", "role", "tabIndex", "contentEditable", "draggable", "translate",
  "aria-label", "aria-describedby", "aria-labelledby", "aria-live", "aria-current", "aria-expanded", "aria-controls", "aria-pressed", "aria-description", "aria-details", "aria-errormessage", "aria-invalid", "aria-placeholder", "aria-required",
] as const;

const commonProps = ["text", "type", "value", "defaultValue", "name", "placeholder", "disabled", "readOnly", "required", "autoComplete", "autoFocus", "min", "max", "minLength", "maxLength", "step", "pattern", "multiple", "checked", "selected", "rows", "cols", "width", "height", "loading", "error"] as const;
const linkProps = ["href", "target", "rel", "download", "src", "alt", "poster"] as const;
const formProps = ["accept", "acceptCharset", "capture", "dirname", "encType", "enterKeyHint", "form", "formEncType", "formMethod", "formNoValidate", "formTarget", "htmlFor", "inputMode", "list", "method", "noValidate", "size", "spellCheck", "wrap"] as const;
const htmlPropAliases = {
  "accept-charset": "acceptCharset",
  autocomplete: "autoComplete",
  autofocus: "autoFocus",
  contenteditable: "contentEditable",
  dirname: "dirname",
  enctype: "encType",
  enterkeyhint: "enterKeyHint",
  for: "htmlFor",
  formenctype: "formEncType",
  formmethod: "formMethod",
  formnovalidate: "formNoValidate",
  formtarget: "formTarget",
  inputmode: "inputMode",
  maxlength: "maxLength",
  minlength: "minLength",
  novalidate: "noValidate",
  readonly: "readOnly",
  spellcheck: "spellCheck",
  tabindex: "tabIndex",
} as const;
const htmlAliases = Object.keys(htmlPropAliases);
const events = ["click", "change", "input", "submit", "focus", "blur"] as const;
const htmlProps = [...globalProps, ...commonProps, ...linkProps, ...formProps, ...htmlAliases] as const;
const shadcnProps = [...htmlProps, "variant", "size", "orientation", "side", "align", "open", "defaultOpen", "modal", "label", "description", "status"] as const;
const domProps = new Set<string>([...globalProps, ...commonProps, ...linkProps, ...formProps]);
const eventProps = new Set(["onClick", "onChange", "onInput", "onSubmit", "onFocus", "onBlur"]);
const urlProps = new Set(["href", "src", "poster"]);

function safeURL(value: string): boolean {
  if (value.startsWith("/") || value.startsWith("#") || value.startsWith("?")) return true;
  try {
    return ["http:", "https:", "mailto:", "tel:"].includes(new URL(value).protocol);
  } catch {
    return false;
  }
}

function sanitizeProps(props: Props): Props {
  const result: Props = {};
  for (const [name, value] of Object.entries(props)) {
    if (value === undefined) continue;
    const normalizedName: string = htmlPropAliases[name as keyof typeof htmlPropAliases] ?? name;
    if (normalizedName in result) throw new Error(`duplicate prop ${normalizedName}`);
    if (normalizedName === "rows") {
      if (!Array.isArray(value)) throw new Error("rows must be an array");
      result[normalizedName] = value;
      continue;
    }
    if (urlProps.has(normalizedName)) {
      if (typeof value !== "string" || !safeURL(value)) throw new Error(`${normalizedName} must be a safe URL`);
      result[normalizedName] = value;
      continue;
    }
    if (["string", "number", "boolean"].includes(typeof value)) {
      result[normalizedName] = value;
      continue;
    }
    throw new Error(`${normalizedName} must be a string, number, or boolean`);
  }
  return result;
}

function DOMElement({ tag, props, className, children }: { tag: string; props: Props; className: string; children?: ReactNode }): ReactNode {
  const elementProps: Props = {};
  for (const [name, value] of Object.entries(props)) {
    if (eventProps.has(name) || (domProps.has(name) && name !== "text" && name !== "rows" && name !== "loading" && name !== "error")) elementProps[name] = value;
  }
  elementProps.className = className;
  const content = props.text ?? children;
  if (tag === "form" && typeof props.onSubmit === "function") {
    elementProps.onSubmit = (event: { preventDefault(): void }) => {
      event.preventDefault();
      void (props.onSubmit as (event: unknown) => Promise<void>)(event);
    };
  }
  return createElement(tag, elementProps, content as ReactNode);
}

function htmlClass(tag: string): string {
  if (tag === "a") return "text-ui-primary underline-offset-4 hover:underline";
  if (tag === "button") return "inline-flex items-center justify-center rounded-md bg-ui-primary px-3 py-2 text-sm font-medium text-ui-primary-foreground hover:bg-ui-primary/90 disabled:pointer-events-none disabled:opacity-50";
  if (["input", "textarea", "select"].includes(tag)) return "flex w-full rounded-md border border-ui-border bg-ui-background px-3 py-2 text-sm text-ui-foreground shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ui-ring disabled:cursor-not-allowed disabled:opacity-50";
  if (tag === "table") return "w-full caption-bottom text-sm";
  if (["th", "td"].includes(tag)) return "border-b border-ui-border px-4 py-3 text-left align-middle";
  if (["h1", "h2", "h3", "h4", "h5", "h6"].includes(tag)) return "font-semibold tracking-tight text-ui-foreground";
  if (tag === "img") return "h-auto max-w-full rounded-md";
  return "text-ui-foreground";
}

export function createDefaultHTMLDefinition(tag: string): ComponentDefinition {
  return {
    component: ((props: Props) => <DOMElement tag={tag} props={props} className={htmlClass(tag)}>{props.children as ReactNode}</DOMElement>) as ComponentType<Props>,
    allowedProps: htmlProps,
    propsSchema: sanitizeProps,
    events,
  };
}

const shadcnElement: Record<string, string> = {
  button: "button", checkbox: "input", input: "input", "native-select": "select", select: "select", slider: "input", switch: "button", textarea: "textarea",
  table: "table", "data-table": "table", progress: "progress", separator: "hr", label: "label", kbd: "kbd", typography: "p", spinner: "output",
  accordion: "details", collapsible: "details", dialog: "dialog", "alert-dialog": "dialog",
};

function shadcnClass(name: string): string {
  if (name === "button") return "inline-flex items-center justify-center rounded-md bg-ui-primary px-3 py-2 text-sm font-medium text-ui-primary-foreground hover:bg-ui-primary/90 disabled:pointer-events-none disabled:opacity-50";
  if (["input", "textarea", "select", "native-select"].includes(name)) return "flex w-full rounded-md border border-ui-border bg-ui-background px-3 py-2 text-sm text-ui-foreground shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ui-ring disabled:cursor-not-allowed disabled:opacity-50";
  if (name === "card") return "rounded-xl border border-ui-border bg-ui-card p-6 text-ui-card-foreground shadow-sm";
  if (["badge", "marker"].includes(name)) return "inline-flex items-center rounded-md bg-ui-muted px-2 py-1 text-xs font-medium text-ui-muted-foreground";
  if (["alert", "empty", "message", "toast"].includes(name)) return "rounded-lg border border-ui-border bg-ui-card p-4 text-ui-card-foreground";
  if (name === "separator") return "my-4 border-ui-border";
  if (name === "spinner") return "inline-flex size-4 animate-spin rounded-full border-2 border-ui-muted border-t-ui-primary";
  if (name === "skeleton") return "animate-pulse rounded-md bg-ui-muted";
  if (name === "table") return "w-full caption-bottom text-sm";
  return "rounded-md border border-ui-border bg-ui-card p-4 text-ui-card-foreground";
}

function DataTable({ rows, loading, error }: Props): ReactNode {
  if (loading) return <output role="status" className="text-sm text-ui-muted-foreground">Loading…</output>;
  if (error) return <output role="alert" className="text-sm text-ui-destructive">{String(error)}</output>;
  const items = Array.isArray(rows) ? rows : [];
  return <table className="w-full caption-bottom text-sm"><tbody>{items.map((row, index) => <tr key={index} className="border-b border-ui-border">{Object.values(row as Record<string, unknown>).map((value, cell) => <td key={cell} className="px-4 py-3 align-middle">{String(value ?? "")}</td>)}</tr>)}</tbody></table>;
}

export function createDefaultShadcnDefinition(name: string): ComponentDefinition {
  if (name === "data-table") return {
    component: DataTable as ComponentType<Props>,
    allowedProps: shadcnProps,
    propsSchema: sanitizeProps,
    events,
    children: false,
  };
  const tag = shadcnElement[name] ?? "div";
  return {
    component: ((props: Props) => <DOMElement tag={tag} props={props} className={shadcnClass(name)}>{props.children as ReactNode}</DOMElement>) as ComponentType<Props>,
    allowedProps: shadcnProps,
    propsSchema: sanitizeProps,
    events,
  };
}
