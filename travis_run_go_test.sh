#!/bin/bash

# Verbose mode and stop on first error
set -ev

if [ "$GIMME_OS" = "linux" ]; then
    # Since Travis uses linux machines for builds, we can only run linux executables.
    go test -v ./...
fi
