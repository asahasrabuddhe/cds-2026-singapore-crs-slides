---
theme: slidev-theme-architectural-console
addons:
  - slidev-addon-live-terminal
title: Build Your Own Container Runtime in Go
info: |
  ContainerDays Singapore 2026, workshop
  Ajitem Sahasrabuddhe
highlighter: shiki
lineNumbers: false
drawings:
  persist: false
transition: none
mdc: true
layout: cover
---

# Build Your Own Container Runtime in Go

::label::
CONTAINERDAYS 2026 · SINGAPORE · WORKSHOP

::subtitle::
Namespaces, cgroups, capabilities and a pivoted root. Forty-five minutes, five tests to make green.

::speaker::
Ajitem Sahasrabuddhe

::role::
Automattic · @asahasrabuddhe

::initials::
AS

<!--
**S01 · COVER · before 0:00**

Up while people find seats. One line about who you are, then S02.
-->

---
layout: default
---

# Before we start

<div class="text-lg space-y-3 pt-2">

**1.** The Multipass VM from `docs/setup.md`. `multipass shell crs`, then:

```bash
git clone https://github.com/asahasrabuddhe/cds-2026-singapore-crs-code.git ~/crs
cd ~/crs && ./scripts/check.sh          # every line OK, last line: ready
```

**2.** A second `multipass shell crs`, with your editor on `~/crs/workshop/start`.

**3.** This, in the VM:

```bash
cd ~/crs/workshop/start && make build rootfs && go test ./...
```

</div>

<div v-click class="pt-2 accent">
Six packages fail. Five are the steps; the sixth goes green with step 4.
</div>

<!--
**S02 · SETUP · from 4:15pm, then 0:00 to 2:00**

Up from the moment the room opens. Walk the room asking one question:
"what did check.sh say?" A VM still provisioning has an empty
`/etc/crs-provisioned`; it is Go downloading, so wait. Anyone with no VM
pairs with a neighbour now, not at step 2.

The sixth package is `cli`, whose tests use the capability set from step 4.
-->

---
layout: default
---

# A container is a process wearing a costume

| Step | The costume | You write | Minutes |
|---|---|---|---|
| 1 | Namespaces: what it can see | `CloneFlags()` | 5 |
| 2 | Rootfs: where it lives | the `pivot_root` sequence | 10 |
| 3 | cgroups: what it can use | `Files()`, `cpuMax()` | 7 |
| 4 | Capabilities: what it can do | `Default()` | 5 |
| 5 | Init: what PID 1 owes | `reap()` | 7 |

<div class="pt-4 text-lg">
Each step: I explain the idea, you write it, the test tells you when it is right.
</div>

<div v-click class="pt-3 accent">
Stuck? The answer is the same file two directories up. <code>cp ../../internal/rootfs/rootfs.go internal/rootfs/</code>, and carry on.
</div>

<!--
**S03 · THE PLAN · 2:00**

Two minutes, then step 1 whatever the clock says. The table is the whole
workshop: five functions, and everything else in the repository is
scaffolding.

The click is the rule that keeps forty-five minutes moving. Copying the
answer is allowed, and saying so now stops people falling silently
behind at step 2.
-->

---
layout: step
---

::title::
Step 1 · Namespaces: what it can see

::default::

<div class="text-base space-y-3">

A namespace is one bit passed to `clone(2)`, and it applies to the child. So `mycontainer run` starts itself again as `mycontainer child`, through `/proc/self/exe`, with the flags on that clone. The child is the container.

**You write** `CloneFlags()` in `internal/namespace/namespace.go`, then one line in `sysProcAttr`, `internal/container/container_linux.go`.

</div>

::right::

<div class="text-base space-y-4">

<div>The seven <code>CloneNew…</code> constants are at the top of the file. OR in each one whose field is on.</div>

<div v-click="1">

```go
var f uintptr
if s.PID {
	f |= CloneNewPID
}
// ...six more
```

</div>

<div v-click="2">

```go
attr.Cloneflags = s.Namespaces.CloneFlags()
```

</div>

</div>

::run::

