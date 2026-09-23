import type { ComponentDefinition, ComponentRegistry } from "./registry.js";
import { createDefaultHTMLDefinition } from "./defaults.js";

export const htmlTagNames = [
  "a", "abbr", "address", "area", "article", "aside", "audio", "b", "bdi", "bdo", "blockquote", "br", "button", "canvas", "caption", "cite", "code", "col", "colgroup", "data", "datalist", "dd", "del", "details", "dfn", "dialog", "div", "dl", "dt", "em", "fieldset", "figcaption", "figure", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "header", "hgroup", "hr", "i", "img", "input", "ins", "kbd", "label", "legend", "li", "main", "map", "mark", "menu", "meter", "nav", "ol", "optgroup", "option", "output", "p", "picture", "pre", "progress", "q", "rp", "rt", "ruby", "s", "samp", "search", "section", "select", "small", "source", "span", "strong", "sub", "summary", "sup", "table", "tbody", "td", "textarea", "tfoot", "th", "thead", "time", "tr", "track", "u", "ul", "var", "video", "wbr",
] as const;

export type HTMLTagName = typeof htmlTagNames[number];

const htmlTagSet = new Set<string>(htmlTagNames);

export function isHTMLTagName(name: string): name is HTMLTagName {
  return htmlTagSet.has(name);
}

// registerHTMLTags registers host-owned implementations for safe HTML tags.
export function registerHTMLTags(registry: ComponentRegistry, definitions: Partial<Record<HTMLTagName, ComponentDefinition>>): void {
  for (const [name, definition] of Object.entries(definitions)) {
    if (!isHTMLTagName(name) || !definition) throw new Error(`unsupported html tag ${name}`);
    registry.register(name, definition);
  }
}

// registerDefaultHTMLComponents registers safe native React implementations for every supported HTML tag.
// Existing host registrations take precedence.
export function registerDefaultHTMLComponents(registry: ComponentRegistry): void {
  for (const name of htmlTagNames) if (!registry.has(name)) registry.register(name, createDefaultHTMLDefinition(name));
}
