package eightball

import (
	"testing"
	"testing/fstest"
)

func TestLoadAndRandom(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/8ball.txt": {Data: []byte("yes\nno\nmaybe\n")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		seen[Random()] = true
	}
	for _, want := range []string{"yes", "no", "maybe"} {
		if !seen[want] {
			t.Errorf("Random never returned %q across 200 draws", want)
		}
	}
}

func TestLoadMissingFile(t *testing.T) {
	if err := Load(fstest.MapFS{}); err == nil {
		t.Fatal("expected error when datasets/8ball.txt absent, got nil")
	}
}

func TestRandomEmptyPool(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/8ball.txt": {Data: []byte("")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "" {
		t.Errorf("Random on empty pool: got %q, want \"\"", got)
	}
}
