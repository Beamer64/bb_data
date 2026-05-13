In-process access to curated text datasets (jokes, facts, roasts, tongue twisters, 8-ball responses, and others), embedded into the Go binary at build time.

bb_data is a small Go module that ships curated text datasets — jokes, facts, roasts, Kanye quotes, etc — embedded directly into the binary via //go:embed. Each dataset lives in its own sub-package, so consumers import only what they use. 

The public surface is two calls: Load() once at startup, Random() thereafter — lock-free and safe for concurrent goroutines after warm-up. No network access, no filesystem reads at runtime. Internal helpers cover plain-text, JSON, and JSONL fixtures; test coverage uses testing/fstest so each dataset's load path is verified against in-memory fixtures.
