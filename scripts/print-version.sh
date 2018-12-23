#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

dir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)

git=$(which git || true)
if [ -z "${git:-}" ]; then
	(>&2 echo "ERROR: git not found")
	(>&2 echo "Please install git and make it available on your PATH.")
	exit 127
fi

num_commits=$("${git}" rev-list --count HEAD) # version is monotone increasing (only works for one branch)
commit_hash=$("${git}" rev-parse --short HEAD) # version is reproducible
build_number="${BUILD_NUMBER:-1}" # version is unique for multiple build runs from same commit (e.g. Jenkins)
version="${num_commits}-${commit_hash}-${build_number}"

echo -n ${version}
