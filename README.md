# Build Your Own Container Runtime in Go

Slides for **Build Your Own Container Runtime in Go**, the forty-five-minute
workshop at ContainerDays Singapore 2026.

Built with [Slidev](https://sli.dev) and the private
`slidev-theme-architectural-console` theme. The workshop's code lives in a
sibling repository,
[`cds-2026-singapore-crs-code`](https://github.com/asahasrabuddhe/cds-2026-singapore-crs-code),
which the deck pulls its Go from by anchor comment. Attendees work in that
repository's `workshop/start/`; the deck is what is on the screen while
they do.

## Setup

```bash
git clone git@github.com:asahasrabuddhe/cds-2026-singapore-crs-slides.git
git clone git@github.com:asahasrabuddhe/cds-2026-singapore-crs-code.git   # as a sibling
cd cds-2026-singapore-crs-slides
cp .env.example .env && cp .npmrc.example .npmrc
npm run install:theme      # op run --env-file=.env -- npm install
npm run dev
```

The theme is a private npm package. `.env` holds a 1Password reference
and `op run` resolves it into the install's environment. Node 24, pinned
in `.nvmrc`. `@slidev/cli` is pinned to 52.14.1; bump it deliberately.

## The shape of the deck

Twenty-five slides. Each of the five steps is the same three or four
slides: a section title with the file and the minutes, one slide on what
the function has to do, and a "your turn" slide with the stub, the test
command and the checkpoint. The room types during the "your turn" slides.
There is one live demo, the strace tally at the end, because forty-five
minutes has no room for a demo after every step.

| Time | Step | Your turn | Minutes |
|---|---|---|---|
| 0:00 | Setup check, the plan, two processes | S02 | 3 |
| 6:00 | 1. Namespaces | S08 | 4 |
| 11:30 | 2. Rootfs | S12 | 8 |
| 21:00 | 3. cgroups | S15 | 6 |
| 29:00 | 4. Capabilities | S18 | 3 |
| 33:30 | 5. Init | S22 | 5 |
| 39:30 | Wrap: the strace tally, takeaways | S23 | |

The checkpoint that matters is 21:00, leaving step 2. Past 22:00, do step
3 together on the screen rather than as an exercise and give its time to
step 5, which is the one people have not seen before. Rootless is
take-home; the last click on S24 points at it.

## The snippets

Two slides show Go pulled from the code repository by anchor comment,
through `cmd/snippets`. `snippets/*.go` is committed so a fresh clone
builds; CI regenerates and fails on drift.

```bash
CRS_CODE_DIR=../cds-2026-singapore-crs-code npm run snippets
```

## The terminal

S23 carries a live terminal into the presenter's VM, landing in the
finished runtime at the root of the sibling code checkout, so
`./scripts/demo.sh 7` works as typed. `scripts/dev-shell.sh` opens it
through OrbStack. Click the terminal to type; `r` restarts a dead session.

## Presenter notes

Every slide has them, keyed by slide number with the cumulative time,
then `SAY`, `DO` and `WATCH FOR`. Press `p` for presenter mode.
[`docs/EVENT-DAY.md`](docs/EVENT-DAY.md) is the runbook for the day.

## Rendering

Checked at 16:9, 1280x720, across all 25 slides and every click state.

## Licence

CC BY-SA 4.0. See [`LICENSE`](LICENSE).
