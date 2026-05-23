package textfonts

import (
	"reflect"
	"testing"
	"testing/fstest"
)

// loadFixture installs a minimal in-memory dataset for testing. Each
// character mapping has exactly one variant so Convert is deterministic
// across runs.
func loadFixture(t *testing.T) {
	t.Helper()
	fsys := fstest.MapFS{
		"datasets/text_fonts.json": &fstest.MapFile{
			Data: []byte(`{
				"cursive": [{
					"a": ["A1"],
					"b": ["B1"],
					"c": ["C1"]
				}],
				"bubble": [{
					"a": ["A2"],
					"b": ["B2"]
				}]
			}`),
		},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}
}

func TestConvert_knownGroupSubstitutesMappedChars(t *testing.T) {
	loadFixture(t)
	got := Convert("abc", "cursive")
	if got != "A1B1C1" {
		t.Errorf("Convert(abc, cursive) = %q, want %q", got, "A1B1C1")
	}
}

func TestConvert_unmappedCharsPassThrough(t *testing.T) {
	loadFixture(t)
	// "z" isn't in the bubble fixture; it should appear as-is.
	got := Convert("abz", "bubble")
	if got != "A2B2z" {
		t.Errorf("Convert(abz, bubble) = %q, want %q", got, "A2B2z")
	}
}

func TestConvert_unknownGroupReturnsTextUnchanged(t *testing.T) {
	loadFixture(t)
	const in = "hello"
	got := Convert(in, "nonexistent")
	if got != in {
		t.Errorf("Convert with unknown group = %q, want %q", got, in)
	}
}

func TestConvert_emptyText(t *testing.T) {
	loadFixture(t)
	if got := Convert("", "cursive"); got != "" {
		t.Errorf("Convert(\"\", cursive) = %q, want empty", got)
	}
}

func TestConvert_picksOneOfMultipleVariants(t *testing.T) {
	fsys := fstest.MapFS{
		"datasets/text_fonts.json": &fstest.MapFile{
			Data: []byte(`{"silly": [{"x": ["X1", "X2", "X3"]}]}`),
		},
	}
	if err := Load(fsys); err != nil {
		t.Fatalf("Load: %v", err)
	}

	// Every output character must be one of the listed variants.
	allowed := map[string]bool{"X1": true, "X2": true, "X3": true}
	for range 50 {
		got := Convert("x", "silly")
		if !allowed[got] {
			t.Errorf("Convert returned %q, want one of X1/X2/X3", got)
		}
	}
}

func TestGroups_returnsSortedNames(t *testing.T) {
	loadFixture(t)
	got := Groups()
	want := []string{"bubble", "cursive"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Groups() = %v, want %v", got, want)
	}
}
