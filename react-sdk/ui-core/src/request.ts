export interface RequestTracker {
  next(key: string): number;
  isCurrent(key: string, requestId: number): boolean;
}

export function createRequestTracker(): RequestTracker {
  const current = new Map<string, number>();
  return {
    next: (key) => {
      const requestId = (current.get(key) ?? 0) + 1;
      current.set(key, requestId);
      return requestId;
    },
    isCurrent: (key, requestId) => current.get(key) === requestId,
  };
}
