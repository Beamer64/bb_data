package affirmations

import (
	"testing"
	"testing/fstest"
)

func TestLoadAndRandom(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/affirmations.jsonl": {Data: []byte(
			`{"prompt":"Category: safety","completion":"alpha"}` + "\n" +
				`{"prompt":"Category: focus","completion":"beta"}` + "\n" +
				`{"prompt":"Category: calm","completion":"gamma"}` + "\n",
		)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		seen[Random()] = true
	}
	for _, want := range []string{"alpha", "beta", "gamma"} {
		if !seen[want] {
			t.Errorf("Random never returned %q across 200 draws", want)
		}
	}
}

func TestLoadSkipsBlankLines(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/affirmations.jsonl": {Data: []byte(
			`{"completion":"alpha"}` + "\n\n   \n" +
				`{"completion":"beta"}` + "\n",
		)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := len(pool); got != 2 {
		t.Errorf("blank lines not skipped: got %d entries, want 2; pool=%v", got, pool)
	}
}

func TestLoadSkipsEmptyCompletions(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/affirmations.jsonl": {Data: []byte(
			`{"completion":"keeper"}` + "\n" +
				`{"completion":""}` + "\n" +
				`{"completion":"   "}` + "\n",
		)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := len(pool); got != 1 {
		t.Errorf("empty completions not skipped: got %d entries, want 1; pool=%v", got, pool)
	}
}

func TestLoadIgnoresUnknownFields(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/affirmations.jsonl": {Data: []byte(
			`{"prompt":"Category: safety","completion":"alpha","extra":42}` + "\n",
		)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "alpha" {
		t.Errorf("got %q, want %q", got, "alpha")
	}
}

func TestLoadMalformedJSONL(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/affirmations.jsonl": {Data: []byte(
			`{"completion":"ok"}` + "\n" + `not json` + "\n",
		)},
	}
	if err := Load(fsys); err == nil {
		t.Fatal("expected error from malformed JSONL line, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if err := Load(fstest.MapFS{}); err == nil {
		t.Fatal("expected error when datasets/affirmations.jsonl absent, got nil")
	}
}

func TestRandomEmptyPool(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/affirmations.jsonl": {Data: []byte("")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "" {
		t.Errorf("Random on empty pool: got %q, want \"\"", got)
	}
}
