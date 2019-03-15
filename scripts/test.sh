#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

dir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
root="${dir}/.."

build_dir="${root}/build"
build_atr_dir="${build_dir}/src/github.com/android-test-runner/atr"

(cd "${build_atr_dir}"; GOPATH=${build_dir} go test ./...)
