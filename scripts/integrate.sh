#!/bin/sh
set -eu
sh scripts/setup-git-credentials.sh
bash scripts/sync-skills.sh
