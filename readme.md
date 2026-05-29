# bb_data

**13 curated datasets for Go programs** — jokes, roasts, facts, Kanye quotes, 8-ball responses, pickup lines, would-you-rathers, and more — embedded into the binary at build time, loaded once at startup, accessed through a tiny per-dataset API. Zero external dependencies.

```
go get github.com/Beamer64/bb_data
```

---

## What it is

bb_data is a small Go module that ships curated text datasets — jokes, facts, roasts, Kanye quotes, etc. — embedded directly into the binary via `//go:embed`. Each dataset lives in its own sub-package, so consumers import only what they use.

The public surface is two calls: **`Load()` once at startup, accessor thereafter** — lock-free and safe for concurrent goroutines after warm-up. No network access, no filesystem reads at runtime, no asset directories to ship alongside your binary.

It started life as the in-process content library behind a Discord bot — every `/get joke`, `/daily affirmation`, `/game wyr`, `/txt cursive`, etc. pulls from one of these subpackages — but the API is fully usable on its own for any Go program that wants embedded text content without a database.

## Quick start

The simplest pattern: load everything at startup, then call subpackage accessors freely.

```go
package main

import (
    "fmt"
    "log"

    "github.com/Beamer64/bb_data"
    "github.com/Beamer64/bb_data/jokes"
    "github.com/Beamer64/bb_data/affirmations"
)

func main() {
    if err := bb_data.Load(); err != nil {
        log.Fatalf("bb_data.Load: %v", err)
    }

    fmt.Println("Joke:", jokes.Random())
    fmt.Println("Affirmation:", affirmations.Random())
}
```

A runnable version is at [`examples/cli`](examples/cli/) — `go run ./examples/cli` prints a random joke + affirmation and reports how long the embedded datasets took to load.

If you only need one or two datasets, you can skip the umbrella `bb_data.Load()` and call the subpackage's own `Load(bb_data.FS)` directly — every subpackage loader takes any `fs.FS`, and `bb_data.FS` (the embedded root) is the natural source:

```go
import (
    "github.com/Beamer64/bb_data"
    "github.com/Beamer64/bb_data/jokes"
)

if err := jokes.Load(bb_data.FS); err != nil {
    log.Fatal(err)
}
```

Before `Load` (or if it errored), accessors return their zero value (`""` for strings, `Poll{}` for `wyr.Random`) — no panics, no required nil-checks at the call site.

## The catalogue

Most subpackages expose a single `Random()` returning a `string`. The two outliers are `wyr` (returns a `Poll` struct with two options + historical vote counts) and `textfonts` (`Convert` and `Groups`, since it's a transformation, not a pick).

| Package | What it serves | Accessor |
|---|---|---|
| [`affirmations`](affirmations/) | One-line affirmations | `Random() string` |
| [`eightball`](eightball/) | Magic 8-ball answers | `Random() string` |
| [`emojis`](emojis/) | Random emoji characters | `Random() string` |
| [`facts`](facts/) | Random "fun fact" entries | `Random() string` |
| [`jokes`](jokes/) | General-purpose short jokes | `Random() string` |
| [`kanye`](kanye/) | Kanye West quotes | `Random() string` |
| [`loadingmessages`](loadingmessages/) | Loading-screen one-liners | `Random() string` |
| [`pickuplines`](pickuplines/) | Pickup lines | `Random() string` |
| [`roasts`](roasts/) | Insult / roast lines | `Random() string` |
| [`tonguetwister`](tonguetwister/) | Tongue twisters | `Random() string` |
| [`yomomma`](yomomma/) | Yo Mama jokes | `Random() string` |
| [`wyr`](wyr/) | Would You Rather polls (two options + vote counts) | `Random() Poll`, `Count() int` |
| [`textfonts`](textfonts/) | ASCII → stylized Unicode (cursive, bubble, leet, flipped, cursed) | `Convert(text, group) string`, `Groups() []string` |

Shared helpers (random selection plus filesystem readers for plain-text, JSON, JSONL, and CSV fixtures) live under [`internal/pick`](internal/pick/) — not part of the public API, but worth a look if you're curious how the subpackages stay tiny.

## Design notes

A few things worth knowing if you're using or extending the library:

### Load-once pattern

Every subpackage exposes `Load(fs.FS) error` and accessors that are safe to call **after a successful `Load`** and zero-valued before. The umbrella `bb_data.Load()` chains every subpackage's loader in order — call it once at startup, check the error, then forget about it. After warm-up the accessors are lock-free; the in-memory pool is written exactly once and read concurrently from then on.

### Embedded, not shipped

The actual data files live under `datasets/` and are bundled into the binary via `//go:embed`. There are no runtime paths to configure, no asset directories to copy next to your binary, no flags to point at content. The embedded `fs.FS` is exposed as `bb_data.FS` if you want to use it as the source for your own custom loaders.

### Zero external dependencies

```go
module github.com/Beamer64/bb_data
go 1.26
```

That's the entire `go.mod` — no `require` directives. Pure standard library: `embed`, `encoding/json`, `encoding/csv`, `io/fs`, `math/rand/v2`, `strings`, `sort`. Smallest possible dependency footprint for a content library.

### Robust loaders

The shared `internal/pick.LoadJSONL` tolerates JSON Lines, concatenated JSON values with no separator, and arbitrary whitespace between values — useful because hand-curated `.jsonl` files occasionally lose a newline between records. `LoadLines` handles both LF and CRLF input and skips empty lines. `wyr.Load` logs and skips malformed CSV rows rather than aborting the whole load. The library tries hard not to be brittle about its own datasets.

### Empty-state safety

`Random()` returns `""` (or the zero value for non-string accessors) when its pool is empty — no panics, no required nil-checks at the call site. Lets you wire it into hot paths without defensive guards on every access.

### Tested against in-memory fixtures

Test coverage uses [`testing/fstest`](https://pkg.go.dev/testing/fstest) so every dataset's load path is verified against in-memory `fs.FS` fixtures — not the real embedded data. That means tests are deterministic regardless of the live dataset's contents, and dataset additions don't break existing test assertions.

## Dependencies

None. The whole library is `go build`able with just the Go toolchain. No CGO, no native libraries, no Docker, no API keys.

## Used by

- **[BuddieBot](https://github.com/Beamer64/BuddieBot)** — Discord bot; everything content-driven (`/get joke`, `/daily affirmation`, `/get yomomma`, `/get 8ball`, `/game wyr`, `/txt *` font transforms, …) pulls from this library.

## Contributing

Adding a new dataset is a tight recipe:

1. **Drop the data file** under `datasets/` (`.json`, `.jsonl`, `.csv`, or line-delimited `.txt` all supported by the helpers in `internal/pick`).
2. **Create a subpackage** at the repo root: `yournewset/yournewset.go` containing:
   - `var pool []T` — the in-memory dataset
   - `func Load(fsys fs.FS) error` — reads the file via `pick.LoadJSON` / `LoadJSONL` / `LoadLines` and populates `pool`
   - `func Random() T` — wraps `pick.Random(pool)` so empty-state behaviour matches the rest of the library
3. **Add a test** at `yournewset/yournewset_test.go` that loads from a small `fstest.MapFS` and asserts the accessor returns sane output. The existing `*_test.go` files are templates worth copying.
4. **Wire the loader into the umbrella `bb_data.Load()`** in `bb_data.go` so consumers picking up the whole library don't have to remember the new package.

Keep the public surface tiny (`Load` + one accessor is the norm). Anything cleverer than `Random()` belongs in a different package.

## License

MIT — see [LICENSE](LICENSE).
