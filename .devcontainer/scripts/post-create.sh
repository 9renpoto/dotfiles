#!/bin/sh
set -eu
# Install mise explicitly; applying dotfiles remains an opt-in action.
curl -fsSL https://mise.jdx.dev/install.sh -o /tmp/install-mise.sh
sh /tmp/install-mise.sh
