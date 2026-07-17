#!/usr/bin/env bash
#
# Shrtr API — curl examples. Free, anonymous, no API key.
# Docs: https://shrtr.top/api  ·  Spec: https://shrtr.top/openapi.json
#
# Requires: curl, jq
set -euo pipefail

BASE="${SHRTR_BASE:-https://shrtr.top/api/v1}"

echo "# 1) Shorten a URL"
resp=$(curl -sS -X POST "$BASE/shorten" \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/some/long/path"}')
echo "$resp" | jq .
code=$(printf '%s' "$resp" | jq -r .code)

echo
echo "# 2) Shorten with a custom alias (5-16 chars, letters/digits/hyphen)"
curl -sS -X POST "$BASE/shorten" \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/pricing","alias":"my-link"}' | jq .

echo
echo "# 3) Aggregate click stats for the code from step 1"
curl -sS "$BASE/stats/$code" | jq .

echo
echo "# 4) Health probe"
curl -sS "$BASE/health" | jq .
