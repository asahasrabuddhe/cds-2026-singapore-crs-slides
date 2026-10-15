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
**S01 · COVER · 0:00**

Up while people find seats and plug in. Say the one thing that matters
before anything else: if `./scripts/check.sh` in your VM does not say
`ready`, put a hand up now. Forty-five minutes has no room for installing
a VM, so anyone without one pairs with a neighbour today.
-->

---
layout: default
---

# Before we start

<div class="text-lg space-y-3 pt-2">

**1.** A Linux VM that passes the pre-flight. In it:

```bash
git clone https://github.com/asahasrabuddhe/cds-2026-singapore-crs-code.git ~/crs
cd ~/crs && ./scripts/check.sh          # every line OK, last line: ready
```

**2.** Your editor open on `~/crs/workshop/start`. Five files have a `TODO(step N)` in them.

**3.** This, in the VM:

```bash
cd ~/crs/workshop/start && make build rootfs && go test ./...
```

</div>

<div v-click class="pt-2 accent">
Six packages fail. Five are the steps; the sixth goes green with step 4.
</div>

<!--
**S02 · SETUP · 0:30**

Three minutes, no more. Walk the room while they type. The two failure
modes: a VM that has not finished provisioning (`cat /etc/crs-provisioned`
is empty; it is Go downloading, wait) and WSL2 on cgroup v1 (pair with a
neighbour; `docs/setup.md` route C is for after).

The sixth package is `cli`, whose tests use the capability set from step 4.
On a Mac the count is five, because the init test only runs on Linux.

**SAY**:
> "Each failing package is a step. When they are all green you have a
> container runtime, and you wrote the parts that make it one."
-->

---
layout: statement
---

# A container is a process<br/>wearing a costume

<!--
**S03 · STATEMENT · 3:30**

The thesis. Everything today is one piece of the costume, and every piece
is one function.
-->

---
layout: default
---

# The costume, and what you write

| Step | Pillar | You write | Minutes | Test |
|---|---|---|---|---|
| 1 | Namespaces | `CloneFlags()` | 5 | `./internal/namespace/` |
| 2 | Rootfs | the `pivot_root` sequence | 10 | `./internal/rootfs/` |
| 3 | cgroups | `Files()`, `cpuMax()` | 8 | `./internal/cgroup/` |
| 4 | Capabilities | `Default()` | 4 | `./internal/caps/` |
| 5 | Init | `reap()` | 7 | `./internal/initproc/` |

<div v-click class="pt-4 opacity-80">
The tests are the finished ones. Stuck? The answer is the same file two directories up:
<code>cp ../../internal/rootfs/rootfs.go internal/rootfs/rootfs.go</code>, and carry on.
</div>

<!--
**S04 · PLAN · 4:00**

Read the "you write" column. Five functions. Everything else in the
repository is scaffolding and it is deliberately boring.

The click is the rule that keeps a 45-minute workshop moving: nobody sits
stuck. Copying the answer is allowed, and saying so now stops people
falling silently behind at step 2.
-->

---
layout: default
---

# Two processes, always

<<< @/snippets/clone.go go

<div v-click class="pt-2 opacity-80 text-base">
Namespaces are set at <code>clone(2)</code> and apply to the child. The Go runtime is several
threads before <code>main</code> runs. So <code>mycontainer run</code> starts <code>mycontainer child</code>
through <code>/proc/self/exe</code>, and the child is the container. <code>runc init</code> is the same trick.
</div>

<!--
**S05 · CODE · 5:00**

Already written, in `internal/container/container_linux.go`. The parent
outside, the child as PID 1 inside, a spec over a pipe between them. Say
it once: every function they write today runs in one or the other.

**CHECKPOINT** 6:00. Start step 1 here whatever the clock says.
-->

---
layout: section
chapter: "01"
---

::title::
Step 1. Namespaces: <span class="ac-accent">what it can see</span>

::default::
`internal/namespace/namespace.go`, then one line in `container_linux.go`. Five minutes.

<!--
**S06 · SECTION · 6:00**
-->

---
layout: default
---

