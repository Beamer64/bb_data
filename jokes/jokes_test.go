package jokes

import (
	"testing"
	"testing/fstest"
)

func TestLoadAndRandom(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/shortjokes.json": {Data: []byte(`{"data":[
			{"id":"1","joke":"alpha"},
			{"id":"2","joke":"beta"},
			{"id":"3","joke":"gamma"}
		]}`)},
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

func TestLoadSkipsEmptyJokes(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/shortjokes.json": {Data: []byte(`{"data":[
			{"id":"1","joke":"keeper"},
			{"id":"2","joke":""},
			{"id":"3","joke":"   "}
		]}`)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := len(pool); got != 1 {
		t.Errorf("empty/whitespace jokes not skipped: got %d entries, want 1; pool=%v", got, pool)
	}
}

func TestLoadTrimsWhitespace(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/shortjokes.json": {Data: []byte(`{"data":[
			{"id":"1","joke":"  padded  "}
		]}`)},
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
		"datasets/shortjokes.json": {Data: []byte(`not json`)},
	}
	if err := Load(fsys); err == nil {
		t.Fatal("expected error from malformed JSON, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if err := Load(fstest.MapFS{}); err == nil {
		t.Fatal("expected error when datasets/shortjokes.json absent, got nil")
	}
}

func TestRandomEmptyPool(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/shortjokes.json": {Data: []byte(`{"data":[]}`)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "" {
		t.Errorf("Random on empty pool: got %q, want \"\"", got)
	}
}