```bash
go test ./internal/namespace/ && make build
sudo ./mycontainer run -ns pid,uts -hostname box ubuntu \
  /bin/sh -c 'hostname; echo $$; ls -d /proc/[0-9]* | wc -l'
```

<div v-click="3" class="accent">
<code>box</code>, a small PID, and <code>/proc</code> still lists the host. Why?
</div>

<!--
**S04 · STEP 1 · 4:00, five minutes**

**SAY**, one minute: the idea on the left. Then "go", and walk the room.

Click 1 at minute two, click 2 at minute three, for the room rather than
for anyone in particular. The test names the missing flag when one is
wrong.

**WATCH FOR** "not in a new UTS namespace": they have `CloneFlags` but not
the wiring line. The refusal is deliberate, because without the flag the
call would rename their VM.

Click 3 at minute four, and leave the question hanging. The shell is
PID 8 or 9 rather than 2 because the Go init's threads took the numbers
in between; it comes back at step 3.

**CHECKPOINT** 9:00. Step 2 whatever the room's state.
-->

---
layout: step
---

::title::
Step 2 · Rootfs: where it lives

::default::

<div class="text-base space-y-3">

`ps` reads `/proc`, and `/proc` is the host's until the container mounts its own. That is step 1's answer. The container's root is a directory that `pivot_root` swaps in, and the order of the calls is the lesson.

**You write** the `steps` list in `Pivot`, `internal/rootfs/rootfs.go`. The comment lists all seven.

</div>

::right::

<div class="text-base space-y-4">

<div>Each step is a name, a call through <code>m</code>, and whether it may fail. The first:</div>

```go
{"make / private", func() error {
	return m.Mount("", "/", "",
		MS_REC|MS_PRIVATE, "")
}, false},
```

<div v-click="1">Write one, run the test. It prints the call it wanted, flags and all.</div>

<div v-click="2">Private first. Bind the rootfs onto itself: <code>pivot_root</code> wants a mount point.</div>

</div>

::run::

```bash
go test ./internal/rootfs/ && make build
touch /tmp/host-was-here
sudo ./mycontainer run ubuntu /bin/sh -c 'ps; ls /; ls /tmp; mount'
```

<div v-click="3" class="accent">
Two processes. The image's tree. An empty <code>/tmp</code>. Only its own mounts.
</div>

<!--
**S05 · STEP 2 · 9:00, ten minutes**

**SAY**, one minute: the answer to step 1's question, then the idea. Step
1 of the seven matters most: systemd makes every mount shared, so
unmounting the old root inside the namespace would propagate out and
unmount the host's. People have done this to their laptops.

Click 1 at minute two, click 2 at minute four. At minute six, say the
`cp` escape hatch out loud again.

If anyone asks why `mount` names the VM's disk for `/`, such as
`/dev/sda1 on / type ext4`: the root is a bind mount of `rootfs/ubuntu`,
which lives on that disk, and the device nodes in `/dev` are bind mounts
of the VM's. What is gone is every other mount the VM has.

**CHECKPOINT** 19:00. Past 20:00, do step 3 together on the screen and
give its time to step 5.
-->

---
layout: step
---

::title::
Step 3 · cgroups: what it can use

::default::

<div class="text-base space-y-3">

cgroup v2 is a filesystem. A limit is a number written to a file in the container's directory, and the kernel enforces it. The child is born inside that directory, so the first fork of a fork bomb already counts.

**You write** `Files()` and `cpuMax()` in `internal/cgroup/cgroup.go`, then two lines in `sysProcAttr`.

<div class="warn">Keep <code>-pids</code> at sixteen or more. It counts the Go init's own threads.</div>

</div>

::right::

<div class="text-base space-y-4">

<div>Four files: <code>memory.max</code> in bytes, <code>memory.swap.max</code> as <code>0</code> beside it, <code>cpu.max</code>, <code>pids.max</code>. A zero limit writes nothing.</div>

<div v-click="1"><code>cpu.max</code> is a quota and a period in microseconds. The period is 100000, so half a CPU is <code>50000 100000</code>.</div>

<div v-click="2">

