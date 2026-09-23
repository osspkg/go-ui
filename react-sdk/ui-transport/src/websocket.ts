import { createRequest, readResponse, type CallOptions, type RPCNotification, type RPCTransport } from "./rpc.js";

interface SocketLike {
  readyState: number;
  onopen: ((event: unknown) => void) | null;
  onclose: ((event: unknown) => void) | null;
  onerror: ((event: unknown) => void) | null;
  onmessage: ((event: { data: string }) => void) | null;
  send(data: string): void;
  close(): void;
}

interface SocketConstructor {
  new (url: string): SocketLike;
}

export interface WebSocketTransportOptions {
  socket?: SocketConstructor;
  timeoutMs?: number;
  reconnect?: boolean;
  reconnectDelayMs?: number;
}

interface PendingCall {
  resolve(value: unknown): void;
  reject(error: unknown): void;
  timer: ReturnType<typeof setTimeout>;
  abort?: () => void;
  signal?: AbortSignal;
}

export class WebSocketTransport implements RPCTransport {
  private socket: SocketLike | undefined;
  private closed = false;
  private connecting: Promise<void> | undefined;
  private readonly pending = new Map<number, PendingCall>();
  private readonly listeners = new Set<(notification: RPCNotification) => void>();
  private readonly reconnectListeners = new Set<() => void>();
  private readonly Socket: SocketConstructor;
  private reconnectTimer: ReturnType<typeof setTimeout> | undefined;
  private hasConnected = false;

  constructor(private readonly endpoint: string, private readonly options: WebSocketTransportOptions = {}) {
    this.Socket = options.socket ?? (WebSocket as unknown as SocketConstructor);
  }

  async call<T>(method: string, params: unknown, callOptions: CallOptions = {}): Promise<T> {
    if (this.closed) throw new Error("transport is closed");
    await this.connect();
    const request = createRequest(method, params);
    const timeoutMs = callOptions.timeoutMs ?? this.options.timeoutMs ?? 30_000;
    return new Promise<T>((resolve, reject) => {
      const pending: PendingCall = { resolve, reject, timer: setTimeout(() => {
        if (this.pending.get(request.id) !== pending) return;
        this.pending.delete(request.id);
        this.cleanupPending(pending);
        reject(new Error("rpc request timed out"));
      }, timeoutMs), ...(callOptions.signal ? { signal: callOptions.signal } : {}) };
      const abort = () => {
        if (this.pending.get(request.id) !== pending) return;
        this.pending.delete(request.id);
        this.cleanupPending(pending);
        reject(new DOMException("request was cancelled", "AbortError"));
      };
      pending.abort = abort;
      this.pending.set(request.id, pending);
      callOptions.signal?.addEventListener("abort", abort, { once: true });
      if (callOptions.signal?.aborted) {
        abort();
        return;
      }
      try {
        this.socket!.send(JSON.stringify(request));
      } catch (error) {
        this.pending.delete(request.id);
        this.cleanupPending(pending);
        reject(error);
      }
    });
  }

  subscribe(handler: (notification: RPCNotification) => void): () => void {
    this.listeners.add(handler);
    return () => this.listeners.delete(handler);
  }

  onReconnect(handler: () => void): () => void {
    this.reconnectListeners.add(handler);
    return () => this.reconnectListeners.delete(handler);
  }

  close(): void {
    this.closed = true;
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.reconnectTimer = undefined;
    this.socket?.close();
    for (const pending of this.pending.values()) {
      this.cleanupPending(pending);
      pending.reject(new Error("transport is closed"));
    }
    this.pending.clear();
  }

  private connect(): Promise<void> {
    if (this.socket?.readyState === 1) return Promise.resolve();
    if (this.connecting) return this.connecting;
    this.connecting = new Promise<void>((resolve, reject) => {
      const socket = new this.Socket(this.endpoint);
      this.socket = socket;
      socket.onopen = () => {
        if (this.socket !== socket) return;
        this.connecting = undefined;
        const reconnect = this.hasConnected;
        this.hasConnected = true;
        resolve();
        if (reconnect) for (const listener of this.reconnectListeners) listener();
      };
      socket.onerror = () => {
        if (this.socket !== socket) return;
        this.socket = undefined;
        this.connecting = undefined;
        reject(new Error("websocket connection failed"));
        if (!this.closed && this.options.reconnect) this.scheduleReconnect();
      };
      socket.onclose = () => {
        if (this.socket !== socket) return;
        this.socket = undefined;
        this.connecting = undefined;
        reject(new Error("websocket closed before connection opened"));
        if (!this.closed && this.options.reconnect) this.scheduleReconnect();
        for (const [id, pending] of this.pending) {
          this.pending.delete(id);
          this.cleanupPending(pending);
          pending.reject(new Error("websocket closed"));
        }
      };
      socket.onmessage = (event) => this.handleMessage(event.data);
    });
    return this.connecting;
  }

  private handleMessage(data: string): void {
    let message: unknown;
    try { message = JSON.parse(data); } catch { return; }
    if (!message || typeof message !== "object") return;
    const record = message as { id?: unknown; method?: unknown };
    if (typeof record.id === "number") {
      const pending = this.pending.get(record.id);
      if (!pending) return;
      this.pending.delete(record.id);
      this.cleanupPending(pending);
      try { pending.resolve(readResponse(message, record.id)); } catch (error) { pending.reject(error); }
      return;
    }
    if (typeof record.method === "string") for (const listener of this.listeners) listener(message as RPCNotification);
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer) return;
    const delay = this.options.reconnectDelayMs ?? 1000;
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = undefined;
      if (!this.closed) void this.connect().catch(() => undefined);
    }, delay);
  }

  private cleanupPending(pending: PendingCall): void {
    clearTimeout(pending.timer);
    if (pending.abort && pending.signal) pending.signal.removeEventListener("abort", pending.abort);
  }
}
