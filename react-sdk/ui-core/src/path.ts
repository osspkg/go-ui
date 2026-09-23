const forbiddenSegments = new Set(["__proto__", "prototype", "constructor"]);

export function resolvePath(value: unknown, path: string): unknown {
  const segments = validatePath(path);
  let current = value;
  for (const segment of segments) {
    if (current === null || current === undefined) return undefined;
    if (typeof current !== "object" && typeof current !== "function") return undefined;
    if (!Object.prototype.hasOwnProperty.call(current, segment)) return undefined;
    current = (current as Record<string, unknown>)[segment];
  }
  return current;
}

export function setPath<T>(value: T, path: string, next: unknown): T {
  const segments = validatePath(path);
  if (segments.length === 0) throw new Error("path is required");
  const root = cloneValue(value) as Record<string, unknown>;
  let current = root;
  for (let index = 0; index < segments.length - 1; index += 1) {
    const segment = segments[index]!;
    const existing = current[segment];
    current[segment] = existing && typeof existing === "object" ? cloneValue(existing) : {};
    current = current[segment] as Record<string, unknown>;
  }
  current[segments.at(-1)!] = next;
  return root as T;
}

export function validatePath(path: string): string[] {
  if (path.length === 0 || path.length > 1024) throw new Error("invalid path");
  const segments = path.split(".");
  if (segments.some((segment) => segment.length === 0 || forbiddenSegments.has(segment))) {
    throw new Error("unsafe path");
  }
  return segments;
}

function cloneValue<T>(value: T): T {
  if (Array.isArray(value)) return value.map(cloneValue) as T;
  if (value && typeof value === "object") {
    return Object.fromEntries(Object.entries(value).map(([key, nested]) => [key, cloneValue(nested)])) as T;
  }
  return value;
}
