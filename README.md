# Shrtr API — client examples

[![CI](https://github.com/alexander-po/shrtr-api/actions/workflows/ci.yml/badge.svg)](https://github.com/alexander-po/shrtr-api/actions/workflows/ci.yml)

Official, zero-dependency examples for the **[Shrtr](https://shrtr.top)** URL-shortener JSON API.

[Shrtr](https://shrtr.top) turns long URLs into short links with instant QR codes. The API is **free, anonymous (no API key, no signup), CORS-enabled, and rate-limited per IP**. Errors follow [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) (`application/problem+json`).

- **Base URL:** `https://shrtr.top/api/v1`
- **Human docs:** <https://shrtr.top/api>
- **Machine-readable spec:** <https://shrtr.top/openapi.json> (canonical; also vendored here as [`openapi.json`](openapi.json))
- **Service status:** <https://status.shrtr.top>

## Quickstart

```bash
curl -sS -X POST https://shrtr.top/api/v1/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/some/long/path"}'
```
```json
{
  "code": "ab3k9xz",
  "short_url": "https://shrtr.top/s/ab3k9xz",
  "original_url": "https://example.com/some/long/path",
  "created_at": "2026-07-16T12:00:00+00:00"
}
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/shorten` | Create a short link. Body: `{"url": "…", "alias": "optional"}`. Returns `201` with the created link. |
| `GET`  | `/api/v1/stats/{code}` | Aggregate click stats for a link (`clicks_count`, `last_clicked_at`, `created_at`, `enabled`). Never per-click data. |
| `GET`  | `/api/v1/health` | Liveness probe. Returns `{"status":"ok"}`. |

A custom `alias` is optional: 5–16 characters, letters/digits/hyphens, starting and ending alphanumeric (case-insensitive — stored lowercase).

## Authentication

**None.** The API is anonymous — no key, no signup, no OAuth. Rate limits are enforced per client IP:

| Endpoint | Limit |
|----------|-------|
| `POST /shorten` | 30 / minute (and 10 / hour when a custom `alias` is set) |
| `GET /stats/{code}` | 120 / minute |
| `GET /health` | unlimited |

A `429` response carries a `Retry-After` header (seconds). `POST /shorten` responses also include `X-RateLimit-Limit` / `X-RateLimit-Remaining` / `X-RateLimit-Reset`.

## Errors

All errors are `application/problem+json` (RFC 7807):

```json
{ "type": "about:blank", "title": "Unprocessable Entity", "status": 422,
  "detail": "The request body failed validation.", "errors": { "url": ["This is not a valid URL."] } }
```

- `422` validation → a top-level `errors` map keyed by field (note the JSON field is `alias` but its error key is `customAlias`); a URL rejected by the service (a blocked or unsafe destination) instead carries only `detail` (e.g. `This URL is not allowed.`).
- `410` on stats → the link exists but is inactive; a `reason` extension is one of `disabled` / `expired` / `exhausted`.
- `429` → includes a `retry_after` (seconds) extension and the `Retry-After` header. Beyond the per-IP limits, `POST /shorten` can also `429` from a per-destination burst limit (the same destination shortened too many times in a day).
- Other statuses: `400` (malformed request), `409` (alias taken), `503` (rare code-generation retry).

## Examples

Each is a single file, **no dependencies**, with a tiny reusable client plus a runnable demo:

| Language | File | Run |
|----------|------|-----|
| curl / bash | [`examples/curl/shorten.sh`](examples/curl/shorten.sh) | `bash examples/curl/shorten.sh` (needs `jq`) |
| Python (stdlib) | [`examples/python/shrtr.py`](examples/python/shrtr.py) | `python3 examples/python/shrtr.py` |
| JavaScript (Node 18+) | [`examples/javascript/shrtr.mjs`](examples/javascript/shrtr.mjs) | `node examples/javascript/shrtr.mjs` |
| Go (1.21+) | [`examples/go/main.go`](examples/go/main.go) | `cd examples/go && go run .` |

Every example honours a **`SHRTR_BASE`** environment variable (default `https://shrtr.top/api/v1`) — point it at a self-hosted instance or a mock server (this repo's CI runs the examples against a [Prism](https://github.com/stoplightio/prism) mock built from `openapi.json`, so tests never touch production).

**PHP, Rust, and others — contributions welcome** (see [Contributing](#contributing)).

## Continuous integration

- **On every push/PR** (hermetic, no network to the live API): lint `openapi.json` as OpenAPI 3.1, syntax-check every example, and **contract-test** the examples against a Prism mock generated from the spec.
- **Daily + on demand** (read-only against production): re-download the live `https://shrtr.top/openapi.json` and fail if the vendored copy has **drifted**, plus a read-only smoke of `GET /health` and the `GET /stats/{unknown}` → `404` problem+json contract. No `POST` is ever made against production, so CI never creates real links.

## Contributing

The examples are intentionally minimal and dependency-free so they never rot. A new-language example should follow the same shape: one file, a small client (`shorten` / `stats`), a runnable demo, and RFC 7807 error handling. Open a PR.

## About

This is the **official** examples repository for [shrtr.top](https://shrtr.top) — maintained by the service operator, not a third-party wrapper. Shrtr is a small, independent, privacy-first project (no tracking, no ads); see the [About page](https://shrtr.top/about).

## License

[MIT](LICENSE).
