# Event day

The runbook for the workshop. What to run, when, and what to do when an
attendee's box is not ready. The teaching script is in the presenter notes.

## The week before

Forty-five minutes has no room for installing a VM, so the setup has to
happen before the day. Send attendees
[`docs/setup.md`](https://github.com/asahasrabuddhe/cds-2026-singapore-crs-code/blob/main/docs/setup.md)
from the code repository, with one line in front of it: run
`./scripts/check.sh` inside the VM and reply with the last line.

On the day, anyone without a ready box pairs with a neighbour. Three
reasons it happens, and the answer to each:

- Multipass has to download an Ubuntu image and the venue wifi is slow.
  Bring the image on a USB stick: `multipass launch` accepts a local file
  path in place of `24.04`. Even then, pair them up; the launch and
  provisioning take longer than step 1.
- Windows attendees on WSL2 find they are on cgroup v1. Pair them.
- Docker Desktop users can follow along in a privileged container, route
  D of the setup doc. Rootful only, which is everything today.

**The day before: update OrbStack.** When an update is pending, the orb
client prints an "OrbStack update available!" banner at the top of every
new shell, and the first in-slide terminal is a new shell on the
projector. Run `orb update` the day before, not the morning of, then
`make vm-shell` once to confirm the banner is gone.

## T-60: the presenter's VM

From the code repository on the laptop:

```console
$ make vm && make vm-shell
$ make build rootfs
$ ./scripts/check.sh          # ready
$ ./scripts/demo.sh 7         # the one live beat
```

Then the starting point, because the room will:

```console
$ cd workshop/start && make build rootfs && go test ./...
```

Six packages fail on Linux. If a seventh does, or one fewer, the stubs
have drifted from the reference; do not present a starting point you have
not run.

## T-30: the deck

```console
$ CRS_CODE_DIR=../cds-2026-singapore-crs-code npm run snippets
$ npm run dev
```

`p` for presenter mode. Click into the terminal on S23 and confirm it
lands in the code repository inside the VM.

## T-10: the room

Put S02 up. It has the commands attendees need and stays up until you
start. Walk the room. The question is "what did check.sh say".

## During

Five "your turn" slides carry the time: 4, 8, 6, 3 and 5 minutes. Walk
the room while they type, and use the click on each to leave the next
question hanging. Anyone stuck gets the finished file, which S04 says out
loud:

```bash
cp ../../internal/rootfs/rootfs.go internal/rootfs/rootfs.go
```

Checkpoints: 6:00 start step 1 whatever the clock says, 21:00 out of step
2, 29:00 out of step 3, 39:30 go to the wrap whatever the room's state.
Past 22:00 at the second, step 3 becomes a demo on the screen and its
time goes to step 5.

## When something dies

**The slide terminal is dead.** `r` in it.

**A test hangs on step 5.** A `wait4` loop without `WNOHANG`. It is in the
S22 notes; say it at minute three.

**Someone's `-pids` run kills the init.** They used a limit under sixteen.
That is the click on S15, and they found it first.

**Someone's `-hostname` is refused.** They have not wired the clone flags
yet. The refusal is deliberate: without `CLONE_NEWUTS` the call would
rename their VM.

**Everything is gone.** The code repository's `docs/` has the expected
output of every checkpoint. The workshop survives without the deck; it
does not survive without the tests, and those are on every attendee's
machine.
