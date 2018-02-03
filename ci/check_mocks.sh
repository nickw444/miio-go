#!/usr/bin/env bash

# Checks if mocks need re-generating

set -euo pipefail

go get github.com/vektra/mockery/.../
export PATH=$PATH:$GOPATH/bin
make mocks
if ! git diff --exit-code HEAD; then
  error "Mocks are not up to date, please run: make mocks"
fi
