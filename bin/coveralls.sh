#!/usr/bin/env bash
set -eu

ROOT_DIR=$(cd "$(dirname "$0")"/.. && pwd)

main() {
  cd ${ROOT_DIR}
  local cover_file=$(mktemp)
  go test -v -coverprofile ${cover_file} ./...
  goveralls -coverprofile ${cover_file}
}

main "$@"
