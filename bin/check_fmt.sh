#!/usr/bin/env bash
set -eu

ROOT_DIR=$(cd "$(dirname "$0")"/.. && pwd)
cd ${ROOT_DIR}

go fmt ./...
if ! git diff --exit-code HEAD; then
  echo "Code is not formatted, please run: go fmt ./..."
  exit 1
fi
