#!/usr/bin/env bash

# Checks if mocks need re-generating

set -euo pipefail

make mocks
if ! git diff --exit-code HEAD; then
  error "Mocks are not up to date, please run: make mocks"
fi
