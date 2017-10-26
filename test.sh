#!/bin/bash

# Script that runs tests, code coverage, and benchmarks all at once.
# Builds a symlink in /tmp, mostly to avoid messing with GOPATH at the user's shell level.

TEMPORARY_PATH="/tmp/gval_test"
SRC_PATH="${TEMPORARY_PATH}/src"
FULL_PATH="${TEMPORARY_PATH}/src/github.com/PaesslerAG/gval"

# set up temporary directory
rm -rf "${SRC_PATH}"
mkdir -p "${FULL_PATH}"
rm -rf "${FULL_PATH}"

ln -s $(pwd) "${FULL_PATH}"
export GOPATH="${TEMPORARY_PATH}"

pushd "${TEMPORARY_PATH}/src"

# run the actual tests.
go test -bench=. -benchmem -coverprofile coverage.out
status=$?

if [ "${status}" != 0 ];
then
	exit $status
fi

# run the actual tests.
go test -bench=Random -benchtime 10m -timeout 30m -benchmem -coverprofile coverage.out
status=$?

if [ "${status}" != 0 ];
then
	exit $status
fi


popd
