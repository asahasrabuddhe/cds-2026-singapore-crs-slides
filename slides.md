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

**1.** The Multipass VM from `docs/setup.md`. On your machine:

```bash
git clone https://github.com/asahasrabuddhe/cds-2026-singapore-crs-code.git crs
multipass mount crs crs:/home/ubuntu/crs
```

**2.** In the VM, `multipass shell crs`, then:

```bash
cd ~/crs && ./scripts/check.sh && make build && go test ./...
```

**3.** Your editor on the `crs` folder, on your machine.

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
Stuck? The answer is on the <code>final</code> branch. <code>git restore --source origin/final internal/rootfs/rootfs.go</code>, and carry on.
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
layout: code-split
lang: go
---

## Namespaces are bits on a clone

A namespace is a flag on `clone(2)`. It applies to the new process, never to the one asking.

<div v-click="1">The kernel has seven: PID, mount, UTS, IPC, network, cgroup and user. One bit each.</div>

<div v-click="2">The container has to set itself up from inside, so the runtime starts itself again as the child.</div>

<div v-click="2" class="accent">Your job: turn the spec's fields into those bits, and hand them to the clone.</div>

::right::

````md magic-move {at:1}
```go
// A process.
cmd := exec.Command("/bin/sh")
cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
cmd.Run()
```
```go
// A process in new namespaces.
cmd := exec.Command("/bin/sh")
cmd.SysProcAttr = &syscall.SysProcAttr{
	Cloneflags: syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWPID,
}
cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
cmd.Run()
```
```go
// The runtime, started again inside them.
cmd := exec.Command("/proc/self/exe", "child")
cmd.SysProcAttr = &syscall.SysProcAttr{
	Cloneflags: syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWPID,
}
cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
cmd.Run()

// In the child, PID 1 in the new namespace:
// the hostname, the rootfs, the limits, then
// the command.
```
````

<!--
**S04 · NAMESPACES · 4:00, one minute**

**SAY**: "A namespace is one bit on a clone, and it is the child that
gets it." Click once per sentence on the left; the code follows.

The last state is the shape of the whole runtime: `run` re-executes
itself through `/proc/self/exe` as `child`, and the child does
everything the next four steps add. Then S05 and "go".
-->

---
layout: step
leftLabel: DO THIS
---

::title::
Step 1 · Namespaces: what it can see

::default::

<div class="steps text-base">

1. Open `internal/namespace/namespace.go` and find `TODO(step 1)`.
2. The PID flag is written for you. Add the other six the same way.
3. Run `go test ./internal/namespace/` until it passes.
4. Open `internal/container/container_linux.go`, find `TODO(step 1)`, and add one line.

</div>

::right::

<div class="text-base space-y-4">

<div v-click="1">

Item 2, once more for each field:

```go
if s.Mount {
	f |= syscall.CLONE_NEWNS
}
```

Mount is the odd one: <code>CLONE_NEWNS</code>. The rest are their own names.

</div>

<div v-click="2">

Item 4 is this line:

```go
attr.Cloneflags = s.Namespaces.CloneFlags()
```

</div>

</div>

::run::

```bash
make build && sudo ./mycontainer run -ns pid,uts -hostname box ubuntu /bin/sh
hostname; echo $$; ls /proc          # type these inside
```

<div v-click="3" class="accent">
<code>box</code>, a small PID, and <code>/proc</code> still lists the VM's processes. Why?
</div>

<!--
**S05 · STEP 1 · 5:00, four minutes**

S04 was the idea. Here, "go", and walk the room.

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
layout: code-split
lang: mounts
---

## The root is a directory, swapped in

The container's `/` is a directory on the VM's disk, swapped in by `pivot_root(2)`. The order is the lesson.

<div v-click="1">Make <code>/</code> private, so nothing leaks out. Bind the image onto itself: <code>pivot_root</code> wants a mount point.</div>

<div v-click="2">Mount <code>/proc</code>, <code>/dev</code> and <code>/sys</code> inside it.</div>

<div v-click="3">Step in and pivot. The VM's root hangs off <code>/.old_root</code>.</div>