```go
if group != nil {
	attr.UseCgroupFD = true
	attr.CgroupFD = group.FD
}
```

</div>

</div>

::run::

```bash
go test ./internal/cgroup/ && make build
sudo ./mycontainer run -pids 20 ubuntu \
  /bin/sh -c 'cat /sys/fs/cgroup/pids.max; for i in $(seq 40); do sleep 2 & done'
```

<div v-click="3" class="accent">
<code>20</code>, then <code>Cannot fork</code>.
</div>

<!--
**S06 · STEP 3 · 19:00, seven minutes**

**SAY**, one minute: the idea. The line that bites: without
`memory.swap.max` a process at `memory.max` carries on into swap and the
limit is a suggestion. Docker sets it for you.

Click 1 at minute two, click 2 at minute three. `UseCgroupFD` is
`CLONE_INTO_CGROUP`.

The warning answers the PID question from step 1: the init's threads
are processes as far as `pids.max` is concerned. On the workshop VM the
floor is about seven, and it moves by one from run to run. Just under
it, `sh` cannot start ("resource temporarily unavailable") or the Go
runtime dies with "failed to create new OS thread". More cores mean more
threads, so sixteen is the margin.

Anyone done early: `-memory 64m` and
`/bin/sh -c 'head -c 100000000 /dev/zero | tail'`, exit 137.

**CHECKPOINT** 26:00.
-->

---
layout: step
---

::title::
Step 4 · Capabilities: what it can do

::default::

<div class="text-base space-y-3">

The kernel does not ask whether you are UID 0. It asks whether the process holds the capability. `mount(2)` wants `CAP_SYS_ADMIN`, and the fourteen Docker grants by default leave it out. Dropping the rest is written; choosing what to keep is yours.

**You write** `Default()` in `internal/caps/caps.go`.

</div>

::right::

<div class="text-base space-y-4">

<div>The constants at the top of <code>caps.go</code> are the fourteen, plus <code>SysAdmin</code>. Add the fourteen to a <code>Set</code>.</div>

<div v-click="1">Then give one back: <code>-caps default,+CAP_SYS_ADMIN</code>. <code>a82425fb</code>, and the mount works. Bit 21.</div>

<div class="warn">Capabilities belong to a thread, and Go moves goroutines between threads. <code>Child</code> calls <code>runtime.LockOSThread()</code> before the drop. Read it.</div>

</div>

::run::

```bash
go test ./internal/caps/ ./internal/cli/ && make build
sudo ./mycontainer run ubuntu \
  /bin/sh -c 'grep CapBnd /proc/self/status; mount -t tmpfs none /tmp'
```

<div v-click="2" class="accent">
<code>a80425fb</code>, and root cannot mount.
</div>

<!--
**S07 · STEP 4 · 26:00, five minutes**

**SAY**, one minute: the idea, then the warning. It happened while
writing this, only on the rootless path, only sometimes: drop on one
thread, fork the command from another, and the command starts with
everything. If there is one Go-specific thing to take home, it is this.

`cli` goes green here too.

**SAY** when the first refusals come up: "mount says permission denied,
to a process that id says is root."

**CHECKPOINT** 31:00.
-->

---
layout: step
---

::title::
Step 5 · Init: what PID 1 owes

::default::

<div class="text-base space-y-3">

When a parent exits before its children, the kernel hands them to PID 1, and only PID 1 can reap them. A PID 1 that never waits collects zombies. The kernel also drops a default signal sent to PID 1, which is why `docker stop` waits ten seconds. Signals are written; reaping is yours.

**You write** `reap()` in `internal/initproc/initproc_linux.go`. Today it waits for its own child and nobody else.

</div>

::right::

<div class="text-base space-y-4">

<div>Loop on <code>syscall.Wait4(-1, &ws, syscall.WNOHANG, nil)</code>: any child, and never block.</div>

<div v-click="1">A pid of <code>0</code> means children exist and none has exited: return. <code>ECHILD</code> means nothing is left: return. <code>EINTR</code>: go round again.</div>