# Seven namespaces, seven bits

| Flag | Isolates | Inside |
|---|---|---|
| `CLONE_NEWPID` | process IDs | the first process is PID 1 |
| `CLONE_NEWNS` | the mount table | a private `/proc`, and the pivot |
| `CLONE_NEWUTS` | hostname | `hostname` is the container id |
| `CLONE_NEWIPC` | shared memory, queues | the host's are invisible |
| `CLONE_NEWNET` | interfaces, routes, ports | a down `lo`, nothing else |
| `CLONE_NEWCGROUP` | the cgroup tree | `/proc/self/cgroup` reads `0::/` |
| `CLONE_NEWUSER` | UIDs and capabilities | root, mapped to you |

<!--
**S07 · TABLE · 6:10**

Thirty seconds. One bit each; the constants are at the top of the file
they are about to edit. `CloneFlags()` ORs together the ones that are on.
-->

---
layout: default
---

# Your turn: `CloneFlags()`, then wire it

```go
// internal/namespace/namespace.go
func (s Set) CloneFlags() uintptr {
	// TODO(step 1): one bit per namespace that is on.
}

// internal/container/container_linux.go, in sysProcAttr
// TODO(step 1): attr.Cloneflags from s.Namespaces.CloneFlags().
```

```bash
go test ./internal/namespace/ && make build
sudo ./mycontainer run -ns pid,uts -hostname box ubuntu \
  /bin/sh -c 'hostname; echo $$; ls -d /proc/[0-9]* | wc -l'
```

<div v-click class="pt-2 accent">
Hostname <code>box</code>, a small PID, and <code>/proc</code> still listing the host. Why?
</div>

<!--
**S08 · YOUR TURN · 6:40**

Four minutes. The test names the missing flag when one is wrong.

Before the wiring line is in, the runtime refuses `-hostname` with
"not in a new UTS namespace". That is deliberate: without the flag, the
call would rename their VM. If someone hits it, that is the hint.

Leave the click hanging while they type.
-->

---
layout: statement
---

# A PID namespace changes the numbers.<br/>`/proc` is what `ps` reads.

<div class="opacity-70 pt-6">
<code>/proc</code> is the host's until the container has its own. That is step 2.
</div>

<!--
**S09 · STATEMENT · 11:00**

The answer to the question on S08, in one line. It is the most common
mistake in this whole subject: people type `unshare -p`, run `ps`, and
wonder why nothing changed.

If anyone asks why the shell is PID 8 and not 2: the Go init's threads
took 2 to 7, because a thread takes an ID from the same namespace. It
comes back at step 3.

**CHECKPOINT** 11:30.
-->

---
layout: section
chapter: "02"
---

::title::
Step 2. Rootfs: <span class="ac-accent">where it lives</span>

::default::
`internal/rootfs/rootfs.go`. Ten minutes, and the order is the lesson.

<!--
**S10 · SECTION · 11:30**
-->

---
layout: default
---

# `pivot_root`, and the order that works

<div class="text-base space-y-1 pt-1">

**1.** `mount("", "/", MS_REC|MS_PRIVATE)`. systemd shares every mount; without this the last step unmounts the host's root.

**2.** `mount(rootfs, rootfs, MS_BIND|MS_REC)`. `pivot_root` wants a mount point. A directory is not one.

**3.** `proc` on `rootfs/proc`, from inside the PID namespace. `tmpfs` on `rootfs/dev`, then `populateDev`.

**4.** `sysfs` on `rootfs/sys`, `cgroup2` on `rootfs/sys/fs/cgroup`. Both optional.

**5.** `chdir(rootfs)`, `pivot_root(".", putOld)`, `chdir("/")`.

**6.** `umount2(putOld, MNT_DETACH)`, then `remove(putOld)`.

</div>

<div v-click class="pt-3 accent">
The test records every call and compares the sequence. It tells you which one is wrong, and how.
</div>

<!--
**S11 · SEQUENCE · 11:40**

