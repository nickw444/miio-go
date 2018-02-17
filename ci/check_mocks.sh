#!/usr/bin/env bash

# Checks if mocks need re-generating

set -euo pipefail

export PATH=$PATH:$GOPATH/bin
make mocks
if ! git diff --exit-code HEAD; then
  echo "Mocks are not up to date, please run: make mocks"
  exit 1
fi
