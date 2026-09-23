import { createRequest, readResponse, type CallOptions, type RPCNotification, type RPCTransport } from "./rpc.js";

export interface HTTPTransportOptions {
  fetch?: typeof globalThis.fetch;
  headers?: Record<string, string>;
  timeoutMs?: number;
}

export class HTTPTransport implements RPCTransport {
  private readonly fetcher: typeof globalThis.fetch;
  private readonly headers: Record<string, string>;
  private readonly timeoutMs: number;

  constructor(private readonly endpoint: string, options: HTTPTransportOptions = {}) {
    this.fetcher = options.fetch ?? globalThis.fetch.bind(globalThis);
    this.headers = { "content-type": "application/json", ...(options.headers ?? {}) };
    this.timeoutMs = options.timeoutMs ?? 30_000;
  }

  async call<T>(method: string, params: unknown, options: CallOptions = {}): Promise<T> {
    const request = createRequest(method, params);
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), options.timeoutMs ?? this.timeoutMs);
    const abort = () => controller.abort();
    options.signal?.addEventListener("abort", abort, { once: true });
    try {
      const response = await this.fetcher(this.endpoint, {
        method: "POST",
        headers: this.headers,
        body: JSON.stringify(request),
        signal: controller.signal,
      });
      const payload = await response.json() as unknown;
      if (!response.ok) throw new Error(`http request failed with status ${response.status}`);
      return readResponse<T>(payload, request.id);
    } finally {
      clearTimeout(timeout);
      options.signal?.removeEventListener("abort", abort);
    }
  }

  subscribe(_handler: (notification: RPCNotification) => void): () => void {
    return () => undefined;
  }

  close(): void {}
}