Ninety seconds, with a beat on 1 and on 5. Step 1: systemd makes every
mount shared, so unmounting the old root inside the namespace propagates
out and unmounts the host's; people have done this to their laptops. Step
5 is why this is not chroot: chroot leaves the old root reachable to
anything holding a file descriptor outside it; `pivot_root` swaps the
mount and the old one is detached.
-->

---
layout: default
---

# Your turn: the sequence

```go
// internal/rootfs/rootfs.go, in Pivot
steps := []struct {
	name     string
	do       func() error
	optional bool
}{
	// TODO(step 2): the sequence, in order. Each step is a mount or a
	// syscall through m. rootfs_test.go pins the exact order and flags.
}
```

```bash
go test ./internal/rootfs/ && make build
touch /tmp/host-was-here
sudo ./mycontainer run ubuntu /bin/sh -c 'ps; ls /; ls /tmp; mount'
```

<div v-click class="pt-2 accent">
Two processes. The image's tree. An empty <code>/tmp</code>. Nothing from the host in <code>mount</code>.
</div>

<!--
**S12 · YOUR TURN · 13:10**

Eight minutes; the longest step. The test prints the call it got against
the call it wanted, flags in hex, so it earns its keep here. Two things
people get wrong: forgetting `MS_REC` on the private remount (the test
says so) and calling `pivot_root` from outside the new root (EINVAL at
runtime; say it out loud at minute four).

At minute six, remind the room about the `cp` escape hatch from S04.

**CHECKPOINT** 21:00. Past 22:00, do step 3 together on the screen rather
than as an exercise, and give the time to step 5.
-->

---
layout: section
chapter: "03"
---

::title::
Step 3. cgroups: <span class="ac-accent">what it can use</span>

::default::
`internal/cgroup/cgroup.go`, then two lines in `container_linux.go`. Eight minutes.

<!--
**S13 · SECTION · 21:00**
-->

---
layout: default
---

# cgroup v2 is a filesystem

```
/sys/fs/cgroup/mycontainer/<id>/
  memory.max        67108864          bytes, then the OOM killer
  memory.swap.max   0                 or memory.max is a suggestion
  cpu.max           50000 100000      quota and period, µs: half a core
  pids.max          20                fork returns EAGAIN past twenty
```

<div v-click class="pt-4 opacity-80 text-base">
<code>Create</code> makes the directory and enables the controllers. You write <code>Files()</code>, which
decides what goes in it, and <code>cpuMax</code>. Then <code>UseCgroupFD</code> and <code>CgroupFD</code>: the child is
born inside its cgroup (<code>CLONE_INTO_CGROUP</code>), so the first fork of a fork bomb already counts.
</div>

<!--
**S14 · CGROUPS · 21:10**

One minute. The line that bites: without `memory.swap.max` a process at
`memory.max` keeps going into swap and the limit is a suggestion; Docker
sets it for you.
-->

---
layout: default
---

# Your turn: `Files()`, `cpuMax()`, and the fd

```go
// internal/cgroup/cgroup.go
func (l Limits) Files() map[string]string {
	// TODO(step 3): memory.max (and memory.swap.max), cpu.max via cpuMax, pids.max.
}
func cpuMax(cpus float64) string {
	// TODO(step 3): "<quota> <period>", period 100000; zero or less is "max".
}
// container_linux.go, when group != nil: attr.UseCgroupFD = true; attr.CgroupFD = group.FD
```

```bash
go test ./internal/cgroup/ && make build
sudo ./mycontainer run -pids 20 ubuntu \
  /bin/sh -c 'cat /sys/fs/cgroup/pids.max; for i in $(seq 40); do sleep 2 & done'
```

<div v-click class="pt-2 warn">
Keep <code>-pids</code> above sixteen. <code>pids.max</code> counts the Go init's own threads, and below that it dies before the bomb does.
</div>

<!--
**S15 · YOUR TURN · 22:10**

Six minutes. Expected: `20`, then `/bin/sh: 0: Cannot fork`.

The click is the warning, and it is the PID 8 question from step 1
answered: the init's threads are processes as far as `pids.max` is
concerned. tini has one thread; that is the trade for writing this in Go.
Someone will try `-pids 8` and see "failed to create new OS thread".
Good. They found it first.