<div v-click="4" class="accent">Detach it and remove it. Nothing points out.</div>

::right::

````md magic-move {at:1}
```text
/                       the VM's root, shared
├── proc                the VM's processes
├── tmp/host-was-here
└── …/rootfs/ubuntu     a directory
```
```text
/                       the VM's root, private
├── proc                the VM's processes
├── tmp/host-was-here
└── …/rootfs/ubuntu     a mount point
```
```text
/                       the VM's root, private
├── proc                the VM's processes
├── tmp/host-was-here
└── …/rootfs/ubuntu     a mount point
    ├── proc            the container's processes
    ├── dev             a tmpfs, a few devices
    └── sys             read-only
```
```text
/                       the image
├── proc                the container's processes
├── dev                 a tmpfs, a few devices
├── sys                 read-only
└── .old_root           the VM's root, still there
    └── tmp/host-was-here
```
```text
/                       the image
├── proc                the container's processes
├── dev                 a tmpfs, a few devices
├── sys                 read-only
└── tmp                 empty
```
````

<!--
**S06 · ROOTFS · 9:00, one minute**

**SAY**, first: "`ps` reads `/proc`, and `/proc` was the host's. That is
step 1's answer." Then one click per picture.

On the first click, the warning: systemd makes every mount shared, so
without the private mount, unmounting the old root inside the namespace
unmounts the VM's. People have done this to their laptops.

These five pictures are the seven steps in `Pivot`: private, bind,
three mounts, chdir and pivot, detach and remove. Then S07 and "go".
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
		syscall.MS_REC|syscall.MS_PRIVATE, "")
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
**S07 · STEP 2 · 10:00, nine minutes**

S06 was the idea and the warning about shared mounts. Here, "go".

Click 1 at minute two, click 2 at minute four. At minute six, say the
`cp` escape hatch out loud again.

If anyone asks why `mount` names the VM's disk for `/`, such as
`/dev/sda1 on / type ext4`: the root is a bind mount of
`/var/lib/mycontainer/rootfs/ubuntu`, which lives on that disk, and the device nodes in `/dev` are bind mounts
of the VM's. What is gone is every other mount the VM has.

**CHECKPOINT** 19:00. Past 20:00, do step 3 together on the screen and
give its time to step 5.
-->

---
layout: code-split
lang: shell
---

## A limit is a number in a file

cgroup v2 is a filesystem. A group is a directory, and the kernel fills it with files.

<div v-click="1">Write a number to a file and the kernel enforces it. Without the swap file, the memory limit is a suggestion.</div>

<div v-click="2">A process joins by writing its PID. The runtime's child is born inside, so its first fork counts.</div>

<div v-click="2" class="accent">Your job: which file gets which number, and the flag that puts the child inside.</div>

::right::

````md magic-move {at:1}
```bash
$ cd /sys/fs/cgroup && mkdir box && cd box
$ ls
cgroup.procs  cpu.max  memory.max  pids.max  …
```
```bash
$ cd /sys/fs/cgroup && mkdir box && cd box
$ echo 20 > pids.max
$ echo 67108864 > memory.max       # 64 MiB
$ echo 0 > memory.swap.max
$ echo "50000 100000" > cpu.max    # half a CPU
```
```bash
$ cd /sys/fs/cgroup && mkdir box && cd box
$ echo 20 > pids.max
$ echo 67108864 > memory.max       # 64 MiB
$ echo 0 > memory.swap.max
$ echo "50000 100000" > cpu.max    # half a CPU
$ echo $$ > cgroup.procs
$ for i in $(seq 40); do sleep 9 & done
# fork fails once the group holds 20
```
````

<!--
**S08 · CGROUPS · 19:00, one minute**

**SAY**: "cgroup v2 is a filesystem." One click per sentence.

`cpu.max` is a quota and a period in microseconds, so half a CPU is
`50000 100000`. The flag on the last click is `CLONE_INTO_CGROUP`, which
Go spells `UseCgroupFD`.

The shell is the idea, not the task: the runtime writes these files from
Go, and the test checks the names and the numbers. Then S09 and "go".
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
**S09 · STEP 3 · 20:00, six minutes**

