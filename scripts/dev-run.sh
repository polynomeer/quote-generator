#!/usr/bin/env bash
set -euo pipefail
QG_CONFIG=${QG_CONFIG:-configs/config.yaml}
go run ./cmd/quote-generator