Anyone done early: `-memory 64m` and a 100 MB string, exit 137.

**CHECKPOINT** 29:00.
-->

---
layout: section
chapter: "04"
---

::title::
Step 4. Capabilities: <span class="ac-accent">what it can do</span>

::default::
`internal/caps/caps.go`. Four minutes, and the shortest function of the day.

<!--
**S16 · SECTION · 29:00**
-->

---
layout: default
---

# Root is fourteen bits

<div class="text-base pt-1">
The kernel does not ask "is this UID 0". It asks whether the process holds the capability.
<code>mount(2)</code> wants <code>CAP_SYS_ADMIN</code>, and Docker's default fourteen do not include it.
</div>

```
CAP_CHOWN CAP_DAC_OVERRIDE CAP_FOWNER CAP_FSETID CAP_KILL CAP_SETGID
CAP_SETUID CAP_SETPCAP CAP_NET_BIND_SERVICE CAP_NET_RAW CAP_SYS_CHROOT
CAP_MKNOD CAP_AUDIT_WRITE CAP_SETFCAP                       = a80425fb
```

<<< @/snippets/capbset-drop.go go

<!--
**S17 · CAPS · 29:10**

The fourteen are the constants at the top of `caps.go`, in that order;
`Default()` adds each to a `Set`. The drop below is already written: one
`prctl(PR_CAPBSET_DROP)` per capability not kept. The bounding set is the
ceiling, so dropping it once covers `sh` and everything `sh` runs.
-->

---
layout: default
---

# Your turn: `Default()`

```go
// internal/caps/caps.go
func Default() Set {
	// TODO(step 4): the fourteen Docker grants by default.
	var s Set
	return s
}
```

```bash
go test ./internal/caps/ ./internal/cli/ && make build
sudo ./mycontainer run ubuntu \
  /bin/sh -c 'grep CapBnd /proc/self/status; mount -t tmpfs none /tmp'
sudo ./mycontainer run -caps default,+CAP_SYS_ADMIN ubuntu \
  /bin/sh -c 'mount -t tmpfs none /tmp && echo MOUNTED'
```

<div v-click class="pt-2 accent">
<code>a80425fb</code> and a refusal. Then <code>a82425fb</code> and <code>MOUNTED</code>. Bit 21.
</div>

<!--
**S18 · YOUR TURN · 30:00**

Three minutes. `cli` goes green here too.

**SAY** when the first refusals come up: "mount says permission denied,
to a process that id says is root. Root just does not mean mount any
more."
-->

---
layout: statement
class: warn
---

# Capabilities belong to a thread.<br/>Go moves goroutines between threads.

<div class="opacity-70 pt-6 text-white">
Drop on one thread, fork the command from another, and the command starts with everything.<br/>
<code>runtime.LockOSThread()</code>, before the drop. Already in <code>Child</code>; read it.
</div>

<!--
**S19 · WARN · 33:00**

Thirty seconds. The Go trap: it happened while writing this, only on the
rootless path, only sometimes. `Child` locks itself to one OS thread before
it does anything. If there is one Go-specific thing to take home, it is
this.
-->

---
layout: section
chapter: "05"
---

::title::
Step 5. Init: <span class="ac-accent">what PID 1 owes</span>

::default::
`internal/initproc/initproc_linux.go`. Seven minutes, and the one people forget.

<!--
**S20 · SECTION · 33:30**
-->

---
layout: default
---

# Two jobs

<div class="grid grid-cols-2 gap-8 pt-4">
<div class="border-l-2 b-accent pl-4">
<div class="accent font-mono text-sm pb-2">REAPING</div>
<div class="text-lg">
When a parent exits before its children, the kernel re-parents them to PID 1.
Only PID 1 can then <code>wait</code> on them. A PID 1 that never does
collects zombies for as long as it runs.
</div>
</div>
<div class="border-l-2 b-accent pl-4">
<div class="accent font-mono text-sm pb-2">SIGNALS</div>
<div class="text-lg">
The kernel drops a default-disposition signal sent to the init of a PID
namespace. That is why <code>docker stop</code> waits ten seconds and sends SIGKILL.
The forwarding is written. The reaping is yours.
</div>
</div>
</div>

