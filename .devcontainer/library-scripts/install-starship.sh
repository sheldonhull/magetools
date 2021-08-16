#!/usr/bin/env bash
# Run with elevated privileges
set -e
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi
curl -fsSL https://starship.rs/install.sh | bash -s -- --bin-dir /usr/local/bin --force && echo "completed setup of starship.rs"
