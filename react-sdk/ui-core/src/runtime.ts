import { resolvePath, setPath } from "./path.js";
import type { RuntimeError, SourceRuntimeState, UIValue } from "./schema.js";

export interface StateRuntime<T> {
  get(): T;
  set(path: string, value: unknown): void;
  merge(path: string, value: Record<string, unknown>): void;
}

export function createStateRuntime<T>(initial: T): StateRuntime<T> {
  let state = initial;
  return {
    get: () => state,
    set: (path, value) => { state = setPath(state, path, value); },
    merge: (path, value) => {
      const existing = resolvePath(state, path);
      const base = existing && typeof existing === "object" ? existing : {};
      state = setPath(state, path, {...base, ...value});
    },
  };
}

export function runtimeError(code: string, message: string, details?: unknown, retryable?: boolean): RuntimeError {
  return { code, message, ...(details === undefined ? {} : { details }), ...(retryable === undefined ? {} : { retryable }) };
}

export function emptySource(): SourceRuntimeState {
  return { status: "idle" };
}
