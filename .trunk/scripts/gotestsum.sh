#!/usr/bin/env bash
set -e
export PATH="${tool}/bin:${PATH}"
export GOTESTS='slow'
# green text on beginning tests
echo -e "\033[32mRunning tests...\033[0m"
gotestsum --format pkgname -- -shuffle=on -tags integration ./...