<div v-click="2">When the pid is <code>main</code>, keep <code>ExitCode(ws)</code> and set <code>done</code>, then go round again. Any other pid was an orphan: reaped, and forgotten.</div>

</div>

::run::

```bash
go test ./internal/initproc/ && go test ./...
```

<div v-click="3" class="accent">
Green means no zombie after an orphan exits, and nothing is red any more.
</div>

<!--
**S08 · STEP 5 · 31:00, seven minutes**

**SAY**, one minute: the idea. This is why Docker has `--init` and tini
exists; most programs were never written to be PID 1.

Click 1 at minute two, click 2 at minute four.

**WATCH FOR** a loop without `WNOHANG`. It hangs when children exist and
none has exited, so the test times out rather than fails. Say it out
loud at minute three.

**CHECKPOINT** 38:00. Go to the wrap whatever the room's state; the test
tells the stragglers when they are done.
-->

---
layout: default
---

# What you just built, counted

<div class="h-64">
  <DemoTerminal />
</div>

<div v-click class="pt-3 font-mono text-sm opacity-80">
1 clone3 · 1 sethostname · 12 mount · 1 pivot_root · 1 umount2 · 27 PR_CAPBSET_DROP · 3 execve
</div>

<!--
**S09 · LIVE · 38:00**

**DO**, in the slide's terminal (the finished runtime):
```bash
./scripts/demo.sh 7
```

One container under `strace -f`, tallied. The tally shows about 75
`prctl`, more than the click. Twenty-seven are capability drops. The rest
are the Go runtime naming its memory with `PR_SET_VMA`.

**SAY**, from the tally on the click:
> "One clone3 carrying every namespace flag and the cgroup fd: steps 1
> and 3. Twelve mounts and one pivot_root: step 2. Twenty-seven capability
> drops: step 4. And an init that reaps: step 5. You wrote each of those
> this afternoon. That is a container. It always was."
-->

---
layout: default
---

# Take home

<v-clicks>

**1.** A container is a process wearing a costume: seven namespace bits, a pivoted mount, a cgroup directory, fourteen capability bits, and an init.

**2.** A PID namespace changes the numbers. `/proc` is what `ps` reads.

**3.** Root is a bounding set. Drop one bit and root cannot mount.

**4.** PID 1 owes the tree reaping and signals. Bring an init.

**5.** In Go: start the child through `/proc/self/exe`, lock the thread before you drop capabilities, and count your own threads against `pids.max`.

</v-clicks>

<div v-click class="pt-6 accent">
At home: three lines of UID mapping make it rootless. <code>docs/workshop.md</code> has them.
</div>

<!--
**S10 · TAKE HOME · 40:00**

Read all five, then leave this up for questions and for the people still
typing.

The last click is the bonus there was no time for: in `sysProcAttr`, one
`UidMappings` and one `GidMappings` entry, and `./mycontainer run` works
without sudo, root inside and them outside.

Questions that come up, in order of likelihood:

1. **Why not `unshare(2)` from Go?** A PID namespace only applies to
   children, and a user namespace needs a single-threaded process. Fork
   and exec is what every runtime does.
2. **Where is the network?** A namespace with a down loopback. Wiring a
   veth pair, a bridge and NAT is its own workshop.
3. **How does this compare to runc?** Same syscalls in the same order.
   runc adds seccomp, AppArmor, hooks, the OCI config and years of edge
   cases. `libcontainer/rootfs_linux.go` will look familiar after today.

Nothing follows in the Workshop Room, so stay with the stragglers until
the drinks at 5:30.
-->

---
layout: end
url: github.com/asahasrabuddhe/cds-2026-singapore-crs-code
qr: /images/qr-slides.svg
handles:
  - label: cds-2026-singapore-crs-slides
    href: https://github.com/asahasrabuddhe/cds-2026-singapore-crs-slides
  - label: cds-2026-singapore-crs-code
    href: https://github.com/asahasrabuddhe/cds-2026-singapore-crs-code
  - label: "@asahasrabuddhe"
    href: https://github.com/asahasrabuddhe
---

<!--
**S11 · END · 45:00**

For the photograph. Questions stay on S10.
-->
