#!/usr/bin/env bash
set -eu

ROOT_DIR=$(cd "$(dirname "$0")"/.. && pwd)
cd ${ROOT_DIR}

find . -name "mocks" | xargs rm -rf
mockery -dir capability -output capability/mocks -all
mockery -dir device/rthrottle -output device/rthrottle/mocks -all
mockery -dir device -output device/mocks -all
mockery -dir protocol -output protocol/mocks -name Protocol
mockery -dir protocol/packet -output protocol/packet/mocks -all
mockery -dir protocol/transport -output protocol/transport/mocks -all
mockery -dir subscription/common -output subscription/common/mocks -all
