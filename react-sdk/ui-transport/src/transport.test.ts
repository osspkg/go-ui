import { describe, expect, it, vi } from "vitest";
import { HTTPTransport } from "./http.js";

describe("HTTPTransport", () => {
  it("sends JSON-RPC requests and returns result", async () => {
    let requestInit: RequestInit | undefined;
    const fetcher = vi.fn(async (_input: RequestInfo | URL, init?: RequestInit) => {
      requestInit = init;
      return new Response(JSON.stringify({ jsonrpc: "2.0", id: 1, result: { ok: true } }), { status: 200 });
    });
    const transport = new HTTPTransport("https://example.test/rpc", { fetch: fetcher });
    await expect(transport.call("ping", { value: 1 })).resolves.toEqual({ ok: true });
    expect(fetcher).toHaveBeenCalledOnce();
    expect(JSON.parse(String(requestInit?.body))).toMatchObject({ method: "ping", params: { value: 1 } });
  });

  it("propagates caller cancellation to fetch", async () => {
    const fetcher = vi.fn((_input: RequestInfo | URL, init?: RequestInit) => new Promise<Response>((_resolve, reject) => {
      init?.signal?.addEventListener("abort", () => reject(new DOMException("aborted", "AbortError")), { once: true });
    }));
    const transport = new HTTPTransport("https://example.test/rpc", { fetch: fetcher });
    const controller = new AbortController();
    const request = transport.call("slow", {}, { signal: controller.signal });

    controller.abort();

    await expect(request).rejects.toMatchObject({ name: "AbortError" });
  });

  it("aborts requests that exceed the timeout", async () => {
    vi.useFakeTimers();
    const fetcher = vi.fn((_input: RequestInfo | URL, init?: RequestInit) => new Promise<Response>((_resolve, reject) => {
      init?.signal?.addEventListener("abort", () => reject(new DOMException("timed out", "AbortError")), { once: true });
    }));
    const transport = new HTTPTransport("https://example.test/rpc", { fetch: fetcher, timeoutMs: 10 });
    const request = transport.call("slow", {});
    const assertion = expect(request).rejects.toMatchObject({ name: "AbortError" });

    await vi.advanceTimersByTimeAsync(10);

    await assertion;
    vi.useRealTimers();
  });

  it("rejects malformed JSON-RPC responses", async () => {
    const fetcher = vi.fn(async () => new Response(JSON.stringify({ jsonrpc: "2.0", id: 999, result: true }), { status: 200 }));
    const transport = new HTTPTransport("https://example.test/rpc", { fetch: fetcher });

    await expect(transport.call("ping", {})).rejects.toMatchObject({ name: "RPCError", code: "INVALID_RESPONSE" });
  });
});
