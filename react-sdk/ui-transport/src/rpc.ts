export interface RPCRequest {
  jsonrpc: "2.0";
  id: number;
  method: string;
  params: unknown;
  meta?: { correlationId: string };
}

export interface RPCResponse<T = unknown> {
  jsonrpc: "2.0";
  id: number;
  result?: T;
  error?: { code: number | string; message: string; data?: unknown };
}

export interface RPCNotification {
  jsonrpc: "2.0";
  method: string;
  params?: unknown;
}

export interface CallOptions {
  signal?: AbortSignal;
  timeoutMs?: number;
}

export interface RPCTransport {
  call<T>(method: string, params: unknown, options?: CallOptions): Promise<T>;
  subscribe(handler: (notification: RPCNotification) => void): () => void;
  onReconnect?(handler: () => void): () => void;
  close(): void;
}

export class RPCError extends Error {
  readonly code: number | string;
  readonly details: unknown;

  constructor(code: number | string, message: string, details?: unknown) {
    super(message);
    this.name = "RPCError";
    this.code = code;
    this.details = details;
  }
}

let nextRequestID = 1;
let nextCorrelationID = 1;

export function createRequest(method: string, params: unknown): RPCRequest {
  const id = nextRequestID++;
  return { jsonrpc: "2.0", id, method, params, meta: { correlationId: `ui-${nextCorrelationID++}` } };
}

export function readResponse<T>(value: unknown, expectedID: number): T {
  if (!value || typeof value !== "object") throw new RPCError("INVALID_RESPONSE", "rpc response is invalid");
  const response = value as RPCResponse<T>;
  if (response.jsonrpc !== "2.0" || response.id !== expectedID) throw new RPCError("INVALID_RESPONSE", "rpc response id is invalid");
  if (response.error) throw new RPCError(response.error.code, response.error.message, response.error.data);
  return response.result as T;
}
