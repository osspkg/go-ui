import type { ComponentType } from "react";

export interface ComponentDefinition<Props = Record<string, unknown>> {
  component: ComponentType<Props>;
  allowedProps?: readonly string[];
  propsSchema?: (props: Record<string, unknown>) => Props;
  events?: readonly string[];
  slots?: readonly string[];
  children?: boolean;
}

export class ComponentRegistry {
  private readonly definitions = new Map<string, ComponentDefinition>();

  register(name: string, definition: ComponentDefinition): void {
    if (!name || this.definitions.has(name)) throw new Error(`component ${name} is already registered`);
    this.definitions.set(name, definition);
  }

  get(name: string): ComponentDefinition | undefined { return this.definitions.get(name); }
  has(name: string): boolean { return this.definitions.has(name); }
}

export function createComponentRegistry(): ComponentRegistry {
  return new ComponentRegistry();
}
