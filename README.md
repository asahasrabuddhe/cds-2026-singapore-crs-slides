# Build Your Own Container Runtime in Go

Slides for **Build Your Own Container Runtime in Go**, the forty-five-minute
workshop at ContainerDays Singapore 2026.

Built with [Slidev](https://sli.dev) and the
[`slidev-theme-architectural-console`](https://github.com/asahasrabuddhe/slidev-theme-architectural-console)
theme. The workshop's code lives in
[`cds-2026-singapore-crs-code`](https://github.com/asahasrabuddhe/cds-2026-singapore-crs-code).
Attendees work in that repository's `workshop/start/`; the deck is what
is on the screen while they do.

## Setup

```bash
git clone https://github.com/asahasrabuddhe/cds-2026-singapore-crs-slides.git
cd cds-2026-singapore-crs-slides
npm ci
npm run dev
```

Node 24, pinned in `.nvmrc`. `@slidev/cli` is pinned to 52.14.1; bump it deliberately.

## The shape of the deck

Eleven slides, and most of the forty-five minutes is the room writing
code. Each of the five steps is one slide: the idea and the task on the
left, with the test and the command that shows it working, and hints on
the right that get more specific with each click. The presenter talks
for a minute, then walks the room. There is one live demo, the strace
tally at the end.

| Time | Slide | Minutes |
|---|---|---|
| before | S01 cover, S02 setup, up while the room fills | |
| 0:00 | S02 setup check, S03 the plan | 4 |
| 4:00 | S04 Step 1, namespaces | 5 |
| 9:00 | S05 Step 2, rootfs | 10 |
| 19:00 | S06 Step 3, cgroups | 7 |
| 26:00 | S07 Step 4, capabilities | 5 |
| 31:00 | S08 Step 5, init | 7 |
| 38:00 | S09 the strace tally, S10 take home, questions | 7 |

The checkpoint that matters is 19:00, leaving step 2. Past 20:00, do step
3 together on the screen rather than as an exercise and give its time to
step 5, which is the one people have not seen before. Rootless is
take-home; the last click on S10 points at it.

## The terminal

S09 carries a live terminal into the presenter's Multipass VM, `crs`,
landing in the finished runtime at `~/crs`, the VM's own clone of the
code repository, so `./scripts/demo.sh 7` works as typed.
`scripts/dev-shell.sh` opens it with `multipass exec`. Click the terminal to type; `r` restarts a dead session.

## Presenter notes

Every slide has them, keyed by slide number with the cumulative time,
then `SAY`, `DO` and `WATCH FOR`. Press `p` for presenter mode.
[`docs/EVENT-DAY.md`](docs/EVENT-DAY.md) is the runbook for the day.

## Rendering

Checked at 16:9, 1280x720, across all 25 slides and every click state.

## Licence

CC BY-SA 4.0. See [`LICENSE`](LICENSE).
