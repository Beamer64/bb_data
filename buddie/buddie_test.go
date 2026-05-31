package buddie

import (
	"testing"
	"testing/fstest"
)

func TestLoadAndRandom(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/buddie/(1).jpg":  &fstest.MapFile{Data: []byte("first")},
		"datasets/buddie/(2).jpeg": &fstest.MapFile{Data: []byte("second")},
		"datasets/buddie/(3).png":  &fstest.MapFile{Data: []byte("third")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got := Count(); got != 3 {
		t.Errorf("Count = %d, want 3", got)
	}

	// Random should return one of the embedded entries, with non-empty fields.
	got := Random()
	if got.Filename == "" || len(got.Data) == 0 {
		t.Errorf("Random returned zero-ish Image: %+v", got)
	}
	known := map[string]bool{"(1).jpg": true, "(2).jpeg": true, "(3).png": true}
	if !known[got.Filename] {
		t.Errorf("Random returned unknown filename %q", got.Filename)
	}
}

func TestLoadIgnoresSubdirectories(t *testing.T) {
	// Files nested under a subdir should NOT be picked up — the loader
	// only reads the top level of datasets/buddie.
	fsys := fstest.MapFS{
		"datasets/buddie/visible.jpg":             &fstest.MapFile{Data: []byte("yes")},
		"datasets/buddie/nested/hidden.jpg":       &fstest.MapFile{Data: []byte("no")},
		"datasets/buddie/deeper/also/hidden.jpeg": &fstest.MapFile{Data: []byte("no")},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := Count(); got != 1 {
		t.Errorf("Count = %d, want 1 (subdirs should be skipped)", got)
	}
	if got := Random(); got.Filename != "visible.jpg" {
		t.Errorf("Random = %q, want visible.jpg", got.Filename)
	}
}

func TestRandomBeforeLoadReturnsZero(t *testing.T) {
	// Reset the package-level pool so this test runs independent of order.
	pool = nil
	got := Random()
	if got.Filename != "" || got.Data != nil {
		t.Errorf("Random with no Load should return zero Image, got %+v", got)
	}
	if Count() != 0 {
		t.Errorf("Count with no Load should be 0, got %d", Count())
	}
}

func TestLoadMissingDirErrors(t *testing.T) {
	// An empty FS has no datasets/buddie directory — Load should surface the
	// underlying fs error rather than silently leaving pool empty.
	if err := Load(fstest.MapFS{}); err == nil {
		t.Error("expected error from Load against empty FS, got nil")
	}
}
