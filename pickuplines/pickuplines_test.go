package pickuplines

import (
	"testing"
	"testing/fstest"
)

func TestLoadAndRandom(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/pickuplines.json": {
			Data: []byte(`{"data":[
				{"category":"cheesy","joke":"alpha"},
				{"category":"cheesy","joke":"beta"},
				{"category":"smooth","joke":"gamma"}
			]}`),
		},
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

func TestLoadMalformedJSON(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/pickuplines.json": {Data: []byte(`not json`)},
	}
	if err := Load(fsys); err == nil {
		t.Fatal("expected error from malformed JSON, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if err := Load(fstest.MapFS{}); err == nil {
		t.Fatal("expected error when datasets/pickuplines.json absent, got nil")
	}
}

func TestRandomEmptyPool(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/pickuplines.json": {Data: []byte(`{"data":[]}`)},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Random(); got != "" {
		t.Errorf("Random on empty pool: got %q, want \"\"", got)
	}
}
