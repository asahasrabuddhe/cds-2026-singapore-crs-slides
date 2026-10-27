#!/usr/bin/env bash
#
# Entry point for the in-slide terminals (slidev-addon-live-terminal).
# Lands in ~/crs inside the presenter's Multipass VM, the same VM every
# attendee makes, as the ordinary user. The speaker notes say `sudo` where
# it is needed.
#
# The dev server runs this as `$SHELL -l -c ./scripts/dev-shell.sh` with
# the slides repo as cwd. ~/crs in the VM is the laptop's checkout of the
# code repository, shared by `make vm` there. Keep it on the final branch.
#
# PROMPT_DIRTRIM keeps bash's prompt to the last path component, so the
# full path does not wrap every command line.
set -euo pipefail

MACHINE="${CRS_VM:-crs}"

if ! command -v multipass >/dev/null; then
  echo "multipass is not installed."
  exit 1
fi

exec multipass exec "$MACHINE" -- bash -lc 'cd ~/crs && PROMPT_DIRTRIM=1 exec bash -l'
