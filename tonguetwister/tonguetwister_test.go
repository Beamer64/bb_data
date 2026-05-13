package tonguetwister

import (
	"testing"
	"testing/fstest"
)

func TestLoadAndRandom(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/tongue_twisters.txt": {Data: []byte("alpha\nbravo\ncharlie\n")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		seen[Random()] = true
	}
	for _, want := range []string{"alpha", "bravo", "charlie"} {
		if !seen[want] {
			t.Errorf("Random never returned %q across 200 draws", want)
		}
	}
}

func TestLoadMissingFile(t *testing.T) {
	if err := Load(fstest.MapFS{}); err == nil {
		t.Fatal("expected error when datasets/tongue_twisters.txt absent, got nil")
	}
}

func TestRandomEmptyPool(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/tongue_twisters.txt": {Data: []byte("")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "" {
		t.Errorf("Random on empty pool: got %q, want \"\"", got)
	}
}
