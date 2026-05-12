package kanye

import (
	"testing"
	"testing/fstest"
)

func TestLoadAndRandom(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/kanyequotes.json": {Data: []byte(`["alpha","beta","gamma"]`)},
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

func TestLoadSkipsEmptyAndWhitespace(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/kanyequotes.json": {Data: []byte(`["keeper","","   "]`)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := len(pool); got != 1 {
		t.Errorf("empty/whitespace entries not skipped: got %d, want 1; pool=%v", got, pool)
	}
}

func TestLoadTrimsWhitespace(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/kanyequotes.json": {Data: []byte(`["  padded  "]`)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "padded" {
		t.Errorf("got %q, want %q", got, "padded")
	}
}

func TestLoadMalformedJSON(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/kanyequotes.json": {Data: []byte(`not json`)},
	}
	if err := Load(fsys); err == nil {
		t.Fatal("expected error from malformed JSON, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if err := Load(fstest.MapFS{}); err == nil {
		t.Fatal("expected error when datasets/kanyequotes.json absent, got nil")
	}
}

func TestRandomEmptyPool(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/kanyequotes.json": {Data: []byte(`[]`)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "" {
		t.Errorf("Random on empty pool: got %q, want \"\"", got)
	}
}
