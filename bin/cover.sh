#!/usr/bin/env bash
set -eu
ROOT_DIR=$(cd "$(dirname "$0")"/.. && pwd)

usage() {
  local prog
  prog=$(basename "$0")
  cat <<EOF
Usage: ${prog} [ OPTIONS... ]

  Run tests & produce coverage

  Options:
    --html             Produce HTML coverage output

  Usage:
    ${prog}
      Will run tests & output coverage information to the terminal.

    ${prog} --html
      Will run tests & output coverage information to a HTML report.
EOF
  exit
}


main() {
  local html=false;

  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      -h|--help)
        usage
        ;;
      --html)
        html=true
        ;;
      -*|--*|*)
        error "unknown option '$1'"
        ;;
    esac
    shift
  done

  cd ${ROOT_DIR}

  local cover_file=$(mktemp);
  local cover_mode;
  go test -v -coverprofile ${cover_file} ./...;

  if [ ${html} = true ]; then
    cover_mode='-html';
  else
    cover_mode='-func';
  fi

  go tool cover ${cover_mode} ${cover_file}
}

main "$@"
