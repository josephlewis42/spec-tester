#!/usr/bin/env sh

set -eux
cd "${0%/*}"/..

make build
./bin/spec-tester $@