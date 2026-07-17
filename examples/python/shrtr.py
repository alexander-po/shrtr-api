#!/usr/bin/env python3
"""Minimal, zero-dependency Python client for the Shrtr URL-shortener API.

Docs:  https://shrtr.top/api
Spec:  https://shrtr.top/openapi.json

The API is free and anonymous (no key, no signup). Errors are RFC 7807
problem+json and surface here as ShrtrError carrying the parsed problem.
Uses only the standard library (urllib) — nothing to install.
"""
from __future__ import annotations

import json
import urllib.error
import urllib.request

BASE = "https://shrtr.top/api/v1"


class ShrtrError(Exception):
    """Raised on a non-2xx response; `.problem` holds the RFC 7807 body."""

    def __init__(self, status: int, problem: dict):
        self.status = status
        self.problem = problem
        detail = problem.get("detail") or problem.get("title") or "request failed"
        super().__init__(f"HTTP {status}: {detail}")


def _request(method: str, path: str, body: dict | None = None) -> dict:
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(
        f"{BASE}{path}",
        data=data,
        method=method,
        headers={"Content-Type": "application/json", "Accept": "application/json"},
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            return json.loads(resp.read() or b"{}")
    except urllib.error.HTTPError as err:
        try:
            problem = json.loads(err.read() or b"{}")
        except json.JSONDecodeError:
            problem = {"title": err.reason}
        raise ShrtrError(err.code, problem) from None


def shorten(url: str, alias: str | None = None) -> dict:
    """Create a short link. Returns {code, short_url, original_url, created_at}."""
    body = {"url": url}
    if alias:
        body["alias"] = alias
    return _request("POST", "/shorten", body)


def stats(code: str) -> dict:
    """Aggregate click stats for a short link."""
    return _request("GET", f"/stats/{code}")


if __name__ == "__main__":
    link = shorten("https://example.com/some/long/path")
    print("short_url:", link["short_url"])
    print("stats:    ", stats(link["code"]))
