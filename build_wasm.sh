#!/bin/bash
set -euo pipefail

# ensure in root
cd "$(dirname "$0")"

GOOS=js GOARCH=wasm go build -ldflags="-w -s" -o ./website/public/SOMAS2020.wasm