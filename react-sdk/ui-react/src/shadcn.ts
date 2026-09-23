import type { ComponentDefinition, ComponentRegistry } from "./registry.js";
import { createDefaultShadcnDefinition } from "./defaults.js";

export const shadcnComponentNames = [
  "accordion", "alert", "alert-dialog", "aspect-ratio", "attachment", "avatar", "badge", "breadcrumb", "bubble", "button", "button-group", "calendar", "card", "carousel", "chart", "checkbox", "collapsible", "combobox", "command", "context-menu", "data-table", "date-picker", "dialog", "direction", "drawer", "dropdown-menu", "empty", "field", "hover-card", "input", "input-group", "input-otp", "item", "kbd", "label", "marker", "menubar", "message", "message-scroller", "native-select", "navigation-menu", "pagination", "popover", "progress", "questionnaire", "radio-group", "resizable", "scroll-area", "select", "separator", "sheet", "sidebar", "skeleton", "slider", "spinner", "switch", "table", "tabs", "textarea", "toast", "toggle", "toggle-group", "tooltip", "typography",
] as const;

export type ShadcnComponentName = typeof shadcnComponentNames[number];

const shadcnComponentSet = new Set<string>(shadcnComponentNames);

export function isShadcnComponentName(name: string): name is ShadcnComponentName {
  return shadcnComponentSet.has(name);
}

// registerShadcnComponents registers host implementations for shadcn/ui names.
// The host deliberately supplies the implementations and their prop allowlists.
export function registerShadcnComponents(registry: ComponentRegistry, definitions: Partial<Record<ShadcnComponentName, ComponentDefinition>>): void {
  for (const [name, definition] of Object.entries(definitions)) {
    if (!isShadcnComponentName(name) || !definition) throw new Error(`unsupported shadcn component ${name}`);
    registry.register(name, definition);
  }
}

// registerDefaultShadcnComponents registers baseline, dependency-free renderers for every shadcn catalog name.
// Install and register host-owned shadcn/ui source components to replace these defaults where richer behavior is needed.
// Existing host registrations take precedence.
export function registerDefaultShadcnComponents(registry: ComponentRegistry): void {
  for (const name of shadcnComponentNames) if (!registry.has(name)) registry.register(name, createDefaultShadcnDefinition(name));
}
