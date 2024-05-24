#!/usr/bin/env bash
set -e
export PATH="${tool}/bin:${PATH}"
export GOTESTS='superslow'
# green text on beginning tests
echo -e "\033[32mRunning tests...\033[0m"
gotestsum --format dots-v2 -- -shuffle=on -tags integration ./...
