#!/usr/bin/env bash
set -eu

ROOT_DIR=$(cd "$(dirname "$0")"/.. && pwd)
cd ${ROOT_DIR}

go test -v --timeout 1s ./...