<!--
**S21 · TWO JOBS · 33:40**

Forty seconds. This is why Docker has `--init` and tini exists; most
programs were never written to be PID 1.
-->

---
layout: default
---

# Your turn: `reap()`

```go
// internal/initproc/initproc_linux.go
func reap(main int) (code int, done bool, err error) {
	// TODO(step 5): this waits for our own child and nobody else. Loop on
	// wait4(-1, WNOHANG) instead: any child, until there is nothing left.
	var ws syscall.WaitStatus
	if _, err := syscall.Wait4(main, &ws, 0, nil); err != nil {
		return 0, false, err
	}
	return ExitCode(ws), true, nil
}
```

```bash
go test ./internal/initproc/
```

<div v-click class="pt-2 accent">
Green means no zombie left behind after an orphan exits. Nothing on this list is red any more.
</div>

<!--
**S22 · YOUR TURN · 34:20**

Five minutes. The shape: loop; `wait4(-1, &ws, WNOHANG)`; retry on
`EINTR`; `ECHILD` means nothing left, return; `pid == 0` means children
exist and none exited, return; `pid == main` means remember its exit code
and that we are done; anything else was an orphan, reaped and forgotten.

**WATCH FOR** a loop without `WNOHANG`: it hangs when children exist and
none has exited, and the test times out rather than fails. Say it at
minute three.

**CHECKPOINT** 39:30. Go to the wrap whatever the room's state; the test
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
1 clone3 · 1 sethostname · 12 mount · 1 pivot_root · 1 umount2 · 27 prctl · 3 execve
</div>

<!--
**S23 · LIVE · 39:30**

**DO**, in the slide's terminal (the finished runtime):
```bash
./scripts/demo.sh 7
```

One container under `strace -f`, tallied.

**SAY**, from the tally on the click:
> "One clone3 carrying every namespace flag and the cgroup fd: steps 1
> and 3. Twelve mounts and one pivot_root: step 2. Twenty-seven prctl:
> step 4. And an init that reaps: step 5. You wrote each of those this
> afternoon. That is a container. It always was."
-->

---
layout: default
---

# Five things to take home

<v-clicks>

**1.** A container is a process wearing a costume: seven namespace bits, a pivoted mount, a cgroup directory, fourteen capability bits, and an init.

**2.** A PID namespace changes the numbers. `/proc` is what `ps` reads.

**3.** Root is a bounding set. Drop one bit and root cannot mount.

**4.** PID 1 owes the tree reaping and signals. Bring an init.

**5.** In Go: fork and exec through `/proc/self/exe`, lock the thread before you drop capabilities, and count your own threads against `pids.max`.

</v-clicks>

<div v-click class="pt-6 accent">
At home: three lines of UID mapping make it rootless. <code>docs/workshop.md</code> has them.
<code>img2rootfs alpine:3.22</code> gives it a real image.
</div>

<!--
**S24 · TAKEAWAYS · 42:00**

Read all five. Leave this up for questions and for the people still
typing.

The last click is the bonus they did not have time for: `sysProcAttr`
with one `UidMappings` and one `GidMappings` line, and
`./mycontainer run` works without sudo, root inside and them outside.
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
**S25 · END · 45:00**

Leave S24 up for questions. This slide is for the photograph.

Questions that come up, in order of likelihood:

1. **My test passes and the container does not run.** The pure packages
   test on macOS; the container needs the VM.
2. **Why not `unshare(2)` from Go?** A PID namespace only applies to
   children, and a user namespace needs a single-threaded process. Fork
   and exec is what every runtime does.
3. **Where is the network?** A namespace with a down loopback. Wiring a
   veth pair, a bridge and NAT is its own workshop.
4. **How does this compare to runc?** Same syscalls in the same order.
   runc adds seccomp, AppArmor, hooks, the OCI config and years of edge
   cases. `libcontainer/rootfs_linux.go` will look familiar after today.
-->
