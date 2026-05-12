// Package pickuplines serves a random pickup line from the embedded
// dataset. The source file pairs each line with a category; this
// package currently exposes only the line text. Add a RandomByCategory
// if filtering becomes useful.
package pickuplines

import (
	"io/fs"

	"github.com/Beamer64/bb_data/internal/pick"
)

type pickupLinesFile struct {
	Data []struct {
		Category string `json:"category"`
		Joke     string `json:"joke"`
	} `json:"data"`
}

var pool []string

// Load reads datasets/pickuplines.json from fsys into memory. Call once
// at startup. After Load returns nil, Random is safe for concurrent use.
func Load(fsys fs.FS) error {
	var f pickupLinesFile
	if err := pick.LoadJSON(fsys, "datasets/pickuplines.json", &f); err != nil {
		return err
	}
	p := make([]string, 0, len(f.Data))
	for _, l := range f.Data {
		p = append(p, l.Joke)
	}
	pool = p
	return nil
}

// Random returns a uniformly random pickup line. Returns "" if Load
// has not been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
