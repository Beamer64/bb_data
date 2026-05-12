// Package facts serves a random fact from the embedded facts dataset.
// Call Load once at startup with the embed.FS from the parent package,
// then call Random from anywhere.
package facts

import (
	"io/fs"

	"github.com/Beamer64/bb_data/internal/pick"
)

var pool []string

// Load reads datasets/facts.txt from fsys into memory. Call once at
// startup. After Load returns nil, Random is safe for concurrent use.
func Load(fsys fs.FS) error {
	lines, err := pick.LoadLines(fsys, "datasets/facts.txt")
	if err != nil {
		return err
	}
	pool = lines
	return nil
}

// Random returns a uniformly random fact. Returns "" if Load has not
// been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
