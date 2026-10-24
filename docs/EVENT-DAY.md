# Event day

The runbook for the workshop. What to run, when, and what to do when an
attendee's box is not ready. The teaching script is in the presenter notes.

## The week before

Forty-five minutes has no room for installing a VM, so the setup has to
happen before the day. Send attendees
[`docs/setup.md`](https://github.com/asahasrabuddhe/cds-2026-singapore-crs-code/blob/main/docs/setup.md)
from the code repository, with one line in front of it: run
`./scripts/check.sh` inside the VM and reply with the last line.

Everyone, the presenter included, uses the same Multipass VM, made the
same way on Windows, macOS and Linux. One tool, one procedure, so every
question in the room has one answer.

On the day, anyone without a ready box pairs with a neighbour. The usual
reason is the Ubuntu image, about 600 MB, on slow venue wifi. Bring it on
a USB stick, both the arm64 image for Apple Silicon and the amd64 one for
everyone else; `multipass launch` takes a `file://` path in place of
`24.04`. Even then, pair them up: the launch and provisioning take longer
than step 1.

## The day before: the presenter's VM

From the code repository on the laptop:

```console
$ make vm                     # multipass launch, provision, clone to ~/crs in the VM
```

The VM has its own clone from GitHub, like every attendee's. Push from
the laptop, then `git pull` in the VM.

## T-60: check the VM

```console
$ make vm-shell               # lands in ~/crs inside the VM
$ git pull
$ ./scripts/check.sh          # ready; pulls ubuntu:26.04 the first time
$ make build
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
$ npm run dev
```

`p` for presenter mode. Click into the terminal on S09 and confirm it
lands in the code repository inside the VM.

## T-10: the room

Put S02 up. It has the commands attendees need and stays up until you
start. Walk the room. The question is "what did check.sh say".

## During

One slide per step, and the room writes code for most of each: 5, 10, 7,
5 and 7 minutes. Explain the idea for a minute, say "go", and walk the
room. The clicks on the right are hints, and the notes say when to give
each one, for the room rather than for anyone in particular. Anyone stuck
gets the finished file, which S03 says out loud:

```bash
cp ../../internal/rootfs/rootfs.go internal/rootfs/rootfs.go
```

Checkpoints: 4:00 start step 1 whatever the clock says, 19:00 out of step
2, 26:00 out of step 3, 31:00 out of step 4, 38:00 go to the wrap
whatever the room's state. Past 20:00 at the second, step 3 becomes a
demo on the screen and its time goes to step 5.

## When something dies

**The slide terminal is dead.** `r` in it.

**A test hangs on step 5.** A `wait4` loop without `WNOHANG`. It is in the
S08 notes; say it at minute three.

**Someone's `-pids` run kills the init.** They went below what the init's threads need, about seven on the workshop VM.
That is the warning on S06, and they found it first.

**Someone's `-hostname` is refused.** They have not wired the clone flags
yet. The refusal is deliberate: without `CLONE_NEWUTS` the call would
rename their VM.

**Everything is gone.** The code repository's `docs/` has the expected
output of every checkpoint. The workshop survives without the deck; it
does not survive without the tests, and those are on every attendee's
machine.
