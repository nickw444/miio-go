#!/usr/bin/env bash
set -eu

ROOT_DIR=$(cd "$(dirname "$0")"/.. && pwd)
cd ${ROOT_DIR}

./bin/regenerate_mocks.sh
if ! git diff --exit-code HEAD -- '*.go'; then
  echo "Mocks are not up to date, please run: make mocks"
  exit 1
fi
