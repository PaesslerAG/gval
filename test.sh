#!/bin/bash

# Script that runs tests, code coverage, and benchmarks all at once.
# Builds a symlink in /tmp, mostly to avoid messing with GOPATH at the user's shell level.

GVAL_PATH=$HOME/gopath/src/github.com/PaesslerAG/gval

# run the actual tests.
cd "${GVAL_PATH}"
go test -bench=. -benchmem -coverprofile coverage.out
status=$?

if [ "${status}" != 0 ];
then
	exit $status
fi

# run random test for a longer period.
go test -bench=Random -benchtime 9m -timeout 30m -benchmem -coverprofile coverage.out
status=$?

if [ "${status}" != 0 ];
then
	exit $status
fi

