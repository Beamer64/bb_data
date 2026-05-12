package roasts

import (
	"testing"
	"testing/fstest"
)

func TestLoadAndRandom(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/roasts.txt": {Data: []byte("alpha\nbeta\ngamma\n")},
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

func TestLoadStripsCRLF(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/roasts.txt": {Data: []byte("alpha\r\nbeta\r\n")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, got := range pool {
		if got == "alpha\r" || got == "beta\r" {
			t.Errorf("CRLF not stripped from %q", got)
		}
	}
}

func TestLoadSkipsBlankLines(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/roasts.txt": {Data: []byte("alpha\n\n\nbeta\n")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := len(pool); got != 2 {
		t.Errorf("blank lines not skipped: got %d entries, want 2", got)
	}
}

func TestLoadMissingFile(t *testing.T) {
	if err := Load(fstest.MapFS{}); err == nil {
		t.Fatal("expected error when datasets/roasts.txt absent, got nil")
	}
}

func TestRandomEmptyPool(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/roasts.txt": {Data: []byte("")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "" {
		t.Errorf("Random on empty pool: got %q, want \"\"", got)
	}
}
