import { afterEach, describe, expect, it, vi } from "vitest";
import { WebSocketTransport } from "./websocket.js";

class FakeSocket {
  static instances: FakeSocket[] = [];
  readyState = 0;
  onopen: ((event: unknown) => void) | null = null;
  onclose: ((event: unknown) => void) | null = null;
  onerror: ((event: unknown) => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  readonly sent: string[] = [];

  constructor(readonly url: string) {
    FakeSocket.instances.push(this);
  }

  send(data: string): void {
    this.sent.push(data);
  }

  close(): void {
    this.readyState = 3;
    this.onclose?.({});
  }

  open(): void {
    this.readyState = 1;
    this.onopen?.({});
  }

  receive(message: unknown): void {
    this.onmessage?.({ data: JSON.stringify(message) });
  }
}

afterEach(() => {
  FakeSocket.instances = [];
  vi.useRealTimers();
});

describe("WebSocketTransport", () => {
  it("sends calls and delivers notifications", async () => {
    const transport = new WebSocketTransport("wss://example.test/rpc", { socket: FakeSocket });
    const notifications: unknown[] = [];
    transport.subscribe((notification) => notifications.push(notification));
    const result = transport.call<{ ok: boolean }>("ping", { value: 1 });
    const socket = FakeSocket.instances[0]!;
    socket.open();
    await Promise.resolve();
    const request = JSON.parse(socket.sent[0]!) as { id: number };
    socket.receive({ jsonrpc: "2.0", id: request.id, result: { ok: true } });
    await expect(result).resolves.toEqual({ ok: true });

    socket.receive({ jsonrpc: "2.0", method: "ui.invalidate", params: { view: "users" } });
    expect(notifications).toEqual([{ jsonrpc: "2.0", method: "ui.invalidate", params: { view: "users" } }]);
    transport.close();
  });

  it("reconnects and notifies the runtime after a new socket opens", async () => {
    vi.useFakeTimers();
    const transport = new WebSocketTransport("wss://example.test/rpc", { socket: FakeSocket, reconnect: true, reconnectDelayMs: 10 });
    const reconnected = vi.fn();
    transport.onReconnect(reconnected);
    const result = transport.call("ping", {});
    const first = FakeSocket.instances[0]!;
    first.open();
    await Promise.resolve();
    const request = JSON.parse(first.sent[0]!) as { id: number };
    first.receive({ jsonrpc: "2.0", id: request.id, result: true });
    await expect(result).resolves.toBe(true);

    first.close();
    await vi.advanceTimersByTimeAsync(10);
    const second = FakeSocket.instances[1]!;
    second.open();

    expect(reconnected).toHaveBeenCalledOnce();
    transport.close();
  });

  it("rejects a call when the socket closes before opening", async () => {
    const transport = new WebSocketTransport("wss://example.test/rpc", { socket: FakeSocket, reconnect: true });
    const result = transport.call("ping", {});
    FakeSocket.instances[0]!.close();

    await expect(result).rejects.toThrow("websocket closed before connection opened");
    transport.close();
  });
});