S08 was the idea. Here, "go". The line that bites, if anyone asks: without
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
layout: code-split
lang: /proc/self/status
---

## Root is a set of bits

The kernel does not ask whether you are UID 0. It asks whether the process holds the capability for what it is doing.

<div v-click="1">Docker keeps fourteen. <code>CAP_SYS_ADMIN</code>, bit 21, is not one of them, and <code>mount(2)</code> wants it.</div>

<div v-click="2">So root in the container cannot mount.</div>

<div v-click="2" class="accent">Your job: the fourteen. Dropping the rest is written.</div>

::right::

````md magic-move {at:1}
```text
# root on the VM
CapBnd: 000001ffffffffff    bits 0 to 40, all of them
```
```text
# root in the container
CapBnd: 00000000a80425fb    fourteen bits

 0 chown          7 setuid             27 mknod
 1 dac_override   8 setpcap            29 audit_write
 3 fowner        10 net_bind_service   31 setfcap
 4 fsetid        13 net_raw
 5 kill          18 sys_chroot
 6 setgid
```
```text
# root in the container
CapBnd: 00000000a80425fb    fourteen bits

 0 chown          7 setuid             27 mknod
 1 dac_override   8 setpcap            29 audit_write
 3 fowner        10 net_bind_service   31 setfcap
 4 fsetid        13 net_raw
 5 kill          18 sys_chroot
 6 setgid        21 sys_admin  ✗

# mount -t tmpfs none /tmp
mount: /tmp: permission denied.
```
````

<!--
**S10 · CAPABILITIES · 26:00, one minute**

**SAY**: "Root is not a number. It is a set of bits." One click per
sentence.

Every bit in `a80425fb` is one of the fourteen, which is why the check
on S11 reads that value. Then S11, the thread warning, and "go".
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
**S11 · STEP 4 · 27:00, four minutes**

S10 was the idea. **SAY** the warning, then "go". It happened while
writing this, only on the rootless path, only sometimes: drop on one
thread, fork the command from another, and the command starts with
everything. If there is one Go-specific thing to take home, it is this.

`cli` goes green here too.

**SAY** when the first refusals come up: "mount says permission denied,
to a process that id says is root."

**CHECKPOINT** 31:00.
-->

---
layout: code-split
lang: ps
---

## PID 1 waits for everyone

When a parent exits before its children, the kernel hands them to PID 1.

<div v-click="1">A child that exits stays in the table, a zombie, until its parent waits for it.</div>

<div v-click="2">Only PID 1 can wait for an orphan. One that waits for its own child and nobody else collects zombies.</div>

<div v-click="3" class="accent">Your job: wait for any child, as often as one exits, and never block.</div>

::right::

````md magic-move {at:1}
```text
  PID  PPID  STAT  CMD
    1     0  S     mycontainer child
    8     1  S     sh
    9     8  S     sh -c 'sleep 2 &'
   10     9  S     sleep 2
```
```text
  PID  PPID  STAT  CMD
    1     0  S     mycontainer child
    8     1  S     sh
   10     1  S     sleep 2

# 9 exited. Its child is PID 1's now.
```
```text
  PID  PPID  STAT  CMD
    1     0  S     mycontainer child
    8     1  S     sh
   10     1  Z     [sleep] <defunct>

# 10 exited. Nobody has waited for it.
```
```text
  PID  PPID  STAT  CMD
    1     0  S     mycontainer child
    8     1  S     sh

# PID 1 called wait4(-1, …, WNOHANG).
```
````

<!--
**S12 · INIT · 31:00, one minute**

**SAY**: "When a parent exits first, its children go to PID 1." One click
per picture.

This is why Docker has `--init` and tini exists: most programs were
never written to be PID 1. Signals are written for them; S13 says which
`wait4` results mean what. Then "go".
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
**S13 · STEP 5 · 32:00, six minutes**

S12 was the idea. Here, "go".

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
**S14 · LIVE · 38:00**

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
**S15 · TAKE HOME · 40:00**

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
**S16 · END · 45:00**

For the photograph. Questions stay on S15.
-->
