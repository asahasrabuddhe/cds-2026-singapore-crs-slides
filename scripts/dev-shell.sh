#!/usr/bin/env bash
#
# Entry point for the in-slide terminals (slidev-addon-live-terminal).
# Lands in the companion code repository inside the crs-dev OrbStack
# machine, as the ordinary user. The speaker notes say `sudo` where it is
# needed.
#
# The dev server runs this as `$SHELL -l -c ./scripts/dev-shell.sh` with
# the slides repo as cwd. The sibling checkout is used, not
# vendor/crs-code: the runbook builds and runs against the sibling, and
# the terminal must land in the same tree.
set -euo pipefail

MACHINE="${CRS_VM:-crs-dev}"
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CODE_REPO="$(cd "$HERE/../cds-2026-singapore-crs-code" 2>/dev/null && pwd || true)"
if [[ -z "$CODE_REPO" || ! -f "$CODE_REPO/Makefile" ]]; then
  echo "cds-2026-singapore-crs-code not found next to this repo."
  echo "clone it as a sibling of $HERE"
  exit 1
fi

# OrbStack mounts /Users at the same path inside the machine, so the code
# repo's path is valid on both sides.
#
# TERM=dumb is for the orb client only; the shell inside gets the real one
# back. On startup orb queries the terminal for its background colour and
# cursor position, and those queries land in the session's scrollback,
# which the addon replays to every terminal that attaches later. With
# TERM=dumb orb asks nothing.
#
# PROMPT_DIRTRIM keeps bash's prompt to the last path component, so the
# full /Users/... path does not wrap every command line.
INNER_TERM=${TERM:-xterm-256color}
export TERM=dumb
exec orb -m "$MACHINE" -w "$CODE_REPO" env TERM="$INNER_TERM" PROMPT_DIRTRIM=1 PATH="$PATH:/usr/local/go/bin" bash -l
