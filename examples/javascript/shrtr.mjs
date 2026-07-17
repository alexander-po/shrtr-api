// Minimal, zero-dependency JavaScript client for the Shrtr URL-shortener API.
//
// Docs:  https://shrtr.top/api
// Spec:  https://shrtr.top/openapi.json
//
// The API is free and anonymous (no key, no signup). Errors are RFC 7807
// problem+json and surface here as ShrtrError. Uses the built-in `fetch`
// (Node 18+ or any modern browser) — nothing to install.

const BASE = "https://shrtr.top/api/v1";

export class ShrtrError extends Error {
  constructor(status, problem) {
    super(`HTTP ${status}: ${problem.detail ?? problem.title ?? "request failed"}`);
    this.name = "ShrtrError";
    this.status = status;
    this.problem = problem; // the parsed RFC 7807 body
  }
}

async function request(method, path, body) {
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: { "Content-Type": "application/json", Accept: "application/json" },
    body: body ? JSON.stringify(body) : undefined,
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new ShrtrError(res.status, data);
  return data;
}

/** Create a short link → { code, short_url, original_url, created_at }. */
export const shorten = (url, alias) =>
  request("POST", "/shorten", alias ? { url, alias } : { url });

/** Aggregate click stats for a short link. */
export const stats = (code) => request("GET", `/stats/${code}`);

// Demo — run: node shrtr.mjs
if (import.meta.url === `file://${process.argv[1]}`) {
  const link = await shorten("https://example.com/some/long/path");
  console.log("short_url:", link.short_url);
  console.log("stats:    ", await stats(link.code));
}
