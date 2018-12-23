#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

dir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
root="${dir}/.."

go=$(which go || true)
if [ -z "${go:-}" ]; then
	(>&2 echo "ERROR: go not found")
	(>&2 echo "Please install go and make it available on your PATH.")
	exit 127
fi


version=$("${dir}/print-version.sh")
build_dir="${root}/build"
build_atr_dir="${build_dir}/src/atr"
arch=${GOARCH:-}
os=${GOOS:-}

rm -rf "${build_dir}"
mkdir -p "${build_atr_dir}"
cp "${root}"/*.go "${build_atr_dir}"
(
	export GOPATH=${build_dir}
	export GOARCH=${arch}
	export GOOS=${os}

	cd ${build_atr_dir}
	"${go}" get -d ./...
	"${go}" build -ldflags "-X main.version=${version}" -o "${build_dir}/bin/${GOOS:-localOS}_${GOARCH:-localARCH}/atr"
)
